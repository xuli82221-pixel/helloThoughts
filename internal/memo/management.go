package memo

import (
	"encoding/json"
	"errors"
	"io"
	"io/fs"
	"os"
	"path"
	"sort"
	"strings"
	"time"
)

const trashDirectory = ".memo/trash"

type Project struct {
	Path     string `json:"path"`
	Name     string `json:"name"`
	Nodes    []Note `json:"nodes"`
	Archived bool   `json:"archived"`
}
type Library struct {
	Projects     []Project `json:"projects"`
	Inbox        []Note    `json:"inbox"`
	AssetVersion string    `json:"assetVersion"`
}
type TrashEntry struct {
	ID           string    `json:"id"`
	OriginalPath string    `json:"originalPath"`
	Name         string    `json:"name"`
	DeletedAt    time.Time `json:"deletedAt"`
	Revision     string    `json:"revision"`
	Reason       string    `json:"reason"`
}

func nodePath(directory, name string) (string, error) {
	if name != strings.TrimSpace(name) || name == "" || strings.ContainsAny(name, "/\\") {
		return "", problem(400, "invalid_path", "Enter a valid node name without path separators.")
	}
	p := path.Join(directory, name+".md")
	if err := validPath(p); err != nil {
		return "", err
	}
	return p, nil
}

// uniqueName rejects case-only duplicates on every supported platform.
func (s *Store) uniqueName(target, except string) error {
	dir, err := s.root.Open(path.Dir(target))
	if err != nil {
		return storageError(err)
	}
	defer dir.Close()
	entries, err := dir.ReadDir(-1)
	if err != nil {
		return storageError(err)
	}
	for _, entry := range entries {
		candidate := path.Join(path.Dir(target), entry.Name())
		if candidate != except && strings.EqualFold(entry.Name(), path.Base(target)) {
			return problem(409, "already_exists", "A file or folder with this name already exists.")
		}
	}
	return nil
}

func classify(p string) (string, string) {
	parts := strings.Split(p, "/")
	if len(parts) < 2 || parts[0] == "Inbox" {
		return "", ""
	}
	if len(parts) >= 3 && (parts[0] == "Projects" || parts[0] == "项目") {
		return strings.Join(parts[:2], "/"), parts[1]
	}
	return parts[0], parts[0]
}

func (s *Store) Library() (Library, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	notes, err := s.list()
	if err != nil {
		return Library{}, err
	}
	library := Library{Projects: []Project{}, Inbox: []Note{}}
	projects := map[string]*Project{}
	add := func(p, name string) *Project {
		if projects[p] == nil {
			projects[p] = &Project{Path: p, Name: name, Nodes: []Note{}}
		}
		return projects[p]
	}
	err = fs.WalkDir(s.root.FS(), ".", func(p string, d fs.DirEntry, e error) error {
		if e != nil {
			return storageError(e)
		}
		if p != "." && strings.HasPrefix(d.Name(), ".") {
			if d.IsDir() {
				return fs.SkipDir
			}
			return nil
		}
		if d.Type()&os.ModeSymlink != 0 {
			return nil
		}
		if !d.IsDir() || p == "." || p == "Inbox" {
			return nil
		}
		if strings.EqualFold(d.Name(), "assets") {
			return fs.SkipDir
		}
		if strings.HasPrefix(p, "Inbox/") {
			return fs.SkipDir
		}
		parts := strings.Split(p, "/")
		if parts[0] == "Projects" || parts[0] == "项目" {
			if len(parts) == 2 {
				add(p, parts[1])
			}
		} else if len(parts) == 1 {
			add(p, parts[0])
		}
		return nil
	})
	if err != nil {
		return Library{}, err
	}
	for _, note := range notes {
		note.Content = ""
		p, name := classify(note.Path)
		if p == "" {
			library.Inbox = append(library.Inbox, note)
		} else {
			add(p, name).Nodes = append(add(p, name).Nodes, note)
		}
	}
	for _, project := range projects {
		project.Archived, err = s.projectArchived(project.Path)
		if err != nil {
			return Library{}, err
		}
		sort.Slice(project.Nodes, func(i, j int) bool { return project.Nodes[i].Name < project.Nodes[j].Name })
		library.Projects = append(library.Projects, *project)
	}
	sort.Slice(library.Projects, func(i, j int) bool { return library.Projects[i].Name < library.Projects[j].Name })
	library.AssetVersion = s.assetVersion()
	return library, nil
}

