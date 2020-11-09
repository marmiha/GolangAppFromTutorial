package handlers

import (
	"fmt"
	"github.com/go-chi/chi"
	"net/http"
)



func (s *Server) setupEndpoints(r *chi.Mux) {
	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/users", func(r chi.Router) {
			r.Post("/register", s.registerUser)
		})
		r.Route("/test", func(r chi.Router) {
			r.Use(Authenticator)
			r.Get("/protected", func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte(fmt.Sprintf("protected area. hi")))
			})
		})
	})
}
