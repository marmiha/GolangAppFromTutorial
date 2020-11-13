package handlers

import (
	"github.com/dgrijalva/jwt-go"
	jwtRequest "github.com/dgrijalva/jwt-go/request"
	"net/http"
	"os"
	"strings"
	"todo/domain"
)


var (
	contextUserKey = "currentUser"	// Request context key for the current logged in user.
	contextTodoKey = "currentTodo"
)

// Extracts the string from the Authorization header and trims it of the prefix 'Bearer '.
var AuthenticationHeaderExtractorFilter = jwtRequest.PostExtractionFilter{
	Extractor: jwtRequest.HeaderExtractor{"Authorization"},
	Filter: func(s string) (string, error) {
		tokenString := strings.TrimPrefix(s, "Bearer ")
		return tokenString, nil
	}}

func ParseToken(r *http.Request, tokenClaims *domain.JWTTokenClaims) (*jwt.Token, error) {
	// This is where we check for our token and save it inside our context.
	token, err := jwtRequest.ParseFromRequest(r, &AuthenticationHeaderExtractorFilter, func(token *jwt.Token) (interface{}, error) {
		key := []byte(os.Getenv("JWT_KEY"))
		return key, nil
	}, jwtRequest.WithClaims(tokenClaims))
	return token, err
}

func (s *Server) currentUserFromContext(r *http.Request) *domain.User {
	return r.Context().Value(contextUserKey).(*domain.User)
}

func (s *Server) todoFromContext(r *http.Request) *domain.Todo {
	return r.Context().Value(contextTodoKey).(*domain.Todo)
}

