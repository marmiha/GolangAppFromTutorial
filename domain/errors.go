package domain

import "errors"

// These will be our business logic errors that can occur. This way code is more readable and it also ensures
// consistency of errors naming across our codebase.
var (
	ErrUserWithEmailAlreadyExists    = errors.New("user with specified email already exists")
	ErrUserWithUsernameAlreadyExists = errors.New("user with specified username already exist")
	ErrNoResult                      = errors.New("no result")
)
