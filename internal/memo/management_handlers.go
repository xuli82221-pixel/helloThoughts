package memo

import "net/http"

func (s *Server) registerManagement(mux *http.ServeMux) {
	register := func(pattern string, handler http.HandlerFunc) { mux.Handle(pattern, s.authorized(handler)) }
	s.registerProjectLifecycle(register)
	register("GET /api/library", func(w http.ResponseWriter, r *http.Request) {
		value, err := s.store.Library()
		if err != nil {
			fail(w, err)
			return
		}
		jsonResponse(w, 200, value)
	})
	register("POST /api/projects", func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Name     string `json:"name"`
			Language string `json:"language"`
			Template bool   `json:"template"`
		}
		if err := decode(w, r, &body); err != nil {
			fail(w, err)
			return
		}
		defer s.broadcast()
		value, err := s.store.CreateProject(body.Name, body.Language, body.Template)
		if err != nil {
			fail(w, err)
			return
		}
		jsonResponse(w, 201, value)
	})
	register("POST /api/notes/rename", func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Path     string `json:"path"`
			Name     string `json:"name"`
			Revision string `json:"revision"`
		}
		if err := decode(w, r, &body); err != nil {
			fail(w, err)
			return
		}
		defer s.broadcast()
		value, err := s.store.Rename(body.Path, body.Name, body.Revision)
		if err != nil {
			fail(w, err)
			return
		}
		jsonResponse(w, 200, value)
	})
	register("POST /api/notes/trash", func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Path     string `json:"path"`
			Revision string `json:"revision"`
		}
		if err := decode(w, r, &body); err != nil {
			fail(w, err)
			return
		}
		defer s.broadcast()
		value, err := s.store.Trash(body.Path, body.Revision)
		if err != nil {
			fail(w, err)
			return
		}
		jsonResponse(w, 200, value)
	})
	register("POST /api/notes/move", func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Path      string `json:"path"`
			Directory string `json:"directory"`
			Revision  string `json:"revision"`
		}
		if err := decode(w, r, &body); err != nil {
			fail(w, err)
			return
		}
		defer s.broadcast()
		note, err := s.store.Move(body.Path, body.Directory, body.Revision)
		if err != nil {
			fail(w, err)
			return
		}
		jsonResponse(w, 200, note)
	})
	register("GET /api/trash", func(w http.ResponseWriter, r *http.Request) {
		value, err := s.store.TrashList()
		if err != nil {
			fail(w, err)
			return
		}
		jsonResponse(w, 200, value)
	})
	register("POST /api/trash/restore", func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			ID       string `json:"id"`
			Revision string `json:"revision"`
		}
		if err := decode(w, r, &body); err != nil {
			fail(w, err)
			return
		}
		defer s.broadcast()
		value, err := s.store.Restore(body.ID, body.Revision)
		if err != nil {
			fail(w, err)
			return
		}
		jsonResponse(w, 200, value)
	})
}
