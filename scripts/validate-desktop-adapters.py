"""Validate local Windows candidates without executing their binaries."""
import hashlib
import json
import os
from pathlib import Path, PurePosixPath
import sys
import zipfile

root = Path(sys.argv[1])
index = json.loads((root / "index.json").read_text(encoding="utf-8-sig"))
assert index["supportedTargets"] == ["windows-amd64"]
assert index["clients"] == ["claude-desktop", "claude-code", "codex"]
assert not index["publicDistributionReady"] and not index["publicSubmission"]
assert len(index["packages"]) == 3
executable_hashes = set()
checked = []
for package in index["packages"]:
    archive = root / package["file"]
    assert hashlib.sha256(archive.read_bytes()).hexdigest() == package["sha256"]
    if package.get("kind") == "offline-current-user-installer":
        info = json.loads((root / "installer-build.json").read_text(encoding="utf-8-sig"))
        assert archive.read_bytes().startswith(b"MZ")
        assert info["file"] == package["file"] and info["version"] == index["version"]
        assert info["installerId"] == "Thoughts.Desktop.Windows"
        assert info["executableSha256"] == package["executableSha256"]
        executable_hashes.add(package["executableSha256"])
        checked.append({"package": package["file"], "status": "PASS", "scope": "structure/hash; execution tested separately"})
        continue
    with zipfile.ZipFile(archive) as bundle:
        assert bundle.testzip() is None
        names = bundle.namelist()
        assert len(names) == len(set(names))
        assert all(not PurePosixPath(n).is_absolute() and ".." not in PurePosixPath(n).parts and "\\" not in n for n in names)
        assert not any(n.lower().endswith(("runtime.json", "config.json", ".log", ".env", ".pfx", ".pem", ".key")) for n in names)
        assert not any(n.startswith(("Inbox/", "Projects/", "library/", ".memo/", "node_modules/", ".git/")) for n in names)
        binary = "server/thoughts.exe" if package["file"].endswith(".mcpb") else "thoughts.exe"
        data = bundle.read(binary)
        assert data.startswith(b"MZ")
        pe = int.from_bytes(data[0x3c:0x40], "little")
        assert data[pe:pe+4] == b"PE\0\0"
        assert int.from_bytes(data[pe+4:pe+6], "little") == 0x8664
        assert int.from_bytes(data[pe+24:pe+26], "little") == 0x20b
        digest = hashlib.sha256(data).hexdigest()
        assert digest == package["executableSha256"]
        executable_hashes.add(digest)
        # trimpath should prevent this development account path from being embedded.
        profile = os.environ.get("USERPROFILE", "")
        if profile:
            for value in (profile, profile.replace("\\", "/")):
                assert value.encode() not in data and value.encode("utf-16-le") not in data, "Workstation path embedded"
        for name in ("LICENSE", "PRIVACY.md", "INSTALL.md", "SECURITY.md", "CHANGELOG.md", "THIRD_PARTY_NOTICES.txt"):
            assert name in names
        notices = bundle.read("THIRD_PARTY_NOTICES.txt")
        assert b"shadcn" in notices and b"github.com/modelcontextprotocol/go-sdk" in notices
        manifest = json.loads(bundle.read("checksums.json"))
        assert manifest["version"] == index["version"]
        for name, expected in manifest["files"].items():
            assert hashlib.sha256(bundle.read(name)).hexdigest() == expected
        if binary.startswith("server/"):
            m = json.loads(bundle.read("manifest.json"))
            assert m["version"] == index["version"]
            assert m["compatibility"]["platforms"] == ["win32"]
            assert m["server"]["type"] == "binary"
            assert m["server"]["mcp_config"] == {"command": "${__dirname}/server/thoughts.exe", "args": ["mcp"], "env": {"THOUGHTS_HOME": "${user_config.thoughts_home}"}}
            assert m["user_config"]["thoughts_home"]["required"] is False
            assert m["user_config"]["thoughts_home"]["default"] == "auto"
            assert m["user_config"]["thoughts_home"]["type"] == "string"
            assert bundle.read(m["icon"]).startswith(b"\x89PNG\r\n\x1a\n")
        else:
            assert "register-mcp.ps1" in names
            assert "configuration/codex.toml.example" in names
            assert "configuration/claude-code.json.example" in names
            assert not any(n.startswith("plugins/") for n in names)
        checked.append({"package": package["file"], "status": "PASS"})
assert len(executable_hashes) == 1, "Clients must ship the same core binary"
print(json.dumps({"packages": checked, "sameCore": True, "status": "PASS"}, indent=2))
