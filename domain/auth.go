package domain

// This is the payload we will get when a user wants to register.
// On the right hand side we have the json attribute names.
type RegisterPayload struct {
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
	Username        string `json:"username"`
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
