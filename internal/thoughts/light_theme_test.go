package thoughts

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestLightThemesOnlyAndLegacyFallback(t *testing.T) {
	if len(builtInThemes()) != 1 || builtInThemes()[0].ID != "classic" {
		t.Fatal("dark bundled theme is still present")
	}
	s, _ := testServer(t)
	s.preferencesPath = filepath.Join(t.TempDir(), "preferences.json")
	files := customThemeFiles(t, "dark-old")
	var tokens map[string]string
	json.Unmarshal(files["tokens.json"], &tokens)
	tokens["background"] = "#14212b"
	files["tokens.json"], _ = json.Marshal(tokens)
	if _, err := s.importTheme(themeZIP(t, files)); err == nil {
		t.Fatal("dark theme import accepted")
	}
	if err := os.WriteFile(s.preferencesPath, []byte(`{"language":"ja","theme":"sample"}`), 0600); err != nil {
		t.Fatal(err)
	}
	p, err := s.readPreferences()
	if err != nil || p.Theme != "classic" || p.Language != "ja" {
		t.Fatal("old dark selection did not fall back without changing language", err)
	}
}
