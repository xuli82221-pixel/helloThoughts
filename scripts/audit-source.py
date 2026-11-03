"""Read-only release audit. Findings contain locations/rule names, never values.

This is a bounded heuristic scan, not proof that arbitrary secrets are absent.
It scans all reachable/reflog Git blobs and the working tree, including ignored
local artifacts. Dependency/cache directories are inventoried but not read.
"""
import argparse
import hashlib
import json
import os
from pathlib import Path
import re
import subprocess

RULES = {
    "private_key": rb"-----BEGIN (?:RSA |EC |OPENSSH |ENCRYPTED )?PRIVATE KEY-----",
    "github_token": rb"\b(?:gh[pousr]_[A-Za-z0-9]{30,}|github_pat_[A-Za-z0-9_]{30,})\b",
    "api_key": rb"\bsk-[A-Za-z0-9_-]{24,}\b",
    "secret_assignment": rb'''(?i)\b(?:password|passwd|api[_-]?key|client[_-]?secret|mnemonic|seed[_ -]?phrase)\b["']?\s*[:=]\s*["'][^"'\r\n]{8,}["']''',
    "labelled_mnemonic": rb'''(?i)(?:mnemonic|seed[_ -]?phrase)\s*[:=]\s*["']?(?:[a-z]+\s+){11,23}[a-z]+''',
    "opening_ticket": rb"(?:#ticket=|[?&]ticket=)[a-fA-F0-9]{32,}",
    "runtime_credential": rb'"(?:token|csrf)"\s*:\s*"[a-f0-9]{64}"',
    "bearer_credential": rb"Bearer [A-Za-z0-9_.-]{24,}",
    "personal_windows_path": rb"[A-Za-z]:[/\\]{1,4}Users[/\\]{1,4}(?!Public\b|USER\b|YourName\b|YOUR_NAME\b)[A-Za-z0-9_.-]+",
}
SKIP = {".git", "node_modules", ".tools", "__pycache__"}
TEXT = {".md", ".json", ".toml", ".yaml", ".yml", ".txt", ".log", ".go", ".ts", ".tsx", ".js", ".mjs", ".py", ".ps1", ".cmd", ".iss", ".html", ".css", ".env", ".sum", ".mod"}


def run(*args):
    return subprocess.check_output(["git", *args])


def findings(data, location, **extra):
    if b"\0" in data[:4096]:
        return []
    result = []
    for rule, pattern in RULES.items():
        matches = list(re.finditer(pattern, data))
        if matches:
            result.append({"file": location, "rule": rule,
                           "lines": sorted({data.count(b"\n", 0, m.start()) + 1 for m in matches}), **extra})
    return result


def audit(root):
    os.chdir(root)
    historical = []
    objects = run("rev-list", "--objects", "--all", "--reflog").decode().splitlines()
    proc = subprocess.Popen(["git", "cat-file", "--batch"], stdin=subprocess.PIPE, stdout=subprocess.PIPE)
    blobs = 0
    try:
        for entry in objects:
            sha, _, name = entry.partition(" ")
            proc.stdin.write((sha + "\n").encode()); proc.stdin.flush()
            header = proc.stdout.readline().split()
            if len(header) != 3:
                raise RuntimeError("Cannot read Git object")
            size = int(header[2]); data = proc.stdout.read(size); proc.stdout.read(1)
            if header[1] == b"blob":
                blobs += 1
                historical.extend(findings(data, name or "(unnamed blob)", blob=sha))
    finally:
        proc.stdin.close(); proc.wait(timeout=10)
    tracked = set(run("ls-files", "-z").decode().split("\0")) - {""}
    workspace = []; inventory = {"files": 0, "bytes": 0, "scannedTextFiles": 0, "skippedDependencyDirectories": 0, "largeOrBinaryFiles": 0, "symlinksSkipped": 0}
    for directory, dirs, files in os.walk(root, followlinks=False):
        for name in list(dirs):
            p = Path(directory) / name
            if name in SKIP or p.is_symlink():
                dirs.remove(name); inventory["skippedDependencyDirectories"] += 1
        for name in files:
            p = Path(directory) / name
            if p.is_symlink():
                inventory["symlinksSkipped"] += 1; continue
            try:
                size = p.stat().st_size
                inventory["files"] += 1; inventory["bytes"] += size
                if size > 5 * 1024 * 1024 or p.suffix.lower() not in TEXT:
                    inventory["largeOrBinaryFiles"] += 1; continue
                inventory["scannedTextFiles"] += 1
                rel = p.relative_to(root).as_posix()
                workspace.extend(findings(p.read_bytes(), rel, tracked=rel in tracked))
            except OSError:
                raise RuntimeError("Unreadable audit input; audit incomplete") from None
    return {"schemaVersion": 1, "historyCommits": len(run("rev-list", "--all", "--reflog").splitlines()),
            "historyBlobs": blobs, "historyFindings": historical, "workspaceInventory": inventory,
            "workspaceFindings": workspace, "historySafeForPublicPush": not historical,
            "limitations": ["Heuristic patterns do not detect every password or mnemonic.", "Dependencies/caches are inventoried, binary contents need separate package inspection.", "Location-only report; no matched values are emitted."]}


if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument("--root", type=Path, default=Path(__file__).resolve().parent.parent)
    parser.add_argument("--output", type=Path, required=True)
    parser.add_argument("--private-rules", type=Path, help="Optional local-only JSON regex rules; never include this file in a release")
    args = parser.parse_args()
    if args.private_rules:
        RULES.update({name: pattern.encode() for name, pattern in json.loads(args.private_rules.read_text(encoding="utf-8")).items()})
    report = audit(args.root.resolve())
    args.output.parent.mkdir(parents=True, exist_ok=True)
    args.output.write_text(json.dumps(report, ensure_ascii=False, indent=2) + "\n", encoding="utf-8")
    print(json.dumps({"historyCommits": report["historyCommits"], "historyBlobs": report["historyBlobs"],
                      "historyFindingFiles": len({x["file"] for x in report["historyFindings"]}),
                      "historyRules": sorted({x["rule"] for x in report["historyFindings"]}),
                      "workspaceFindingCount": len(report["workspaceFindings"]),
                      "historySafeForPublicPush": report["historySafeForPublicPush"]}))
