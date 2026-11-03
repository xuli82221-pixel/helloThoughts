package memo

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSharedHomeAcrossHostProfiles(t *testing.T) {
	requireWindows(t)
	home := filepath.Join(t.TempDir(), "Thoughts")
	t.Setenv("THOUGHTS_HOME", home)
	t.Setenv("MEMO_CONFIG_DIR", filepath.Join(t.TempDir(), "legacy-config"))
	if _, err := setupMCP("", ""); err != nil {
		t.Fatal(err)
	}
	first, err := Resolve("")
	if err != nil || first.Vault != filepath.Join(home, "library") {
		t.Fatal("shared home was not used", err)
	}
	store, err := NewStore(first.Vault)
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	note, err := store.Write("Inbox/shared.md", "Shared local Markdown", "", true)
	if err != nil {
		t.Fatal(err)
	}
	// Host-local Windows profile paths can differ. Explicit home must win.
	t.Setenv("MEMO_CONFIG_DIR", filepath.Join(t.TempDir(), "other-host-config"))
	if _, err = setupMCP("", ""); err != nil {
		t.Fatal(err)
	}
	second, err := Resolve("")
	if err != nil || second != first {
		t.Fatal("host used a separate library or service directory", err)
	}
	if body, err := os.ReadFile(filepath.Join(second.Vault, "Inbox", "shared.md")); err != nil || string(body) != note.Content {
		t.Fatal("reinitialization lost the existing file", err)
	}
	if _, err = setupMCP("", t.TempDir()); err == nil {
		t.Fatal("reinitialization switched libraries")
	}
}

func TestSharedHomeRejectsRelativePath(t *testing.T) {
	t.Setenv("THOUGHTS_HOME", "relative-library")
	if _, err := configDir(); err == nil {
		t.Fatal("relative home would depend on host working directory")
	}
}
