"""Real stdio/CLI/filesystem acceptance, isolated from all user libraries.

Does not constitute acceptance inside Claude Desktop, Claude Code or Codex.
Only the verified service created in this test directory may be terminated.
"""
import argparse
import hashlib
import json
import os
from pathlib import Path
import queue
import signal
import subprocess
import threading
import time
import urllib.parse
import urllib.request
import uuid


def main(exe, root, image):
    root.mkdir(parents=True, exist_ok=False)
    home = root / "home"
    env = dict(os.environ, THOUGHTS_HOME=str(home))
    env.pop("MEMO_CONFIG_DIR", None)
    steps = []
    pid = None
    proc = None

    def cli(*args):
        return json.loads(subprocess.check_output([str(exe), *args, "--json"], env=env, timeout=20, stderr=subprocess.PIPE))

    def stop_owned():
        nonlocal pid
        if pid is None:
            return
        state = cli("doctor")
        if state.get("running") and state.get("pid") == pid and Path(state["vault"]) == home / "library" and state["version"] == version:
            os.kill(pid, signal.SIGTERM)
            for _ in range(20):
                if not cli("doctor").get("running"):
                    break
                time.sleep(.1)
        pid = None

    try:
        version = json.loads(subprocess.check_output([str(exe), "--version"], env=env, timeout=10, stderr=subprocess.PIPE))["version"]
        proc = subprocess.Popen([str(exe), "mcp"], env=env, stdin=subprocess.PIPE, stdout=subprocess.PIPE, stderr=subprocess.PIPE,
                                text=True, encoding="utf-8", creationflags=getattr(subprocess, "CREATE_NO_WINDOW", 0))
        incoming = queue.Queue()

        def reader():
            for line in proc.stdout:
                try:
                    incoming.put(json.loads(line))
                except ValueError:
                    incoming.put({"invalid_protocol_output": True})

        threading.Thread(target=reader, daemon=True).start()
        seq = 0

        def rpc(method, params):
            nonlocal seq
            seq += 1
            proc.stdin.write(json.dumps({"jsonrpc": "2.0", "id": seq, "method": method, "params": params}) + "\n")
            proc.stdin.flush()
            deadline = time.monotonic() + 20
            while time.monotonic() < deadline:
                msg = incoming.get(timeout=max(.1, deadline - time.monotonic()))
                assert not msg.get("invalid_protocol_output"), "stdout contains non-protocol data"
                if msg.get("id") == seq:
                    assert "error" not in msg, "JSON-RPC failed"
                    return msg["result"]
            raise TimeoutError(method)

        def tool(name, args=None, error=False):
            result = rpc("tools/call", {"name": name, "arguments": args or {}})
            assert bool(result.get("isError")) == error, "Unexpected MCP error state: " + name
            if error:
                return None
            return result["structuredContent"]["data"]

        init = rpc("initialize", {"protocolVersion": "2025-03-26", "capabilities": {}, "clientInfo": {"name": "thoughts-isolated-acceptance", "version": "1"}})
        assert init["serverInfo"]["version"] == version
        proc.stdin.write('{"jsonrpc":"2.0","method":"notifications/initialized"}\n'); proc.stdin.flush()
        assert len(rpc("tools/list", {})["tools"]) == 14
        assert not home.exists(), "Handshake initialized a library"
        steps.append("stdio handshake; 14 tools; no implicit initialization")
        tool("thoughts_setup")
        first_config = (home / "config/config.json").read_bytes()
        tool("thoughts_setup")
        assert (home / "config/config.json").read_bytes() == first_config
        steps.append("shared home initialization and repeated setup")
        note = tool("thoughts_create", {"path": "Inbox/acceptance-can-delete.md", "content": "# Thoughts acceptance\n\nIsolated local Markdown integration test.\n"})
        read = tool("thoughts_read", {"path": note["path"]})
        assert (home / "library" / note["path"]).read_text(encoding="utf-8") == read["content"]
        note = tool("thoughts_update", {"path": read["path"], "content": read["content"] + "\nUpdated through real stdio MCP.\n", "revision": read["revision"]})
        tool("thoughts_update", {"path": read["path"], "content": "stale", "revision": read["revision"]}, error=True)
        assert tool("thoughts_read", {"path": note["path"]})["content"] == note["content"]
        assert tool("thoughts_search", {"query": "real stdio MCP"})
        steps.append("real Markdown create/read/update/search and stale revision rejection")
        state = tool("thoughts_status")
        pid = state["pid"]
        assert Path(state["vault"]) == home / "library"
        runtime = json.loads(next((home / "config/instances").glob("*/runtime.json")).read_text())
        assert runtime["pid"] == pid and runtime["url"] == state["url"]
        assert urllib.parse.urlsplit(runtime["url"]).hostname == "127.0.0.1"

        def http(method, route, body=None):
            req = urllib.request.Request(runtime["url"] + route, data=body, method=method,
                                         headers={"Authorization": "Bearer " + runtime["token"], "Content-Type": "application/octet-stream"})
            with urllib.request.urlopen(req, timeout=10) as response:
                return response.read()

        uploaded = json.loads(http("POST", "/api/images?note=" + urllib.parse.quote(note["path"]), image.read_bytes()))
        image_bytes = http("GET", "/api/images?note=" + urllib.parse.quote(note["path"]) + "&src=" + urllib.parse.quote(uploaded["source"]))
        assert hashlib.sha256(image_bytes).digest() == hashlib.sha256(image.read_bytes()).digest()
        note = tool("thoughts_update", {"path": note["path"], "content": note["content"] + "\n![Local acceptance image](" + uploaded["source"] + ")\n", "revision": note["revision"]})
        note = tool("thoughts_rename", {"path": note["path"], "name": "acceptance-renamed", "revision": note["revision"]})
        project = tool("thoughts_create_project", {"name": "Acceptance", "template": False})
        note = tool("thoughts_move", {"path": note["path"], "directory": project["path"], "revision": note["revision"]})
        trashed = tool("thoughts_trash", {"path": note["path"], "revision": note["revision"]})
        assert tool("thoughts_trash_list")
        note = tool("thoughts_restore", {"id": trashed["id"], "revision": trashed["revision"]})
        tool("thoughts_read", {"path": "../outside.md"}, error=True)
        steps.append("real local image upload/read; rename/move/trash/restore; path isolation")
        proc.stdin.close(); proc.wait(timeout=5); proc = None
        assert cli("doctor")["running"], "MCP exit stopped webpage service"
        assert cli("notes", "read", "--path", note["path"])["content"] == note["content"]
        steps.append("MCP process exit leaves owned service and files available")
        stop_owned()
        opening = cli("open", "--browser", "none")
        assert cli("notes", "read", "--path", note["path"])["content"] == note["content"]
        pid = cli("doctor")["pid"]
        steps.append("verified isolated service termination/reopen; Markdown preserved")
        # Credentials remain only inside the ignored acceptance directory.
        (root / "browser-open-private.json").write_text(json.dumps(opening), encoding="utf-8")
        result = {"status": "PASS", "version": version, "checks": steps, "realClientAcceptance": False,
                  "note": note["path"], "serviceLeftForBrowserQA": True}
        (root / "results.json").write_text(json.dumps(result, indent=2), encoding="utf-8")
        pid = None  # Deliberately leave only this test service for browser QA.
        return result
    except Exception as exc:
        result = {"status": "BLOCKED_OR_FAILED", "completedChecks": steps, "failureType": type(exc).__name__, "realClientAcceptance": False}
        (root / "results.json").write_text(json.dumps(result, indent=2), encoding="utf-8")
        return result
    finally:
        if proc is not None:
            proc.stdin.close()
            try:
                proc.wait(timeout=5)
            except subprocess.TimeoutExpired:
                proc.terminate(); proc.wait(timeout=5)
        if pid is not None:
            stop_owned()


if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument("--exe", type=Path, required=True)
    parser.add_argument("--root", type=Path, required=True, help="New, nonexistent isolated directory")
    parser.add_argument("--image", type=Path, required=True)
    args = parser.parse_args()
    result = main(args.exe.resolve(), args.root.resolve(), args.image.resolve())
    print(json.dumps(result))
    raise SystemExit(result["status"] != "PASS")
