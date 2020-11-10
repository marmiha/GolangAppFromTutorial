package handlers

import (
	"net/http"
	"todo/domain"
)

type authenticationResponse struct {
	User  *domain.User `json:"user"`
	Token string       `json:"token"`
}

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

		// Generate jwt token.
		token, err := user.GenerateToken()
		if err != nil {
			badRequestResponse(writer, err)
			return
		}

		// Fill our response struct with the user and the token.
		responseBody := authenticationResponse{
			User:  user,
			Token: *token,
		}

		// Return the newly created object.
		jsonResponse(writer, responseBody, http.StatusCreated)

	}, &payload)

	// Serve the returned
	next.ServeHTTP(writer, request)
}

func (s *Server) loginUser(writer http.ResponseWriter, request *http.Request) {
	var payload domain.LoginPayload

	next := validatePayload(func(writer http.ResponseWriter, request *http.Request) {
		// Call the registration domain logic.
		user, err := s.Domain.Login(payload)

		// Return the error if it occurs in the domain logic.
		if err != nil {
			badRequestResponse(writer, err)
			return
		}

		// Generate jwt token.
		token, err := user.GenerateToken()
		if err != nil {
			badRequestResponse(writer, err)
			return
		}

		// Fill our response struct with the user and the token.
		responseBody := authenticationResponse{
			User:  user,
			Token: *token,
		}

		// Return the newly created object.
		jsonResponse(writer, responseBody, http.StatusCreated)

	}, &payload)
	// Serve the returned
	next.ServeHTTP(writer, request)
}
