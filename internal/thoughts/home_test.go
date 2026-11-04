package thoughts

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func isolatedDefaultProfile(t *testing.T) string {
	t.Helper()
	profile := t.TempDir()
	previous := profileDirectory
	profileDirectory = func() (string, error) { return profile, nil }
	t.Cleanup(func() { profileDirectory = previous })
	t.Setenv("THOUGHTS_HOME", "")
	t.Setenv("MEMO_CONFIG_DIR", "")
	return profile
}

func TestDefaultHomeIndependentOfHostAndInstallationOrder(t *testing.T) {
	requireWindows(t)
	for _, order := range [][]string{{"codex", "claude-desktop", "claude-code"}, {"claude-desktop", "claude-code", "codex"}} {
		t.Run(order[0]+"-first", func(t *testing.T) {
			profile := isolatedDefaultProfile(t)
			var first Location
			var last Note
			for i, host := range order {
				t.Setenv("APPDATA", filepath.Join(t.TempDir(), host))
				t.Setenv("HOME", filepath.Join(t.TempDir(), "host-home"))
				value := ""
				if host == "claude-desktop" {
					value = "auto"
				}
				t.Setenv("THOUGHTS_HOME", value)
				if _, e := setupMCP("", ""); e != nil {
					t.Fatal(e)
				}
				loc, e := Resolve("")
				if e != nil || loc.Vault != filepath.Join(profile, "Thoughts", "library") {
					t.Fatal("host split the library", e)
				}
				if i == 0 {
					first = loc
				} else if loc != first {
					t.Fatal("host split service metadata")
				}
				store, e := NewStore(loc.Vault)
				if e != nil {
					t.Fatal(e)
				}
				if i > 0 {
					read, e := store.Read("Inbox/shared.md")
					if e != nil || read.Content != last.Content {
						store.Close()
						t.Fatal("later host did not read earlier save", e)
					}
				}
				last, e = store.Write("Inbox/shared.md", "Real isolated file saved by "+host, last.Revision, i == 0)
				store.Close()
				if e != nil {
					t.Fatal(e)
				}
			}
		})
	}
}

func TestDefaultAliasesAndCustomOverride(t *testing.T) {
	profile := isolatedDefaultProfile(t)
	for _, value := range []string{"", "auto", "AUTO", "${HOME}/Thoughts", `${HOME}\Thoughts`} {
		t.Setenv("THOUGHTS_HOME", value)
		home, e := sharedHome()
		if e != nil || home != filepath.Join(profile, "Thoughts") {
			t.Fatal("default alias depends on host substitution", e)
		}
	}
	custom := filepath.Join(t.TempDir(), "chosen-home")
	t.Setenv("THOUGHTS_HOME", custom)
	base, e := configDir()
	if e != nil || base != filepath.Join(custom, "config") {
		t.Fatal("explicit selection ignored", e)
	}
	for _, invalid := range []string{"relative", `\\server\share`, `C:\`} {
		t.Setenv("THOUGHTS_HOME", invalid)
		if _, e := sharedHome(); e == nil {
			t.Fatal("unsafe home accepted")
		}
	}
}

func TestDefaultDoesNotSilentlyAbandonLegacyConfiguration(t *testing.T) {
	requireWindows(t)
	profile := isolatedDefaultProfile(t)
	legacy := filepath.Join(profile, "AppData", "Roaming", "Memo")
	os.MkdirAll(legacy, 0700)
	old := []byte(`{"vault":"old selection retained"}`)
	os.WriteFile(filepath.Join(legacy, "config.json"), old, 0600)
	if _, e := Init(""); e == nil {
		t.Fatal("created an empty default over legacy configuration")
	}
	if _, e := os.Stat(filepath.Join(profile, "Thoughts")); !os.IsNotExist(e) {
		t.Fatal("failed initialization created another home")
	}
	t.Setenv("MEMO_CONFIG_DIR", legacy)
	if base, e := configDir(); e != nil || base != legacy {
		t.Fatal("explicit legacy override lost", e)
	}
	t.Setenv("MEMO_CONFIG_DIR", "")
	home := filepath.Join(profile, "Thoughts")
	t.Setenv("THOUGHTS_HOME", home)
	if _, e := Init(""); e != nil {
		t.Fatal("explicit approved home could not be initialized", e)
	}
	t.Setenv("THOUGHTS_HOME", "")
	if base, e := configDir(); e != nil || base != filepath.Join(home, "config") {
		t.Fatal("existing shared home did not win", e)
	}
	b, _ := os.ReadFile(filepath.Join(legacy, "config.json"))
	if string(b) != string(old) {
		t.Fatal("legacy config modified")
	}
}

func TestPathsIsReadOnlyAndInitUsesOneSharedLock(t *testing.T) {
	requireWindows(t)
	profile := isolatedDefaultProfile(t)
	paths, e := Paths()
	if e != nil || paths["initialized"] != false || paths["home"] != filepath.Join(profile, "Thoughts") {
		t.Fatal(paths, e)
	}
	if _, e := os.Stat(filepath.Join(profile, "Thoughts")); !os.IsNotExist(e) {
		t.Fatal("path inspection wrote to disk")
	}
	base := filepath.Join(profile, "Thoughts", "config")
	if e := privateDir(base); e != nil {
		t.Fatal(e)
	}
	lock, e := acquire(Location{Dir: base}, "init.lock")
	if e != nil {
		t.Fatal(e)
	}
	if _, e := Init(""); e == nil {
		lock.Close()
		t.Fatal("concurrent initialization was allowed")
	}
	lock.Close()
	vault, e := Init("")
	if e != nil {
		t.Fatal(e)
	}
	paths, e = Paths()
	if e != nil || paths["vault"] != vault || paths["source"] != "windows-user-default" {
		t.Fatal(paths, e)
	}
	b, e := os.ReadFile(filepath.Join(base, "config.json"))
	var saved map[string]string
	if e != nil || json.Unmarshal(b, &saved) != nil || saved["vault"] != vault {
		t.Fatal("configuration not persisted")
	}
}
