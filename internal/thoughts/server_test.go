package thoughts

import (
	"bufio"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func testServer(t *testing.T) (*Server, *httptest.Server) {
	t.Helper()
	s := newServer(testStore(t), State{Product: "memo", Version: Version, Instance: randomID(), Token: randomID()})
	ts := httptest.NewServer(s.Handler())
	s.state.URL = ts.URL
	t.Cleanup(ts.Close)
	return s, ts
}
func request(t *testing.T, ts *httptest.Server, method, route, body string, headers map[string]string, cookie *http.Cookie) *http.Response {
	t.Helper()
	r, e := http.NewRequest(method, ts.URL+route, strings.NewReader(body))
	if e != nil {
		t.Fatal(e)
	}
	for k, v := range headers {
		if k == "Host" {
			r.Host = v
		} else {
			r.Header.Set(k, v)
		}
	}
	if cookie != nil {
		r.AddCookie(cookie)
	}
	res, e := ts.Client().Do(r)
	if e != nil {
		t.Fatal(e)
	}
	t.Cleanup(func() { res.Body.Close() })
	return res
}
func TestRequestSecurity(t *testing.T) {
	s, ts := testServer(t)
	auth := "Bearer " + s.state.Token
	cases := []struct {
		name, method, route string
		headers             map[string]string
		status              int
	}{
		{"unauthenticated", "GET", "/api/notes", nil, 401},
		{"wrong host", "GET", "/api/notes", map[string]string{"Authorization": auth, "Host": "evil.example"}, 403},
		{"evil origin", "POST", "/api/notes", map[string]string{"Authorization": auth, "Origin": "https://evil.example"}, 403},
		{"cross site", "GET", "/api/notes", map[string]string{"Authorization": auth, "Sec-Fetch-Site": "cross-site"}, 403},
		{"authenticated", "GET", "/api/health", map[string]string{"Authorization": auth}, 200},
		{"encoded traversal", "GET", "/api/notes?path=%2e%2e%2fsecret.md", map[string]string{"Authorization": auth}, 400},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			r := request(t, ts, c.method, c.route, "", c.headers, nil)
			if r.StatusCode != c.status {
				t.Fatalf("status %d", r.StatusCode)
			}
		})
	}
	r := request(t, ts, "GET", "/", "", nil, nil)
	if r.StatusCode != 200 || !strings.Contains(r.Header.Get("Content-Security-Policy"), "frame-ancestors 'none'") {
		t.Fatal("missing static assets or CSP")
	}
}
func TestBrowserTicketCSRFAndConflict(t *testing.T) {
	s, ts := testServer(t)
	r := request(t, ts, "POST", "/api/ticket", "{}", map[string]string{"Authorization": "Bearer " + s.state.Token}, nil)
	var ticket struct {
		Ticket string `json:"ticket"`
	}
	if e := json.NewDecoder(r.Body).Decode(&ticket); e != nil {
		t.Fatal(e)
	}
	body := `{"ticket":"` + ticket.Ticket + `"}`
	r = request(t, ts, "POST", "/api/session", body, map[string]string{"Origin": ts.URL}, nil)
	if r.StatusCode != 200 {
		t.Fatalf("session: %d", r.StatusCode)
	}
	cookies := r.Cookies()
	if len(cookies) != 1 || !cookies[0].HttpOnly || cookies[0].SameSite != http.SameSiteStrictMode {
		t.Fatal("unsafe cookie")
	}
	cookie := cookies[0]
	var sess struct {
		CSRF string `json:"csrf"`
	}
	json.NewDecoder(r.Body).Decode(&sess)
	r = request(t, ts, "POST", "/api/session", body, map[string]string{"Origin": ts.URL}, nil)
	if r.StatusCode != 401 {
		t.Fatal("ticket replay allowed")
	}
	noteBody := `{"path":"中文.md","content":"# hello"}`
	r = request(t, ts, "POST", "/api/notes", noteBody, map[string]string{"Origin": ts.URL}, cookie)
	if r.StatusCode != 403 {
		t.Fatal("missing CSRF accepted")
	}
	h := map[string]string{"Origin": ts.URL, "X-Memo-CSRF": sess.CSRF}
	r = request(t, ts, "POST", "/api/notes", noteBody, h, cookie)
	if r.StatusCode != 200 {
		b, _ := io.ReadAll(r.Body)
		t.Fatalf("create %d %s", r.StatusCode, b)
	}
	r = request(t, ts, "PUT", "/api/notes", `{"path":"中文.md","content":"stale","revision":"wrong"}`, h, cookie)
	if r.StatusCode != 409 {
		t.Fatal("conflict accepted")
	}
}
func TestExternalEditSSEWithinTwoSeconds(t *testing.T) {
	s, ts := testServer(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go s.watch(ctx)
	req, _ := http.NewRequestWithContext(ctx, "GET", ts.URL+"/api/events", nil)
	req.Header.Set("Authorization", "Bearer "+s.state.Token)
	r, e := ts.Client().Do(req)
	if e != nil {
		t.Fatal(e)
	}
	defer r.Body.Close()
	scanner := bufio.NewScanner(r.Body)
	if !scanner.Scan() || scanner.Text() != "event: connected" {
		t.Fatal("no connection event")
	}
	// Consume initial inventory notification before measuring a new external write.
	lines := make(chan string, 16)
	go func() {
		for scanner.Scan() {
			select {
			case lines <- scanner.Text():
			case <-ctx.Done():
				return
			}
		}
	}()
	time.Sleep(650 * time.Millisecond)
	for len(lines) > 0 {
		<-lines
	}
	started := time.Now()
	if e = os.WriteFile(filepath.Join(s.store.Vault, "external.md"), []byte("# 外部写入"), 0600); e != nil {
		t.Fatal(e)
	}
	timeout := time.NewTimer(2 * time.Second)
	defer timeout.Stop()
	for {
		select {
		case line := <-lines:
			if line == "event: changed" {
				n, e := s.store.Read("external.md")
				if e != nil || n.Content != "# 外部写入" {
					t.Fatal("wrong content")
				}
				t.Logf("external write to SSE: %s", time.Since(started))
				return
			}
		case <-timeout.C:
			t.Fatal("refresh exceeded two seconds")
		}
	}
}

func TestPageStyleNonceKeepsScriptPolicyStrict(t *testing.T) {
	s, ts := testServer(t)
	for _, route := range []string{"/", "/index.html"} {
		r := request(t, ts, "GET", route, "", nil, nil)
		b, e := io.ReadAll(r.Body)
		if e != nil {
			t.Fatal(e)
		}
		if !strings.Contains(string(b), `name="thoughts-style-nonce" content="`+s.state.Instance+`"`) {
			t.Fatal("missing trusted style nonce")
		}
		policy := r.Header.Get("Content-Security-Policy")
		if !strings.Contains(policy, "script-src 'self';") || !strings.Contains(policy, "style-src 'self' 'nonce-"+s.state.Instance+"'") {
			t.Fatal("unexpected CSP")
		}
		if strings.Contains(string(b), s.state.Token) {
			t.Fatal("API credential exposed in HTML")
		}
	}
}
