package handlers

import (
	"github.com/go-chi/chi"
)

func (s *Server) setupEndpoints(r *chi.Mux) {
	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/users", func(r chi.Router) {
			r.Post("/register", s.registerUser)
			r.Post("/login", s.loginUser)
		})

		// Routes for todos.
		r.Route("/todos", func(r chi.Router) {
			r.Use(s.WithUserAuthentication)
			r.Post("/", s.createTodo)
			r.Get("/", s.getTodos)

			r.Route("/{todo_id}", func(r chi.Router) {
				r.Use(s.WithTodo)
				r.Use(s.WithOwner(contextTodoKey))
				r.Delete("/", s.deleteTodo)
				r.Patch("/", s.patchTodo)
			})
		})
	})
}
