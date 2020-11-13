package handlers

import (
	"context"
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	validation "github.com/go-ozzo/ozzo-validation/v4"
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

func successfulDeleteResponse(writer http.ResponseWriter) {
	// Our response, what we reply back. Using the map[string]string we can
	// define json properties and their values.
	response := map[string]string {}
	jsonResponse(writer, response, http.StatusNoContent)
}


// For better code readability as we will use this a lot.
func internalServerErrorResponse(writer http.ResponseWriter, error error) {
	// Our response, what we reply back. Using the map[string]string we can
	// define json properties and their values.
	response := map[string]string {
		"error": error.Error(),
	}
	jsonResponse(writer, response, http.StatusInternalServerError)
}

// For better code readability as we will use this a lot.
func unauthorizedResponse(writer http.ResponseWriter, error error) {
	// Our response, what we reply back. Using the map[string]string we can
	// define json properties and their values.
	response := map[string]string {
		"error": error.Error(),
	}
	jsonResponse(writer, response, http.StatusUnauthorized)
}

func forbiddenResponse(writer http.ResponseWriter) {
	response := map[string]string {
		"error": "Forbidden",
	}
	jsonResponse(writer, response, http.StatusForbidden)
}

func notFoundResponse(writer http.ResponseWriter, error error) {
	response := map[string]string {
		"error": error.Error(),
	}
	jsonResponse(writer, response, http.StatusNotFound)
}

// Universal payload tester for Http endpoints. Throws validation errors or decoding errors.
func validatePayload(next http.HandlerFunc, payload validation.Validatable) http.HandlerFunc{
	return func(writer http.ResponseWriter, request *http.Request) {
		// We will try to decode our request.body with our RegisterPayload struct.
		// We will save it inside the payload pointer.
		decodingError := json.NewDecoder(request.Body).Decode(&payload)

		if decodingError != nil {
			// Call our defined handler function for Http BadRequest status code.
			badRequestResponse(writer, decodingError)
			return
		}

		// We close the body after this function exits.
		defer  request.Body.Close()

		// Validate the payload.
		if validationErrors := payload.Validate(); validationErrors != nil {
			// Bad request response.
			badRequestResponse(writer, validationErrors)
			return
		}

		// We will add a new filed to the context of our request.
		// As if we would do: req.payload = payload on a objet in javascript.
		// We can access it later.
		newPayloadContext := context.WithValue(request.Context(), "payload", payload)
		next.ServeHTTP(writer, request.WithContext(newPayloadContext))
	}
}