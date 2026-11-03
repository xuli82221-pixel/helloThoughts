package memo

import (
	"archive/zip"
	"bytes"
	"embed"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"unicode/utf8"
)

//go:embed themes/*/*.json
var bundledThemes embed.FS

const maxThemeArchive = 3 * 1024 * 1024

var themeID = regexp.MustCompile(`^[a-z][a-z0-9-]{1,39}$`)
var themeVersion = regexp.MustCompile(`^[0-9]+\.[0-9]+\.[0-9]+$`)
var themeColor = regexp.MustCompile(`^#[0-9a-fA-F]{6}([0-9a-fA-F]{2})?$`)
var themeRadius = regexp.MustCompile(`^([0-9]|1[0-6])px$`)
var themeColors = strings.Fields("background foreground card card-foreground popover popover-foreground primary primary-foreground secondary secondary-foreground muted muted-foreground accent accent-foreground destructive border input ring sidebar sidebar-foreground sidebar-primary sidebar-primary-foreground sidebar-accent sidebar-accent-foreground sidebar-border sidebar-ring workspace-tab-background editor-background code-background")

type ThemeManifest struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Version     string `json:"version"`
	ThemeSchema int    `json:"themeSchema"`
}
type Theme struct {
	ThemeManifest
	Tokens  map[string]string `json:"tokens"`
	BuiltIn bool              `json:"builtIn"`
	Preview bool              `json:"preview"`
}
type themePackage struct {
	Theme
	preview     []byte
	previewMIME string
}

