package handlers

import (
	"fmt"
	"github.com/go-chi/chi"
	"net/http"
	"todo/domain"
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
		})

		// Functionality testing routes.
		r.Route("/test", func(r chi.Router) {
			r.Use(s.WithUserAuthentication)
			r.Get("/protected", func(w http.ResponseWriter, r *http.Request) {
				user := r.Context().Value(contextUserKey).(*domain.User)
				_, _ = w.Write([]byte(fmt.Sprintf("Welcome %v!\n", user.Username)))
			})
		})
	})
}
