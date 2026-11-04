package thoughts

import "net/http"

func (s *Server) registerProjectLifecycle(register func(string, http.HandlerFunc)) {
	register("GET /api/project", func(w http.ResponseWriter, r *http.Request) {
		value, e := s.store.ProjectStatus(r.URL.Query().Get("path"))
		if e != nil {
			fail(w, e)
			return
		}
		jsonResponse(w, 200, value)
	})
	register("POST /api/projects/archive", func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Path     string `json:"path"`
			Revision string `json:"revision"`
			Archived bool   `json:"archived"`
		}
		if e := decode(w, r, &body); e != nil {
			fail(w, e)
			return
		}
		value, e := s.store.ArchiveProject(body.Path, body.Revision, body.Archived)
		if e != nil {
			fail(w, e)
			return
		}
		s.broadcast()
		jsonResponse(w, 200, value)
	})
	register("POST /api/projects/trash", func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Path     string `json:"path"`
			Revision string `json:"revision"`
		}
		if e := decode(w, r, &body); e != nil {
			fail(w, e)
			return
		}
		value, e := s.store.TrashProject(body.Path, body.Revision)
		if e != nil {
			fail(w, e)
			return
		}
		s.broadcast()
		jsonResponse(w, 200, value)
	})
	register("GET /api/project-trash", func(w http.ResponseWriter, r *http.Request) {
		value, e := s.store.ProjectTrashList()
		if e != nil {
			fail(w, e)
			return
		}
		jsonResponse(w, 200, value)
	})
	register("POST /api/project-trash/restore", func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			ID       string `json:"id"`
			Revision string `json:"revision"`
		}
		if e := decode(w, r, &body); e != nil {
			fail(w, e)
			return
		}
		value, e := s.store.RestoreProject(body.ID, body.Revision)
		if e != nil {
			fail(w, e)
			return
		}
		s.broadcast()
		jsonResponse(w, 200, value)
	})
}
