"""Apply real publisher metadata; leave unknown public URLs absent."""
import argparse
import json
import re
from pathlib import Path
from urllib.parse import urlparse
from legal_notices import render_notices

ROOT = Path(__file__).resolve().parent.parent


def load(file):
    return json.loads(file.read_text(encoding="utf-8-sig"))


def save(file, data):
    file.write_text(json.dumps(data, ensure_ascii=False, indent=2) + "\n", encoding="utf-8")


def publisher_issues(profile):
    issues = []
    for field in ("name", "copyrightHolder", "supportEmail"):
        if not profile.get(field) or any(ord(c) < 32 for c in profile[field]):
            issues.append("Missing or invalid " + field)
    if not re.fullmatch(r"[^\s@]+@[^\s@]+\.[^\s@]+", profile.get("supportEmail", "")):
        issues.append("Invalid supportEmail")
    for field in ("homepageUrl", "supportPageUrl", "privacyPolicyUrl", "termsOfServiceUrl"):
        value = profile.get(field)
        parsed = urlparse(value or "")
        if not value:
            issues.append("Public HTTPS URL not configured: " + field)
        elif parsed.scheme != "https" or not parsed.hostname or parsed.username or parsed.password or parsed.hostname in ("localhost", "127.0.0.1", "example.com"):
            issues.append("Invalid public HTTPS URL: " + field)
    if not profile.get("contactVerified"):
        issues.append("Publisher contact has not been verified")
    if not profile.get("publicPagesVerified"):
        issues.append("Public pages have not been checked")
    if not profile.get("signingProvider"):
        issues.append("Production signing provider is not configured")
    return issues


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("--check-public", action="store_true")
    args = parser.parse_args()
    profile = load(ROOT / "release/publisher.json")
    if args.check_public:
        issues = publisher_issues(profile)
        print(json.dumps({"publicMetadataReady": not issues, "blockers": issues}, ensure_ascii=False, indent=2))
        raise SystemExit(bool(issues))
    if not profile.get("name") or not re.fullmatch(r"[^\s@]+@[^\s@]+\.[^\s@]+", profile.get("supportEmail", "")):
        raise ValueError("Real publisher name and support email are required")
    author = {"name": profile["name"], "email": profile["supportEmail"]}
    # Legacy plugin manifests are maintained privately when present, but are not
    # dependencies of the public stdio distribution.
    for relative in ("plugins/memo/plugin.json", "plugins/memo/.codex-plugin/plugin.json"):
        file = ROOT / relative
        if not file.is_file():
            continue
        manifest = load(file)
        manifest["author"] = author
        manifest["repository"] = profile["repositoryUrl"]
        if profile.get("homepageUrl"):
            manifest["homepage"] = profile["homepageUrl"]
        save(file, manifest)
    claude_file = ROOT / "adapters/claude-desktop/manifest.json"
    claude = load(claude_file)
    claude["author"] = author
    claude["repository"] = {"type": "git", "url": profile["repositoryUrl"]}
    claude["support"] = profile.get("supportPageUrl") or "mailto:" + profile["supportEmail"]
    if profile.get("homepageUrl"):
        claude["homepage"] = profile["homepageUrl"]
    if profile.get("privacyPolicyUrl"):
        claude["privacy_policies"] = [profile["privacyPolicyUrl"]]
    else:
        claude.pop("privacy_policies", None)
    save(claude_file, claude)
    license_file = ROOT / "LICENSE"
    license_text = re.sub(r"^Copyright \(c\) 2026 .+$", "Copyright (c) 2026 " + profile["copyrightHolder"], license_file.read_text(encoding="utf-8"), flags=re.M)
    license_file.write_text(license_text, encoding="utf-8")
    render_notices(ROOT, profile)
    print("Publisher metadata and local privacy notice synchronized; public readiness checked separately.")


if __name__ == "__main__":
    main()
