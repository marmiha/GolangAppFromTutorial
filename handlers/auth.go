package handlers

import (
	"context"
	"github.com/dgrijalva/jwt-go"
	jwtRequest "github.com/dgrijalva/jwt-go/request"
	"net/http"
	"os"
	"strings"
	"todo/domain"
)


var (
	contextUserKey = "currentUser"	// Request context key for the current logged in user.
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
