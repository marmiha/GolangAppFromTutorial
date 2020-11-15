package domain

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Todo struct {
	Entities
	tableName struct{} `pg:"todos,alias:todo"` // Default values would be the same.

	Title       string `json:"title" pg:"title"`
	Description string `json:"description" pg:"description"`
	Completed   bool   `json:"completed" pg:"completed,default:false"`

	UserId int64 `json:"user_id" pg:"user_id"`
}

func (t Todo) IsOwner(user *User) bool {
	return t.UserId == user.Id
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

type PatchTodoPayload struct {
	// These are pointers so that we can detect if any of these were omitted or not.
	// If these would not be pointers than Completed would default to false if omitted
	// thus we would not know if the user specified this field to be patched or not.
	Title       *string `json:"title,omitempty"`
	Description *string `json:"description,omitempty"`
	Completed   *bool   `json:"completed,omitempty"`
}

func (p PatchTodoPayload) Validate() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.Title, validation.Length(1, 30)),
		validation.Field(&p.Description, validation.Length(0, 256)),
	)
}

func (d *Domain) PatchTodo(payload PatchTodoPayload, todo *Todo) error {

	if payload.Completed != nil {
		todo.Completed = *payload.Completed
	}

	if payload.Description != nil {
		todo.Description = *payload.Description
	}

	if payload.Title != nil {
		todo.Title = *payload.Title
	}

	err := d.DB.TodoRepository.Patch(todo)
	if err != nil {
		return err
	}

	return nil
}

func (d *Domain) CreateTodo(payload CreateTodoPayload, user *User) (*Todo, error) {
	data := &Todo{
		Title:       payload.Title,
		Description: payload.Description,
		Completed:   false,
		UserId:      user.Id,
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

func (d *Domain) GetTodoById(id int64) (*Todo, error) {
	todo, err := d.DB.TodoRepository.GetById(id)
	if err != nil {
		return nil, err
	}
	return todo, nil
}

func (d *Domain) DeleteTodo(todo *Todo) error {
	err := d.DB.TodoRepository.Delete(todo)
	if err != nil {
		return err
	}
	return nil
}
