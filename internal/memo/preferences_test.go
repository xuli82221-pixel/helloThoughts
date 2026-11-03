package memo

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"testing"
)

func TestLanguagePreferencesAuthorizationAndPersistence(t *testing.T) {
	s, ts := testServer(t)
	s.preferencesPath = filepath.Join(t.TempDir(), "preferences.json")
	if r := request(t, ts, "GET", "/api/preferences", "", nil, nil); r.StatusCode != http.StatusUnauthorized {
		t.Fatal("unauthorized preferences accepted")
	}
	h := map[string]string{"Authorization": "Bearer " + s.state.Token}
	r := request(t, ts, "GET", "/api/preferences", "", h, nil)
	var p Preferences
	if json.NewDecoder(r.Body).Decode(&p) != nil || p.Language != "en" {
		t.Fatal("default language must be English")
	}
	for _, language := range []string{"ko", "ja", "zh-Hans", "zh-Hant", "de", "fr", "es", "en"} {
		r = request(t, ts, "PUT", "/api/preferences", `{"language":"`+language+`"}`, h, nil)
		if r.StatusCode != 200 {
			t.Fatalf("language %s rejected", language)
		}
		// Recreate the server object to verify that preferences survive a service restart.
		next := newServer(s.store, s.state)
		next.preferencesPath = s.preferencesPath
		recorded, err := os.ReadFile(next.preferencesPath)
		if err != nil || json.Unmarshal(recorded, &p) != nil || p.Language != language {
			t.Fatal("language not persisted")
		}
	}
	original, _ := os.ReadFile(s.preferencesPath)
	for _, body := range []string{`{"language":"it"}`, `{"language":"en","path":"outside"}`} {
		r = request(t, ts, "PUT", "/api/preferences", body, h, nil)
		if r.StatusCode != 400 {
			t.Fatal("invalid preferences accepted")
		}
	}
	after, _ := os.ReadFile(s.preferencesPath)
	if string(original) != string(after) {
		t.Fatal("invalid update changed settings")
	}
}
