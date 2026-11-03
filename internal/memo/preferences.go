package memo

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
)

type Preferences struct {
	Language string `json:"language"`
	Theme    string `json:"theme"`
}

func validLanguage(value string) bool {
	switch value {
	case "en", "ko", "ja", "zh-Hans", "zh-Hant", "de", "fr", "es":
		return true
	}
	return false
}

func (s *Server) readPreferences() (Preferences, error) {
	p := Preferences{Language: "en", Theme: "classic"}
	body, err := os.ReadFile(s.preferencesPath)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return p, problem(422, "preferences_unavailable", "Cannot read preferences.")
	}
	if err == nil && (len(body) > 8192 || json.Unmarshal(body, &p) != nil || !validLanguage(p.Language)) {
		return p, problem(422, "preferences_unavailable", "Invalid preferences file.")
	}
	if p.Theme == "" {
		p.Theme = "classic"
	}
	if _, err := s.loadTheme(p.Theme); err != nil {
		p.Theme = "classic"
	}
	return p, nil
}
func (s *Server) preferences(w http.ResponseWriter, r *http.Request) {
	s.preferencesMu.Lock()
	defer s.preferencesMu.Unlock()
	if s.preferencesPath == "" {
		fail(w, problem(503, "preferences_unavailable", "Preferences storage is unavailable."))
		return
	}
	switch r.Method {
	case "GET":
		preferences, err := s.readPreferences()
		if err != nil {
			fail(w, err)
			return
		}
		jsonResponse(w, 200, preferences)
	case "PUT":
		var update struct {
			Language *string `json:"language"`
			Theme    *string `json:"theme"`
		}
		if err := decode(w, r, &update); err != nil {
			fail(w, err)
			return
		}
		preferences, err := s.readPreferences()
		if err != nil {
			fail(w, err)
			return
		}
		if update.Language != nil {
			preferences.Language = *update.Language
		}
		if update.Theme != nil {
			if _, err := s.loadTheme(*update.Theme); err != nil {
				fail(w, themeError())
				return
			}
			preferences.Theme = *update.Theme
		}
		if !validLanguage(preferences.Language) {
			fail(w, problem(400, "invalid_language", "Unsupported interface language."))
			return
		}
		body, _ := json.Marshal(preferences)
		if err := atomicConfig(s.preferencesPath, body); err != nil {
			fail(w, problem(422, "preferences_unavailable", "Cannot save preferences."))
			return
		}
		s.broadcastEvent("preferences")
		jsonResponse(w, 200, preferences)
	default:
		fail(w, problem(405, "method_not_allowed", "Method not allowed."))
	}
}
