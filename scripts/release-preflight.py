"""Read-only checks of the proposed source allowlist, not release approval."""
import importlib.util
import json
from pathlib import Path, PurePosixPath
import re
from urllib.parse import unquote, urlparse

ROOT = Path(__file__).resolve().parent.parent


def module(name, file):
    spec = importlib.util.spec_from_file_location(name, file)
    loaded = importlib.util.module_from_spec(spec)
    spec.loader.exec_module(loaded)
    return loaded


def safe_relative(name):
    path = PurePosixPath(name)
    forbidden = {".git", "node_modules", ".tools", ".artifacts", "dist", "library", "backups", "Inbox", "Projects", "__pycache__"}
    return (bool(name) and not path.is_absolute() and "\\" not in name and ":" not in name
            and ".." not in path.parts and not any(part in forbidden for part in path.parts)
            and not any(part.startswith(".env") for part in path.parts)
            and path.suffix.lower() not in {".exe", ".log", ".pfx", ".p12", ".pem", ".key", ".pyc"})


def inspect_versions(root):
    files = ["web/package.json", "web/package-lock.json", "adapters/claude-desktop/manifest.json"]
    versions = [json.loads((root / name).read_text(encoding="utf-8-sig"))["version"] for name in files]
    source = (root / "internal/thoughts/runtime.go").read_text(encoding="utf-8")
    match = re.search(r'(?:const|var) Version = "([^"]+)"', source)
    if not match or any(value != match.group(1) for value in versions):
        return [{"rule": "version_mismatch", "file": "package metadata"}]
    return []


def inspect_links(root, files):
    issues = []
    for relative in files:
        if relative.suffix != ".md":
            continue
        text = (root / relative).read_text(encoding="utf-8-sig")
        for link in re.findall(r'!?\[[^\]]*\]\(([^)]+)\)', text):
            target = link.split(' "', 1)[0].strip("<>")
            if target.startswith("#") or urlparse(target).scheme:
                continue
            path = unquote(target.split("#", 1)[0])
            if path and not (root / relative.parent / path).is_file():
                issues.append({"rule": "missing_document_link", "file": relative.as_posix(), "target": path})
    return issues


def inspect(root=ROOT):
    export = module("source_export", root / "scripts/prepare-public-source.py")
    audit = module("source_audit", root / "scripts/audit-source.py")
    files = export.source_files()
    issues = []
    for relative in files:
        if not safe_relative(relative.as_posix()):
            issues.append({"rule": "unsafe_export_path", "file": relative.as_posix()})
            continue
        file = root / relative
        if not file.is_file() or file.is_symlink():
            issues.append({"rule": "missing_or_symlink_source", "file": relative.as_posix()})
            continue
        issues.extend(audit.findings(file.read_bytes(), relative.as_posix()))
    issues.extend(inspect_versions(root))
    if not any(issue["rule"] == "missing_or_symlink_source" for issue in issues):
        issues.extend(inspect_links(root, files))
    for locale in ("", "-ko", "-ja", "-zh-Hans", "-zh-Hant", "-de", "-fr", "-es"):
        if not (root / ("web/public/privacy" + locale + ".html")).is_file():
            issues.append({"rule": "missing_localized_notice", "file": "privacy" + locale + ".html"})
    return {"status": "PASS" if not issues else "BLOCKED", "sourceFileCount": len(files),
            "findings": issues, "historyIncluded": False, "publicDistributionReady": False,
            "scope": "current-source allowlist, versions, notices and document links; no publication or history approval"}


if __name__ == "__main__":
    result = inspect()
    print(json.dumps(result, indent=2))
    raise SystemExit(result["status"] != "PASS")
