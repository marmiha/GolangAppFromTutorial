package handlers

import (
	"fmt"
	"net/http"
	"todo/domain"
)

func (s *Server) registerUser(writer http.ResponseWriter, request *http.Request) {
	// This is the payload that we need.
	var payload domain.RegisterPayload

	// Call the validate payload.
	next := validatePayload(func(writer http.ResponseWriter, request *http.Request) {
		// If successful then this will be executed with next.ServeHTTP().
		fmt.Println(payload)
	}, &payload)

	// Serve the returned
	next.ServeHTTP(writer, request)
}
