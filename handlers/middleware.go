package handlers

import (
	"context"
	"github.com/go-chi/chi"
	"net/http"
	"strconv"
	"todo/domain"
)

// Injects User from the JWT token into the context.
func (s *Server) WithUserAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// This is where we check for our token and save it inside our context.
		var tokenClaims domain.JWTTokenClaims
		_, err := ParseToken(r, &tokenClaims)

		// If signature is invalid or the token does not exist.
		if err != nil {
			unauthorizedResponse(w, err)
			return
		}

		// If the token is not valid (expired...).
		if err := tokenClaims.Valid(); err != nil {
			unauthorizedResponse(w, err)
			return
		}

		// Get the user from the database and insert it into the context.
		user, err := s.Domain.GetUserById(tokenClaims.UserId)
		if err != nil {
			internalServerErrorResponse(w, err)
			return
		}

		ctx := context.WithValue(r.Context(), contextUserKey, user)
		// Token is authenticated, pass the request with the user in context to the next handler.
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Injects the _Todo into the context.
func (s *Server) WithTodo(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		todo := new(domain.Todo)

		if stringId := chi.URLParam(r, "todo_id"); stringId != "" {
			todoId, err := strconv.ParseInt(stringId, 0, 0)
			if err != nil {
				badRequestResponse(w, err)
				return
			}

			todo, err = s.Domain.GetTodoById(todoId)
			if err != nil {
				notFoundResponse(w, domain.ErrNoResult)
				return
			}
		}

		ctx := context.WithValue(r.Context(), contextTodoKey, todo)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (s *Server) WithOwner(subjectContextKey string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			currentUser := s.currentUserFromContext(request)
			// Get the subject from the context.
			subject := request.Context().Value(subjectContextKey).(domain.HaveOwner)
			// If it's not the owner then reply with forbidden response.
			if !subject.IsOwner(currentUser) {
				forbiddenResponse(writer)
				return
			}
			next.ServeHTTP(writer, request)
		})
	}
}