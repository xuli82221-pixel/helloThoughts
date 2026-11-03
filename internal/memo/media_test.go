package memo

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
)

func testPNG(t *testing.T) []byte {
	t.Helper()
	var b bytes.Buffer
	m := image.NewRGBA(image.Rect(0, 0, 8, 6))
	m.Set(2, 2, color.RGBA{R: 220, A: 255})
	if err := png.Encode(&b, m); err != nil {
		t.Fatal(err)
	}
	return b.Bytes()
}
func themeZIP(t *testing.T, files map[string][]byte) []byte {
	t.Helper()
	var b bytes.Buffer
	z := zip.NewWriter(&b)
	for name, data := range files {
		f, e := z.Create(name)
		if e != nil {
			t.Fatal(e)
		}
		if _, e = f.Write(data); e != nil {
			t.Fatal(e)
		}
	}
	if err := z.Close(); err != nil {
		t.Fatal(err)
	}
	return b.Bytes()
}
func customThemeFiles(t *testing.T, id string) map[string][]byte {
	t.Helper()
	v := builtInThemes()[0]
	v.ID = id
	v.Name = "Test theme"
	m, _ := json.Marshal(v.ThemeManifest)
	tokens, _ := json.Marshal(v.Tokens)
	return map[string][]byte{"manifest.json": m, "tokens.json": tokens}
}

