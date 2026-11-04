"""Publish only bounded test summaries, never raw QA homes, sessions or logs."""
import argparse
import hashlib
import json
from pathlib import Path


def main(packages, qa_root):
    results = {}
    for name, relative in {"stdioMcp": "mcp/results.json", "offlineInstaller": "installer/result.json"}.items():
        result = json.loads((qa_root / relative).read_text(encoding="utf-8"))
        if result.get("status") != "PASS":
            raise ValueError("Preview acceptance did not pass: " + name)
        results[name] = {"status": "PASS", "checks": result["checks"]}
    evidence = {"scope": "GitHub-hosted Windows runner and synthetic test data",
                "checks": results, "realAiClientAcceptance": False,
                "cleanSmartAppControlAcceptance": False, "trustedSignature": False,
                "stableReleaseAccepted": False}
    (packages / "ACCEPTANCE.json").write_text(json.dumps(evidence, indent=2) + "\n", encoding="utf-8")
    files = sorted(p for p in packages.iterdir() if p.is_file() and
                   (p.suffix in {".exe", ".zip", ".mcpb"} or p.name in {"BUILDINFO.json", "UNSIGNED_NOTICE.txt", "ACCEPTANCE.json"}))
    (packages / "SHA256SUMS.txt").write_text("".join(hashlib.sha256(p.read_bytes()).hexdigest() + "  " + p.name + "\n" for p in files), encoding="ascii")


if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument("--packages", type=Path, required=True)
    parser.add_argument("--qa-root", type=Path, required=True)
    args = parser.parse_args()
    main(args.packages.resolve(), args.qa_root.resolve())
