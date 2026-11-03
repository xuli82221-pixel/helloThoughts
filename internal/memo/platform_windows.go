//go:build windows

package memo

import (
	"os"
	"os/exec"
	"os/user"
	"syscall"
	"unsafe"
)

var kernel32 = syscall.NewLazyDLL("kernel32.dll")
var advapi32 = syscall.NewLazyDLL("advapi32.dll")

func windowsProfileDirectory() (string, error) {
	u, err := user.Current()
	if err != nil {
		return "", err
	}
	if u.HomeDir == "" {
		return "", syscall.EINVAL
	}
	return u.HomeDir, nil
}

func documentsDir() (string, error) {
	id := syscall.GUID{Data1: 0xFDD39AD0, Data2: 0x238F, Data3: 0x46AF, Data4: [8]byte{0xAD, 0xB4, 0x6C, 0x85, 0x48, 0x03, 0x69, 0xC7}}
	var p *uint16
	hresult, _, _ := syscall.NewLazyDLL("shell32.dll").NewProc("SHGetKnownFolderPath").Call(uintptr(unsafe.Pointer(&id)), 0, 0, uintptr(unsafe.Pointer(&p)))
	if hresult != 0 {
		return "", syscall.Errno(hresult)
	}
	defer syscall.NewLazyDLL("ole32.dll").NewProc("CoTaskMemFree").Call(uintptr(unsafe.Pointer(p)))
	var text []uint16
	for i := uintptr(0); i < 32768; i++ {
		c := *(*uint16)(unsafe.Pointer(uintptr(unsafe.Pointer(p)) + 2*i))
		if c == 0 {
			return syscall.UTF16ToString(text), nil
		}
		text = append(text, c)
	}
	return "", syscall.EINVAL
}

func lockFile(f *os.File) error {
	var ov syscall.Overlapped
	r, _, e := kernel32.NewProc("LockFileEx").Call(f.Fd(), 3, 0, 1, 0, uintptr(unsafe.Pointer(&ov)))
	if r == 0 {
		return e
	}
	return nil
}

func privateDir(p string) error {
	if err := os.MkdirAll(p, 0700); err != nil {
		return err
	}
	u, err := user.Current()
	if err != nil {
		return err
	}
	sddl, err := syscall.UTF16PtrFromString("D:P(A;OICI;FA;;;SY)(A;OICI;FA;;;" + u.Uid + ")")
	if err != nil {
		return err
	}
	var descriptor uintptr
	r, _, e := advapi32.NewProc("ConvertStringSecurityDescriptorToSecurityDescriptorW").Call(uintptr(unsafe.Pointer(sddl)), 1, uintptr(unsafe.Pointer(&descriptor)), 0)
	if r == 0 {
		return e
	}
	defer kernel32.NewProc("LocalFree").Call(descriptor)
	pw, err := syscall.UTF16PtrFromString(p)
	if err != nil {
		return err
	}
	r, _, e = advapi32.NewProc("SetFileSecurityW").Call(uintptr(unsafe.Pointer(pw)), 0x80000004, descriptor)
	if r == 0 {
		return e
	}
	return nil
}

func startDetached(exe string, args []string, log *os.File) error {
	start := func(flags uint32) error {
		cmd := exec.Command(exe, args...)
		cmd.Stdout = log
		cmd.Stderr = log
		cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true, CreationFlags: flags}
		if err := cmd.Start(); err != nil {
			return err
		}
		return cmd.Process.Release()
	}
	if err := start(0x01000208); err != nil {
		return start(0x00000208)
	}
	return nil
}

func systemBrowser(url string) error {
	cmd := exec.Command("rundll32.exe", "url.dll,FileProtocolHandler", url)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	if err := cmd.Start(); err != nil {
		return err
	}
	return cmd.Process.Release()
}
