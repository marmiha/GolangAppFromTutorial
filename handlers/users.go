package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"todo/domain"
)

func (s *Server) registerUser() http.HandlerFunc	{
	// Handler function is used fro HTTP response. It handles the requests.
	return func(writer http.ResponseWriter, request *http.Request) {
		// This is the payload that we need.
		payload := domain.RegisterPayload{}
		// We will try to decode our request.body with our RegisterPayload struct.
		decodingError := json.NewDecoder(request.Body).Decode(&payload)

		if decodingError != nil {
			// Call our defined handler function for Http BadRequest status code.
			badRequestResponse(writer, decodingError)
			return
		}

		fmt.Printf("payload %v", payload)
		return
	}
}