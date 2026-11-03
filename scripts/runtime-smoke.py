"""Actual local runtime install/upgrade/rollback and concurrent MCP acceptance.

Requires two explicitly supplied reviewed builds with managed-runtime capabilities.
Creates a new isolated home; never points to an existing user library.
"""
import argparse
import hashlib
import json
import os
from pathlib import Path
import queue
import subprocess
import threading
import traceback
import urllib.request


class MCP:
    def __init__(self, exe, env):
        self.proc = subprocess.Popen([str(exe), "mcp"], env=env, stdin=subprocess.PIPE, stdout=subprocess.PIPE, stderr=subprocess.PIPE, text=True, encoding="utf-8", creationflags=getattr(subprocess, "CREATE_NO_WINDOW", 0))
        self.incoming = queue.Queue(); self.seq = 0
        def read():
            for line in self.proc.stdout:
                try: self.incoming.put(json.loads(line))
                except ValueError: self.incoming.put({"invalid": True})
        threading.Thread(target=read, daemon=True).start()
        self.version = self.rpc("initialize", {"protocolVersion": "2025-03-26", "capabilities": {}, "clientInfo": {"name": "isolated-shared-runtime-test", "version": "1"}})["serverInfo"]["version"]
        self.proc.stdin.write('{"jsonrpc":"2.0","method":"notifications/initialized"}\n'); self.proc.stdin.flush()
    def rpc(self, method, params):
        self.seq += 1
        self.proc.stdin.write(json.dumps({"jsonrpc": "2.0", "id": self.seq, "method": method, "params": params}) + "\n"); self.proc.stdin.flush()
        while True:
            message = self.incoming.get(timeout=20)
            assert not message.get("invalid"), "non-protocol stdout"
            if message.get("id") == self.seq:
                assert "error" not in message, "JSON-RPC error"
                return message["result"]
    def tool(self, name, arguments=None, error=False):
        result = self.rpc("tools/call", {"name": name, "arguments": arguments or {}})
        assert bool(result.get("isError")) == error, "unexpected tool error state: " + name
        return None if error else result["structuredContent"]["data"]
    def close(self):
        self.proc.stdin.close(); self.proc.wait(timeout=10)


