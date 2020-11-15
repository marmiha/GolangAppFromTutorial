package handlers

import (
	"net/http"
	"todo/domain"
)

type getTodosResponse struct {
	Todos []*domain.Todo `json:"todos"`
}

type createTodoResponse struct {
	Todo *domain.Todo `json:"todo"`
}

type patchTodoResponse struct {
	Todo *domain.Todo `json:"todo"`
}

func (s *Server) patchTodo(writer http.ResponseWriter, request *http.Request) {
	var payload domain.PatchTodoPayload

	next := validatePayload(func(writer http.ResponseWriter, request *http.Request) {

		todo := s.todoFromContext(request)
		err := s.Domain.PatchTodo(payload, todo)

		if err != nil {
			badRequestResponse(writer, err)
			return
		}

		response := patchTodoResponse{
			Todo: todo,
		}
		jsonResponse(writer, response, http.StatusOK)

	}, &payload)

	next.ServeHTTP(writer, request)
}

func (s *Server) deleteTodo(writer http.ResponseWriter, request *http.Request) {
	todo := s.todoFromContext(request)
	err := s.Domain.DeleteTodo(todo)
	if err != nil {
		badRequestResponse(writer, err)
		return
	}
	successfulDeleteResponse(writer)
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

func (s *Server) getTodos(writer http.ResponseWriter, request *http.Request) {

	user := s.currentUserFromContext(request)
	todos, err := s.Domain.GetTodosOfUser(user)

	if err != nil {
		badRequestResponse(writer, err)
		return
	}

	responseBody := getTodosResponse{
		Todos: todos,
	}

	jsonResponse(writer, responseBody, http.StatusOK)
}