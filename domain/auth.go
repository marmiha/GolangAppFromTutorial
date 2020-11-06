package domain

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

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

	// Getting the password.
	passwordHash, passwordError := domain.setPassword(payload.Password)
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

func (domain *Domain) setPassword(password string) (*string, error) {
	return nil, nil
}
