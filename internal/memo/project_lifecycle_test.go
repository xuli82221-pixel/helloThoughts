package memo

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestProjectArchiveSearchTrashRestoreKeepsAllFiles(t *testing.T) {
	s := testStore(t)
	p, e := s.CreateProject("Lifecycle", "en", false)
	if e != nil {
		t.Fatal(e)
	}
	note, e := s.Write(p.Path+"/note.md", "# Needle\n![image](assets/picture.png)\n", "", true)
	if e != nil {
		t.Fatal(e)
	}
	s.root.MkdirAll(p.Path+"/assets", 0700)
	image := []byte{137, 80, 78, 71, 1, 2, 3}
	os.WriteFile(filepath.Join(s.Vault, filepath.FromSlash(p.Path), "assets/picture.png"), image, 0600)
	snap, e := s.ProjectStatus(p.Path)
	if e != nil || snap.NoteCount != 1 {
		t.Fatal(snap, e)
	}
	archived, e := s.ArchiveProject(p.Path, snap.Revision, true)
	if e != nil || !archived.Archived {
		t.Fatal(archived, e)
	}
	rows, e := s.Search("Needle")
	if e != nil || len(rows) != 0 {
		t.Fatal("archived project leaked into default search", e)
	}
	rows, e = s.SearchScoped("Needle", true)
	if e != nil || len(rows) != 1 {
		t.Fatal("explicit archived search failed", e)
	}
	lib, e := s.Library()
	if e != nil || !lib.Projects[0].Archived {
		t.Fatal(lib, e)
	}
	if n, e := s.Read(note.Path); e != nil || n.Content != note.Content {
		t.Fatal("archive changed Markdown", e)
	}
	entry, e := s.TrashProject(p.Path, archived.Revision)
	if e != nil {
		t.Fatal(e)
	}
	if _, e = s.root.Stat(p.Path); !os.IsNotExist(e) {
		t.Fatal("project not moved to Trash")
	}
	if rows, e := s.SearchScoped("Needle", true); e != nil || len(rows) != 0 {
		t.Fatal("Trash leaked into search", e)
	}
	deleted, e := s.ProjectTrashList()
	if e != nil || len(deleted) != 1 || deleted[0].Revision != entry.Revision {
		t.Fatal(deleted, e)
	}
	restored, e := s.RestoreProject(entry.ID, entry.Revision)
	if e != nil || !restored.Archived {
		t.Fatal(restored, e)
	}
	b, e := os.ReadFile(filepath.Join(s.Vault, filepath.FromSlash(p.Path), "assets/picture.png"))
	if e != nil || string(b) != string(image) {
		t.Fatal("image changed", e)
	}
	n, e := s.Read(note.Path)
	if e != nil || n.Content != note.Content {
		t.Fatal("restored note changed", e)
	}
	if _, e = s.ArchiveProject(p.Path, restored.Revision, false); e != nil {
		t.Fatal(e)
	}
	if rows, e := s.Search("Needle"); e != nil || len(rows) != 1 {
		t.Fatal("unarchive failed", e)
	}
}
func TestProjectRevisionIncludesAttachmentsAndRestoreRefusesCollision(t *testing.T) {
	s := testStore(t)
	p, e := s.CreateProject("Conflict", "en", false)
	if e != nil {
		t.Fatal(e)
	}
	snap, _ := s.ProjectStatus(p.Path)
	os.WriteFile(filepath.Join(s.Vault, filepath.FromSlash(p.Path), "attachment.bin"), []byte("new external attachment"), 0600)
	if _, e = s.TrashProject(p.Path, snap.Revision); e == nil {
		t.Fatal("stale project deletion accepted")
	}
	snap, _ = s.ProjectStatus(p.Path)
	entry, e := s.TrashProject(p.Path, snap.Revision)
	if e != nil {
		t.Fatal(e)
	}
	if _, e = s.CreateProject("Conflict", "en", false); e != nil {
		t.Fatal(e)
	}
	if _, e = s.RestoreProject(entry.ID, entry.Revision); e == nil {
		t.Fatal("restore overwrote existing project")
	}
	if _, e = s.root.Stat(projectTrashDir + "/" + entry.ID + "/attachment.bin"); e != nil {
		t.Fatal("collision lost trashed attachment", e)
	}
}
func TestProjectLifecycleRejectsContainersTraversalAndTamperedMetadata(t *testing.T) {
	s := testStore(t)
	for _, p := range []string{"Projects", "项目", "Inbox", "../outside", "Projects/a/b", ".memo", "assets", ""} {
		if _, e := s.ProjectStatus(p); e == nil {
			t.Fatal("unsafe project accepted", p)
		}
	}
	p, _ := s.CreateProject("Metadata", "en", false)
	snap, _ := s.ProjectStatus(p.Path)
	entry, e := s.TrashProject(p.Path, snap.Revision)
	if e != nil {
		t.Fatal(e)
	}
	entry.OriginalPath = "../outside"
	b, _ := json.Marshal(entry)
	os.WriteFile(filepath.Join(s.Vault, filepath.FromSlash(projectTrashDir), entry.ID+".json"), b, 0600)
	if _, e = s.RestoreProject(entry.ID, entry.Revision); e == nil {
		t.Fatal("tampered restore accepted")
	}
	if _, e = s.root.Stat(projectTrashDir + "/" + entry.ID); e != nil {
		t.Fatal("tampered payload lost", e)
	}
}
func TestProjectLifecycleHTTPRequiresLocalAuthorization(t *testing.T) {
	s, ts := testServer(t)
	for _, route := range []string{"/api/project?path=Projects/a", "/api/project-trash"} {
		if r := request(t, ts, "GET", route, "", nil, nil); r.StatusCode != 401 {
			t.Fatal("unauthorized project endpoint", r.StatusCode)
		}
	}
	p, _ := s.store.CreateProject("HTTP", "en", false)
	snap, _ := s.store.ProjectStatus(p.Path)
	b, _ := json.Marshal(map[string]any{"path": p.Path, "revision": snap.Revision, "archived": true})
	r := request(t, ts, "POST", "/api/projects/archive", string(b), map[string]string{"Authorization": "Bearer " + s.state.Token}, nil)
	if r.StatusCode != 200 {
		t.Fatal("authorized archive failed", r.StatusCode)
	}
}
