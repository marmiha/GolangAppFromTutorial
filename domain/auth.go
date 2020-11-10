package domain

import (
	"github.com/dgrijalva/jwt-go"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"golang.org/x/crypto/bcrypt"
)

// Our custom claims for the JWT token field.
type JWTTokenClaims struct {
	UserId         int64  `json:"user_id"`
	Username       string `json:"username"`
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