func themeError() error {
	return problem(422, "invalid_theme", "Invalid theme package. Only schema 1 JSON tokens and a PNG/WebP preview are supported.")
}
func strictJSON(body []byte, value any) error {
	d := json.NewDecoder(bytes.NewReader(body))
	d.DisallowUnknownFields()
	if err := d.Decode(value); err != nil {
		return err
	}
	var extra any
	if err := d.Decode(&extra); err != io.EOF {
		return themeError()
	}
	return nil
}
func validateTheme(manifest, tokens []byte) (Theme, error) {
	var result Theme
	if len(manifest) > 4096 || len(tokens) > 16384 || strictJSON(manifest, &result.ThemeManifest) != nil || strictJSON(tokens, &result.Tokens) != nil || !themeID.MatchString(result.ID) || validPath(result.ID+".md") != nil || !themeVersion.MatchString(result.Version) || result.ThemeSchema != 1 || strings.TrimSpace(result.Name) != result.Name || result.Name == "" || utf8.RuneCountInString(result.Name) > 60 || strings.IndexFunc(result.Name, func(r rune) bool { return r < 32 }) >= 0 {
		return result, themeError()
	}
	if len(result.Tokens) != len(themeColors)+1 {
		return result, themeError()
	}
	for _, key := range themeColors {
		if !themeColor.MatchString(result.Tokens[key]) {
			return result, themeError()
		}
	}
	if !themeRadius.MatchString(result.Tokens["radius"]) {
		return result, themeError()
	}
	for _, key := range strings.Fields("background card popover sidebar workspace-tab-background editor-background code-background") {
		color := result.Tokens[key]
		if len(color) == 9 && !strings.EqualFold(color[7:], "ff") {
			return result, problem(422, "invalid_theme", "Only opaque light themes are supported.")
		}
		rgb, _ := strconv.ParseUint(color[1:7], 16, 32)
		brightness := float64(rgb>>16)*.2126 + float64((rgb>>8)&255)*.7152 + float64(rgb&255)*.0722
		if brightness < 180 {
			return result, problem(422, "invalid_theme", "Only light themes are supported.")
		}
	}
	return result, nil
}
func builtInThemes() []Theme {
	values := []Theme{}
	for _, id := range []string{"classic"} {
		m, _ := bundledThemes.ReadFile("themes/" + id + "/manifest.json")
		t, _ := bundledThemes.ReadFile("themes/" + id + "/tokens.json")
		v, err := validateTheme(m, t)
		if err != nil {
			panic("Invalid bundled theme: " + id)
		}
		v.BuiltIn = true
		values = append(values, v)
	}
	return values
}
func parseThemeArchive(body []byte) (themePackage, error) {
	var result themePackage
	if len(body) > maxThemeArchive {
		return result, themeError()
	}
	z, err := zip.NewReader(bytes.NewReader(body), int64(len(body)))
	if err != nil || len(z.File) < 2 || len(z.File) > 3 {
		return result, themeError()
	}
	files := map[string][]byte{}
	var total uint64
	for _, entry := range z.File {
		if entry.Name != "manifest.json" && entry.Name != "tokens.json" && entry.Name != "preview.png" && entry.Name != "preview.webp" {
			return result, themeError()
		}
		if !entry.Mode().IsRegular() || entry.Mode()&os.ModeSymlink != 0 || files[entry.Name] != nil {
			return result, themeError()
		}
		limit := uint64(maxThemeArchive)
		if entry.Name == "manifest.json" {
			limit = 4096
		}
		if entry.Name == "tokens.json" {
			limit = 16384
		}
		total += entry.UncompressedSize64
		if entry.UncompressedSize64 > limit || total > 5*1024*1024 {
			return result, themeError()
		}
		f, e := entry.Open()
		if e != nil {
			return result, themeError()
		}
		b, e := io.ReadAll(io.LimitReader(f, int64(limit)+1))
		f.Close()
		if e != nil || uint64(len(b)) > limit {
			return result, themeError()
		}
		files[entry.Name] = b
	}
	result.Theme, err = validateTheme(files["manifest.json"], files["tokens.json"])
	if err != nil {
		return result, err
	}
	for _, name := range []string{"preview.png", "preview.webp"} {
		if b := files[name]; b != nil {
			_, mime, c, e := inspectImage(b)
			if e != nil || c.Width > 2048 || c.Height > 2048 || (name == "preview.png" && mime != "image/png") || (name == "preview.webp" && mime != "image/webp") {
				return result, themeError()
			}
			result.preview = b
			result.previewMIME = mime
			result.Preview = true
		}
	}
	return result, nil
}
func (s *Server) openThemeRoot() (*os.Root, error) {
	if s.preferencesPath == "" {
		return nil, problem(503, "preferences_unavailable", "Theme storage is unavailable.")
	}
	dir := filepath.Join(filepath.Dir(s.preferencesPath), "themes")
	if err := os.Mkdir(dir, 0700); err != nil && !errors.Is(err, os.ErrExist) {
		return nil, storageError(err)
	}
	info, err := os.Lstat(dir)
	if err != nil || !info.IsDir() || info.Mode()&os.ModeSymlink != 0 {
		return nil, problem(422, "invalid_theme", "Theme directory must be a real local directory.")
	}
	return os.OpenRoot(dir)
}
func (s *Server) loadTheme(id string) (themePackage, error) {
	for _, theme := range builtInThemes() {
		if theme.ID == id {
			return themePackage{Theme: theme}, nil
		}
	}
	if !themeID.MatchString(id) {
		return themePackage{}, themeError()
	}
	r, err := s.openThemeRoot()
	if err != nil {
		return themePackage{}, err
	}
	defer r.Close()
	info, err := r.Lstat(id + ".thoughts-theme")
	if err != nil || !info.Mode().IsRegular() {
		return themePackage{}, themeError()
	}
	f, err := r.Open(id + ".thoughts-theme")
	if err != nil {
		return themePackage{}, themeError()
	}
	defer f.Close()
	body, err := io.ReadAll(io.LimitReader(f, maxThemeArchive+1))
	if err != nil {
		return themePackage{}, themeError()
	}
	value, err := parseThemeArchive(body)
	if err != nil || value.ID != id {
		return themePackage{}, themeError()
	}
	return value, nil
}
func (s *Server) themeCatalog() []Theme {
	values := builtInThemes()
	r, err := s.openThemeRoot()
	if err != nil {
		return values
	}
	defer r.Close()
	f, err := r.Open(".")
	if err != nil {
		return values
	}
	defer f.Close()
	entries, err := f.ReadDir(-1)
	if err != nil {
		return values
	}
	for _, entry := range entries {
		if !entry.Type().IsRegular() || !strings.HasSuffix(entry.Name(), ".thoughts-theme") {
			continue
		}
		id := strings.TrimSuffix(entry.Name(), ".thoughts-theme")
		if id == "classic" || id == "sample" {
			continue
		}
		if value, e := s.loadTheme(id); e == nil {
			values = append(values, value.Theme)
		}
	}
	sort.SliceStable(values, func(i, j int) bool {
		if values[i].BuiltIn != values[j].BuiltIn {
			return values[i].BuiltIn
		}
		return values[i].Name < values[j].Name
	})
	return values
}
func (s *Server) importTheme(body []byte) (Theme, error) {
	value, err := parseThemeArchive(body)
	if err != nil {
		return Theme{}, err
	}
	if value.ID == "classic" || value.ID == "sample" {
		return Theme{}, problem(409, "theme_exists", "Bundled themes cannot be replaced.")
	}
	s.themeMu.Lock()
	defer s.themeMu.Unlock()
	r, err := s.openThemeRoot()
	if err != nil {
		return Theme{}, err
	}
	defer r.Close()
	tmp := ".import-" + randomID() + ".tmp"
	f, err := r.OpenFile(tmp, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
	if err != nil {
		return Theme{}, storageError(err)
	}
	defer r.Remove(tmp)
	if _, err = f.Write(body); err == nil {
		err = f.Sync()
	}
	closeErr := f.Close()
	if err != nil {
		return Theme{}, storageError(err)
	}
	if closeErr != nil {
		return Theme{}, storageError(closeErr)
	}
	if err = r.Link(tmp, value.ID+".thoughts-theme"); err != nil {
		if errors.Is(err, os.ErrExist) {
			return Theme{}, problem(409, "theme_exists", "A theme with this ID is already installed.")
		}
		return Theme{}, storageError(err)
	}
	return value.Theme, nil
}
func (s *Server) registerThemes(mux *http.ServeMux) {
	mux.Handle("GET /api/themes/sample/package", s.authorized(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		value := builtInThemes()[0]
		value.ID = "sample-light"
		value.Name = "Sample Light (import)"
		manifest, _ := json.Marshal(value.ThemeManifest)
		tokens, _ := json.Marshal(value.Tokens)
		var body bytes.Buffer
		z := zip.NewWriter(&body)
		for name, data := range map[string][]byte{"manifest.json": manifest, "tokens.json": tokens} {
			f, _ := z.Create(name)
			f.Write(data)
		}
		z.Close()
		w.Header().Set("Content-Type", "application/zip")
		w.Header().Set("Content-Disposition", `attachment; filename="sample-light.thoughts-theme"`)
		w.Write(body.Bytes())
	})))
	mux.Handle("GET /api/themes", s.authorized(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { jsonResponse(w, 200, s.themeCatalog()) })))
	mux.Handle("POST /api/themes", s.authorized(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, maxThemeArchive)
		body, err := io.ReadAll(r.Body)
		if err != nil {
			fail(w, themeError())
			return
		}
		value, err := s.importTheme(body)
		if err != nil {
			fail(w, err)
			return
		}
		jsonResponse(w, 201, value)
	})))
	mux.Handle("GET /api/themes/{id}/preview", s.authorized(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		value, err := s.loadTheme(r.PathValue("id"))
		if err != nil || !value.Preview {
			fail(w, problem(404, "not_found", "No preview is available."))
			return
		}
		w.Header().Set("Content-Type", value.previewMIME)
		w.Write(value.preview)
	})))
	mux.Handle("DELETE /api/themes/{id}", s.authorized(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if !themeID.MatchString(id) || id == "classic" || id == "sample" {
			fail(w, themeError())
			return
		}
		s.themeMu.Lock()
		defer s.themeMu.Unlock()
		root, err := s.openThemeRoot()
		if err != nil {
			fail(w, err)
			return
		}
		defer root.Close()
		s.preferencesMu.Lock()
		defer s.preferencesMu.Unlock()
		p, err := s.readPreferences()
		if err != nil {
			fail(w, err)
			return
		}
		if p.Theme == id {
			p.Theme = "classic"
			body, _ := json.Marshal(p)
			if err = atomicConfig(s.preferencesPath, body); err != nil {
				fail(w, storageError(err))
				return
			}
		}
		if err = root.Remove(id + ".thoughts-theme"); err != nil {
			fail(w, storageError(err))
			return
		}
		s.broadcastEvent("preferences")
		jsonResponse(w, 200, map[string]bool{"removed": true})
	})))
}