func TestLocalImagesFormatsRenameTrashAndRestore(t *testing.T) {
	s := testStore(t)
	n, err := s.Write("Projects/P/Details.md", "before\n", "", true)
	if err != nil {
		t.Fatal(err)
	}
	var jpg bytes.Buffer
	jpeg.Encode(&jpg, image.NewRGBA(image.Rect(0, 0, 4, 4)), nil)
	webp, err := os.ReadFile("testdata/gopher.webp")
	if err != nil {
		t.Fatal(err)
	}
	for _, body := range [][]byte{testPNG(t), jpg.Bytes(), webp} {
		asset, e := s.AddImage(n.Path, body)
		if e != nil {
			t.Fatal(e)
		}
		got, mime, e := s.ReadImage(n.Path, asset.Source)
		if e != nil || !bytes.Equal(got, body) || mime != asset.MIME {
			t.Fatal("Image bytes/MIME did not round trip", e)
		}
	}
	a, _ := s.AddImage(n.Path, testPNG(t))
	b, _ := s.AddImage(n.Path, testPNG(t))
	if a.Path == b.Path {
		t.Fatal("Repeated image upload overwrote an attachment")
	}
	n, err = s.Write(n.Path, "before\n\n![local]("+a.Source+")\n\nafter\n", n.Revision, false)
	if err != nil {
		t.Fatal(err)
	}
	renamed, err := s.Rename(n.Path, "More detail", n.Revision)
	if err != nil {
		t.Fatal(err)
	}
	if _, _, err = s.ReadImage(renamed.Path, a.Source); err != nil {
		t.Fatal("Rename broke image", err)
	}
	trash, err := s.Trash(renamed.Path, renamed.Revision)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = os.Stat(filepath.Join(s.Vault, filepath.FromSlash(a.Path))); err != nil {
		t.Fatal("Trash deleted shared attachment")
	}
	restored, err := s.Restore(trash.ID, trash.Revision)
	if err != nil {
		t.Fatal(err)
	}
	if _, _, err = s.ReadImage(restored.Path, a.Source); err != nil {
		t.Fatal("Restore broke image", err)
	}
}
func TestImagesRejectOversizeFakesEscapesAndKeepNote(t *testing.T) {
	s := testStore(t)
	n, _ := s.Write("Inbox/A.md", "untouched", "", true)
	for _, body := range [][]byte{[]byte("<svg onload='x'></svg>"), []byte("<html>fake.png</html>"), testPNG(t)[:30], make([]byte, MaxImageSize+1)} {
		if _, err := s.AddImage(n.Path, body); err == nil {
			t.Fatal("Invalid image accepted")
		}
	}
	for _, source := range []string{"../../outside.png", "../../../x.png", "/absolute.png", "file:///secret.png", "https://example.com/x.png", "//example.com/x.png", "%2e%2e/%2e%2e/outside.png", "..\\outside.png", "assets/x.svg", "assets/x.png?url=secret", ".memo/trash/hidden.png"} {
		if _, _, err := s.ReadImage(n.Path, source); err == nil {
			t.Fatal("Unsafe source accepted", source)
		}
	}
	after, _ := s.read(n.Path)
	if after.Revision != n.Revision {
		t.Fatal("Image rejection modified the note")
	}
}
func TestImageAPIRequiresAuthorizationAndCSRF(t *testing.T) {
	s, ts := testServer(t)
	s.store.Write("Inbox/A.md", "body", "", true)
	if r := request(t, ts, "POST", "/api/images?note=Inbox/A.md", string(testPNG(t)), nil, nil); r.StatusCode != 401 {
		t.Fatal("Anonymous upload accepted")
	}
	h := map[string]string{"Authorization": "Bearer " + s.state.Token}
	r := request(t, ts, "POST", "/api/images?note=Inbox/A.md", string(testPNG(t)), h, nil)
	if r.StatusCode != 201 {
		t.Fatal("Authorized upload failed", r.StatusCode)
	}
	var asset ImageAsset
	json.NewDecoder(r.Body).Decode(&asset)
	if r = request(t, ts, "GET", "/api/images?note=Inbox/A.md&src="+url.QueryEscape(asset.Source), "", nil, nil); r.StatusCode != 401 {
		t.Fatal("Image exposed without session")
	}
	if r = request(t, ts, "GET", "/api/images?note=Inbox/A.md&src="+url.QueryEscape(asset.Source), "", h, nil); r.StatusCode != 200 || r.Header.Get("Content-Type") != "image/png" {
		t.Fatal("Local image route failed")
	}
}
func TestImageAssetsRejectSymbolicLinkEscape(t *testing.T) {
	s := testStore(t)
	s.Write("Projects/P/A.md", "body", "", true)
	outside := t.TempDir()
	link := filepath.Join(s.Vault, "Projects", "P", "assets")
	if runtime.GOOS == "windows" {
		cmd := exec.Command("powershell.exe", "-NoProfile", "-NonInteractive", "-Command", "$ErrorActionPreference='Stop'; $null=New-Item -ItemType Junction -Path $env:THOUGHTS_TEST_LINK -Target $env:THOUGHTS_TEST_TARGET")
		cmd.Env = append(os.Environ(), "THOUGHTS_TEST_LINK="+link, "THOUGHTS_TEST_TARGET="+outside)
		if out, e := cmd.CombinedOutput(); e != nil {
			t.Fatalf("Junction creation failed: %v %s", e, out)
		}
	} else if e := os.Symlink(outside, link); e != nil {
		t.Fatal(e)
	}
	if _, err := s.AddImage("Projects/P/A.md", testPNG(t)); err == nil {
		t.Fatal("Upload followed an escaping link")
	}
	entries, _ := os.ReadDir(outside)
	if len(entries) != 0 {
		t.Fatal("Image escaped library")
	}
}
func TestThemeImportPersistenceFallbackAndNotesUntouched(t *testing.T) {
	s, ts := testServer(t)
	s.preferencesPath = filepath.Join(t.TempDir(), "preferences.json")
	n, _ := s.store.Write("Inbox/A.md", "do not touch", "", true)
	before, _ := os.Stat(filepath.Join(s.store.Vault, "Inbox", "A.md"))
	h := map[string]string{"Authorization": "Bearer " + s.state.Token}
	files := customThemeFiles(t, "custom-test")
	files["preview.png"] = testPNG(t)
	theme, err := s.importTheme(themeZIP(t, files))
	if err != nil || !theme.Preview {
		t.Fatal("Valid theme rejected", err)
	}
	if _, err = s.importTheme(themeZIP(t, files)); err == nil {
		t.Fatal("Theme overwrite accepted")
	}
	if r := request(t, ts, "PUT", "/api/preferences", `{"theme":"custom-test"}`, h, nil); r.StatusCode != 200 {
		t.Fatal("Cannot activate theme", r.StatusCode)
	}
	if r := request(t, ts, "PUT", "/api/preferences", `{"language":"ja"}`, h, nil); r.StatusCode != 200 {
		t.Fatal("Language update failed")
	}
	next := newServer(s.store, s.state)
	next.preferencesPath = s.preferencesPath
	p, err := next.readPreferences()
	if err != nil || p.Theme != "custom-test" || p.Language != "ja" {
		t.Fatal("Theme did not persist independently of language", p, err)
	}
	after, _ := os.Stat(filepath.Join(s.store.Vault, "Inbox", "A.md"))
	note, _ := s.store.read(n.Path)
	if !before.ModTime().Equal(after.ModTime()) || note.Revision != n.Revision {
		t.Fatal("Theme changed note data or mtime")
	}
	os.Remove(filepath.Join(filepath.Dir(s.preferencesPath), "themes", "custom-test.thoughts-theme"))
	p, err = next.readPreferences()
	if err != nil || p.Theme != "classic" || p.Language != "ja" {
		t.Fatal("Missing theme did not fall back", p, err)
	}
}
func TestThemePackagesRejectCodeTraversalBombsAndBadTokens(t *testing.T) {
	s, _ := testServer(t)
	s.preferencesPath = filepath.Join(t.TempDir(), "preferences.json")
	for _, badName := range []string{"../escape.txt", "/absolute", "tokens.css", "script.js", "sub/preview.png"} {
		files := customThemeFiles(t, "safe-theme")
		files[badName] = []byte("bad")
		if _, err := s.importTheme(themeZIP(t, files)); err == nil {
			t.Fatal("Unsafe archive accepted", badName)
		}
	}
	for _, value := range []string{"url(https://example.com/x)", "red; background:url(x)", "expression(alert(1))"} {
		files := customThemeFiles(t, "safe-theme")
		var tokens map[string]string
		json.Unmarshal(files["tokens.json"], &tokens)
		tokens["background"] = value
		files["tokens.json"], _ = json.Marshal(tokens)
		if _, err := s.importTheme(themeZIP(t, files)); err == nil {
			t.Fatal("Unsafe CSS token accepted")
		}
	}
	files := customThemeFiles(t, "safe-theme")
	files["preview.png"] = bytes.Repeat([]byte{'x'}, maxThemeArchive+1)
	if _, err := s.importTheme(themeZIP(t, files)); err == nil {
		t.Fatal("Oversize decompression accepted")
	}
	files = customThemeFiles(t, "classic")
	if _, err := s.importTheme(themeZIP(t, files)); err == nil {
		t.Fatal("Bundled theme overwrite accepted")
	}
	var b bytes.Buffer
	z := zip.NewWriter(&b)
	for name, body := range customThemeFiles(t, "safe-theme") {
		file, _ := z.Create(name)
		file.Write(body)
	}
	header := &zip.FileHeader{Name: "preview.png", Method: zip.Store}
	header.SetMode(os.ModeSymlink | 0777)
	f, _ := z.CreateHeader(header)
	f.Write([]byte("../../file"))
	z.Close()
	if _, err := s.importTheme(b.Bytes()); err == nil {
		t.Fatal("Symlink archive accepted")
	}
	for _, radius := range []string{"00px", "+1px", "1.5px", "17px", "-1px"} {
		files := customThemeFiles(t, "safe-theme")
		var tokens map[string]string
		json.Unmarshal(files["tokens.json"], &tokens)
		tokens["radius"] = radius
		files["tokens.json"], _ = json.Marshal(tokens)
		if _, err := s.importTheme(themeZIP(t, files)); err == nil {
			t.Fatal("Noncanonical radius accepted", radius)
		}
	}
	if _, err := s.importTheme(themeZIP(t, customThemeFiles(t, "con"))); err == nil {
		t.Fatal("Windows reserved theme ID accepted")
	}
}
func TestThemeAPIUninstallResetsActiveWithoutChangingLanguage(t *testing.T) {
	s, ts := testServer(t)
	s.preferencesPath = filepath.Join(t.TempDir(), "preferences.json")
	h := map[string]string{"Authorization": "Bearer " + s.state.Token}
	if r := request(t, ts, "GET", "/api/themes", "", nil, nil); r.StatusCode != 401 {
		t.Fatal("Theme API exposed without authorization")
	}
	s.importTheme(themeZIP(t, customThemeFiles(t, "removable")))
	r := request(t, ts, "PUT", "/api/preferences", `{"language":"ko","theme":"removable"}`, h, nil)
	if r.StatusCode != 200 {
		t.Fatal("Theme selection failed")
	}
	r = request(t, ts, "DELETE", "/api/themes/removable", "", h, nil)
	if r.StatusCode != 200 {
		t.Fatal("Theme removal failed")
	}
	p, _ := s.readPreferences()
	if p.Theme != "classic" || p.Language != "ko" {
		t.Fatal("Theme removal damaged preferences")
	}
	r = request(t, ts, "DELETE", "/api/themes/classic", "", h, nil)
	if r.StatusCode < 400 {
		t.Fatal("Built-in theme removed")
	}
}
func TestExternalBrowserGetsIndependentOneUseSessionAndSharedFiles(t *testing.T) {
	s, ts := testServer(t)
	var opened string
	s.openBrowser = func(value string) error { opened = value; return nil }
	h := map[string]string{"Authorization": "Bearer " + s.state.Token}
	r := request(t, ts, "POST", "/api/ticket", "{}", h, nil)
	var ticket struct {
		Ticket string `json:"ticket"`
	}
	json.NewDecoder(r.Body).Decode(&ticket)
	r = request(t, ts, "POST", "/api/session", `{"ticket":"`+ticket.Ticket+`"}`, map[string]string{"Origin": ts.URL}, nil)
	first := r.Cookies()[0]
	var session struct {
		CSRF string `json:"csrf"`
	}
	json.NewDecoder(r.Body).Decode(&session)
	r = request(t, ts, "POST", "/api/open-browser", "{}", map[string]string{"Origin": ts.URL}, first)
	if r.StatusCode != 403 || opened != "" {
		t.Fatal("Browser launch skipped CSRF")
	}
	r = request(t, ts, "POST", "/api/open-browser", "{}", map[string]string{"Origin": ts.URL, "X-Memo-CSRF": session.CSRF}, first)
	if r.StatusCode != 200 || !strings.HasPrefix(opened, ts.URL+"/#ticket=") {
		t.Fatal("Browser launch did not issue a local ticket")
	}
	u, _ := url.Parse(opened)
	params, _ := url.ParseQuery(u.Fragment)
	payload := `{"ticket":"` + params.Get("ticket") + `"}`
	r = request(t, ts, "POST", "/api/session", payload, map[string]string{"Origin": ts.URL}, nil)
	if r.StatusCode != 200 {
		t.Fatal("Second browser cannot connect")
	}
	second := r.Cookies()[0]
	if second.Value == first.Value {
		t.Fatal("Browsers share a session credential")
	}
	if r = request(t, ts, "POST", "/api/session", payload, map[string]string{"Origin": ts.URL}, nil); r.StatusCode != 401 {
		t.Fatal("Opening ticket could be reused")
	}
	s.store.Write("Inbox/Shared.md", "same local files", "", true)
	for _, cookie := range []*http.Cookie{first, second} {
		r = request(t, ts, "GET", "/api/notes?path=Inbox/Shared.md", "", nil, cookie)
		body, _ := io.ReadAll(r.Body)
		if r.StatusCode != 200 || !bytes.Contains(body, []byte("same local files")) {
			t.Fatal("Browsers are not connected to the same vault")
		}
	}
}

