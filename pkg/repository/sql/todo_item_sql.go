package sql

import (
	"context"
	_ "embed"
	"fmt"
	todo "github.com/dafuqqqyunglean/todoRestAPI"
	"github.com/jmoiron/sqlx"
	"strings"
)

type TodoItemRepository interface {
	Create(ctx context.Context, listId int, item todo.TodoItem) (int, error)
	GetAll(ctx context.Context, userId, listId int) ([]todo.TodoItem, error)
	GetById(ctx context.Context, userId, itemId int) (todo.TodoItem, error)
	Delete(ctx context.Context, userId, itemId int) error
	Update(ctx context.Context, userId, itemId int, input todo.UpdateItemInput) error
}

type TodoItemPostgres struct {
	db *sqlx.DB
}

func NewTodoItemPostgres(db *sqlx.DB) *TodoItemPostgres {
	return &TodoItemPostgres{
		db: db,
	}
}

//go:embed query/CreateItem.sql
var createItem string

//go:embed query/CreateListsItems.sql
var createListsItems string

func (r *TodoItemPostgres) Create(ctx context.Context, listId int, item todo.TodoItem) (int, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}

	var itemId int
	row := tx.QueryRow(createItem, item.Title, item.Description)
	if err := row.Scan(&itemId); err != nil {
		tx.Rollback()
		return 0, err
	}

	_, err = tx.Exec(createListsItems, listId, itemId)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	return itemId, tx.Commit()
}

//go:embed query/GetAllItems.sql
var getAllItems string

func (r *TodoItemPostgres) GetAll(ctx context.Context, userId, listId int) ([]todo.TodoItem, error) {
	var items []todo.TodoItem

	if err := r.db.SelectContext(ctx, &items, getAllItems, listId, userId); err != nil {
		return nil, err
	}

	return items, nil
}

//go:embed query/GetItemById.sql
var getItemById string

func (r *TodoItemPostgres) GetById(ctx context.Context, userId, itemId int) (todo.TodoItem, error) {
	var item todo.TodoItem

	if err := r.db.GetContext(ctx, &item, getItemById, itemId, userId); err != nil {
		return item, err
	}

	return item, nil
}

//go:embed query/DeleteItem.sql
var deleteItem string

func (r *TodoItemPostgres) Delete(ctx context.Context, userId, itemId int) error {
	_, err := r.db.ExecContext(ctx, deleteItem, userId, itemId)

	return err
}

//go:embed query/UpdateItem.sql
var updateItem string

func (r *TodoItemPostgres) Update(ctx context.Context, userId, itemId int, input todo.UpdateItemInput) error {
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

	if input.Done != nil {
		setValues = append(setValues, fmt.Sprintf("done = $%d", argId))
		args = append(args, *input.Done)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf(updateItem, setQuery, argId, argId+1)
	args = append(args, userId, itemId)

	_, err := r.db.ExecContext(ctx, query, args...)
	return err
}
