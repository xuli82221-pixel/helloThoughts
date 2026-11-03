package memo

import (
	"encoding/json"
	"errors"
	"math"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

type ReaderState struct {
	Schema   int                `json:"schema"`
	Tabs     []string           `json:"tabs"`
	Selected string             `json:"selected"`
	Page     string             `json:"page"`
	Scrolls  map[string]float64 `json:"scrolls"`
}
type editorLease struct {
	Busy     bool
	Sequence int64
	Expires  time.Time
}

var editorID = regexp.MustCompile(`^[a-f0-9-]{16,64}$`)

func emptyReaderState() ReaderState {
	return ReaderState{Schema: 1, Tabs: []string{}, Page: "empty", Scrolls: map[string]float64{}}
}
func validReaderCollection(id string) bool {
	if id == "inbox" || id == "trash" || id == "search" || id == "archive" || id == "projects" {
		return true
	}
	return strings.HasPrefix(id, "project:") && validProjectPath(strings.TrimPrefix(id, "project:"))
}
func validReaderTab(id string) bool {
	if validReaderCollection(id) {
		return true
	}
	return strings.HasPrefix(id, "note:") && validPath(strings.TrimPrefix(id, "note:")) == nil
}
func validateReaderState(p ReaderState) bool {
	if p.Schema != 1 || len(p.Tabs) > 200 || len(p.Scrolls) > 500 || p.Selected != "" && validPath(p.Selected) != nil {
		return false
	}
	if p.Page != "empty" && p.Page != "document" && !validReaderCollection(p.Page) {
		return false
	}
	seen := map[string]bool{}
	for _, id := range p.Tabs {
		if !validReaderTab(id) || seen[id] {
			return false
		}
		seen[id] = true
	}
	for id, pos := range p.Scrolls {
		if !validReaderTab(id) || math.IsNaN(pos) || math.IsInf(pos, 0) || pos < 0 || pos > 1e8 {
			return false
		}
	}
	return true
}
func (s *Server) readerState(w http.ResponseWriter, r *http.Request) {
	s.readerMu.Lock()
	defer s.readerMu.Unlock()
	if s.readerPath == "" {
		fail(w, problem(503, "preferences_unavailable", "Reader state storage is unavailable."))
		return
	}
	if r.Method == "GET" {
		p := emptyReaderState()
		b, e := os.ReadFile(s.readerPath)
		if e != nil && !errors.Is(e, os.ErrNotExist) {
			fail(w, problem(422, "preferences_unavailable", "Cannot read reader state."))
			return
		}
		if e == nil && (len(b) > 128*1024 || json.Unmarshal(b, &p) != nil || !validateReaderState(p)) {
			fail(w, problem(422, "preferences_unavailable", "Invalid reader state."))
			return
		}
		jsonResponse(w, 200, map[string]any{"exists": e == nil, "state": p})
		return
	}
	var p ReaderState
	if e := decode(w, r, &p); e != nil {
		fail(w, e)
		return
	}
	if !validateReaderState(p) {
		fail(w, problem(400, "invalid_request", "Invalid reader state; no note content is accepted."))
		return
	}
	b, _ := json.Marshal(p)
	if len(b) > 128*1024 || atomicConfig(s.readerPath, b) != nil {
		fail(w, problem(422, "preferences_unavailable", "Cannot save reader state."))
		return
	}
	jsonResponse(w, 200, map[string]bool{"saved": true})
}
func (s *Server) editorState(w http.ResponseWriter, r *http.Request) {
	var p struct {
		ID       string `json:"id"`
		Busy     bool   `json:"busy"`
		Sequence int64  `json:"sequence"`
	}
	if e := decode(w, r, &p); e != nil {
		fail(w, e)
		return
	}
	if !editorID.MatchString(p.ID) || p.Sequence < 0 {
		fail(w, problem(400, "invalid_request", "Invalid editor identity."))
		return
	}
	s.readerMu.Lock()
	defer s.readerMu.Unlock()
	if s.editorLeases == nil {
		s.editorLeases = map[string]editorLease{}
	}
	for id, lease := range s.editorLeases {
		if time.Now().After(lease.Expires) {
			delete(s.editorLeases, id)
		}
	}
	if old, ok := s.editorLeases[p.ID]; ok && p.Sequence < old.Sequence {
		jsonResponse(w, 200, map[string]bool{"saved": false})
		return
	}
	if _, ok := s.editorLeases[p.ID]; !ok && len(s.editorLeases) >= 256 {
		fail(w, problem(429, "too_many_sessions", "Too many active editors."))
		return
	}
	s.editorLeases[p.ID] = editorLease{p.Busy, p.Sequence, time.Now().Add(35 * time.Second)}
	jsonResponse(w, 200, map[string]bool{"saved": true})
}
func (s *Server) stopForUpgrade(w http.ResponseWriter, r *http.Request) {
	// Runtime control is restricted to the local bearer credential, not a webpage.
	if !s.bearer(r) {
		fail(w, problem(403, "forbidden", "Runtime control requires the local service credential."))
		return
	}
	var p struct {
		Instance string `json:"instance"`
	}
	if e := decode(w, r, &p); e != nil {
		fail(w, e)
		return
	}
	if !equal(p.Instance, s.state.Instance) {
		fail(w, problem(409, "instance_changed", "The service instance changed."))
		return
	}
	s.readerMu.Lock()
	defer s.readerMu.Unlock()
	for _, lease := range s.editorLeases {
		if lease.Busy && time.Now().Before(lease.Expires) {
			fail(w, problem(409, "editor_busy", "Save browser drafts and finish uploads before updating."))
			return
		}
	}
	if s.shutdown == nil {
		fail(w, problem(503, "unavailable", "Shutdown is unavailable."))
		return
	}
	jsonResponse(w, 200, map[string]bool{"stopping": true})
	go func() { time.Sleep(100 * time.Millisecond); s.shutdown() }()
}