func TestMoveCopiesInlineAndReferenceImagesWithoutChangingMarkdown(t *testing.T) {
	s := testStore(t)
	s.CreateProject("A", "en", false)
	s.CreateProject("B", "en", false)
	n, _ := s.Write("Projects/A/Note.md", "initial", "", true)
	asset, err := s.AddImage(n.Path, testPNG(t))
	if err != nil {
		t.Fatal(err)
	}
	body := "---\r\nsource: test\r\n---\r\n# Exact bytes\r\n\r\n![inline](" + asset.Source + " \"title\")\r\n\r\n![reference][picture]\r\n\r\n[picture]: " + asset.Source + "\r\n\r\n```md\r\n![example](assets/not-real.png)\r\n```\r\n"
	n, err = s.Write(n.Path, body, n.Revision, false)
	if err != nil {
		t.Fatal(err)
	}
	moved, err := s.Move(n.Path, "Projects/B", n.Revision)
	if err != nil {
		t.Fatal(err)
	}
	if moved.Content != body || moved.Revision != n.Revision {
		t.Fatal("Move changed Markdown bytes")
	}
	if _, _, err = s.ReadImage(moved.Path, asset.Source); err != nil {
		t.Fatal("Moved image unavailable", err)
	}
	if _, err = os.Stat(filepath.Join(s.Vault, filepath.FromSlash(asset.Path))); err != nil {
		t.Fatal("Move deleted shared original")
	}
	entry, err := s.Trash(moved.Path, moved.Revision)
	if err != nil {
		t.Fatal(err)
	}
	restored, err := s.Restore(entry.ID, entry.Revision)
	if err != nil {
		t.Fatal(err)
	}
	if _, _, err = s.ReadImage(restored.Path, asset.Source); err != nil {
		t.Fatal("Restored moved image unavailable", err)
	}
}
func TestMoveRejectsImageCollisionsAndStaleRevisions(t *testing.T) {
	s := testStore(t)
	s.CreateProject("A", "en", false)
	s.CreateProject("B", "en", false)
	n, _ := s.Write("Projects/A/Note.md", "initial", "", true)
	asset, _ := s.AddImage(n.Path, testPNG(t))
	n, _ = s.Write(n.Path, "![image]("+asset.Source+")", n.Revision, false)
	destination, _ := resolveImagePath("Projects/B/Note.md", asset.Source)
	s.root.MkdirAll("Projects/B/assets", 0755)
	var jpg bytes.Buffer
	jpeg.Encode(&jpg, image.NewRGBA(image.Rect(0, 0, 2, 2)), nil)
	if err := s.publishImage(destination, jpg.Bytes()); err != nil {
		t.Fatal(err)
	}
	if _, err := s.Move(n.Path, "Projects/B", n.Revision); err == nil {
		t.Fatal("Different image overwritten")
	}
	kept, _ := s.read(n.Path)
	if kept.Revision != n.Revision {
		t.Fatal("Failed move changed source")
	}
	got, _, _ := s.readImageFile(destination)
	if !bytes.Equal(got, jpg.Bytes()) {
		t.Fatal("Collision changed target image")
	}
	s.Write(n.Path, n.Content+"newer", n.Revision, false)
	if _, err := s.Move(n.Path, "Inbox", n.Revision); err == nil {
		t.Fatal("Stale move accepted")
	}
}
func TestExternallyCopiedImagesRefreshAssetVersionWithoutTouchingNote(t *testing.T) {
	s := testStore(t)
	n, _ := s.Write("Inbox/A.md", "![external](assets/external.png)", "", true)
	before, _ := s.Library()
	s.root.MkdirAll("Inbox/assets", 0755)
	if err := os.WriteFile(filepath.Join(s.Vault, "Inbox", "assets", "external.png"), testPNG(t), 0600); err != nil {
		t.Fatal(err)
	}
	after, _ := s.Library()
	if before.AssetVersion == after.AssetVersion {
		t.Fatal("External image not detected")
	}
	if _, _, err := s.ReadImage(n.Path, "assets/external.png"); err != nil {
		t.Fatal("Agent-copied image unavailable", err)
	}
	note, _ := s.read(n.Path)
	if note.Revision != n.Revision || !note.Modified.Equal(n.Modified) {
		t.Fatal("Image refresh changed note")
	}
}

