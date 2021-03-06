package domain

import (
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"os"
	"strconv"
	"time"
)

// User struct will be used for our postgres database. Right hand side we have our
// json marshal definitions. If we write json:"-" that field won't be marshalled.
type User struct {
	Entities
	// tableName is an optional field that specifies custom table name and alias.
	// By default go-pg generates table name and alias from struct name.
	tableName struct{} `pg:"users,alias:u"` // Default values would be the same.

	Username string `json:"username" pg:",unique"`
	Email    string `json:"email" pg:",unique"`
	Password string `json:"-" pg:""`

	Todos []*Todo `json:"todos" pg:"rel:has-many"`
}

func (user *User) GenerateToken() (*string, error) {
	// The tokens will expire in one day. Unix function converts the
	// date to the seconds passed so int64.
	expiresAt := time.Now().Add(time.Hour * 10)

	// Populate the claims.
	claims := JWTTokenClaims{
		UserId:   user.Id,
		Username: user.Username,
		StandardClaims: jwt.StandardClaims{
			Id:        strconv.FormatInt(user.Id, 10),
			Issuer:    "TodoApp",
			ExpiresAt: expiresAt.Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Sign the token with the JWT_KEY environment variable.
	key := []byte(os.Getenv("JWT_KEY"))
	signedString, err := token.SignedString(key)

	if err != nil {
		return nil, err
	}

	return &signedString, nil
}

func (user *User) CheckPassword(password string) error {
	bytePassword, byteHashedPassword := []byte(password), []byte(user.Password)
	return bcrypt.CompareHashAndPassword(byteHashedPassword, bytePassword)
}

// Function for setting the user password.
func (user *User) SetPassword(password string) error {
	hashedPassword, err := hashPassword(password)
	if err != nil {
		return err
	}
	user.Password = *hashedPassword
	return nil
}

// Used for getting the user with id from database. We use this for our inserting the user
// in our HTTP request context.
func (domain *Domain) GetUserById(id int64) (*User, error) {
	user, err := domain.DB.UserRepository.GetById(id)
	if err != nil {
		return nil, err
	}
	return user, nil
}


