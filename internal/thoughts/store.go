package thoughts

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
	"unicode/utf8"
)

const MaxNoteSize = 1024 * 1024

type Error struct {
	Status  int    `json:"-"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e *Error) Error() string                       { return e.Message }
func problem(status int, code, message string) error { return &Error{status, code, message} }

type Note struct {
	Path     string    `json:"path"`
	Name     string    `json:"name"`
	Title    string    `json:"title"`
	Content  string    `json:"content"`
	Revision string    `json:"revision"`
	Modified time.Time `json:"modified"`
	Preview  string    `json:"preview,omitempty"`
	Snippet  string    `json:"snippet,omitempty"`
}

type Store struct {
	root  *os.Root
	Vault string
	mu    sync.Mutex
}

func NewStore(vault string) (*Store, error) {
	r, err := os.OpenRoot(vault)
	if err != nil {
		return nil, err
	}
	return &Store{root: r, Vault: vault}, nil
}
func (s *Store) Close() error { return s.root.Close() }
func digest(b []byte) string  { h := sha256.Sum256(b); return hex.EncodeToString(h[:]) }
func randomID() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	return hex.EncodeToString(b)
}

func validPath(p string) error {
	if len(p) > 1000 || !utf8.ValidString(p) || !fs.ValidPath(p) || strings.ContainsAny(p, "\\:<>\"|?*\x00") || !strings.EqualFold(path.Ext(p), ".md") {
		return problem(400, "invalid_path", "请使用资料库内的相对 Markdown 路径，例如 产品创意/想法.md。")
	}
	for _, part := range strings.Split(p, "/") {
		stem := strings.ToUpper(strings.SplitN(part, ".", 2)[0])
		if strings.IndexFunc(part, func(r rune) bool { return r < 32 }) >= 0 || strings.HasPrefix(part, ".") || strings.TrimRight(part, ". ") != part || len(part) > 240 || stem == "CON" || stem == "PRN" || stem == "AUX" || stem == "NUL" || (len(stem) == 4 && (strings.HasPrefix(stem, "COM") || strings.HasPrefix(stem, "LPT")) && stem[3] >= '0' && stem[3] <= '9') {
			return problem(400, "invalid_path", "路径包含保留名称、隐藏目录或无效字符。")
		}
	}
	return nil
}

func storageError(err error) error {
	if errors.Is(err, os.ErrNotExist) {
		return problem(404, "not_found", "文件不存在或已在外部移动。")
	}
	if errors.Is(err, os.ErrExist) {
		return problem(409, "already_exists", "同名笔记已存在，请换一个名称。")
	}
	return problem(422, "file_unavailable", "文件无法访问，请检查目录权限和符号链接；操作仅允许在资料库内进行。")
}

func (s *Store) read(p string) (Note, error) {
	if err := validPath(p); err != nil {
		return Note{}, err
	}
	return s.readFile(p)
}

// readFile is only used with a validated public path or a generated internal path.
func (s *Store) readFile(p string) (Note, error) {
	f, err := s.root.Open(p)
	if err != nil {
		return Note{}, storageError(err)
	}
	defer f.Close()
	st, err := f.Stat()
	if err != nil || !st.Mode().IsRegular() {
		return Note{}, problem(422, "not_regular", "只支持普通 Markdown 文件。")
	}
	b, err := io.ReadAll(io.LimitReader(f, MaxNoteSize+1))
	if err != nil {
		return Note{}, storageError(err)
	}
	if len(b) > MaxNoteSize || !utf8.Valid(b) || strings.ContainsRune(string(b), 0) {
		return Note{}, problem(422, "invalid_content", "M0 仅支持 1 MiB 以内的 UTF-8 Markdown 文本。")
	}
	title := strings.TrimSuffix(path.Base(p), path.Ext(p))
	for _, line := range strings.Split(string(b), "\n") {
		if strings.HasPrefix(line, "# ") {
			title = strings.TrimSpace(strings.TrimPrefix(line, "# "))
			break
		}
	}
	return Note{Path: p, Name: strings.TrimSuffix(path.Base(p), path.Ext(p)), Title: title, Content: string(b), Revision: digest(b), Modified: st.ModTime(), Preview: notePreview(string(b)), Snippet: trimPreview(string(b), 900)}, nil
}

func notePreview(content string) string {
	var fallback string
	for _, line := range strings.Split(content, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, "```") {
			continue
		}
		line = strings.Join(strings.Fields(line), " ")
		if fallback == "" {
			fallback = line
		}
		if len([]rune(line)) >= 20 {
			return trimPreview(line, 120)
		}
	}
	return trimPreview(fallback, 120)
}
func trimPreview(s string, limit int) string {
	r := []rune(s)
	if len(r) > limit {
		return string(r[:limit]) + "…"
	}
	return s
}
func (s *Store) Search(query string) ([]Note, error) {
	return s.SearchScoped(query, false)
}
func (s *Store) SearchScoped(query string, includeArchived bool) ([]Note, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	notes, err := s.list()
	if err != nil {
		return nil, err
	}
	q := strings.ToLower(strings.TrimSpace(query))
	results := []Note{}
	for _, n := range notes {
		if !includeArchived {
			if p, _ := classify(n.Path); p != "" {
				archived, e := s.projectArchived(p)
				if e != nil {
					return nil, e
				}
				if archived {
					continue
				}
			}
		}
		body := strings.ToLower(n.Content)
		if q != "" && !strings.Contains(strings.ToLower(n.Name), q) && !strings.Contains(strings.ToLower(n.Title), q) && !strings.Contains(body, q) {
			continue
		}
		if q != "" {
			if at := strings.Index(body, q); at >= 0 {
				start := utf8.RuneCountInString(body[:at])
				r := []rune(n.Content)
				lo := max(0, start-32)
				hi := min(len(r), start+len([]rune(q))+96)
				n.Preview = strings.Join(strings.Fields(string(r[lo:hi])), " ")
				if lo > 0 {
					n.Preview = "…" + n.Preview
				}
				if hi < len(r) {
					n.Preview += "…"
				}
			}
		}
		n.Content = ""
		results = append(results, n)
	}
	return results, nil
}
func (s *Store) Read(p string) (Note, error) { s.mu.Lock(); defer s.mu.Unlock(); return s.read(p) }

