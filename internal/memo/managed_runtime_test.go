package memo

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestManagedPinsRejectMissingTamperedAndEscapingExecutables(t *testing.T) {
	home := t.TempDir()
	t.Setenv("THOUGHTS_HOME", home)
	p := runtimePin{Version: "1.0.0", SHA256: digest([]byte("test bytes"))}
	p.File = filepath.Join("app", "versions", p.Version+"-"+p.SHA256[:12], "memo.exe")
	if e := os.MkdirAll(filepath.Dir(filepath.Join(home, p.File)), 0700); e != nil {
		t.Fatal(e)
	}
	if e := os.WriteFile(filepath.Join(home, p.File), []byte("test bytes"), 0600); e != nil {
		t.Fatal(e)
	}
	if _, e := pinPath(home, p); e != nil {
		t.Fatal(e)
	}
	base := filepath.Join(home, "config")
	os.MkdirAll(base, 0700)
	b, _ := json.Marshal(runtimePins{Current: p})
	os.WriteFile(filepath.Join(base, "managed-runtime.json"), b, 0600)
	os.WriteFile(filepath.Join(home, p.File), []byte("changed"), 0600)
	if _, e := readPins(home); e == nil {
		t.Fatal("tampered runtime accepted")
	}
	os.Remove(filepath.Join(home, p.File))
	if forwarded, e := ForwardManaged([]string{"--version"}); forwarded || e == nil {
		t.Fatal("missing managed executable silently fell back")
	}
	p.File = filepath.Join("..", "outside.exe")
	if _, e := pinPath(home, p); e == nil {
		t.Fatal("runtime escape accepted")
	}
}
func TestInitReinstallKeepsConfiguredLibraryAndRejectsSilentSwitch(t *testing.T) {
	requireWindows(t)
	t.Setenv("THOUGHTS_HOME", filepath.Join(t.TempDir(), "home"))
	first, e := Init("")
	if e != nil {
		t.Fatal(e)
	}
	os.WriteFile(filepath.Join(first, "keep.md"), []byte("keep existing data"), 0600)
	again, e := Init("")
	if e != nil || again != first {
		t.Fatal("reinstall changed library", e)
	}
	other := filepath.Join(t.TempDir(), "not-created")
	if _, e = Init(other); e == nil {
		t.Fatal("silent library switch accepted")
	}
	if _, e = os.Stat(other); !os.IsNotExist(e) {
		t.Fatal("conflicting initialization created another library")
	}
	b, e := os.ReadFile(filepath.Join(first, "keep.md"))
	if e != nil || string(b) != "keep existing data" {
		t.Fatal("reinstall modified data")
	}
}

func TestDirectManagedEntryKeepsHomeWithoutEnvironment(t *testing.T) {
	home := t.TempDir()
	b := []byte("managed entry test bytes")
	p := runtimePin{Version: Version, SHA256: digest(b)}
	p.File = filepath.Join("app", "versions", p.Version+"-"+p.SHA256[:12], "memo.exe")
	file := filepath.Join(home, p.File)
	if e := os.MkdirAll(filepath.Dir(file), 0700); e != nil {
		t.Fatal(e)
	}
	if e := os.WriteFile(file, b, 0600); e != nil {
		t.Fatal(e)
	}
	os.MkdirAll(filepath.Join(home, "config"), 0700)
	metadata, _ := json.Marshal(runtimePins{Current: p})
	os.WriteFile(filepath.Join(home, "config", "managed-runtime.json"), metadata, 0600)
	if actual := discoverManagedHome(file); actual != home {
		t.Fatal("direct entry lost its stable library home")
	}
	if actual := discoverManagedHome(filepath.Join(t.TempDir(), "memo.exe")); actual != "" {
		t.Fatal("unmanaged portable program acquired another library")
	}
}
