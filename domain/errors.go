package domain

import (
	"errors"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// These will be our business logic errors that can occur other than validation. This way code is more readable and it
// also ensures consistency of errors naming across our codebase.
var (
	ErrUserWithEmailAlreadyExists    = errors.New("user with specified email already exists")
	ErrUserWithUsernameAlreadyExists = errors.New("user with specified username already exist")
	ErrNoResult                      = errors.New("no result")
	ErrInvalidLoginCredentials       = errors.New("invalid login credentials")
	ErrUserNotFound					= errors.New("user not found")
)

// These errors are used for go-ozzo validation in our business logic, mostly for our payloads.
var (
	ErrUsernameOrEmailRequired = validation.NewError("validation_key_username_or_email_required", "username or email is required")
	ErrPasswordsDoNotMatch     = validation.NewError("validation_key_passwords_do_not_match", "passwords don't match")
)
