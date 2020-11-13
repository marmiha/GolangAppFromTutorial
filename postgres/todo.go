package postgres

import (
	"errors"
	"github.com/go-pg/pg/v10"
	"todo/domain"
)

type TodoRepository struct {
	DB *pg.DB
}

func (t TodoRepository) Delete(todo *domain.Todo) error {
	_, err := t.DB.Model(todo).WherePK().Delete()
	if err != nil {
		return err
	}
	return nil
}

func NewTodoRepository(DB *pg.DB) *TodoRepository {
	return &TodoRepository{DB: DB}
}

func (t TodoRepository) GetById(id int64) (*domain.Todo, error) {
	todo := new(domain.Todo)
	err := t.DB.Model(todo).Where("id = ?", id).First()
	if err != nil {
		// Check if error err is equal to postgres error ErrNoRows.
		if errors.Is(err, pg.ErrNoRows) {
			// Return our custom error which represent that user is not found.
			return nil, domain.ErrNoResult
		}
		return nil, err
	}
	return todo, nil
}

func (t TodoRepository) Update(todo *domain.Todo) (*domain.Todo, error) {
	_, err := t.DB.Model(todo).WherePK().Update()
	if err != nil {
		return nil, err
	}
	return todo, nil
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