func (s *Store) List() ([]Note, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.list()
}
func (s *Store) list() ([]Note, error) {
	notes := []Note{}
	err := fs.WalkDir(s.root.FS(), ".", func(p string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return storageError(walkErr)
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
		if d.IsDir() || !strings.EqualFold(path.Ext(p), ".md") {
			return nil
		}
		n, err := s.read(p)
		if err != nil {
			return fmt.Errorf("%s: %w", p, err)
		}
		notes = append(notes, n)
		return nil
	})
	sort.Slice(notes, func(i, j int) bool { return notes[i].Modified.After(notes[j].Modified) })
	return notes, err
}

func (s *Store) Write(p, content, revision string, create bool) (Note, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.write(p, content, revision, create)
}
func (s *Store) write(p, content, revision string, create bool) (Note, error) {
	if err := validPath(p); err != nil {
		return Note{}, err
	}
	if len(content) > MaxNoteSize || !utf8.ValidString(content) || strings.ContainsRune(content, 0) {
		return Note{}, problem(400, "invalid_content", "正文必须是 1 MiB 以内的 UTF-8 文本。")
	}
	if !create {
		old, err := s.read(p)
		if err != nil {
			return Note{}, err
		}
		if revision == "" || old.Revision != revision {
			return Note{}, problem(409, "revision_conflict", "磁盘内容已变化，草稿未被覆盖。请重新读取或另存副本。")
		}
	}
	if err := s.root.MkdirAll(path.Dir(p), 0755); err != nil {
		return Note{}, storageError(err)
	}
	if create {
		if err := s.uniqueName(p, ""); err != nil {
			return Note{}, err
		}
	}
	tmp := path.Join(path.Dir(p), ".memo-"+randomID()+".tmp")
	f, err := s.root.OpenFile(tmp, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		return Note{}, storageError(err)
	}
	defer s.root.Remove(tmp)
	_, writeErr := io.WriteString(f, content)
	if writeErr == nil {
		writeErr = f.Sync()
	}
	closeErr := f.Close()
	if writeErr != nil {
		return Note{}, storageError(writeErr)
	}
	if closeErr != nil {
		return Note{}, storageError(closeErr)
	}
	if create {
		// A hard link publishes the complete file without replacing an existing target.
		err = s.root.Link(tmp, p)
	} else {
		// Recheck immediately before replacement, including external writes with unchanged mtime.
		old, readErr := s.read(p)
		if readErr != nil {
			return Note{}, readErr
		}
		if old.Revision != revision {
			return Note{}, problem(409, "revision_conflict", "磁盘内容已变化，请保留草稿并重新读取。")
		}
		err = s.root.Rename(tmp, p)
	}
	if err != nil {
		return Note{}, storageError(err)
	}
	return s.read(p)
}

func CanonicalVault(vault string) (string, error) {
	abs, err := filepath.Abs(vault)
	if err != nil {
		return "", err
	}
	return filepath.EvalSymlinks(abs)
}
