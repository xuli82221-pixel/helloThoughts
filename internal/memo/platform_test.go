package memo

import (
	"runtime"
	"testing"
)

// Keep Windows integration assertions active on Windows while allowing the
// platform-independent tests to run in Linux build containers.
func requireWindows(t *testing.T) {
	t.Helper()
	if runtime.GOOS != "windows" {
		t.Skip("requires Windows filesystem permissions, file locking, or desktop runtime")
	}
}
