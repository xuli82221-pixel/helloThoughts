package memo

import (
	"bytes"
	"context"
	"crypto/subtle"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

var Version = "0.3.0-preview.20"

type State struct {
	Product  string `json:"product"`
	Version  string `json:"version"`
	Vault    string `json:"vault"`
	Instance string `json:"instance"`
	Token    string `json:"token"`
	URL      string `json:"url"`
	PID      int    `json:"pid"`
}
type Location struct {
	Dir   string
	Log   string
	Vault string
}

func configDir() (string, error) {
	if legacyConfigOverride() {
		return filepath.Abs(os.Getenv("MEMO_CONFIG_DIR"))
	}
	home, err := sharedHome()
	if err != nil {
		return "", err
	}
	base := filepath.Join(home, "config")
	// Never silently start an empty default library over a known old installation.
	if automaticHome() {
		if _, e := os.Stat(filepath.Join(base, "config.json")); os.IsNotExist(e) {
			profile, e := profileDirectory()
			if e != nil {
				return "", e
			}
			legacy := filepath.Join(profile, "AppData", "Roaming", "Memo", "config.json")
			if _, e := os.Stat(legacy); e == nil {
				return "", fmt.Errorf("legacy Memo configuration found; back up and migrate it to the shared Thoughts home, or explicitly select MEMO_CONFIG_DIR; no new library was created")
			} else if !os.IsNotExist(e) {
				return "", e
			}
		} else if e != nil {
			return "", e
		}
	}
	return base, nil
}
func Init(vault string) (string, error) {
	base, err := configDir()
	if err != nil {
		return "", err
	}
	if err = privateDir(base); err != nil {
		return "", err
	}
	lock, err := acquire(Location{Dir: base}, "init.lock")
	if err != nil {
		return "", fmt.Errorf("another initialization is in progress; existing configuration was kept")
	}
	defer lock.Close()
	if _, e := os.Stat(filepath.Join(base, "config.json")); e == nil {
		loc, e := Resolve("")
		if e != nil {
			return "", e
		}
		if vault != "" && !sameMCPVault(loc.Vault, vault) {
			return "", fmt.Errorf("existing library was kept; use the explicit migration workflow to switch")
		}
		return loc.Vault, nil
	} else if !os.IsNotExist(e) {
		return "", e
	}
	if vault == "" {
		if !legacyConfigOverride() {
			home, e := sharedHome()
			if e != nil {
				return "", e
			}
			vault = filepath.Join(home, "library")
		} else {
			documents, e := documentsDir()
			if e != nil {
				return "", e
			}
			vault = filepath.Join(documents, "Memo")
		}
	}
	if err = os.MkdirAll(vault, 0755); err != nil {
		return "", err
	}
	vault, err = CanonicalVault(vault)
	if err != nil {
		return "", err
	}
	b, _ := json.Marshal(map[string]string{"vault": vault})
	if err = atomicConfig(filepath.Join(base, "config.json"), b); err != nil {
		return "", err
	}
	return vault, nil
}
func Resolve(vault string) (Location, error) {
	base, err := configDir()
	if err != nil {
		return Location{}, err
	}
	if vault == "" {
		b, e := os.ReadFile(filepath.Join(base, "config.json"))
		if e != nil {
			return Location{}, fmt.Errorf("尚未初始化，请运行 memo init（可用 --vault 指定资料库）")
		}
		var c struct {
			Vault string `json:"vault"`
		}
		if e = json.Unmarshal(b, &c); e != nil {
			return Location{}, fmt.Errorf("配置文件损坏，请重新初始化；笔记未被修改")
		}
		vault = c.Vault
	}
	vault, err = CanonicalVault(vault)
	if err != nil {
		return Location{}, fmt.Errorf("资料目录不存在或不可访问")
	}
	key := vault
	if runtime.GOOS == "windows" {
		key = strings.ToLower(key)
	}
	dir := filepath.Join(base, "instances", digest([]byte(key)))
	if err = privateDir(dir); err != nil {
		return Location{}, err
	}
	return Location{dir, filepath.Join(dir, "service.log"), vault}, nil
}
func atomicConfig(p string, b []byte) error {
	tmp := p + "." + randomID() + ".tmp"
	if err := os.WriteFile(tmp, b, 0600); err != nil {
		return err
	}
	defer os.Remove(tmp)
	return os.Rename(tmp, p)
}
func readState(loc Location) (State, error) {
	var s State
	b, e := os.ReadFile(filepath.Join(loc.Dir, "runtime.json"))
	if e != nil {
		return s, e
	}
	e = json.Unmarshal(b, &s)
	if e != nil {
		return s, e
	}
	u, e := url.Parse(s.URL)
	if e != nil || u.Scheme != "http" || u.Hostname() != "127.0.0.1" || u.Port() == "" || u.User != nil || u.Path != "" || u.RawQuery != "" || u.Fragment != "" || s.Vault != loc.Vault || s.Product != "memo" || len(s.Token) != 64 || len(s.Instance) != 64 {
		return s, fmt.Errorf("服务记录无效")
	}
	return s, nil
}
func client() *http.Client {
	return &http.Client{Timeout: 2 * time.Second, Transport: &http.Transport{Proxy: nil}, CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}
}
func authenticate(ctx context.Context, loc Location) (State, error) {
	s, e := readState(loc)
	if e != nil {
		return State{}, e
	}
	req, e := http.NewRequestWithContext(ctx, "GET", s.URL+"/api/health", nil)
	if e != nil {
		return State{}, e
	}
	req.Header.Set("Authorization", "Bearer "+s.Token)
	r, e := client().Do(req)
	if e != nil {
		return State{}, e
	}
	defer r.Body.Close()
	var health struct {
		State
		ProtocolVersion int `json:"protocolVersion"`
	}
	if r.StatusCode != 200 {
		return State{}, fmt.Errorf("本地服务身份验证失败")
	}
	if e = json.NewDecoder(io.LimitReader(r.Body, 8192)).Decode(&health); e != nil {
		return State{}, e
	}
	if health.Product != "memo" || health.Vault != loc.Vault || health.Version != s.Version || health.PID != s.PID || (health.Version != Version && health.ProtocolVersion != 1) || subtle.ConstantTimeCompare([]byte(health.Instance), []byte(s.Instance)) != 1 {
		return State{}, fmt.Errorf("端口归属或服务版本不匹配，拒绝复用")
	}
	return s, nil
}
func acquire(loc Location, name string) (*os.File, error) {
	f, e := os.OpenFile(filepath.Join(loc.Dir, name), os.O_CREATE|os.O_RDWR, 0600)
	if e != nil {
		return nil, e
	}
	if e = lockFile(f); e != nil {
		f.Close()
		return nil, e
	}
	return f, nil
}
func Ensure(ctx context.Context, loc Location) (State, error) {
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()
	if s, e := authenticate(ctx, loc); e == nil {
		return s, nil
	}
	var launch *os.File
	for {
		var e error
		launch, e = acquire(loc, "launch.lock")
		if e == nil {
			break
		}
		select {
		case <-ctx.Done():
			return State{}, fmt.Errorf("启动等待超时（15 秒）；日志：%s", loc.Log)
		case <-time.After(100 * time.Millisecond):
		}
	}
	defer launch.Close()
	if s, e := authenticate(ctx, loc); e == nil {
		return s, nil
	}
	// A live unverified process must never be treated as Memo or replaced.
	probe, e := acquire(loc, "service.lock")
	if e != nil {
		return State{}, fmt.Errorf("资料库服务锁被占用且身份无法验证，拒绝连接；日志：%s", loc.Log)
	}
	probe.Close()
	log, e := os.OpenFile(loc.Log, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if e != nil {
		return State{}, e
	}
	defer log.Close()
	exe, e := os.Executable()
	if e != nil {
		return State{}, e
	}
	if e = startDetached(exe, []string{"serve", "--vault", loc.Vault}, log); e != nil {
		return State{}, fmt.Errorf("服务启动失败：%w；日志：%s", e, loc.Log)
	}
	for {
		if s, e := authenticate(ctx, loc); e == nil {
			return s, nil
		}
		select {
		case <-ctx.Done():
			return State{}, fmt.Errorf("健康检查未在 15 秒内通过；日志：%s", loc.Log)
		case <-time.After(100 * time.Millisecond):
		}
	}
}
func Call(ctx context.Context, loc Location, method, route string, payload any, out any) error {
	s, e := Ensure(ctx, loc)
	if e != nil {
		return e
	}
	var b bytes.Buffer
	if payload != nil {
		if e = json.NewEncoder(&b).Encode(payload); e != nil {
			return e
		}
	}
	req, e := http.NewRequestWithContext(ctx, method, s.URL+route, &b)
	if e != nil {
		return e
	}
	req.Header.Set("Authorization", "Bearer "+s.Token)
	req.Header.Set("Content-Type", "application/json")
	r, e := client().Do(req)
	if e != nil {
		return fmt.Errorf("服务连接中断，请重试；笔记不会被自动重写")
	}
	defer r.Body.Close()
	if r.StatusCode >= 400 {
		var body struct {
			Error Error `json:"error"`
		}
		if e = json.NewDecoder(r.Body).Decode(&body); e != nil {
			return fmt.Errorf("服务返回错误 %d", r.StatusCode)
		}
		body.Error.Status = r.StatusCode
		return &body.Error
	}
	return json.NewDecoder(r.Body).Decode(out)
}
func Open(ctx context.Context, loc Location, browser string) (map[string]any, error) {
	s, e := Ensure(ctx, loc)
	if e != nil {
		return nil, e
	}
	var ticket struct {
		Ticket string `json:"ticket"`
	}
	if e = Call(ctx, loc, "POST", "/api/ticket", struct{}{}, &ticket); e != nil {
		return nil, e
	}
	u := s.URL + "/#ticket=" + ticket.Ticket
	if browser == "system" {
		if e = systemBrowser(u); e != nil {
			return nil, fmt.Errorf("系统浏览器启动失败；本地服务仍可用：%w", e)
		}
	}
	return map[string]any{"url": u, "vault": loc.Vault, "version": Version, "browser": browser, "log": loc.Log}, nil
}
func Doctor(ctx context.Context, loc Location) map[string]any {
	s, e := authenticate(ctx, loc)
	m := map[string]any{"version": Version, "vault": loc.Vault, "log": loc.Log, "running": e == nil, "runtimeRequired": false}
	if e == nil {
		m["clientVersion"] = Version
		m["version"] = s.Version
		m["url"] = s.URL
		m["pid"] = s.PID
		m["instance"] = s.Instance
	} else {
		m["detail"] = "未运行或服务身份未验证"
	}
	return m
}
