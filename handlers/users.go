package handlers

import (
	"fmt"
	"net/http"
	"todo/domain"
)

func (s *Server) registerUser() http.HandlerFunc	{
	// This is the payload that we need.
	var payload domain.RegisterPayload

	// Call the validate payload.
	return validatePayload(func(writer http.ResponseWriter, request *http.Request) {
		fmt.Printf("payload %v", payload)
	}, &payload)
}