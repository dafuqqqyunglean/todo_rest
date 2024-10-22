package sql

import (
	"context"
	_ "embed"
	"fmt"
	todo "github.com/dafuqqqyunglean/todoRestAPI"
	"github.com/jmoiron/sqlx"
	"strings"
)

type TodoListRepository interface {
	Create(ctx context.Context, userId int, list todo.TodoList) (int, error)
	GetAll(ctx context.Context, userId int) ([]todo.TodoList, error)
	GetById(ctx context.Context, userId, listId int) (todo.TodoList, error)
	Delete(ctx context.Context, userId, listId int) error
	Update(ctx context.Context, userId, listId int, input todo.UpdateListInput) error
}

type TodoListPostgres struct {
	db *sqlx.DB
}

func NewTodoListPostgres(db *sqlx.DB) *TodoListPostgres {
	return &TodoListPostgres{db: db}
}

//go:embed query/CreateList.sql
var createList string

//go:embed query/CreateUsersLists.sql
var createUsersLists string

func (r *TodoListPostgres) Create(ctx context.Context, userId int, list todo.TodoList) (int, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}

	var id int
	row := tx.QueryRow(createList, list.Title, list.Description) // stores information about the returned row from db
	if err := row.Scan(&id); err != nil {
		tx.Rollback()
		return 0, err
	}

	_, err = tx.Exec(createUsersLists, userId, id)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	return id, tx.Commit()
}

//go:embed query/GetAllLists.sql
var getAllLists string

func (r *TodoListPostgres) GetAll(ctx context.Context, userId int) ([]todo.TodoList, error) {
	var lists []todo.TodoList

	err := r.db.SelectContext(ctx, &lists, getAllLists, userId)

	return lists, err
}

//go:embed query/GetListById.sql
var getListById string

func (r *TodoListPostgres) GetById(ctx context.Context, userId, listId int) (todo.TodoList, error) {
	var list todo.TodoList

	err := r.db.GetContext(ctx, &list, getListById, userId, listId)

	return list, err
}

//go:embed query/DeleteList.sql
var deleteList string

func (r *TodoListPostgres) Delete(ctx context.Context, userId, listId int) error {
	_, err := r.db.ExecContext(ctx, deleteList, userId, listId)

	return err
}

//go:embed query/UpdateList.sql
var updateList string

func (r *TodoListPostgres) Update(ctx context.Context, userId, listId int, input todo.UpdateListInput) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if input.Title != nil {
		setValues = append(setValues, fmt.Sprintf("title = $%d", argId))
		args = append(args, *input.Title)
		argId++
	}

	if input.Description != nil {
		setValues = append(setValues, fmt.Sprintf("description = $%d", argId))
		args = append(args, *input.Description)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf(updateList, setQuery, argId, argId+1)
	args = append(args, listId, userId)

	_, err := r.db.ExecContext(ctx, query, args...)
	return err
}
