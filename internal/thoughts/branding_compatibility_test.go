package thoughts

import (
	"os"
	"path/filepath"
	"testing"
)

func TestBothManagedExecutableNamesKeepHashAndPathChecks(t *testing.T) {
	for _, name := range []string{"thoughts.exe", "memo.exe"} {
		t.Run(name, func(t *testing.T) {
			home := t.TempDir()
			data := []byte("isolated compatibility fixture")
			p := runtimePin{Version: "0.3.0-preview.20", SHA256: digest(data)}
			p.File = filepath.Join("app", "versions", p.Version+"-"+p.SHA256[:12], name)
			file := filepath.Join(home, p.File)
			if err := os.MkdirAll(filepath.Dir(file), 0700); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(file, data, 0600); err != nil {
				t.Fatal(err)
			}
			if _, err := pinPath(home, p); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(file, []byte("tampered"), 0600); err != nil {
				t.Fatal(err)
			}
			if _, err := pinPath(home, p); err == nil {
				t.Fatal("tampered core accepted")
			}
			p.File = filepath.Join("app", "versions", p.Version+"-"+p.SHA256[:12], "other.exe")
			if _, err := pinPath(home, p); err == nil {
				t.Fatal("unrelated executable name accepted")
			}
		})
	}
}
