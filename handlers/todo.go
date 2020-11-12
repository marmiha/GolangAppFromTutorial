package handlers

import (
	"net/http"
	"todo/domain"
)

type createTodoResponse struct {
	Todo *domain.Todo `json:"todo"`
}

func (s *Server) createTodo(writer http.ResponseWriter, request *http.Request) {
	var payload domain.CreateTodoPayload

	next := validatePayload(func(writer http.ResponseWriter, request *http.Request) {
		// Get the current user.
		user := s.currentUserFromContext(request)
		todo, err := s.Domain.CreateTodo(payload, user)


		if err != nil {
			badRequestResponse(writer, err)
			return
		}

		responseBody := createTodoResponse{
			Todo: todo,
		}

		// Return the newly created object.
		jsonResponse(writer, responseBody, http.StatusCreated)

	}, &payload)

	// Serve the returned
	next.ServeHTTP(writer, request)
}

