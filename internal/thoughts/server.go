package thoughts

import (
	"context"
	"crypto/subtle"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"thoughts.local/thoughts/internal/webui"
)

type session struct {
	CSRF    string
	Expires time.Time
}
type Server struct {
	store           *Store
	state           State
	mu              sync.Mutex
	tickets         map[string]time.Time
	sessions        map[string]session
	clients         map[chan string]bool
	preferencesPath string
	preferencesMu   sync.Mutex
	themeMu         sync.Mutex
	readerMu        sync.Mutex
	readerPath      string
	editorLeases    map[string]editorLease
	shutdown        func()
	openBrowser     func(string) error
}

func newServer(store *Store, state State) *Server {
	return &Server{store: store, state: state, tickets: map[string]time.Time{}, sessions: map[string]session{}, clients: map[chan string]bool{}, openBrowser: systemBrowser}
}
func jsonResponse(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
func fail(w http.ResponseWriter, e error) {
	var p *Error
	if !errors.As(e, &p) {
		p = &Error{500, "internal_error", "操作失败，请检查本地日志或目录权限。"}
	}
	jsonResponse(w, p.Status, map[string]any{"error": p})
}
func decode(w http.ResponseWriter, r *http.Request, v any) error {
	r.Body = http.MaxBytesReader(w, r.Body, MaxNoteSize+16384)
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if e := d.Decode(v); e != nil {
		return problem(400, "invalid_request", "请求格式或大小不符合要求。")
	}
	var extra any
	if e := d.Decode(&extra); e != io.EOF {
		return problem(400, "invalid_request", "仅允许一个 JSON 请求。")
	}
	return nil
}
func equal(a, b string) bool { return a != "" && subtle.ConstantTimeCompare([]byte(a), []byte(b)) == 1 }
func (s *Server) bearer(r *http.Request) bool {
	return equal(r.Header.Get("Authorization"), "Bearer "+s.state.Token)
}
func (s *Server) cleanupLocked() {
	now := time.Now()
	for k, t := range s.tickets {
		if now.After(t) {
			delete(s.tickets, k)
		}
	}
	for k, v := range s.sessions {
		if now.After(v.Expires) {
			delete(s.sessions, k)
		}
	}
}
func (s *Server) sessionCookieName() string {
	return "memo_session_" + digest([]byte(s.state.Instance))[:16]
}
func (s *Server) browserSession(r *http.Request) (session, bool) {
	c, e := r.Cookie(s.sessionCookieName())
	if e != nil {
		return session{}, false
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.cleanupLocked()
	v, ok := s.sessions[c.Value]
	return v, ok
}
func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/health", func(w http.ResponseWriter, r *http.Request) {
		if !s.bearer(r) {
			fail(w, problem(401, "unauthorized", "需要本地授权。"))
			return
		}
		h := s.state
		h.Token = ""
		jsonResponse(w, 200, struct {
			State
			ProtocolVersion int `json:"protocolVersion"`
		}{h, 1})
	})
	mux.HandleFunc("POST /api/ticket", func(w http.ResponseWriter, r *http.Request) {
		if !s.bearer(r) {
			fail(w, problem(401, "unauthorized", "需要本地授权。"))
			return
		}
		s.mu.Lock()
		defer s.mu.Unlock()
		s.cleanupLocked()
		if len(s.tickets) >= 64 {
			fail(w, problem(429, "too_many_tickets", "请稍后重试。"))
			return
		}
		ticket := randomID()
		s.tickets[ticket] = time.Now().Add(2 * time.Minute)
		jsonResponse(w, 200, map[string]string{"ticket": ticket})
	})
	mux.HandleFunc("POST /api/session", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Origin") != s.state.URL {
			fail(w, problem(403, "invalid_origin", "请从本地 Thoughts 页面打开。"))
			return
		}
		var body struct {
			Ticket string `json:"ticket"`
		}
		if e := decode(w, r, &body); e != nil {
			fail(w, e)
			return
		}
		s.mu.Lock()
		defer s.mu.Unlock()
		s.cleanupLocked()
		expires, ok := s.tickets[body.Ticket]
		if !ok || time.Now().After(expires) {
			fail(w, problem(401, "ticket_expired", "打开链接已失效，请让 Agent 重新打开 Thoughts。"))
			return
		}
		delete(s.tickets, body.Ticket)
		if len(s.sessions) >= 128 {
			fail(w, problem(429, "too_many_sessions", "本地会话过多，请重启服务。"))
			return
		}
		key := randomID()
		v := session{randomID(), time.Now().Add(12 * time.Hour)}
		s.sessions[key] = v
		http.SetCookie(w, &http.Cookie{Name: s.sessionCookieName(), Value: key, Path: "/", HttpOnly: true, SameSite: http.SameSiteStrictMode, MaxAge: 12 * 60 * 60})
		jsonResponse(w, 200, map[string]string{"csrf": v.CSRF, "vault": filepath.Base(s.state.Vault), "vaultPath": s.state.Vault, "version": Version})
	})
	mux.HandleFunc("GET /api/session", func(w http.ResponseWriter, r *http.Request) {
		v, ok := s.browserSession(r)
		if !ok {
			fail(w, problem(401, "unauthorized", "请让 Agent 打开 Thoughts，获取此浏览器的本地会话。"))
			return
		}
		jsonResponse(w, 200, map[string]string{"csrf": v.CSRF, "vault": filepath.Base(s.state.Vault), "vaultPath": s.state.Vault, "version": Version})
	})
	mux.Handle("/api/notes", s.authorized(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			if p := r.URL.Query().Get("path"); p != "" {
				n, e := s.store.Read(p)
				if e != nil {
					fail(w, e)
					return
				}
				jsonResponse(w, 200, n)
				return
			}
			notes, e := s.store.List()
			if e != nil {
				fail(w, e)
				return
			}
			for i := range notes {
				notes[i].Content = ""
			}
			jsonResponse(w, 200, notes)
		case "POST", "PUT":
			var body struct {
				Path     string `json:"path"`
				Content  string `json:"content"`
				Revision string `json:"revision"`
			}
			if e := decode(w, r, &body); e != nil {
				fail(w, e)
				return
			}
			n, e := s.store.Write(body.Path, body.Content, body.Revision, r.Method == "POST")
			if e != nil {
				fail(w, e)
				return
			}
			s.broadcast()
			jsonResponse(w, 200, n)
		default:
			fail(w, problem(405, "method_not_allowed", "此操作不属于 M0。"))
		}
	})))
	mux.Handle("/api/events", s.authorized(http.HandlerFunc(s.events)))
	mux.Handle("/api/preferences", s.authorized(http.HandlerFunc(s.preferences)))
	mux.Handle("GET /api/reader-state", s.authorized(http.HandlerFunc(s.readerState)))
	mux.Handle("PUT /api/reader-state", s.authorized(http.HandlerFunc(s.readerState)))
	mux.Handle("POST /api/editor-state", s.authorized(http.HandlerFunc(s.editorState)))
	mux.Handle("POST /api/control/shutdown", s.authorized(http.HandlerFunc(s.stopForUpgrade)))
	s.registerManagement(mux)
	s.registerThemes(mux)
	s.registerImages(mux)
	mux.Handle("POST /api/open-browser", s.authorized(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ticket := randomID()
		s.mu.Lock()
		s.cleanupLocked()
		if len(s.tickets) >= 64 {
			s.mu.Unlock()
			fail(w, problem(429, "too_many_tickets", "Please try again later."))
			return
		}
		s.tickets[ticket] = time.Now().Add(2 * time.Minute)
		s.mu.Unlock()
		if err := s.openBrowser(s.state.URL + "/#ticket=" + ticket); err != nil {
			s.mu.Lock()
			delete(s.tickets, ticket)
			s.mu.Unlock()
			fail(w, problem(422, "browser_unavailable", "Could not open the external browser. Ask your agent to open Thoughts in that browser."))
			return
		}
		jsonResponse(w, 200, map[string]bool{"opened": true})
	})))
	mux.Handle("GET /api/search", s.authorized(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := strings.TrimSpace(r.URL.Query().Get("q"))
		if len(query) > 512 {
			fail(w, problem(400, "query_too_long", "搜索内容请控制在 512 字节以内。"))
			return
		}
		results, e := s.store.SearchScoped(query, r.URL.Query().Get("includeArchived") == "1")
		if e != nil {
			fail(w, e)
			return
		}
		jsonResponse(w, 200, results)
	})))
	mux.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) { fail(w, problem(404, "not_found", "接口不存在。")) })
	assets := http.FileServer(http.FS(webui.Files()))
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" && r.Method != "HEAD" {
			w.WriteHeader(405)
			return
		}
		if r.URL.Path == "/" || r.URL.Path == "/index.html" {
			file, e := webui.Files().Open("index.html")
			if e != nil {
				fail(w, e)
				return
			}
			defer file.Close()
			body, e := io.ReadAll(file)
			if e != nil {
				fail(w, e)
				return
			}
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			if r.Method == "GET" {
				io.WriteString(w, strings.Replace(string(body), "<head>", "<head><meta name=\"thoughts-style-nonce\" content=\""+s.state.Instance+"\">", 1))
			}
			return
		}
		assets.ServeHTTP(w, r)
	})
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Referrer-Policy", "no-referrer")
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("Cross-Origin-Resource-Policy", "same-origin")
		w.Header().Set("Content-Security-Policy", "default-src 'none'; script-src 'self'; style-src 'self' 'nonce-"+s.state.Instance+"'; style-src-attr 'unsafe-inline'; font-src 'self'; img-src 'self' data:; connect-src 'self'; base-uri 'none'; form-action 'self'; frame-ancestors 'none'")
		if r.Host != strings.TrimPrefix(s.state.URL, "http://") {
			fail(w, problem(403, "invalid_host", "拒绝非预期的本地主机名。"))
			return
		}
		if origin := r.Header.Get("Origin"); origin != "" && origin != s.state.URL {
			fail(w, problem(403, "invalid_origin", "拒绝跨站请求。"))
			return
		}
		if r.Header.Get("Sec-Fetch-Site") == "cross-site" {
			fail(w, problem(403, "cross_site", "拒绝跨站请求。"))
			return
		}
		mux.ServeHTTP(w, r)
	})
}
func (s *Server) authorized(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if s.bearer(r) {
			next.ServeHTTP(w, r)
			return
		}
		v, ok := s.browserSession(r)
		if !ok {
			fail(w, problem(401, "unauthorized", "本地会话已失效，请让 Agent 重新打开 Thoughts。"))
			return
		}
		if r.Method != "GET" && r.Method != "HEAD" {
			if r.Header.Get("Origin") != s.state.URL || !equal(r.Header.Get("X-Memo-CSRF"), v.CSRF) {
				fail(w, problem(403, "invalid_csrf", "本地写入授权校验失败。"))
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}
func (s *Server) broadcast() {
	s.broadcastEvent("changed")
}
func (s *Server) broadcastEvent(kind string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for c := range s.clients {
		select {
		case c <- kind:
		default:
		}
	}
}
func (s *Server) events(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(405)
		return
	}
	f, ok := w.(http.Flusher)
	if !ok {
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("X-Accel-Buffering", "no")
	c := make(chan string, 1)
	s.mu.Lock()
	s.clients[c] = true
	s.mu.Unlock()
	defer func() { s.mu.Lock(); delete(s.clients, c); s.mu.Unlock() }()
	fmt.Fprint(w, "event: connected\ndata: {}\n\n")
	f.Flush()
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-r.Context().Done():
			return
		case <-ticker.C:
			fmt.Fprint(w, ": heartbeat\n\n")
			f.Flush()
		case kind := <-c:
			fmt.Fprintf(w, "event: %s\ndata: {}\n\n", kind)
			f.Flush()
		}
	}
}
func (s *Server) watch(ctx context.Context) {
	last := ""
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			library, e := s.store.Library()
			fingerprint := "unavailable"
			if e == nil {
				parts := []string{"assets:" + library.AssetVersion}
				notes := append([]Note{}, library.Inbox...)
				for _, project := range library.Projects {
					parts = append(parts, fmt.Sprintf("project:%s:%t", project.Path, project.Archived))
					notes = append(notes, project.Nodes...)
				}
				for _, n := range notes {
					parts = append(parts, n.Path+":"+n.Revision)
				}
				if rows, err := s.store.TrashList(); err == nil {
					for _, entry := range rows {
						parts = append(parts, "trash:"+entry.ID+":"+entry.Revision)
					}
				}
				sort.Strings(parts)
				fingerprint = digest([]byte(strings.Join(parts, "\n")))
			}
			if fingerprint != last {
				last = fingerprint
				s.broadcast()
			}
		}
	}
}
func Serve(loc Location) error {
	lock, e := acquire(loc, "service.lock")
	if e != nil {
		return fmt.Errorf("资料库已有服务实例；请使用 thoughts open")
	}
	defer lock.Close()
	store, e := NewStore(loc.Vault)
	if e != nil {
		return e
	}
	defer store.Close()
	// Config directories can differ between host profiles. The vault itself also
	// owns a lock, so misconfigured homes cannot create two writers for one vault.
	if e = store.root.MkdirAll(".memo/runtime", 0700); e != nil {
		return e
	}
	vaultGuard, e := store.root.OpenFile(".memo/runtime/service.lock", os.O_CREATE|os.O_RDWR, 0600)
	if e != nil {
		return e
	}
	if e = lockFile(vaultGuard); e != nil {
		vaultGuard.Close()
		return fmt.Errorf("this library already has a service; use the same THOUGHTS_HOME in every client")
	}
	defer vaultGuard.Close()
	listener, e := net.Listen("tcp4", "127.0.0.1:0")
	if e != nil {
		return e
	}
	defer listener.Close()
	state := State{"memo", Version, loc.Vault, randomID(), randomID(), "http://" + listener.Addr().String(), os.Getpid()}
	app := newServer(store, state)
	app.preferencesPath = filepath.Join(loc.Dir, "preferences.json")
	app.readerPath = filepath.Join(loc.Dir, "reader-state.json")
	server := &http.Server{Handler: app.Handler(), ReadHeaderTimeout: 5 * time.Second, IdleTimeout: 60 * time.Second, MaxHeaderBytes: 16 * 1024}
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	app.shutdown = cancel
	b, _ := json.Marshal(state)
	if e = atomicConfig(filepath.Join(loc.Dir, "runtime.json"), b); e != nil {
		return e
	}
	defer os.Remove(filepath.Join(loc.Dir, "runtime.json"))
	go app.watch(ctx)
	go func() {
		<-ctx.Done()
		c, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		server.Shutdown(c)
	}()
	fmt.Fprintln(os.Stderr, time.Now().UTC().Format(time.RFC3339), "Thoughts", Version, "service started")
	e = server.Serve(listener)
	if errors.Is(e, http.ErrServerClosed) {
		return nil
	}
	return e
}
