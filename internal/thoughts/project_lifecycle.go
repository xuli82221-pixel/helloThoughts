package thoughts

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"sort"
	"strings"
	"time"
)

const projectMarker = ".memo-project.json"
const projectTrashDir = ".memo/project-trash"

type ProjectSnapshot struct {
	Path      string `json:"path"`
	Name      string `json:"name"`
	Revision  string `json:"revision"`
	NoteCount int    `json:"noteCount"`
	Archived  bool   `json:"archived"`
}
type ProjectTrashEntry struct {
	ID           string    `json:"id"`
	OriginalPath string    `json:"originalPath"`
	Name         string    `json:"name"`
	DeletedAt    time.Time `json:"deletedAt"`
	Revision     string    `json:"revision"`
	NoteCount    int       `json:"noteCount"`
}

func validProjectPath(p string) bool {
	if validPath(p+"/check.md") != nil {
		return false
	}
	parts := strings.Split(p, "/")
	if len(parts) == 2 {
		return parts[0] == "Projects" || parts[0] == "项目"
	}
	return len(parts) == 1 && !strings.EqualFold(p, "Inbox") && p != "Projects" && p != "项目" && !strings.EqualFold(p, "assets")
}
func (s *Store) projectArchived(p string) (bool, error) {
	f, e := s.root.Open(p + "/" + projectMarker)
	if errors.Is(e, os.ErrNotExist) {
		return false, nil
	}
	if e != nil {
		return false, storageError(e)
	}
	defer f.Close()
	st, e := f.Stat()
	if e != nil || !st.Mode().IsRegular() || st.Size() > 1024 {
		return false, problem(422, "file_unavailable", "Project metadata is invalid; files were kept.")
	}
	var value struct {
		Schema   int  `json:"schema"`
		Archived bool `json:"archived"`
	}
	b, e := io.ReadAll(io.LimitReader(f, 1025))
	if e != nil || strictJSON(b, &value) != nil || value.Schema != 1 {
		return false, problem(422, "file_unavailable", "Project metadata is invalid; files were kept.")
	}
	return value.Archived, nil
}

