package domain

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"golang.org/x/crypto/bcrypt"
)

// Our custom claims for the JWT token field.
type JWTTokenClaims struct {
	UserId         int64              `json:"user_id"`
	Username       string             `json:"username"`
	StandardClaims jwt.StandardClaims `json:"standard_claims"`
}

// This will always return valid for now. In the future you could
// validate the claims before signing them.
func (J JWTTokenClaims) Valid() error {
	return J.StandardClaims.Valid()
}

// This is the payload we will get when a user wants to register.
// On the right hand side we have the json attribute names.
type RegisterPayload struct {
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
	Username        string `json:"username"`
}

// go-ozzo validator implementation for our register payload package.
func (rp RegisterPayload) Validate() error {
	return validation.ValidateStruct(&rp,
		validation.Field(&rp.Email, validation.Required, validation.NotNil, is.Email),
		validation.Field(&rp.Username, validation.Required, validation.NotNil, validation.Length(3, 15)),
		validation.Field(&rp.Password, validation.Required, validation.NotNil, validation.Length(5, 50)),
		validation.Field(&rp.ConfirmPassword, validation.Required, validation.NotNil, validation.In(rp.Password).ErrorObject(ErrPasswordsDoNotMatch)),
	)
}

// This will be used for logging the user in.
type LoginPayload struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// Validation for our LoginPayload
func (lp LoginPayload) Validate() error {
	return validation.ValidateStruct(&lp,
		validation.Field(&lp.Email, validation.When(lp.Username == "", validation.Required.ErrorObject(ErrUsernameOrEmailRequired), validation.NotNil, is.Email)),
		validation.Field(&lp.Username, validation.Length(3, 15)),
		validation.Field(&lp.Password, validation.Required, validation.NotNil, validation.Length(5, 50)),
	)
}

func (domain *Domain) Login(payload LoginPayload) (*User, error) {
	var user *User
	var err error

	// Get user by email.
	email := payload.Email
	if len(email) > 0 {
		user, err = domain.DB.UserRepository.GetByEmail(email)
		if err != nil {
			return nil, ErrInvalidLoginCredentials
		}
	}

	// Get user by username.
	if user == nil {
		user, err = domain.DB.UserRepository.GetByUsername(payload.Username)
		if err != nil {
			return nil, ErrInvalidLoginCredentials
		}
	}

	err = user.CheckPassword(payload.Password)
	// Invalid password.
	if err != nil {
		fmt.Printf("Invalid password")
		return nil, ErrInvalidLoginCredentials
	}

	// Correct password.
	return user, nil
}

// Business logic for registration
func (domain *Domain) Register(payload RegisterPayload) (*User, error) {
	// Unique email checkpoint.
	userExist, _ := domain.DB.UserRepository.GetByEmail(payload.Email)
	if userExist != nil {
		// Email is already taken, reject the registration with the corresponding passwordError.
		return nil, ErrUserWithEmailAlreadyExists
	}

	// Unique username checkpoint.
	userExist, _ = domain.DB.UserRepository.GetByUsername(payload.Username)
	if userExist != nil {
		// Username is already taken, reject the registration with the corresponding passwordError.
		return nil, ErrUserWithUsernameAlreadyExists
	}

	// Getting the password hash.
	passwordHash, passwordError := hashPassword(payload.Password)
	if passwordError != nil {
		return nil, passwordError
	}

	// User struct using our verified data.
	userData := &User{
		Username: payload.Username,
		Email:    payload.Email,
		Password: *passwordHash,
	}

	// Try to create our newly registered user in the database.
	user, createUserError := domain.DB.UserRepository.Create(userData)
	if createUserError != nil {
		return nil, createUserError
	}

	return user, nil
}

// Function for hashing passwords.
func hashPassword(password string) (*string, error) {
	bytePassword := []byte(password)
	passwordHash, err := bcrypt.GenerateFromPassword(bytePassword, bcrypt.DefaultCost)

	if err != nil {
		return nil, err
	}

	password = string(passwordHash)
	return &password, nil
}
