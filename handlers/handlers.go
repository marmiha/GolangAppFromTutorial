package handlers

import (
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"net/http"
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

func jsonResponse(writer http.ResponseWriter, responseData interface{}, httpStatusCode int) {
	// Let's respond with a error corresponding to httpStatusCode.
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(httpStatusCode)

	// Our response, what we reply back. Using the map[string]string we can
	// define json properties and their values.
	// If data is empty, then send an empty object.
	if responseData == nil {
		responseData = map[string]string{}
	}

	// Encode the response and if something goes wrong with the encoding then
	// return a internal server error response.
	if encodingError := json.NewEncoder(writer).Encode(responseData); encodingError != nil {
		http.Error(writer, encodingError.Error(), http.StatusInternalServerError)
		return
	}
	return
}

// For better code readability as we will use this a lot.
func badRequestResponse(writer http.ResponseWriter, error error) {
	// Our response, what we reply back. Using the map[string]string we can
	// define json properties and their values.
	response := map[string]string {
		"error": error.Error(),
	}
	jsonResponse(writer, response, http.StatusBadRequest)
}
