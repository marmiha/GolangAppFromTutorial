package domain

// UserRepository is used as an interface for our domain business logic, which will be
// implemented in the handlers package, this way we decouple the business logic from
// database implementations with specifying the arguments and return types.
type UserRepository interface {
	GetById(id int64) (*User, error)
	GetByEmail(email string) (*User, error)
	GetByUsername(username string) (*User, error)
	Create(user *User) (*User, error)
}

type TodoRepository interface {
	GetById(id int64) (*Todo, error)
	Create(todo *Todo) (*Todo, error)
	Delete(todo *Todo) error
	Patch(todo *Todo) error
	GetTodosOfUser(user *User) ([]*Todo, error)
}

// Used for implementing on struct that have an owner.
type HaveOwner interface {
	IsOwner(user *User) bool
}

// DB struct encapsulates all the interfaces handling the database gateway interfaces.
type DB struct {
	UserRepository UserRepository
	TodoRepository TodoRepository
}

// Struct used for common variables across the business logic (domain package).
type Domain struct {
	DB DB
}
