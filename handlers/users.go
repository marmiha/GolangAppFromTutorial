package handlers

import (
	"net/http"
	"todo/domain"
)

func (s *Server) registerUser(writer http.ResponseWriter, request *http.Request) {
	// This is the payload that we need.
	var payload domain.RegisterPayload

	// Call the validate payload.
	next := validatePayload(func(writer http.ResponseWriter, request *http.Request) {
		// Call the registration domain logic.
		user, err := s.Domain.Register(payload)

		// Return the error if it occurs in the domain logic.
		if err != nil {
			badRequestResponse(writer, err)
			return
		}

		// Return the newly created object.
		jsonResponse(writer, user, http.StatusCreated)
	}, &payload)

	// Serve the returned
	next.ServeHTTP(writer, request)
}
