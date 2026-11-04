//go:build !windows

package thoughts

import (
	"errors"
	"os"
)

var errPlatformDeferred = errors.New("This Thoughts preview supports Windows only; macOS work is deferred")

func documentsDir() (string, error)                  { return "", errPlatformDeferred }
func windowsProfileDirectory() (string, error)       { return "", errPlatformDeferred }
func lockFile(*os.File) error                        { return errPlatformDeferred }
func privateDir(string) error                        { return errPlatformDeferred }
func startDetached(string, []string, *os.File) error { return errPlatformDeferred }
func systemBrowser(string) error                     { return errPlatformDeferred }
