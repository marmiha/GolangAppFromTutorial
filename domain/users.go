package domain

import "time"

// User struct will be used for our postgres database. Right hand side we have our
// json marshal definitions. If we write json:"-" that field won't be marshalled.
type User struct {
	// tableName is an optional field that specifies custom table name and alias.
	// By default go-pg generates table name and alias from struct name.
	tableName struct{} `pg:"user, alias:user"` // Default values would be the same.

	ID       int64  `json:"id" pg:"id,pk"`
	Username string `json:"username" pg:"username"`
	Email    string `json:"email" pg:"email"`
	Password string `json:"-" pg:"password"`

	CreatedAt time.Time `json:"created_at" pg:"created_at,default:now()"`
	UpdatedAt time.Time `json:"updated_at" pg:"updated_at.default:now()"`
}
