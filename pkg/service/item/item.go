package item

import (
	"context"
	todo "github.com/dafuqqqyunglean/todoRestAPI"
	"github.com/dafuqqqyunglean/todoRestAPI/pkg/repository/cache"
	"github.com/dafuqqqyunglean/todoRestAPI/pkg/repository/sql"
)

type TodoItemService interface {
	Create(ctx context.Context, userId, listId int, item todo.TodoItem) (int, error)
	GetAll(ctx context.Context, userId, listId int) ([]todo.TodoItem, error)
	GetById(ctx context.Context, userId, itemId int) (todo.TodoItem, error)
	Delete(ctx context.Context, userId, itemId int) error
	Update(ctx context.Context, userId, itemId int, input todo.UpdateItemInput) error
}

type ImplTodoItem struct {
	repo     sql.TodoItemRepository
	listRepo sql.TodoListRepository
	cache    cache.RedisCache
}

func NewTodoItemService(repo sql.TodoItemRepository, listRepo sql.TodoListRepository, cache cache.RedisCache) *ImplTodoItem {
	return &ImplTodoItem{
		repo:     repo,
		listRepo: listRepo,
		cache:    cache,
	}
}

func (s *ImplTodoItem) Create(ctx context.Context, userId, listId int, item todo.TodoItem) (int, error) {
	_, err := s.listRepo.GetById(ctx, userId, listId)
	if err != nil {
		return 0, err
	}

	return s.repo.Create(ctx, listId, item)
}

func (s *ImplTodoItem) GetAll(ctx context.Context, userId, listId int) ([]todo.TodoItem, error) {
	items, err := s.repo.GetAll(ctx, userId, listId)
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (s *ImplTodoItem) GetById(ctx context.Context, userId, itemId int) (todo.TodoItem, error) {
	item, err := s.cache.GetItem(ctx, userId, itemId)
	if err == nil {
		return item, nil
	}

	item, err = s.repo.GetById(ctx, userId, itemId)
	if err != nil {
		return item, err
	}

	s.cache.SetItem(ctx, userId, itemId, item)

	return item, nil
}

func (s *ImplTodoItem) Delete(ctx context.Context, userId, itemId int) error {
	err := s.repo.Delete(ctx, userId, itemId)
	if err != nil {
		return err
	}
	s.cache.Delete(ctx, userId, itemId)

	return nil
}

func (s *ImplTodoItem) Update(ctx context.Context, userId, itemId int, input todo.UpdateItemInput) error {
	if err := input.Validate(); err != nil {
		return err
	}

	err := s.repo.Update(ctx, userId, itemId, input)
	if err != nil {
		return err
	}

	s.cache.Delete(ctx, userId, itemId)
	return nil
}