func (s *Store) CreateProject(name, language string, template bool) (Project, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !validLanguage(language) {
		return Project{}, problem(400, "invalid_language", "Unsupported template language.")
	}
	checked, err := nodePath("Projects", name)
	if err != nil {
		return Project{}, err
	}
	target := strings.TrimSuffix(checked, ".md")
	if err = s.root.MkdirAll("Projects", 0755); err != nil {
		return Project{}, storageError(err)
	}
	if err = s.uniqueName(target, ""); err != nil {
		return Project{}, err
	}
	// Mkdir is exclusive: repeated project creation never refills removed template nodes.
	if err = s.root.Mkdir(target, 0755); err != nil {
		return Project{}, storageError(err)
	}
	project := Project{Path: target, Name: path.Base(target), Nodes: []Note{}}
	if !template {
		return project, nil
	}
	names := map[string][]string{
		"en":      {"Summary", "Rules", "Progress", "Details", "Ideas", "Questions", "Tasks"},
		"ko":      {"요약", "규칙", "진행 상황", "세부 사항", "아이디어", "질문", "할 일"},
		"ja":      {"概要", "ルール", "進捗", "詳細", "アイデア", "質問", "タスク"},
		"zh-Hans": {"摘要", "规则", "进度", "细节", "想法", "问题", "待办"},
		"zh-Hant": {"摘要", "規則", "進度", "細節", "想法", "問題", "待辦"},
		"de":      {"Zusammenfassung", "Regeln", "Fortschritt", "Details", "Ideen", "Fragen", "Aufgaben"},
		"fr":      {"Résumé", "Règles", "Avancement", "Détails", "Idées", "Questions", "Tâches"},
		"es":      {"Resumen", "Reglas", "Progreso", "Detalles", "Ideas", "Preguntas", "Tareas"},
	}[language]
	for _, nodeName := range names {
		note, e := s.write(path.Join(target, nodeName+".md"), "", "", true)
		if e != nil {
			return Project{}, problem(422, "project_incomplete", "The project was created, but not every starter node could be written. Existing files are preserved; inspect the current tree before continuing.")
		}
		note.Content = ""
		project.Nodes = append(project.Nodes, note)
	}
	return project, nil
}

func (s *Store) ensureTrash() error {
	for _, p := range []string{".memo", trashDirectory} {
		info, err := s.root.Lstat(p)
		if err == nil && (!info.IsDir() || info.Mode()&os.ModeSymlink != 0) {
			return problem(422, "file_unavailable", "Internal storage must be a real directory.")
		}
		if err != nil && !errors.Is(err, os.ErrNotExist) {
			return storageError(err)
		}
		if err = s.root.MkdirAll(p, 0700); err != nil {
			return storageError(err)
		}
	}
	return nil
}
func trashPaths(id string) (string, string) {
	return path.Join(trashDirectory, id+".md"), path.Join(trashDirectory, id+".json")
}
func validTrashID(id string) bool {
	if len(id) != 64 {
		return false
	}
	for _, r := range id {
		if !strings.ContainsRune("0123456789abcdef", r) {
			return false
		}
	}
	return true
}

// Metadata is synced before moving the source. An interrupted operation remains
// visible and recoverable in Trash rather than silently losing its only copy.
func (s *Store) stage(p, revision, reason string) (TrashEntry, error) {
	note, err := s.read(p)
	if err != nil {
		return TrashEntry{}, err
	}
	if revision == "" || note.Revision != revision {
		return TrashEntry{}, problem(409, "revision_conflict", "The file changed. Read its current revision before this operation.")
	}
	info, err := s.root.Lstat(p)
	if err != nil {
		return TrashEntry{}, storageError(err)
	}
	if !info.Mode().IsRegular() {
		return TrashEntry{}, problem(422, "not_regular", "Only regular Markdown files can be renamed or moved to Trash.")
	}
	if err = s.ensureTrash(); err != nil {
		return TrashEntry{}, err
	}
	entry := TrashEntry{ID: randomID(), OriginalPath: p, Name: note.Name, DeletedAt: time.Now(), Revision: revision, Reason: reason}
	payload, metadata := trashPaths(entry.ID)
	body, _ := json.Marshal(entry)
	file, err := s.root.OpenFile(metadata, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		return TrashEntry{}, storageError(err)
	}
	_, err = file.Write(body)
	if err == nil {
		err = file.Sync()
	}
	closeErr := file.Close()
	if err == nil {
		err = closeErr
	}
	if err != nil {
		s.root.Remove(metadata)
		return TrashEntry{}, storageError(err)
	}
	current, err := s.read(p)
	if err != nil || current.Revision != revision {
		s.root.Remove(metadata)
		return TrashEntry{}, problem(409, "revision_conflict", "The file changed before the operation.")
	}
	if err = s.root.Rename(p, payload); err != nil {
		s.root.Remove(metadata)
		return TrashEntry{}, storageError(err)
	}
	captured, err := s.readFile(payload)
	if err != nil || captured.Revision != revision {
		s.rollbackStage(entry)
		return TrashEntry{}, problem(409, "operation_conflict", "A concurrent change was preserved. Check the original path and Trash before retrying.")
	}
	return entry, nil
}
func (s *Store) rollbackStage(entry TrashEntry) {
	payload, metadata := trashPaths(entry.ID)
	if err := s.root.Link(payload, entry.OriginalPath); err == nil {
		s.root.Remove(payload)
		s.root.Remove(metadata)
	}
}

