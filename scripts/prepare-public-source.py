"""Create an audited working-copy export; never copy or alter Git history.

No GitHub writes, git init/commit/push, credentials or release binaries.
The original private history is deliberately not included because it has findings.
"""
import argparse
import hashlib
import importlib.util
import json
from pathlib import Path
import shutil
import zipfile

ROOT = Path(__file__).resolve().parent.parent
ROOT_FILES = [".gitignore", ".gitattributes", ".github/workflows/ci.yml", "README.md", "CONTRIBUTING.md", "LICENSE", "SECURITY.md", "INSTALL.md", "CHANGELOG.md", "PRIVACY.md", "THIRD_PARTY_NOTICES.txt", "go.mod", "go.sum"]
SOURCE_DIRS = ["cmd", "internal", "web", "adapters", "branding", "licenses", "legal", "release", "scripts"]
DOC_FILES = ["docs/README.md", "docs/WORKSPACE_LAYOUT.md", "docs/DEPENDENCIES.md", "docs/PUBLIC_RELEASE_REVIEW.md", "docs/PUBLIC_SCOPE.md", "docs/SECURITY_AUDIT.md", "docs/WINDOWS_SECURITY_COMPATIBILITY.md", "docs/UI_TABS_IMAGES_ACCEPTANCE.md", "docs/SHARED_LIBRARY_UPGRADE_ACCEPTANCE.md", "docs/PROJECT_WORKSPACE_POLISH.md", "docs/PUBLIC_SOURCE_FILES.md", "docs/WINDOWS_NATIVE_TESTS_20261104.md", "docs/RELEASE_READINESS_20261104.md", "docs/assets/preview.png"]
EXCLUDE_DIRS = {"node_modules", "dist", "__pycache__", ".git", ".artifacts", ".tools"}


def source_files():
    result = {Path(name) for name in ROOT_FILES + DOC_FILES}
    for directory in SOURCE_DIRS:
        for p in (ROOT / directory).rglob("*"):
            if not p.is_file():
                continue
            relative = p.relative_to(ROOT)
            if any(part in EXCLUDE_DIRS for part in relative.parts) or p.suffix in {".exe", ".pyc", ".tsbuildinfo", ".log"}:
                continue
            if relative.as_posix() == "scripts/install-local.ps1":
                continue  # Private legacy desktop-plugin registration, not stdio.
            if p.is_symlink():
                raise ValueError("Source symlink rejected")
            result.add(relative)
    return sorted(result, key=lambda p: p.as_posix())


def prepare(output, private_rules=None):
    if output.exists():
        raise ValueError("Choose a new review directory; previous exports are never overwritten")
    files = source_files()
    inventory = "# Proposed public source files\n\nOnly this reviewed working-copy snapshot is proposed. No original Git objects, local libraries, personal screenshots, logs, caches, credentials, binaries or historical local reports are included. Public upload and history handling require owner approval.\n\n" + "\n".join("- `" + p.as_posix() + "`" for p in files) + "\n"
    (ROOT / "docs/PUBLIC_SOURCE_FILES.md").write_text(inventory, encoding="utf-8")
    spec = importlib.util.spec_from_file_location("audit_source", ROOT / "scripts/audit-source.py")
    audit = importlib.util.module_from_spec(spec); spec.loader.exec_module(audit)
    if private_rules:
        audit.RULES.update({name: pattern.encode() for name, pattern in json.loads(private_rules.read_text(encoding="utf-8")).items()})
    findings = []
    manifest = []
    for relative in files:
        source = ROOT / relative
        if not source.is_file():
            raise ValueError("Missing required source/review file: " + relative.as_posix())
        data = source.read_bytes()
        findings.extend(audit.findings(data, relative.as_posix()))
        manifest.append({"path": relative.as_posix(), "bytes": len(data), "sha256": hashlib.sha256(data).hexdigest()})
    if findings:
        # Locations/rules only; never matched data.
        raise ValueError(json.dumps({"publicationBlocked": True, "findings": findings}))
    source_root = output / "source"
    source_root.mkdir(parents=True)
    for entry in manifest:
        dest = source_root / entry["path"]
        dest.parent.mkdir(parents=True, exist_ok=True)
        shutil.copy2(ROOT / entry["path"], dest)
        assert hashlib.sha256(dest.read_bytes()).hexdigest() == entry["sha256"]
    result = {"target": "https://github.com/xuli82221-pixel/hello-Thoughts-.git", "purpose": "REVIEW_ONLY_NOT_APPROVED_FOR_PUSH",
              "historyIncluded": False, "historyRewritten": False, "fileCount": len(manifest), "sourceFindings": [], "files": manifest}
    (output / "upload-manifest.json").write_text(json.dumps(result, indent=2) + "\n", encoding="utf-8")
    (output / "SHA256SUMS-source.txt").write_text("".join(e["sha256"] + "  " + e["path"] + "\n" for e in manifest), encoding="utf-8")
    archive = output / "thoughts-public-source-review.zip"
    with zipfile.ZipFile(archive, "x", compression=zipfile.ZIP_DEFLATED) as bundle:
        for entry in manifest:
            bundle.write(source_root / entry["path"], entry["path"])
    return {"fileCount": len(manifest), "sourceFindings": 0, "historyIncluded": False, "historyRewritten": False,
            "archiveSha256": hashlib.sha256(archive.read_bytes()).hexdigest(), "pushed": False}


if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument("--output", type=Path, required=True)
    parser.add_argument("--private-rules", type=Path)
    args = parser.parse_args()
    output = args.output.resolve()
    if not output.is_relative_to(ROOT):
        raise ValueError("Review export must stay inside this project")
    print(json.dumps(prepare(output, args.private_rules)))
