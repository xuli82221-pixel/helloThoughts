package thoughts

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func mcpTestClient(t *testing.T, server *mcp.Server) *mcp.ClientSession {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	t.Cleanup(cancel)
	ct, st := mcp.NewInMemoryTransports()
	ss, err := server.Connect(ctx, st, nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { ss.Close() })
	client := mcp.NewClient(&mcp.Implementation{Name: "acceptance", Version: "1"}, nil)
	cs, err := client.Connect(ctx, ct, nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { cs.Close() })
	return cs
}
func TestMCPConnectDoesNotInitializeLibrary(t *testing.T) {
	base := filepath.Join(t.TempDir(), "config")
	t.Setenv("MEMO_CONFIG_DIR", base)
	cs := mcpTestClient(t, NewMCPServer(""))
	tools, err := cs.ListTools(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(tools.Tools) != 14 {
		t.Fatalf("unexpected tool count %d", len(tools.Tools))
	}
	for _, tool := range tools.Tools {
		if tool.InputSchema == nil || tool.Annotations == nil {
			t.Fatal("missing schema or annotations", tool.Name)
		}
	}
	if _, err := os.Stat(base); !os.IsNotExist(err) {
		t.Fatal("MCP handshake changed configuration")
	}
}
func TestMCPSetupPreservesConfigurationAndExistingFiles(t *testing.T) {
	requireWindows(t)
	base := t.TempDir()
	t.Setenv("MEMO_CONFIG_DIR", filepath.Join(base, "config"))
	vault := filepath.Join(base, "资料库")
	if _, err := setupMCP("", vault); err != nil {
		t.Fatal(err)
	}
	note := filepath.Join(vault, "Existing.md")
	os.WriteFile(note, []byte("keep"), 0600)
	if _, err := setupMCP("", vault); err != nil {
		t.Fatal(err)
	}
	if _, err := setupMCP("", filepath.Join(base, "Another")); err == nil {
		t.Fatal("setup switched library")
	}
	loc, err := Resolve("")
	if err != nil || loc.Vault != vault {
		t.Fatal("configured library changed", err)
	}
	if body, _ := os.ReadFile(note); string(body) != "keep" {
		t.Fatal("setup changed note")
	}
	os.WriteFile(filepath.Join(base, "config", "config.json"), []byte("broken"), 0600)
	if _, err := setupMCP("", vault); err == nil {
		t.Fatal("setup silently replaced corrupt config")
	}
}
func TestMCPUsesSameServiceAndRefusesStaleWrites(t *testing.T) {
	requireWindows(t)
	t.Setenv("MEMO_CONFIG_DIR", t.TempDir())
	s, _ := testServer(t)
	s.state.Vault = s.store.Vault
	loc, err := Resolve(s.store.Vault)
	if err != nil {
		t.Fatal(err)
	}
	body, _ := json.Marshal(s.state)
	if err = os.WriteFile(filepath.Join(loc.Dir, "runtime.json"), body, 0600); err != nil {
		t.Fatal(err)
	}
	cs := mcpTestClient(t, NewMCPServer(s.store.Vault))
	call := func(name string, args map[string]any) *mcp.CallToolResult {
		t.Helper()
		result, err := cs.CallTool(context.Background(), &mcp.CallToolParams{Name: name, Arguments: args})
		if err != nil {
			t.Fatal(err)
		}
		return result
	}
	created := call("thoughts_create", map[string]any{"path": "Inbox/MCP.md", "content": "真实本地 Markdown"})
	if created.IsError {
		t.Fatal(created.Content)
	}
	note, err := s.store.Read("Inbox/MCP.md")
	if err != nil {
		t.Fatal("MCP did not write real file", err)
	}
	if got := call("thoughts_read", map[string]any{"path": note.Path}); got.IsError {
		t.Fatal(got.Content)
	}
	s.store.Write(note.Path, "external edit", note.Revision, false)
	if got := call("thoughts_update", map[string]any{"path": note.Path, "content": "stale", "revision": note.Revision}); !got.IsError {
		t.Fatal("MCP overwrote external change")
	}
	if got, _ := s.store.Read(note.Path); got.Content != "external edit" {
		t.Fatal("conflict lost data")
	}
	if got := call("thoughts_read", map[string]any{"path": "../outside.md"}); !got.IsError {
		t.Fatal("MCP path escaped vault")
	}
	if got := call("thoughts_open", map[string]any{"browser": "https://external.invalid"}); !got.IsError {
		t.Fatal("arbitrary browser argument allowed")
	}
	if got := call("thoughts_update", map[string]any{"path": note.Path, "content": "missing revision"}); !got.IsError {
		t.Fatal("required revision schema not enforced")
	}
}
