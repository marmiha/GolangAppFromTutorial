package domain

import "time"

type Todo struct {
	tableName struct{} `pg:"\"todos\",alias:todo"` // Default values would be the same.

	Id        int64  `json:"id" pg:"id,pk"`
	Title     string `json:"title" pg:"title"`
	Completed bool   `json:"completed" pg:"completed"`

	User *User `json:"user" pg:"rel:has-one"`

	CreatedAt time.Time `json:"created_at" pg:"default:now()"`
	UpdatedAt time.Time `json:"updated_at" pg:"default:now()"`
}