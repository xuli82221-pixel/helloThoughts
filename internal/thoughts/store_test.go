package thoughts

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"testing"
	"time"
)

func testStore(t *testing.T) *Store {
	t.Helper()
	s, e := NewStore(t.TempDir())
	if e != nil {
		t.Fatal(e)
	}
	t.Cleanup(func() { s.Close() })
	return s
}
func statusIs(t *testing.T, e error, status int) {
	t.Helper()
	var p *Error
	if !errors.As(e, &p) || p.Status != status {
		t.Fatalf("expected status %d, got %v", status, e)
	}
}
func TestCreateReadUpdateUnicode(t *testing.T) {
	s := testStore(t)
	p := "产品 创意/想法 📝.md"
	n, e := s.Write(p, "# 中文标题\n\n实际内容\n", "", true)
	if e != nil {
		t.Fatal(e)
	}
	if n.Title != "中文标题" || len(n.Revision) != 64 {
		t.Fatal("bad note metadata")
	}
	updated, e := s.Write(p, "# 修改\n\n第二版", n.Revision, false)
	if e != nil {
		t.Fatal(e)
	}
	b, e := os.ReadFile(filepath.Join(s.Vault, filepath.FromSlash(p)))
	if e != nil || string(b) != updated.Content {
		t.Fatal("file content differs")
	}
	rows, e := s.List()
	if e != nil || len(rows) != 1 || rows[0].Path != p {
		t.Fatalf("list: %v %v", rows, e)
	}
	files, e := os.ReadDir(filepath.Join(s.Vault, "产品 创意"))
	if e != nil || len(files) != 1 {
		t.Fatal("temporary file leaked")
	}
}
func TestCreateNeverOverwrites(t *testing.T) {
	s := testStore(t)
	_, e := s.Write("a.md", "original", "", true)
	if e != nil {
		t.Fatal(e)
	}
	_, e = s.Write("a.md", "replacement", "", true)
	statusIs(t, e, 409)
	n, _ := s.Read("a.md")
	if n.Content != "original" {
		t.Fatal("overwritten")
	}
}
func TestSameMtimeConflict(t *testing.T) {
	s := testStore(t)
	n, e := s.Write("a.md", "first", "", true)
	if e != nil {
		t.Fatal(e)
	}
	p := filepath.Join(s.Vault, "a.md")
	if e = os.WriteFile(p, []byte("external"), 0600); e != nil {
		t.Fatal(e)
	}
	if e = os.Chtimes(p, time.Now(), n.Modified); e != nil {
		t.Fatal(e)
	}
	_, e = s.Write("a.md", "stale draft", n.Revision, false)
	statusIs(t, e, 409)
	current, _ := s.Read("a.md")
	if current.Content != "external" {
		t.Fatal("external edit lost")
	}
	_, e = s.Write("a.md", "missing revision", "", false)
	statusIs(t, e, 409)
}
func TestInvalidPaths(t *testing.T) {
	s := testStore(t)
	for _, p := range []string{"../secret.md", "/secret.md", "C:/secret.md", "a\\b.md", "a/../../s.md", "a.md:stream", "a.txt", ".hidden/a.md", "CON.md", "a/NUL.md", "a./b.md", "a /b.md", "a//b.md", "COM1.md"} {
		t.Run(p, func(t *testing.T) { _, e := s.Write(p, "body", "", true); statusIs(t, e, 400) })
	}
}
func TestSymlinkEscape(t *testing.T) {
	s := testStore(t)
	outside := t.TempDir()
	if e := os.WriteFile(filepath.Join(outside, "secret.md"), []byte("private"), 0600); e != nil {
		t.Fatal(e)
	}
	if e := os.Symlink(outside, filepath.Join(s.Vault, "escape")); e != nil {
		if runtime.GOOS != "windows" {
			t.Fatal(e)
		}
		cmd := exec.Command("powershell.exe", "-NoProfile", "-NonInteractive", "-Command", "$ErrorActionPreference='Stop'; $null = New-Item -ItemType Junction -Path $env:MEMO_TEST_LINK -Target $env:MEMO_TEST_TARGET")
		cmd.Env = append(os.Environ(), "MEMO_TEST_LINK="+filepath.Join(s.Vault, "escape"), "MEMO_TEST_TARGET="+outside)
		if output, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("cannot create symlink or junction: %v %s", err, output)
		}
		t.Log("ordinary symlink privilege unavailable; testing Windows junction escape instead")
	}
	if _, e := s.Read("escape/secret.md"); e == nil {
		t.Fatal("read escaped root")
	}
	if _, e := s.Write("escape/write.md", "body", "", true); e == nil {
		t.Fatal("write escaped root")
	}
	if _, e := os.Stat(filepath.Join(outside, "write.md")); !os.IsNotExist(e) {
		t.Fatal("outside file written")
	}
}
func TestConcurrentCreate(t *testing.T) {
	s := testStore(t)
	var wg sync.WaitGroup
	success := make(chan bool, 12)
	for i := 0; i < 12; i++ {
		wg.Add(1)
		go func() { defer wg.Done(); _, e := s.Write("once.md", "body", "", true); success <- e == nil }()
	}
	wg.Wait()
	close(success)
	count := 0
	for ok := range success {
		if ok {
			count++
		}
	}
	if count != 1 {
		t.Fatalf("%d creates succeeded", count)
	}
}
func TestContentLimits(t *testing.T) {
	s := testStore(t)
	for _, body := range []string{strings.Repeat("x", MaxNoteSize+1), string([]byte{0xff}), "a\x00b"} {
		_, e := s.Write("a.md", body, "", true)
		statusIs(t, e, 400)
	}
}
func TestServiceLock(t *testing.T) {
	requireWindows(t)
	loc := Location{Dir: t.TempDir()}
	first, e := acquire(loc, "service.lock")
	if e != nil {
		t.Fatal(e)
	}
	if second, e := acquire(loc, "service.lock"); e == nil {
		second.Close()
		t.Fatal("second lock succeeded")
	}
	first.Close()
	second, e := acquire(loc, "service.lock")
	if e != nil {
		t.Fatal(e)
	}
	second.Close()
}
func TestInitIdempotent(t *testing.T) {
	requireWindows(t)
	t.Setenv("MEMO_CONFIG_DIR", t.TempDir())
	vault := filepath.Join(t.TempDir(), "资料库")
	v, e := Init(vault)
	if e != nil {
		t.Fatal(e)
	}
	if e = os.WriteFile(filepath.Join(v, "existing.md"), []byte("keep"), 0600); e != nil {
		t.Fatal(e)
	}
	if _, e = Init(vault); e != nil {
		t.Fatal(e)
	}
	b, e := os.ReadFile(filepath.Join(v, "existing.md"))
	if e != nil || string(b) != "keep" {
		t.Fatal("init modified existing data")
	}
}
