"""Exercise a real Windows installer in a NEW explicit temporary home.

Refuses to replace an existing installer registration. No browser, real notes,
AI-client configuration, execution policy or Windows protection is changed.
"""
import argparse
import hashlib
import json
import os
from pathlib import Path
import subprocess
import time
import urllib.request
import uuid
import winreg


UNINSTALL_KEY = r"Software\Microsoft\Windows\CurrentVersion\Uninstall\Thoughts.Desktop.Windows_is1"


def registered():
    try:
        with winreg.OpenKey(winreg.HKEY_CURRENT_USER, UNINSTALL_KEY, 0, winreg.KEY_READ | winreg.KEY_WOW64_64KEY):
            return True
    except FileNotFoundError:
        return False


def run(setup, previous, root):
    if registered():
        raise ValueError("Existing Thoughts installer registration found; use a clean test account")
    root.mkdir(parents=True, exist_ok=False)
    home = root / "home"
    program = root / "program"
    scratch = root / "temp"
    scratch.mkdir()
    env = dict(os.environ, THOUGHTS_HOME=str(home), TEMP=str(scratch), TMP=str(scratch))
    env.pop("MEMO_CONFIG_DIR", None)
    env.pop("THOUGHTS_NO_DELEGATE", None)
    checks = []
    results = {"status": "RUNNING", "checks": checks, "realLibraryTouched": False}
    core = None
    shortcut = None
    shortcut_group = "Thoughts Installer QA " + uuid.uuid4().hex[:12]

    def process(exe, args, ok=True):
        p = subprocess.run([str(exe), *args], env=env, cwd=root, capture_output=True,
                           text=True, encoding="utf-8", errors="replace", timeout=120,
                           creationflags=subprocess.CREATE_NO_WINDOW)
        if ok != (p.returncode == 0):
            raise RuntimeError("Unexpected process exit: " + str(p.returncode) + " (" + exe.name + ")")
        return p

    def install(label, ok=True, icons=False):
        args = ["/VERYSILENT", "/SUPPRESSMSGBOXES", "/NORESTART", "/LANG=en",
                "/DIR=" + str(program), "/LOG=" + str(root / (label + ".log"))]
        args.extend(["/GROUP=" + shortcut_group, "/NoDesktop=true"] if icons else ["/NOICONS"])
        return process(setup, args, ok)

    def cli(*args, ok=True):
        p = process(core, [*args, "--json"], ok)
        return json.loads(p.stdout) if ok else p

    def pins():
        return (home / "config/managed-runtime.json").read_bytes()

    def state():
        return json.loads(next((home / "config/instances").glob("*/runtime.json")).read_text(encoding="utf-8"))

    def api(route, body):
        s = state()
        req = urllib.request.Request(s["url"] + route, data=json.dumps(body).encode(), method="POST",
                                     headers={"Authorization": "Bearer " + s["token"], "Content-Type": "application/json"})
        with urllib.request.urlopen(req, timeout=10) as response:
            return json.load(response)

    def stop():
        try:
            s = state()
            api("/api/editor-state", {"id": "00000000-0000-0000-0000-000000000001", "busy": False, "sequence": 99})
            api("/api/control/shutdown", {"instance": s["instance"]})
        except (OSError, StopIteration):
            pass

    try:
        install("first-install", icons=True)
        assert registered()
        core = next((program / "launchers").glob("*/memo.exe"))
        expected = hashlib.sha256(core.read_bytes()).hexdigest()
        assert json.loads(pins())["current"]["sha256"] == expected
        assert (home / "library").is_dir()
        checks.append("fresh offline install initializes the isolated home and activates the expected core")
        script = ("$w = New-Object -ComObject WScript.Shell; "
                  "$p = Join-Path ([Environment]::GetFolderPath('Programs')) $env:THOUGHTS_TEST_GROUP; "
                  "$f = Join-Path $p 'Thoughts.lnk'; $s = $w.CreateShortcut($f); "
                  "[pscustomobject]@{File=$f; Exists=(Test-Path -LiteralPath $f); Target=$s.TargetPath; Args=$s.Arguments} | ConvertTo-Json")
        env["THOUGHTS_TEST_GROUP"] = shortcut_group
        link = json.loads(process(Path(os.environ["SystemRoot"]) / "System32/WindowsPowerShell/v1.0/powershell.exe",
                                  ["-NoProfile", "-NonInteractive", "-Command", script]).stdout)
        shortcut = Path(link["File"])
        target = Path(link["Target"])
        assert link["Exists"] and target == home / json.loads(pins())["current"]["file"]
        assert link["Args"] == "open --browser system --json"
        unbound_env = env.copy()
        unbound_env.pop("THOUGHTS_HOME", None)
        resolved = subprocess.run([str(target), "paths", "--json"], env=unbound_env, capture_output=True,
                                  text=True, encoding="utf-8", check=True, timeout=20,
                                  creationflags=subprocess.CREATE_NO_WINDOW)
        assert Path(json.loads(resolved.stdout)["home"]) == home
        checks.append("actual Start menu shortcut targets the installed shared core and retains its home without environment overrides")
        note = home / "library/Inbox/installer-test.md"
        image = home / "library/installer-test.png"
        note.parent.mkdir(parents=True, exist_ok=True)
        note.write_text("# Installer fixture\n\nOnly isolated test content.\n", encoding="utf-8")
        image.write_bytes(b"\x89PNG\r\n\x1a\ninstaller-test-fixture")
        config = home / "config/config.json"
        preserved = {file: file.read_bytes() for file in (note, image, config)}

        def check_preserved():
            assert all(file.read_bytes() == content for file, content in preserved.items())

        first_pin = pins()
        install("repeat-install")
        check_preserved()
        assert pins() == first_pin
        checks.append("repeat install preserves the note, image, config and current runtime pin byte-for-byte")
        cli("runtime", "install", "--source", str(core), "--sha256", "0" * 64, ok=False)
        assert pins() == first_pin
        checks.append("incorrect executable hash is rejected before changing the runtime pin")

        if previous:
            old_sha = hashlib.sha256(previous.read_bytes()).hexdigest()
            assert old_sha != expected
            cli("runtime", "install", "--source", str(previous), "--sha256", old_sha)
            old_pin = pins()
            api("/api/editor-state", {"id": "00000000-0000-0000-0000-000000000001", "busy": True, "sequence": 1})
            install("busy-editor-refused", ok=False)
            assert pins() == old_pin
            check_preserved()
            api("/api/editor-state", {"id": "00000000-0000-0000-0000-000000000001", "busy": False, "sequence": 2})
            checks.append("installer refuses an upgrade while the isolated editor is busy; existing pin/data remain")
            install("upgrade")
            assert json.loads(pins())["current"]["sha256"] == expected
            assert json.loads(pins())["previous"]["sha256"] == old_sha
            check_preserved()
            cli("runtime", "rollback")
            assert json.loads(pins())["current"]["sha256"] == old_sha
            check_preserved()
            install("after-rollback")
            check_preserved()
            checks.append("real older-to-current upgrade, runtime rollback and reinstall preserve fixture data/config")

        before_uninstall = pins()
        process(program / "unins000.exe", ["/VERYSILENT", "/SUPPRESSMSGBOXES", "/NORESTART",
                                           "/LOG=" + str(root / "uninstall.log")])
        for _ in range(100):
            if not registered():
                break
            time.sleep(0.1)
        assert not registered()
        assert not core.exists()
        assert not shortcut.exists()
        check_preserved()
        assert pins() == before_uninstall
        checks.append("real uninstall removes launcher/registration and retains shared data/config/managed versions")
        results["status"] = "PASS"
    except Exception as error:
        results.update(status="FAILED_OR_BLOCKED", errorType=type(error).__name__, detail=str(error))
        raise
    finally:
        stop()
        if registered() and (program / "unins000.exe").is_file():
            try:
                process(program / "unins000.exe", ["/VERYSILENT", "/SUPPRESSMSGBOXES", "/NORESTART",
                                                   "/LOG=" + str(root / "cleanup-uninstall.log")])
            except (OSError, RuntimeError, subprocess.TimeoutExpired) as cleanup_error:
                results["cleanupError"] = type(cleanup_error).__name__
        (root / "result.json").write_text(json.dumps(results, indent=2) + "\n", encoding="utf-8")
    return results


if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument("--setup", type=Path, required=True)
    parser.add_argument("--previous", type=Path)
    parser.add_argument("--root", type=Path, required=True)
    args = parser.parse_args()
    print(json.dumps(run(args.setup.resolve(), args.previous.resolve() if args.previous else None, args.root.resolve())))
