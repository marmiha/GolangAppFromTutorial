package domain

import "time"

// User struct will be used for our postgres database. Right hand side we have our
// json marshal definitions. If we write json:"-" that field won't be marshalled.
type User struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"-"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