// Snapshot every regular file, including images, before a whole-project mutation.
// Reparse points and links are refused, rather than silently omitting their data.
func (s *Store) projectSnapshotAt(p string) (ProjectSnapshot, error) {
	info, e := s.root.Lstat(p)
	if e != nil {
		return ProjectSnapshot{}, storageError(e)
	}
	if !info.IsDir() || info.Mode()&os.ModeSymlink != 0 {
		return ProjectSnapshot{}, problem(422, "file_unavailable", "Project must be a real local directory.")
	}
	h := sha256.New()
	count, files := 0, 0
	var total int64
	e = fs.WalkDir(s.root.FS(), p, func(file string, d fs.DirEntry, e error) error {
		if e != nil {
			return storageError(e)
		}
		if d.Type()&os.ModeSymlink != 0 {
			return problem(422, "file_unavailable", "Project links must be removed before this operation; files were kept.")
		}
		rel := strings.TrimPrefix(file, p)
		if d.IsDir() {
			fmt.Fprintf(h, "dir:%q\n", rel)
			return nil
		}
		files++
		if files > 10000 {
			return problem(422, "project_too_large", "Project has too many files for this operation.")
		}
		f, e := s.root.Open(file)
		if e != nil {
			return storageError(e)
		}
		defer f.Close()
		st, e := f.Stat()
		if e != nil || !st.Mode().IsRegular() {
			return problem(422, "file_unavailable", "Project contains an unsupported file.")
		}
		total += st.Size()
		if total > 1<<30 {
			return problem(422, "project_too_large", "Project exceeds the 1 GiB operation limit.")
		}
		sum := sha256.New()
		if _, e = io.Copy(sum, f); e != nil {
			return storageError(e)
		}
		fmt.Fprintf(h, "file:%q:%x\n", rel, sum.Sum(nil))
		if strings.EqualFold(path.Ext(file), ".md") {
			count++
		}
		return nil
	})
	if e != nil {
		return ProjectSnapshot{}, e
	}
	archived, e := s.projectArchived(p)
	if e != nil {
		return ProjectSnapshot{}, e
	}
	return ProjectSnapshot{Path: p, Name: path.Base(p), Revision: hex.EncodeToString(h.Sum(nil)), NoteCount: count, Archived: archived}, nil
}
func (s *Store) ProjectStatus(p string) (ProjectSnapshot, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !validProjectPath(p) {
		return ProjectSnapshot{}, problem(400, "invalid_path", "Choose one project, not its container.")
	}
	return s.projectSnapshotAt(p)
}
func (s *Store) ArchiveProject(p, revision string, archived bool) (ProjectSnapshot, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !validProjectPath(p) {
		return ProjectSnapshot{}, problem(400, "invalid_path", "Invalid project path.")
	}
	old, e := s.projectSnapshotAt(p)
	if e != nil {
		return old, e
	}
	if revision == "" || revision != old.Revision {
		return old, problem(409, "revision_conflict", "Project files changed; review the project before retrying.")
	}
	b, _ := json.Marshal(map[string]any{"schema": 1, "archived": archived})
	tmp := p + "/.project-" + randomID() + ".tmp"
	f, e := s.root.OpenFile(tmp, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
	if e != nil {
		return old, storageError(e)
	}
	_, e = f.Write(b)
	if e == nil {
		e = f.Sync()
	}
	f.Close()
	defer s.root.Remove(tmp)
	if e != nil {
		return old, storageError(e)
	}
	if e = s.root.Rename(tmp, p+"/"+projectMarker); e != nil {
		return old, storageError(e)
	}
	return s.projectSnapshotAt(p)
}
func (s *Store) TrashProject(p, revision string) (ProjectTrashEntry, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !validProjectPath(p) {
		return ProjectTrashEntry{}, problem(400, "invalid_path", "Invalid project path.")
	}
	snap, e := s.projectSnapshotAt(p)
	if e != nil {
		return ProjectTrashEntry{}, e
	}
	if revision == "" || snap.Revision != revision {
		return ProjectTrashEntry{}, problem(409, "revision_conflict", "Project files changed; review before deleting.")
	}
	if e = s.ensureTrash(); e != nil {
		return ProjectTrashEntry{}, e
	}
	if e = s.root.MkdirAll(projectTrashDir, 0700); e != nil {
		return ProjectTrashEntry{}, storageError(e)
	}
	st, e := s.root.Lstat(projectTrashDir)
	if e != nil || !st.IsDir() || st.Mode()&os.ModeSymlink != 0 {
		return ProjectTrashEntry{}, problem(422, "file_unavailable", "Project Trash must be a real directory.")
	}
	entry := ProjectTrashEntry{ID: randomID(), OriginalPath: p, Name: snap.Name, DeletedAt: time.Now().UTC(), Revision: snap.Revision, NoteCount: snap.NoteCount}
	b, _ := json.Marshal(entry)
	meta := projectTrashDir + "/" + entry.ID + ".json"
	f, e := s.root.OpenFile(meta, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
	if e != nil {
		return entry, storageError(e)
	}
	_, e = f.Write(b)
	if e == nil {
		e = f.Sync()
	}
	f.Close()
	if e != nil {
		return entry, storageError(e)
	}
	if e = s.root.Rename(p, projectTrashDir+"/"+entry.ID); e != nil {
		s.root.Remove(meta)
		return entry, storageError(e)
	}
	return entry, nil
}
func (s *Store) projectTrashEntry(id string) (ProjectTrashEntry, error) {
	if !validTrashID(id) {
		return ProjectTrashEntry{}, problem(400, "invalid_path", "Invalid project Trash entry.")
	}
	f, e := s.root.Open(projectTrashDir + "/" + id + ".json")
	if e != nil {
		return ProjectTrashEntry{}, storageError(e)
	}
	defer f.Close()
	b, e := io.ReadAll(io.LimitReader(f, 8193))
	var entry ProjectTrashEntry
	if e != nil || len(b) > 8192 || strictJSON(b, &entry) != nil || entry.ID != id || !validProjectPath(entry.OriginalPath) {
		return entry, problem(422, "file_unavailable", "Invalid project Trash metadata; files were kept.")
	}
	snap, e := s.projectSnapshotAt(projectTrashDir + "/" + id)
	if e != nil {
		return entry, e
	}
	entry.Revision = snap.Revision
	entry.NoteCount = snap.NoteCount
	return entry, nil
}
func (s *Store) ProjectTrashList() ([]ProjectTrashEntry, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	rows := []ProjectTrashEntry{}
	f, e := s.root.Open(projectTrashDir)
	if errors.Is(e, os.ErrNotExist) {
		return rows, nil
	}
	if e != nil {
		return nil, storageError(e)
	}
	defer f.Close()
	entries, e := f.ReadDir(-1)
	if e != nil {
		return nil, storageError(e)
	}
	for _, file := range entries {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".json") {
			continue
		}
		id := strings.TrimSuffix(file.Name(), ".json")
		if !validTrashID(id) {
			continue
		}
		if _, e := s.root.Lstat(projectTrashDir + "/" + id); errors.Is(e, os.ErrNotExist) {
			continue
		}
		entry, e := s.projectTrashEntry(id)
		if e != nil {
			return nil, e
		}
		rows = append(rows, entry)
	}
	sort.Slice(rows, func(i, j int) bool { return rows[i].DeletedAt.After(rows[j].DeletedAt) })
	return rows, nil
}
func (s *Store) RestoreProject(id, revision string) (ProjectSnapshot, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	entry, e := s.projectTrashEntry(id)
	if e != nil {
		return ProjectSnapshot{}, e
	}
	if revision == "" || revision != entry.Revision {
		return ProjectSnapshot{}, problem(409, "revision_conflict", "Trashed project changed; reload before restoring.")
	}
	if _, e = s.root.Lstat(entry.OriginalPath); !errors.Is(e, os.ErrNotExist) {
		if e == nil {
			return ProjectSnapshot{}, problem(409, "already_exists", "A project occupies the original path; both were kept.")
		}
		return ProjectSnapshot{}, storageError(e)
	}
	if e = s.root.MkdirAll(path.Dir(entry.OriginalPath), 0755); e != nil {
		return ProjectSnapshot{}, storageError(e)
	}
	if e = s.uniqueName(entry.OriginalPath, ""); e != nil {
		return ProjectSnapshot{}, e
	}
	if e = s.root.Rename(projectTrashDir+"/"+id, entry.OriginalPath); e != nil {
		return ProjectSnapshot{}, storageError(e)
	}
	_ = s.root.Remove(projectTrashDir + "/" + id + ".json")
	return s.projectSnapshotAt(entry.OriginalPath)
}
