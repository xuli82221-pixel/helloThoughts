package thoughts

import (
	"bytes"
	"context"
	"debug/pe"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

type runtimePin struct {
	File    string `json:"file"`
	SHA256  string `json:"sha256"`
	Version string `json:"version"`
}
type runtimePins struct {
	Current  runtimePin  `json:"current"`
	Previous *runtimePin `json:"previous,omitempty"`
}

var runtimeVersion = regexp.MustCompile(`^[0-9]+\.[0-9]+\.[0-9]+(?:-[a-zA-Z0-9.-]+)?$`)

func managedHome() (string, error) {
	if legacyConfigOverride() {
		return "", fmt.Errorf("managed runtime requires a Thoughts home, not MEMO_CONFIG_DIR")
	}
	return sharedHome()
}
func pinPath(home string, p runtimePin) (string, error) {
	if !runtimeVersion.MatchString(p.Version) || !regexp.MustCompile(`^[a-f0-9]{64}$`).MatchString(p.SHA256) {
		return "", fmt.Errorf("invalid runtime pin")
	}
	// The former executable name is accepted only for existing immutable pins.
	name := filepath.Base(p.File)
	if name != "thoughts.exe" && name != "memo.exe" {
		return "", fmt.Errorf("unsupported managed executable name")
	}
	expected := filepath.Join("app", "versions", p.Version+"-"+p.SHA256[:12], name)
	if p.File != expected {
		return "", fmt.Errorf("runtime path is not an immutable managed version")
	}
	root, e := os.OpenRoot(home)
	if e != nil {
		return "", e
	}
	defer root.Close()
	f, e := root.Open(p.File)
	if e != nil {
		return "", e
	}
	defer f.Close()
	b, e := io.ReadAll(io.LimitReader(f, 150*1024*1024))
	if e != nil || digest(b) != p.SHA256 {
		return "", fmt.Errorf("managed executable hash mismatch; no fallback was executed")
	}
	full := filepath.Join(home, p.File)
	canonical, e := filepath.EvalSymlinks(full)
	if e != nil || !strings.EqualFold(canonical, full) {
		return "", fmt.Errorf("managed runtime reparse points are not allowed")
	}
	return full, nil
}
func readPins(home string) (runtimePins, error) {
	var p runtimePins
	b, e := os.ReadFile(filepath.Join(home, "config", "managed-runtime.json"))
	if e != nil {
		return p, e
	}
	if len(b) > 16384 || strictJSON(b, &p) != nil {
		return p, fmt.Errorf("invalid managed runtime metadata")
	}
	if _, e = pinPath(home, p.Current); e != nil {
		return p, fmt.Errorf("managed runtime pin is invalid: %v", e)
	}
	return p, nil
}

// Bundled MCPB/stdio executables become stable launchers after one managed install.
// The destination is a hash-pinned executable under this same explicit home.
func ForwardManaged(args []string) (bool, error) {
	if os.Getenv("THOUGHTS_HOME") == "" && !legacyConfigOverride() {
		if self, e := os.Executable(); e == nil {
			if home := discoverManagedHome(self); home != "" {
				_ = os.Setenv("THOUGHTS_HOME", home)
			}
		}
	}
	if os.Getenv("THOUGHTS_NO_DELEGATE") == "1" || legacyConfigOverride() || len(args) > 0 && (args[0] == "runtime" || args[0] == "capabilities" || args[0] == "paths") {
		return false, nil
	}
	home, e := managedHome()
	if e != nil {
		return false, e
	}
	pins, e := readPins(home)
	if errors.Is(e, os.ErrNotExist) {
		return false, nil
	}
	if e != nil {
		return false, e
	}
	target, e := pinPath(home, pins.Current)
	if e != nil {
		return false, e
	}
	self, e := os.Executable()
	if e != nil {
		return false, e
	}
	if strings.EqualFold(filepath.Clean(self), target) {
		return false, nil
	}
	cmd := exec.Command(target, args...)
	cmd.Env = append(os.Environ(), "THOUGHTS_HOME="+home, "THOUGHTS_NO_DELEGATE=1")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return true, cmd.Run()
}

// A managed executable opened directly must retain its own library binding.
// Portable/unmanaged programs keep the explicit-env/legacy resolution rules.
func discoverManagedHome(executable string) string {
	directory := filepath.Dir(executable)
	versions := filepath.Dir(directory)
	app := filepath.Dir(versions)
	name := filepath.Base(executable)
	if (name != "thoughts.exe" && name != "memo.exe") || filepath.Base(versions) != "versions" || filepath.Base(app) != "app" {
		return ""
	}
	home := filepath.Dir(app)
	if _, e := readPins(home); e != nil {
		return ""
	}
	b, e := os.ReadFile(executable)
	if e != nil {
		return ""
	}
	relative, e := filepath.Rel(home, executable)
	if e != nil {
		return ""
	}
	self := runtimePin{File: relative, SHA256: digest(b), Version: Version}
	if _, e := pinPath(home, self); e != nil {
		return ""
	}
	return home
}
func stageRuntime(home, source, expected string) (runtimePin, error) {
	var p runtimePin
	if !regexp.MustCompile(`^[a-fA-F0-9]{64}$`).MatchString(expected) {
		return p, fmt.Errorf("provide the reviewed executable SHA256")
	}
	source, e := filepath.Abs(source)
	if e != nil {
		return p, e
	}
	f, e := os.Open(source)
	if e != nil {
		return p, e
	}
	defer f.Close()
	b, e := io.ReadAll(io.LimitReader(f, 150*1024*1024))
	if e != nil || digest(b) != strings.ToLower(expected) {
		return p, fmt.Errorf("source executable SHA256 mismatch; nothing was installed")
	}
	image, e := pe.NewFile(bytes.NewReader(b))
	if e != nil {
		return p, fmt.Errorf("source is not a Windows executable")
	}
	machine := image.Machine
	image.Close()
	if machine != pe.IMAGE_FILE_MACHINE_AMD64 {
		return p, fmt.Errorf("Windows x64 executable required")
	}
	staging := filepath.Join(home, "app", "staging", randomID())
	if e = privateDir(staging); e != nil {
		return p, e
	}
	staged := filepath.Join(staging, "thoughts.exe")
	if e = os.WriteFile(staged, b, 0600); e != nil {
		return p, e
	}
	defer os.Remove(staging)
	defer os.Remove(staged)
	// Preserve downloaded-file provenance; never remove Mark of the Web to execute.
	zone, zoneErr := os.ReadFile(source + ":Zone.Identifier")
	if zoneErr == nil {
		if len(zone) > 16384 {
			return p, fmt.Errorf("invalid download provenance metadata")
		}
		if e = os.WriteFile(staged+":Zone.Identifier", zone, 0600); e != nil {
			return p, e
		}
	} else if !errors.Is(zoneErr, os.ErrNotExist) {
		return p, zoneErr
	}
	// Verify identity/capabilities before stopping an existing service or changing pins.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, staged, "capabilities", "--json")
	cmd.Env = append(os.Environ(), "THOUGHTS_NO_DELEGATE=1")
	output, e := cmd.Output()
	if e != nil {
		return p, fmt.Errorf("candidate cannot run or does not support safe upgrades; original runtime kept: %w", e)
	}
	var caps struct {
		Product  string `json:"product"`
		Version  string `json:"version"`
		Protocol int    `json:"protocolVersion"`
		Managed  int    `json:"managedRuntime"`
	}
	if json.Unmarshal(output, &caps) != nil || caps.Product != "thoughts" || !runtimeVersion.MatchString(caps.Version) || caps.Protocol != 1 || caps.Managed != 1 {
		return p, fmt.Errorf("candidate identity or runtime protocol is unsupported")
	}
	p = runtimePin{SHA256: strings.ToLower(expected), Version: caps.Version}
	name := "thoughts.exe"
	// Older reviewed cores require their original name for rollback compatibility.
	if filepath.Base(source) == "memo.exe" {
		name = "memo.exe"
	}
	p.File = filepath.Join("app", "versions", p.Version+"-"+p.SHA256[:12], name)
	target := filepath.Join(home, p.File)
	if _, e = os.Stat(target); e == nil {
		_, e = pinPath(home, p)
		return p, e
	} else if !errors.Is(e, os.ErrNotExist) {
		return p, e
	}
	if e = privateDir(filepath.Dir(target)); e != nil {
		return p, e
	}
	out, e := os.OpenFile(target, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
	if e != nil {
		return p, e
	}
	_, writeErr := out.Write(b)
	if writeErr == nil {
		writeErr = out.Sync()
	}
	closeErr := out.Close()
	if writeErr != nil {
		return p, writeErr
	}
	if closeErr != nil {
		return p, closeErr
	}
	if zoneErr == nil {
		if e = os.WriteFile(target+":Zone.Identifier", zone, 0600); e != nil {
			return p, e
		}
	}
	_, e = pinPath(home, p)
	return p, e
}
func stopManagedService(ctx context.Context, loc Location) (bool, error) {
	s, e := authenticate(ctx, loc)
	if e != nil {
		probe, lockErr := acquire(loc, "service.lock")
		if lockErr != nil {
			return false, fmt.Errorf("unverified or incompatible service is running; no process was stopped")
		}
		probe.Close()
		return false, nil
	}
	var result any
	if e = Call(ctx, loc, "POST", "/api/control/shutdown", map[string]string{"instance": s.Instance}, &result); e != nil {
		return false, e
	}
	deadline := time.Now().Add(6 * time.Second)
	for time.Now().Before(deadline) {
		probe, e := acquire(loc, "service.lock")
		if e == nil {
			probe.Close()
			return true, nil
		}
		time.Sleep(50 * time.Millisecond)
	}
	return false, fmt.Errorf("verified service did not stop in time; runtime pin was not switched")
}
func startPinned(ctx context.Context, loc Location, home string, p runtimePin) error {
	target, e := pinPath(home, p)
	if e != nil {
		return e
	}
	log, e := os.OpenFile(loc.Log, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if e != nil {
		return e
	}
	defer log.Close()
	if e = startDetached(target, []string{"serve", "--vault", loc.Vault}, log); e != nil {
		return e
	}
	deadline := time.Now().Add(12 * time.Second)
	for time.Now().Before(deadline) {
		if s, e := authenticate(ctx, loc); e == nil && s.Version == p.Version {
			return nil
		}
		time.Sleep(100 * time.Millisecond)
	}
	return fmt.Errorf("candidate service did not become healthy; inspect local service log")
}
func ManageRuntime(ctx context.Context, action, source, expected string) (any, error) {
	home, e := managedHome()
	if e != nil {
		return nil, e
	}
	base, e := configDir()
	if e != nil {
		return nil, e
	}
	if e = privateDir(base); e != nil {
		return nil, e
	}
	old, e := readPins(home)
	hasOld := e == nil
	if e != nil && !errors.Is(e, os.ErrNotExist) {
		return nil, e
	}
	if action == "status" {
		return map[string]any{"installed": hasOld, "runtime": old}, nil
	}
	if action != "install" && action != "upgrade" && action != "rollback" {
		return nil, fmt.Errorf("use runtime install, upgrade, rollback or status")
	}
	guard, e := acquire(Location{Dir: base}, "upgrade.lock")
	if e != nil {
		return nil, fmt.Errorf("another runtime update is in progress")
	}
	defer guard.Close()
	// Re-read after the lock; never switch from an outdated snapshot.
	old, e = readPins(home)
	hasOld = e == nil
	if e != nil && !errors.Is(e, os.ErrNotExist) {
		return nil, e
	}
	var next runtimePin
	if action == "rollback" {
		if !hasOld || old.Previous == nil {
			return nil, fmt.Errorf("no previous managed version is available")
		}
		next = *old.Previous
		if _, e = pinPath(home, next); e != nil {
			return nil, e
		}
	} else {
		if source == "" {
			return nil, fmt.Errorf("--source and --sha256 are required")
		}
		next, e = stageRuntime(home, source, expected)
		if e != nil {
			return nil, e
		}
	}
	if hasOld && old.Current == next {
		return map[string]any{"installed": true, "reused": true, "version": next.Version}, nil
	}
	loc, e := Resolve("")
	if e != nil {
		return nil, e
	}
	launch, e := acquire(loc, "launch.lock")
	if e != nil {
		return nil, fmt.Errorf("service launch is in progress; retry after it finishes")
	}
	defer launch.Close()
	wasRunning, e := stopManagedService(ctx, loc)
	if e != nil {
		return nil, e
	}
	pins := runtimePins{Current: next}
	if hasOld {
		previous := old.Current
		pins.Previous = &previous
	}
	pinFile := filepath.Join(base, "managed-runtime.json")
	if hasOld {
		backup := filepath.Join(home, "backups", "runtime")
		if e = privateDir(backup); e != nil {
			return nil, e
		}
		b, _ := json.Marshal(old)
		if e = os.WriteFile(filepath.Join(backup, time.Now().UTC().Format("20060102T150405.000000000")+".json"), b, 0600); e != nil {
			return nil, e
		}
	}
	b, _ := json.Marshal(pins)
	if e = atomicConfig(pinFile, b); e != nil {
		return nil, e
	}
	if e = startPinned(ctx, loc, home, next); e != nil {
		if hasOld {
			previous, _ := json.Marshal(old)
			_ = atomicConfig(pinFile, previous)
			// Never kill an unknown process. Only restart the old pin when no service holds the lock.
			probe, lockErr := acquire(loc, "service.lock")
			if lockErr == nil {
				probe.Close()
				if wasRunning {
					_ = startPinned(ctx, loc, home, old.Current)
				}
			}
		} else {
			_ = os.Remove(pinFile)
		}
		return nil, fmt.Errorf("runtime activation failed; prior pin restored when available: %w", e)
	}
	return map[string]any{"installed": true, "version": next.Version, "executable": filepath.Join(home, next.File), "sha256": next.SHA256, "library": loc.Vault, "rollbackAvailable": hasOld}, nil
}
