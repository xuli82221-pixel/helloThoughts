package thoughts

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestReaderStatePersistsAndRejectsContentAndEscape(t *testing.T) {
	s, ts := testServer(t)
	s.readerPath = filepath.Join(t.TempDir(), "reader-state.json")
	h := map[string]string{"Authorization": "Bearer " + s.state.Token}
	if request(t, ts, "GET", "/api/reader-state", "", nil, nil).StatusCode != 401 {
		t.Fatal("unauthorized state exposed")
	}
	body := `{"schema":1,"tabs":["note:Inbox/中文.md","search"],"selected":"Inbox/中文.md","page":"document","scrolls":{"note:Inbox/中文.md":427}}`
	if request(t, ts, "PUT", "/api/reader-state", body, h, nil).StatusCode != 200 {
		t.Fatal("state save failed")
	}
	next := newServer(s.store, s.state)
	next.readerPath = s.readerPath
	var stored ReaderState
	b, e := os.ReadFile(next.readerPath)
	if e != nil || json.Unmarshal(b, &stored) != nil || stored.Selected != "Inbox/中文.md" || stored.Scrolls["note:Inbox/中文.md"] != 427 {
		t.Fatal("reading state did not persist")
	}
	for _, bad := range []string{`{"schema":1,"tabs":[],"selected":"../private.md","page":"document","scrolls":{}}`, `{"schema":1,"tabs":[],"selected":"","page":"empty","scrolls":{},"draft":"must not be stored"}`} {
		if request(t, ts, "PUT", "/api/reader-state", bad, h, nil).StatusCode != 400 {
			t.Fatal("unsafe reader state accepted")
		}
	}
	after, _ := os.ReadFile(next.readerPath)
	if string(after) != string(b) {
		t.Fatal("invalid input changed saved state")
	}
}
func TestUpgradeShutdownRequiresIdentityAndNoBusyEditors(t *testing.T) {
	s, ts := testServer(t)
	h := map[string]string{"Authorization": "Bearer " + s.state.Token}
	stopped := make(chan bool, 1)
	s.shutdown = func() { stopped <- true }
	id := "01234567-0123-0123-0123-012345678901"
	if request(t, ts, "POST", "/api/editor-state", `{"id":"`+id+`","busy":true,"sequence":2}`, h, nil).StatusCode != 200 {
		t.Fatal("editor lease failed")
	}
	stop := `{"instance":"` + s.state.Instance + `"}`
	if request(t, ts, "POST", "/api/control/shutdown", stop, h, nil).StatusCode != 409 {
		t.Fatal("busy editor did not block upgrade")
	}
	request(t, ts, "POST", "/api/editor-state", `{"id":"`+id+`","busy":false,"sequence":1}`, h, nil)
	if request(t, ts, "POST", "/api/control/shutdown", stop, h, nil).StatusCode != 409 {
		t.Fatal("stale heartbeat cleared busy lease")
	}
	request(t, ts, "POST", "/api/editor-state", `{"id":"`+id+`","busy":false,"sequence":3}`, h, nil)
	if request(t, ts, "POST", "/api/control/shutdown", `{"instance":"wrong"}`, h, nil).StatusCode != 409 {
		t.Fatal("wrong instance accepted")
	}
	if request(t, ts, "POST", "/api/control/shutdown", stop, h, nil).StatusCode != 200 {
		t.Fatal("idle service did not accept shutdown")
	}
	select {
	case <-stopped:
	case <-time.After(time.Second):
		t.Fatal("shutdown not invoked")
	}
}

func TestReaderProjectParentsPersistAndRejectEscape(t *testing.T) {
	s, ts := testServer(t)
	s.readerPath = filepath.Join(t.TempDir(), "reader-state.json")
	h := map[string]string{"Authorization": "Bearer " + s.state.Token}
	body := `{"schema":1,"tabs":["projects","project:Projects/工作"],"selected":"","page":"project:Projects/工作","scrolls":{"project:Projects/工作":128}}`
	if request(t, ts, "PUT", "/api/reader-state", body, h, nil).StatusCode != 200 {
		t.Fatal("project parent state rejected")
	}
	b, e := os.ReadFile(s.readerPath)
	var saved ReaderState
	if e != nil || json.Unmarshal(b, &saved) != nil || saved.Page != "project:Projects/工作" || saved.Scrolls[saved.Page] != 128 {
		t.Fatal("project location not preserved")
	}
	for _, id := range []string{"project:../escape", "project:C:/private", "project:Inbox", "project:Projects", "project:.memo", "project:Projects/CON"} {
		bad := saved
		bad.Page = id
		bad.Tabs = []string{id}
		encoded, _ := json.Marshal(bad)
		if request(t, ts, "PUT", "/api/reader-state", string(encoded), h, nil).StatusCode != 400 {
			t.Fatal("unsafe project state accepted")
		}
	}
	after, _ := os.ReadFile(s.readerPath)
	if string(after) != string(b) {
		t.Fatal("invalid state replaced valid parent location")
	}
}
