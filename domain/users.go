package domain

import (
	"encoding/base64"
	"github.com/dgrijalva/jwt-go"
	"os"
	"strconv"
	"time"
)

// User struct will be used for our postgres database. Right hand side we have our
// json marshal definitions. If we write json:"-" that field won't be marshalled.
type User struct {
	// tableName is an optional field that specifies custom table name and alias.
	// By default go-pg generates table name and alias from struct name.
	tableName struct{} `pg:"\"user\",alias:u"` // Default values would be the same.

	ID       int64  `json:"id" pg:"id,pk"`
	Username string `json:"username" pg:",unique"`
	Email    string `json:"email" pg:",unique"`
	Password string `json:"-" pg:""`

	CreatedAt time.Time `json:"created_at" pg:"default:now()"`
	UpdatedAt time.Time `json:"updated_at" pg:"default:now()"`
}

func (user *User) GenerateToken() (*string, error) {
	// The tokens will expire in one day. Unix function converts the
	// date to the seconds passed so int64.
	expiresAt := time.Now().Add(time.Minute * 1)

	// Populate the claims.
	claims := JWTTokenClaims{
		UserId:   user.ID,
		Username: user.Username,
		StandardClaims: jwt.StandardClaims{
			Id:        strconv.FormatInt(user.ID, 10),
			Issuer:    "TodoApp",
			ExpiresAt: expiresAt.Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Sign the token with the JWT_SECRET environment variable.
	key, _ := base64.URLEncoding.DecodeString(os.Getenv("JWT_KEY"))
	signedString, err := token.SignedString(key)

	if err != nil {
		return nil, err
	}

	return &signedString, nil
}

// Function for setting the user password.
func (user *User) setPassword(password string) error {
	hashedPassword, err := hashPassword(password)
	if err != nil {
		return err
	}
	user.Password = *hashedPassword
	return nil
}