func TestSeparateVaultCookiesDoNotCollideOnLocalhostPorts(t *testing.T) {
	one, oneHTTP := testServer(t)
	two, twoHTTP := testServer(t)
	jar, _ := cookiejar.New(nil)
	for _, item := range []struct {
		s      *Server
		origin string
	}{{one, oneHTTP.URL}, {two, twoHTTP.URL}} {
		ticket := randomID()
		item.s.tickets[ticket] = time.Now().Add(time.Minute)
		res, err := http.Post(item.origin+"/api/session", "application/json", strings.NewReader(`{"ticket":"`+ticket+`"}`))
		if err != nil {
			t.Fatal(err)
		}
		res.Body.Close()
		if res.StatusCode != 403 {
			t.Fatal("Missing Origin unexpectedly accepted")
		}
		req, _ := http.NewRequest("POST", item.origin+"/api/session", strings.NewReader(`{"ticket":"`+ticket+`"}`))
		req.Header.Set("Origin", item.origin)
		res, err = http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		u, _ := url.Parse(item.origin)
		jar.SetCookies(u, res.Cookies())
		res.Body.Close()
		if res.StatusCode != 200 {
			t.Fatal("Session failed")
		}
	}
	if one.sessionCookieName() == two.sessionCookieName() {
		t.Fatal("Vault cookies collide")
	}
	client := &http.Client{Jar: jar}
	for _, origin := range []string{oneHTTP.URL, twoHTTP.URL} {
		res, err := client.Get(origin + "/api/library")
		if err != nil {
			t.Fatal(err)
		}
		res.Body.Close()
		if res.StatusCode != 200 {
			t.Fatal("Opening another vault invalidated existing browser session")
		}
	}
}
