package handlers

import (
	"fmt"
	"net/http"
	"todo/domain"
)

func (s *Server) registerUser(writer http.ResponseWriter, request *http.Request) {
	// This is the payload that we need.
	payload := new(domain.RegisterPayload)

	// Call the validate payload.
	next := validatePayload(func(writer http.ResponseWriter, request *http.Request) {
		// If successful then this will be executed.
		fmt.Println(payload)
	}, payload)

	next.ServeHTTP(writer, request)
}