def run(exe, candidate, root):
    root.mkdir(parents=True, exist_ok=False)
    home = root / "home"; env = dict(os.environ, THOUGHTS_HOME=str(home))
    env.pop("MEMO_CONFIG_DIR", None); env.pop("THOUGHTS_NO_DELEGATE", None)
    checks = []; clients = []; stage = "initialize"
    def cli(*args, content=None, error=False):
        p = subprocess.run([str(exe), *args, "--json"], input=content, text=True, encoding="utf-8", env=env, capture_output=True, timeout=40, creationflags=getattr(subprocess, "CREATE_NO_WINDOW", 0))
        assert (p.returncode != 0) == error, "unexpected CLI result"
        if not p.stdout.strip():
            raise RuntimeError("CLI returned empty JSON output: " + " ".join(args[:2]))
        return json.loads(p.stdout)
    def state():
        value = json.loads(next((home / "config/instances").glob("*/runtime.json")).read_text(encoding="utf-8"))
        assert Path(value["vault"]) == home / "library"
        return value
    def http(method, route, body=None):
        s = state(); data = None if body is None else json.dumps(body).encode()
        request = urllib.request.Request(s["url"] + route, data=data, method=method, headers={"Authorization": "Bearer " + s["token"], "Content-Type": "application/json"})
        with urllib.request.urlopen(request, timeout=15) as r: return json.load(r)
    try:
        initial = cli("init"); original_config = (home / "config/config.json").read_bytes()
        assert cli("init")["vault"] == initial["vault"]
        assert (home / "config/config.json").read_bytes() == original_config
        stage = "create isolated note"
        note = cli("notes", "create", "--path", "Inbox/shared.md", content="# Isolated shared note\n")
        stage = "connect two MCP processes"
        a = MCP(exe, env); clients.append(a); b = MCP(exe, env); clients.append(b)
        first = a.tool("thoughts_read", {"path": note["path"]})
        second = b.tool("thoughts_update", {"path": note["path"], "content": first["content"] + "MCP B update\n", "revision": first["revision"]})
        a.tool("thoughts_update", {"path": note["path"], "content": "stale", "revision": first["revision"]}, error=True)
        web = http("PUT", "/api/notes", {"path": note["path"], "content": second["content"] + "Web API update\n", "revision": second["revision"]})
        assert a.tool("thoughts_read", {"path": note["path"]})["content"] == web["content"]
        assert a.tool("thoughts_status")["pid"] == b.tool("thoughts_status")["pid"]
        reader = {"schema": 1, "tabs": ["note:Inbox/shared.md"], "selected": "Inbox/shared.md", "page": "document", "scrolls": {"note:Inbox/shared.md": 432}}
        http("PUT", "/api/reader-state", reader)
        original_note = (home / "library/Inbox/shared.md").read_bytes()
        checks.append("repeat init; two actual MCP processes and web API share one PID/library; stale write rejected")
        stage = "install"
        first_hash = hashlib.sha256(exe.read_bytes()).hexdigest(); next_hash = hashlib.sha256(candidate.read_bytes()).hexdigest()
        cli("runtime", "install", "--source", str(exe), "--sha256", first_hash)
        first_pins = (home / "config/managed-runtime.json").read_bytes()
        assert cli("runtime", "install", "--source", str(exe), "--sha256", first_hash)["reused"]
        cli("runtime", "upgrade", "--source", str(candidate), "--sha256", "0" * 64, error=True)
        assert (home / "config/managed-runtime.json").read_bytes() == first_pins
        busy = {"id": "12345678-1234-1234-1234-123456789012", "busy": True, "sequence": 1}
        http("POST", "/api/editor-state", busy)
        cli("runtime", "upgrade", "--source", str(candidate), "--sha256", next_hash, error=True)
        assert (home / "config/managed-runtime.json").read_bytes() == first_pins
        http("POST", "/api/editor-state", {**busy, "busy": False, "sequence": 2})
        checks.append("managed install/reinstall; wrong SHA and busy editor leave active pin unchanged")
        stage = "upgrade"
        installed = cli("runtime", "upgrade", "--source", str(candidate), "--sha256", next_hash)
        assert installed["rollbackAvailable"]
        assert a.tool("thoughts_status")["version"] == installed["version"]
        assert a.tool("thoughts_status")["pid"] == b.tool("thoughts_status")["pid"]
        c = MCP(exe, env); clients.append(c)
        assert c.version == installed["version"], "bundled old entry did not delegate"
        assert http("GET", "/api/reader-state")["state"] == reader
        assert (home / "library/Inbox/shared.md").read_bytes() == original_note
        checks.append("real service upgrade; old MCP connections stay compatible; new bundled entry delegates; notes/reader state preserved")
        stage = "rollback"
        rolled_back = cli("runtime", "rollback")
        assert rolled_back["version"] == a.version
        assert c.tool("thoughts_read", {"path": note["path"]})["content"] == web["content"]
        assert http("GET", "/api/reader-state")["state"] == reader
        assert (home / "library/Inbox/shared.md").read_bytes() == original_note
        assert (home / "config/config.json").read_bytes() == original_config
        checks.append("real rollback; newer MCP connection stays compatible; library/config/reader state unchanged")
        result = {"status": "PASS", "checks": checks, "realClientAcceptance": False, "userLibraryTouched": False}
    except Exception as e:
        result = {"status": "FAILED_OR_BLOCKED", "stage": stage, "failureType": type(e).__name__, "failureFrames": [{"function": f.name, "line": f.lineno} for f in traceback.extract_tb(e.__traceback__)], "checks": checks, "realClientAcceptance": False}
    finally:
        for client in clients:
            try: client.close()
            except Exception: pass
        try:
            s = state()
            http("POST", "/api/editor-state", {"id": "12345678-1234-1234-1234-123456789012", "busy": False, "sequence": 99})
            http("POST", "/api/control/shutdown", {"instance": s["instance"]})
        except Exception: pass
    (root / "results.json").write_text(json.dumps(result, indent=2))
    return result

if __name__ == "__main__":
    p = argparse.ArgumentParser(); p.add_argument("--exe", type=Path, required=True); p.add_argument("--candidate", type=Path, required=True); p.add_argument("--root", type=Path, required=True)
    args = p.parse_args(); result = run(args.exe.resolve(), args.candidate.resolve(), args.root.resolve()); print(json.dumps(result)); raise SystemExit(result["status"] != "PASS")
