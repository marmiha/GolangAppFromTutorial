package handlers

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"time"
	"todo/domain"
)

type Server struct {
	Domain *domain.Domain
}

func setupMiddleWare(router *chi.Mux) {
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Compress(6, "application/json"))
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
	router.Use(middleware.Timeout(60 * time.Second))
}

func NewServer(domain *domain.Domain) *Server {
	return &Server{Domain: domain}
}

// By passing the domain we have access to all of the domain functions so basically
// all of our business logic.
func SetupRouter(domain *domain.Domain) *chi.Mux {
	// Pass our server the business logic struct (the domain package).
	server := NewServer(domain)

	// Make a new router and attach the middleware.
	router := chi.NewRouter()
	setupMiddleWare(router)

	// Register our endpoints on the router.
	server.setupEndpoints(router)
	return router
}
