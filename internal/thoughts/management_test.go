package thoughts

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestEuropeanProjectTemplatesCreateLocalizedNotes(t *testing.T) {
	for _, entry := range []struct{ language, firstNote string }{
		{"de", "Zusammenfassung"},
		{"fr", "Résumé"},
		{"es", "Resumen"},
	} {
		t.Run(entry.language, func(t *testing.T) {
			s := testStore(t)
			project, err := s.CreateProject("Locale-"+entry.language, entry.language, true)
			if err != nil || len(project.Nodes) != 7 {
				t.Fatalf("localized starter notes missing: %v %v", project, err)
			}
			note, err := s.Read(project.Path + "/" + entry.firstNote + ".md")
			if err != nil || note.Content != "" {
				t.Fatal("localized Markdown file was not created", err)
			}
			if _, err := os.Stat(filepath.Join(s.Vault, filepath.FromSlash(project.Path), "Summary.md")); !os.IsNotExist(err) {
				t.Fatal("localized project silently used the English template")
			}
		})
	}
}

func TestEditableStarterNodesAreNotRecreated(t *testing.T) {
	s := testStore(t)
	project, err := s.CreateProject("SPXIPO", "zh-Hans", true)
	if err != nil || len(project.Nodes) != 7 {
		t.Fatalf("project: %v %v", project, err)
	}
	detail, err := s.Read("Projects/SPXIPO/细节.md")
	if err != nil {
		t.Fatal(err)
	}
	if _, err = s.Rename(detail.Path, "更细", detail.Revision); err != nil {
		t.Fatal(err)
	}
	rule, _ := s.Read("Projects/SPXIPO/规则.md")
	if _, err = s.Trash(rule.Path, rule.Revision); err != nil {
		t.Fatal(err)
	}
	_, err = s.CreateProject("SPXIPO", "zh-Hans", true)
	statusIs(t, err, 409)
	library, err := s.Library()
	if err != nil || len(library.Projects) != 1 || len(library.Projects[0].Nodes) != 6 {
		t.Fatal("template nodes recreated", library, err)
	}
	if _, err = s.Read(detail.Path); !os.IsNotExist(err) {
		statusIs(t, err, 404)
	}
}
func TestRenameKeepsExactContentAndRejectsStaleOrDuplicate(t *testing.T) {
	s := testStore(t)
	body := "---\r\ncustom: keep\r\n---\r\n# Original heading\r\n\r\n한국어 日本語 正文\r\n"
	original, err := s.Write("Projects/Test/细节.md", body, "", true)
	if err != nil {
		t.Fatal(err)
	}
	renamed, err := s.Rename(original.Path, "更细", original.Revision)
	if err != nil {
		t.Fatal(err)
	}
	if renamed.Name != "更细" || renamed.Content != body || renamed.Revision != original.Revision {
		t.Fatal("rename changed body")
	}
	_, err = s.Read(original.Path)
	statusIs(t, err, 404)
	latest, err := s.Write(renamed.Path, body+"external write", renamed.Revision, false)
	if err != nil {
		t.Fatal(err)
	}
	_, err = s.Rename(renamed.Path, "Lost", renamed.Revision)
	statusIs(t, err, 409)
	current, _ := s.Read(renamed.Path)
	if current.Content != latest.Content {
		t.Fatal("concurrent content lost")
	}
	if _, err = s.Write("Projects/Test/Taken.md", "do not overwrite", "", true); err != nil {
		t.Fatal(err)
	}
	_, err = s.Rename(renamed.Path, "taken", latest.Revision)
	statusIs(t, err, 409)
	taken, _ := s.Read("Projects/Test/Taken.md")
	if taken.Content != "do not overwrite" {
		t.Fatal("duplicate overwritten")
	}
}
func TestTrashRestoreAndCollisionPreserveBothFiles(t *testing.T) {
	s := testStore(t)
	note, _ := s.Write("Projects/P/规则.md", "saved rules", "", true)
	entry, err := s.Trash(note.Path, note.Revision)
	if err != nil {
		t.Fatal(err)
	}
	notes, _ := s.List()
	if len(notes) != 0 {
		t.Fatal("trashed content leaked into notes")
	}
	rows, err := s.TrashList()
	if err != nil || len(rows) != 1 || rows[0].OriginalPath != note.Path {
		t.Fatal(rows, err)
	}
	if _, err = s.Write(note.Path, "new independent rules", "", true); err != nil {
		t.Fatal(err)
	}
	_, err = s.Restore(entry.ID, entry.Revision)
	statusIs(t, err, 409)
	rows, _ = s.TrashList()
	if len(rows) != 1 {
		t.Fatal("failed restore lost trash")
	}
	replacement, _ := s.Read(note.Path)
	if replacement.Content != "new independent rules" {
		t.Fatal("restoration overwrote file")
	}
	second, err := s.Trash(replacement.Path, replacement.Revision)
	if err != nil {
		t.Fatal(err)
	}
	restored, err := s.Restore(entry.ID, entry.Revision)
	if err != nil || restored.Content != "saved rules" {
		t.Fatal(restored, err)
	}
	rows, _ = s.TrashList()
	if len(rows) != 1 || rows[0].ID != second.ID {
		t.Fatal("wrong entry restored")
	}
}
func TestInterruptedRenameStagingIsRecoverableAfterReopen(t *testing.T) {
	s := testStore(t)
	note, _ := s.Write("Projects/P/Notes.md", "only copy", "", true)
	entry, err := s.stage(note.Path, note.Revision, "rename")
	if err != nil {
		t.Fatal(err)
	}
	reopened, err := NewStore(s.Vault)
	if err != nil {
		t.Fatal(err)
	}
	defer reopened.Close()
	rows, err := reopened.TrashList()
	if err != nil || len(rows) != 1 {
		t.Fatal(rows, err)
	}
	restored, err := reopened.Restore(entry.ID, entry.Revision)
	if err != nil || restored.Content != "only copy" {
		t.Fatal("interrupted content not recoverable", err)
	}
}
func TestRecoveredRenameDoesNotAliasPublishedNote(t *testing.T) {
	s := testStore(t)
	original, _ := s.Write("Projects/P/Old.md", "preserve published content", "", true)
	entry, err := s.stage(original.Path, original.Revision, "rename")
	if err != nil {
		t.Fatal(err)
	}
	payload, _ := trashPaths(entry.ID)
	if err = s.root.Link(payload, "Projects/P/New.md"); err != nil {
		t.Fatal(err)
	}
	if _, err = s.Restore(entry.ID, entry.Revision); err != nil {
		t.Fatal(err)
	}
	if err = os.WriteFile(filepath.Join(s.Vault, "Projects", "P", "Old.md"), []byte("changed restored content"), 0600); err != nil {
		t.Fatal(err)
	}
	published, err := s.Read("Projects/P/New.md")
	if err != nil || published.Content != "preserve published content" {
		t.Fatal("recovered copy aliased published note", err)
	}
}
func TestManagementPathsAndTrashMetadataCannotEscape(t *testing.T) {
	s := testStore(t)
	for _, name := range []string{"../escape", "CON", "a/b", ".hidden", "NUL", "a:bad"} {
		if _, err := s.CreateProject(name, "en", true); err == nil {
			t.Fatalf("accepted %s", name)
		}
	}
	note, _ := s.Write("Inbox/Note.md", "keep this", "", true)
	entry, err := s.Trash(note.Path, note.Revision)
	if err != nil {
		t.Fatal(err)
	}
	entry.OriginalPath = "../outside.md"
	body, _ := json.Marshal(entry)
	_, metadata := trashPaths(entry.ID)
	if err = os.WriteFile(filepath.Join(s.Vault, filepath.FromSlash(metadata)), body, 0600); err != nil {
		t.Fatal(err)
	}
	_, err = s.Restore(entry.ID, entry.Revision)
	statusIs(t, err, 422)
	payload, _ := trashPaths(entry.ID)
	preserved, err := s.readFile(payload)
	if err != nil || preserved.Content != "keep this" {
		t.Fatal("corrupt metadata caused content loss")
	}
	_, err = s.Restore("../../escape", entry.Revision)
	statusIs(t, err, 400)
}
func TestLibraryRetainsEmptyProjectsAndClassifiesExistingFiles(t *testing.T) {
	s := testStore(t)
	if _, err := s.CreateProject("Empty", "en", false); err != nil {
		t.Fatal(err)
	}
	for _, p := range []string{"Root.md", "Inbox/Quick.md", "项目/Memo/总览.md", "Legacy/Child/Node.md"} {
		if _, err := s.Write(p, "body", "", true); err != nil {
			t.Fatal(err)
		}
	}
	library, err := s.Library()
	if err != nil || len(library.Inbox) != 2 || len(library.Projects) != 3 {
		t.Fatal(library, err)
	}
	for _, project := range library.Projects {
		for _, node := range project.Nodes {
			if node.Content != "" {
				t.Fatal("library leaked full content")
			}
		}
	}
}
func TestManagementAPIAuthorizationAndRevision(t *testing.T) {
	s, ts := testServer(t)
	r := request(t, ts, "POST", "/api/projects", `{"name":"P","language":"en","template":true}`, nil, nil)
	if r.StatusCode != 401 {
		t.Fatal("unauthorized project creation")
	}
	h := map[string]string{"Authorization": "Bearer " + s.state.Token}
	r = request(t, ts, "POST", "/api/projects", `{"name":"P","language":"en","template":true}`, h, nil)
	if r.StatusCode != 201 {
		t.Fatal("project API failed", r.StatusCode)
	}
	r = request(t, ts, "POST", "/api/notes/trash", `{"path":"Projects/P/Details.md","revision":"stale"}`, h, nil)
	if r.StatusCode != 409 {
		t.Fatal("stale deletion accepted")
	}
	r = request(t, ts, "GET", "/api/library", "", h, nil)
	var library Library
	if json.NewDecoder(r.Body).Decode(&library) != nil || len(library.Projects) != 1 || len(library.Projects[0].Nodes) != 7 {
		t.Fatal("library response incorrect")
	}
}
