package domain

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"time"
)

type Todo struct {
	tableName struct{} `pg:"todos,alias:todo"` // Default values would be the same.

	Id          int64  `json:"id" pg:"id,pk"`
	Title       string `json:"title" pg:"title"`
	Description string `json:"description" pg:"description"`
	Completed   bool   `json:"completed" pg:"completed"`

	UserId int64 `json:"user_id" pg:"user_id"`

	CreatedAt time.Time `json:"created_at" pg:"default:now()"`
	UpdatedAt time.Time `json:"updated_at" pg:"default:now()"`
}

type CreateTodoPayload struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

func (c CreateTodoPayload) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Title, validation.Required, validation.Length(1, 30)),
		validation.Field(&c.Description, validation.Length(0, 256)),
	)
}

func (d *Domain) CreateTodo(payload CreateTodoPayload, user *User) (*Todo, error) {
	data := &Todo {
		Title: payload.Title,
		Description: payload.Description,
		Completed: false,
		UserId: user.Id,
	}

	todo, err := d.DB.TodoRepository.Create(data)
	if err != nil {
		return nil, err
	}

	return todo, nil
}

func (d *Domain) GetTodosOfUser(user *User) ([]*Todo, error) {
	todos, err := d.DB.TodoRepository.GetTodosOfUser(user)

	if err != nil {
		return nil, err
	}

	return todos, nil
}