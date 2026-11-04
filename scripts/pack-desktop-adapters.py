"""Build-only packager: explicit files, portable ZIP modes, no user data."""
import argparse
import hashlib
import json
import re
import shutil
import zipfile
from pathlib import Path
from release_metadata import publisher_issues


def read_json(file):
    return json.loads(file.read_text(encoding="utf-8-sig"))


def write_json(file, value):
    file.parent.mkdir(parents=True, exist_ok=True)
    file.write_text(json.dumps(value, ensure_ascii=False, indent=2) + "\n", encoding="utf-8")


def hash_file(file):
    return hashlib.sha256(file.read_bytes()).hexdigest()


def archive(source, target):
    with zipfile.ZipFile(target, "x", compression=zipfile.ZIP_DEFLATED, compresslevel=6) as bundle:
        for file in sorted(source.rglob("*")):
            if not file.is_file():
                continue
            if file.is_symlink():
                raise ValueError("Package symlink rejected")
            info = zipfile.ZipInfo(file.relative_to(source).as_posix(), (2026, 10, 28, 0, 0, 0))
            info.create_system = 3
            mode = 0o100755 if file.name in ("memo", "thoughts.exe", "install-local.sh") else 0o100644
            info.external_attr = mode << 16
            info.compress_type = zipfile.ZIP_DEFLATED
            bundle.writestr(info, file.read_bytes())


