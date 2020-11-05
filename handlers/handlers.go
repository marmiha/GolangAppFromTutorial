package handlers

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"time"
)

type Server struct {

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

func NewServer() *Server {
	return &Server{}
}

func SetupRouter() *chi.Mux {
	server := NewServer()
	router := chi.NewRouter()
	setupMiddleWare(router)
	server.setupEndpoints(router)
	return router
}
