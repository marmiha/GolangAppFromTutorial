package handlers

import (
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
			s = strings.TrimPrefix(s, "Bearer ")
			fmt.Printf(s)
			return s, nil
		},
	}
)

func Authenticator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// This is where we check for our token and save it inside our context.
		tokenClaims := domain.JWTTokenClaims{}
		token, err := jwtRequest.ParseFromRequest(r, &AuthenticationHeaderExtractorFilter, func(token *jwt.Token) (interface{}, error) {
			pass := []byte(os.Getenv("JWT_SECRET"))
			return pass, nil
		}, jwtRequest.WithClaims(&tokenClaims))


		if err != nil {
			internalServerErrorResponse(w, err)
			return
		}

		fmt.Printf("token: %v", token)

		// Token is authenticated, pass it through
		next.ServeHTTP(w, r)
	})
}