func (s *Store) Rename(p, name, revision string) (Note, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := validPath(p); err != nil {
		return Note{}, err
	}
	target, err := nodePath(path.Dir(p), name)
	if err != nil {
		return Note{}, err
	}
	if err = s.uniqueName(target, p); err != nil {
		return Note{}, err
	}
	if p == target {
		n, e := s.read(p)
		if e == nil && (revision == "" || revision != n.Revision) {
			return Note{}, problem(409, "revision_conflict", "The file changed.")
		}
		return n, e
	}
	entry, err := s.stage(p, revision, "rename")
	if err != nil {
		return Note{}, err
	}
	payload, metadata := trashPaths(entry.ID)
	if err = s.root.Link(payload, target); err != nil {
		s.rollbackStage(entry)
		return Note{}, storageError(err)
	}
	note, err := s.read(target)
	if err != nil {
		return Note{}, problem(409, "operation_conflict", "The captured content remains in Trash and at the new path. Inspect both files.")
	}
	s.root.Remove(payload)
	s.root.Remove(metadata)
	return note, nil
}

func (s *Store) Trash(p, revision string) (TrashEntry, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.stage(p, revision, "deleted")
}

func (s *Store) trashEntry(id string) (TrashEntry, error) {
	if !validTrashID(id) {
		return TrashEntry{}, problem(400, "invalid_path", "Invalid Trash entry.")
	}
	payload, metadata := trashPaths(id)
	info, err := s.root.Lstat(payload)
	if err != nil {
		return TrashEntry{}, storageError(err)
	}
	if !info.Mode().IsRegular() {
		return TrashEntry{}, problem(422, "not_regular", "Trash payload must be a regular Markdown file.")
	}
	file, err := s.root.Open(metadata)
	if err != nil {
		return TrashEntry{}, storageError(err)
	}
	defer file.Close()
	body, err := io.ReadAll(io.LimitReader(file, 8193))
	if err != nil || len(body) > 8192 {
		return TrashEntry{}, problem(422, "file_unavailable", "Trash metadata is unreadable.")
	}
	var entry TrashEntry
	if json.Unmarshal(body, &entry) != nil || entry.ID != id || validPath(entry.OriginalPath) != nil {
		return TrashEntry{}, problem(422, "file_unavailable", "Trash metadata is invalid; its payload was preserved.")
	}
	note, err := s.readFile(payload)
	if err != nil {
		return TrashEntry{}, err
	}
	entry.Revision = note.Revision
	return entry, nil
}
func (s *Store) TrashList() ([]TrashEntry, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	rows := []TrashEntry{}
	dir, err := s.root.Open(trashDirectory)
	if errors.Is(err, os.ErrNotExist) {
		return rows, nil
	}
	if err != nil {
		return nil, storageError(err)
	}
	defer dir.Close()
	files, err := dir.ReadDir(-1)
	if err != nil {
		return nil, storageError(err)
	}
	for _, file := range files {
		if file.IsDir() || file.Type()&os.ModeSymlink != 0 || !strings.HasSuffix(file.Name(), ".json") {
			continue
		}
		id := strings.TrimSuffix(file.Name(), ".json")
		if !validTrashID(id) {
			continue
		}
		payload, _ := trashPaths(id)
		if _, e := s.root.Lstat(payload); errors.Is(e, os.ErrNotExist) {
			continue
		}
		entry, e := s.trashEntry(id)
		if e != nil {
			return nil, e
		}
		rows = append(rows, entry)
	}
	sort.Slice(rows, func(i, j int) bool { return rows[i].DeletedAt.After(rows[j].DeletedAt) })
	return rows, nil
}
func (s *Store) Restore(id, revision string) (Note, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	entry, err := s.trashEntry(id)
	if err != nil {
		return Note{}, err
	}
	if revision == "" || entry.Revision != revision {
		return Note{}, problem(409, "revision_conflict", "The Trash content changed. Reload it before restoring.")
	}
	target := entry.OriginalPath
	if err = s.root.MkdirAll(path.Dir(target), 0755); err != nil {
		return Note{}, storageError(err)
	}
	if err = s.uniqueName(target, ""); err != nil {
		return Note{}, err
	}
	payload, metadata := trashPaths(id)
	snapshot, err := s.readFile(payload)
	if err != nil {
		return Note{}, err
	}
	if snapshot.Revision != revision {
		return Note{}, problem(409, "revision_conflict", "The Trash content changed before restoration.")
	}
	// Restoring creates an independent file. An interrupted rename may already
	// have published another hard link to this payload; do not alias that note.
	note, err := s.write(target, snapshot.Content, "", true)
	if err != nil {
		return Note{}, err
	}
	latest, err := s.readFile(payload)
	if err != nil || latest.Revision != revision {
		return Note{}, problem(409, "operation_conflict", "The restored snapshot and changed Trash content were both preserved.")
	}
	s.root.Remove(payload)
	s.root.Remove(metadata)
	return note, nil
}
