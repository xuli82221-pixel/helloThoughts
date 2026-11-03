package memo

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Read the Windows account profile, not a client's HOME/APPDATA or working dir.
// Tests replace this function with a fresh isolated profile.
var profileDirectory = windowsProfileDirectory

func sharedHome() (string, error) {
	home := os.Getenv("THOUGHTS_HOME")
	if home == "" || strings.EqualFold(home, "auto") || home == "${HOME}/Thoughts" || home == `${HOME}\Thoughts` {
		profile, err := profileDirectory()
		if err != nil {
			return "", fmt.Errorf("cannot locate the Windows account profile: %w", err)
		}
		home = filepath.Join(profile, "Thoughts")
	}
	if !filepath.IsAbs(home) || strings.HasPrefix(home, `\\`) || strings.HasPrefix(home, "//") || filepath.Dir(filepath.Clean(home)) == filepath.Clean(home) {
		return "", fmt.Errorf("THOUGHTS_HOME must be auto or an absolute local non-root directory")
	}
	return filepath.Clean(home), nil
}

func legacyConfigOverride() bool {
	return os.Getenv("THOUGHTS_HOME") == "" && os.Getenv("MEMO_CONFIG_DIR") != ""
}

func automaticHome() bool {
	v := os.Getenv("THOUGHTS_HOME")
	return v == "" || strings.EqualFold(v, "auto") || v == "${HOME}/Thoughts" || v == `${HOME}\Thoughts`
}

// Inspect paths without initializing a library, starting a service or reading notes.
func Paths() (map[string]any, error) {
	base, err := configDir()
	if err != nil {
		return nil, err
	}
	home := ""
	source := "legacy-override"
	if !legacyConfigOverride() {
		home, err = sharedHome()
		if err != nil {
			return nil, err
		}
		source = "explicit-home"
		if automaticHome() {
			source = "windows-user-default"
		}
	}
	var config struct {
		Vault string `json:"vault"`
	}
	b, err := os.ReadFile(filepath.Join(base, "config.json"))
	initialized := err == nil
	if initialized {
		if json.Unmarshal(b, &config) != nil || config.Vault == "" {
			return nil, fmt.Errorf("library configuration is invalid; no files were changed")
		}
	} else if !os.IsNotExist(err) {
		return nil, err
	}
	return map[string]any{"home": home, "config": base, "vault": config.Vault, "initialized": initialized, "source": source}, nil
}
