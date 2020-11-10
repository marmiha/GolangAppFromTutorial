package handlers

import (
	"encoding/base64"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	jwtRequest "github.com/dgrijalva/jwt-go/request"
	"net/http"
	"os"
	"strings"
	"todo/domain"
)

var (
	// Extracts the string from the Authorization header and trims it of the prefix 'Bearer '.
	AuthenticationHeaderExtractorFilter = jwtRequest.PostExtractionFilter{
		Extractor: jwtRequest.HeaderExtractor{"Authorization"},
		Filter: func(s string) (string, error) {
			tokenString := strings.TrimPrefix(s, "Bearer ")
			return tokenString, nil
		},
	}
)

func Authenticator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// This is where we check for our token and save it inside our context.
		var tokenClaims domain.JWTTokenClaims
		token, err := jwtRequest.ParseFromRequest(r, &AuthenticationHeaderExtractorFilter, func(token *jwt.Token) (interface{}, error) {
			key, _ := base64.URLEncoding.DecodeString(os.Getenv("JWT_KEY"))
			return key, nil
		}, jwtRequest.WithClaims(&tokenClaims))


		if err != nil {
			internalServerErrorResponse(w, err)
			return
		}

		if err := tokenClaims.Valid(); err != nil {
			unauthorizedResponse(w, err)
		}

		fmt.Printf("token: %v", token)

		// Token is authenticated, pass it through
		next.ServeHTTP(w, r)
	})
}