def checksums(folder, version):
    files = {file.relative_to(folder).as_posix(): hash_file(file)
             for file in sorted(folder.rglob("*")) if file.is_file() and file.name != "checksums.json"}
    write_json(folder / "checksums.json", {"version": version, "files": files})
    files["checksums.json"] = hash_file(folder / "checksums.json")
    (folder / "SHA256SUMS.txt").write_text("".join(digest + "  " + name + "\n" for name, digest in files.items()), encoding="utf-8")


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("--bin-dir", type=Path, required=True)
    parser.add_argument("--output", type=Path, required=True)
    parser.add_argument("--setup-file", type=Path, required=True)
    parser.add_argument("--public-release", action="store_true")
    parser.add_argument("--unsigned-preview", action="store_true")
    parser.add_argument("--source-commit", default="")
    parser.add_argument("--build-url", default="")
    args = parser.parse_args()
    if args.unsigned_preview:
        if args.public_release or not re.fullmatch(r"[a-f0-9]{40}", args.source_commit):
            raise ValueError("Unsigned preview requires an exact source commit and cannot be a signed public release")
        if not re.fullmatch(r"https://github.com/xuli82221-pixel/hello-Thoughts-/actions/runs/[0-9]+", args.build_url):
            raise ValueError("Unsigned preview requires its public GitHub Actions build URL")
    root = Path(__file__).resolve().parent.parent
    output = args.output.resolve()
    if output.exists():
        raise ValueError("Choose a new output directory; existing artifacts are never overwritten")
    manifest = read_json(root / "adapters/claude-desktop/manifest.json")
    version = manifest["version"]
    runtime = re.search(r'(?:const|var) Version = "([^"]+)"', (root / "internal/thoughts/runtime.go").read_text(encoding="utf-8"))
    assert runtime and runtime.group(1) == version, "Executable/package version mismatch"
    publisher = read_json(root / "release/publisher.json")
    signatures = read_json(args.bin_dir / "common/signatures.json")
    signature_map = {item["file"]: item for item in signatures}
    binary = args.bin_dir / "windows-amd64/thoughts.exe"
    installer = args.bin_dir / "common/register-mcp.ps1"
    setup_name = "thoughts-" + version + "-windows-amd64-setup.exe"
    setup_info = read_json(args.setup_file.parent / "installer-build.json")
    if (args.setup_file.name != setup_name or setup_info["version"] != version
            or setup_info["installerId"] != "Thoughts.Desktop.Windows"
            or setup_info["executableSha256"] != hash_file(binary)):
        raise ValueError("Installer is not the production-ID build of this exact core")
    for file in (binary, installer, args.setup_file):
        if hash_file(file) != signature_map[file.name]["sha256"]:
            raise ValueError("File changed after signature verification")
    blockers = publisher_issues(publisher)
    blockers += ["Missing valid Authenticode signature: " + name for name in ("thoughts.exe", "register-mcp.ps1", setup_name) if signature_map.get(name, {}).get("status") != "Valid"]
    if args.public_release and blockers:
        raise ValueError("Signed public-candidate checks failed: " + "; ".join(blockers))
    # Producing a signed candidate still does not publish it or attest client acceptance.
    blockers += ["Public source/history review and all three real-client acceptance results require separate approval"]
    output.mkdir(parents=True)
    result = {"version": version, "purpose": "LOCAL_UNSIGNED_RELEASE_CANDIDATE", "publicSubmission": False,
              "platformHostAcceptance": "PENDING", "supportedTargets": ["windows-amd64"],
              "clients": ["claude-desktop", "claude-code", "codex"], "macOS": "NOT_SUPPORTED", "linux": "NOT_SUPPORTED",
              "publicDistributionReady": False, "blockers": blockers, "signatures": signatures, "packages": []}
    if args.unsigned_preview:
        result.update(purpose="PUBLIC_UNSIGNED_PREVIEW", unsignedPreviewAuthorized=True,
                      sourceCommit=args.source_commit, buildUrl=args.build_url, signingStatus="UNSIGNED")
    notice = ("UNSIGNED EXPERIMENTAL PREVIEW / 未签名实验预览版\n"
              "Thoughts has no trusted publisher signature. SignPath Foundation has not approved the application.\n"
              "Windows SmartScreen, SAC or organization policy may block these files. Stop if blocked; do not disable protection or bypass warnings.\n"
              "Not a stable release. Native AI-client and clean-SAC installation acceptance remain incomplete. Back up your library.\n"
              "Thoughts does not encrypt files. Authorized online AI clients can receive note contents under their own policies.\n")
    common = ("LICENSE", "PRIVACY.md", "INSTALL.md", "SECURITY.md", "CHANGELOG.md", "THIRD_PARTY_NOTICES.txt")
    portable = output / "sources/windows-amd64"
    portable.mkdir(parents=True)
    shutil.copy2(binary, portable / "thoughts.exe")
    shutil.copy2(installer, portable / "register-mcp.ps1")
    for name in common:
        shutil.copy2(root / name, portable / name)
    if args.unsigned_preview:
        (portable / "UNSIGNED_NOTICE.txt").write_text(notice, encoding="utf-8")
    # Configuration templates contain no username, credentials or actual vault.
    shutil.copytree(root / "adapters/stdio", portable / "configuration")
    write_json(portable / "PUBLICATION.json", {"publicDistributionReady": False, "blockers": blockers, "signatures": signatures})
    checksums(portable, version)
    target_zip = output / ("thoughts-" + version + "-windows-amd64.zip")
    archive(portable, target_zip)
    claude = output / "sources/claude-desktop-windows-amd64"
    (claude / "server").mkdir(parents=True)
    write_json(claude / "manifest.json", manifest)
    shutil.copy2(binary, claude / "server/thoughts.exe")
    shutil.copytree(root / "adapters/claude-desktop/assets", claude / "assets")
    for name in common:
        shutil.copy2(root / name, claude / name)
    if args.unsigned_preview:
        (claude / "UNSIGNED_NOTICE.txt").write_text(notice, encoding="utf-8")
    checksums(claude, version)
    target_mcpb = output / ("thoughts-claude-desktop-" + version + "-windows-amd64.mcpb")
    archive(claude, target_mcpb)
    for clients, package in ((["claude-code", "codex"], target_zip), (["claude-desktop"], target_mcpb)):
        result["packages"].append({"clients": clients, "target": "windows-amd64", "file": package.name,
                                   "sha256": hash_file(package), "bytes": package.stat().st_size,
                                   "runtimeAccepted": False, "hostAccepted": False, "executableSha256": hash_file(binary)})
    target_setup = output / setup_name
    shutil.copy2(args.setup_file, target_setup)
    shutil.copy2(args.setup_file.parent / "installer-build.json", output / "installer-build.json")
    result["packages"].append({"clients": ["standalone", "claude-code", "codex"], "target": "windows-amd64",
                               "kind": "offline-current-user-installer", "file": target_setup.name,
                               "sha256": hash_file(target_setup), "bytes": target_setup.stat().st_size,
                               "runtimeAccepted": False, "hostAccepted": False, "executableSha256": hash_file(binary)})
    write_json(output / "index.json", result)
    if args.unsigned_preview:
        (output / "UNSIGNED_NOTICE.txt").write_text(notice, encoding="utf-8")
        write_json(output / "BUILDINFO.json", {"product": "Thoughts", "version": version,
                   "sourceCommit": args.source_commit, "buildUrl": args.build_url,
                   "repository": "https://github.com/xuli82221-pixel/hello-Thoughts-",
                   "signingStatus": "UNSIGNED", "signpathApproved": False,
                   "stableReleaseAccepted": False, "executableSha256": hash_file(binary)})
    (output / "SHA256SUMS.txt").write_text("".join(p["sha256"] + "  " + p["file"] + "\n" for p in result["packages"]), encoding="ascii")
    shutil.copy2(root / "INSTALL.md", output / "INSTALL.md")
    print(json.dumps({"output": str(output), "packages": len(result["packages"]), "clients": result["clients"], "hostAcceptance": "PENDING", "publicDistributionReady": False}))


if __name__ == "__main__":
    main()
