package postgres

import (
	"github.com/go-pg/pg/v10"
	"todo/domain"
)

type TodoRepository struct {
	DB *pg.DB
}

func NewTodoRepository(DB *pg.DB) *TodoRepository {
	return &TodoRepository{DB: DB}
}

func (t TodoRepository) Create(todo *domain.Todo) (*domain.Todo, error) {
	_, err := t.DB.Model(todo).Returning("*").Insert()
	if err != nil {
		return nil, err
	}
	return todo, err
}
