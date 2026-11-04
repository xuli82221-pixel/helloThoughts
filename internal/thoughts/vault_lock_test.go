package thoughts

import (
	"os"
	"testing"
)

func TestVaultLockRejectsASecondConfigHome(t *testing.T) {
	requireWindows(t)
	vault := t.TempDir()
	a, e := NewStore(vault)
	if e != nil {
		t.Fatal(e)
	}
	defer a.Close()
	b, e := NewStore(vault)
	if e != nil {
		t.Fatal(e)
	}
	defer b.Close()
	if e = a.root.MkdirAll(".memo/runtime", 0700); e != nil {
		t.Fatal(e)
	}
	first, e := a.root.OpenFile(".memo/runtime/service.lock", os.O_CREATE|os.O_RDWR, 0600)
	if e != nil {
		t.Fatal(e)
	}
	if e = lockFile(first); e != nil {
		t.Fatal(e)
	}
	second, e := b.root.OpenFile(".memo/runtime/service.lock", os.O_CREATE|os.O_RDWR, 0600)
	if e != nil {
		t.Fatal(e)
	}
	defer second.Close()
	if e = lockFile(second); e == nil {
		t.Fatal("second writer for same physical library was accepted")
	}
	first.Close()
	if e = lockFile(second); e != nil {
		t.Fatal("clean close did not release library lock", e)
	}
}
