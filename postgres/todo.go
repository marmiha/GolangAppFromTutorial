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

func (t TodoRepository) GetTodosOfUser(user *domain.User) ([]*domain.Todo, error) {
	var todos = make([]*domain.Todo, 0)
	err := t.DB.Model(&todos).Where("user_id = ?", user.Id).Select()
	if err != nil {
		return nil, err
	}
	return todos, nil
}


