package domain

// UserRepository is used as an interface for our domain business logic, which will be
// implemented in the handlers package, this way we decouple the business logic from
// database implementations with specifying the arguments and return types.
type UserRepository interface {
	GetByEmail(email string) (*User, error)
	GetByUsername(username string) (*User, error)
	Create(user *User) (*User, error)
}

// DB struct encapsulates all the interfaces handling the database gateway interfaces.
type DB struct {
	UserRepository UserRepository
}

// Struct used for common variables across the business logic (domain package).
type Domain struct {
	DB DB
}
