package list

import (
	"context"
	todo "github.com/dafuqqqyunglean/todoRestAPI"
	"github.com/dafuqqqyunglean/todoRestAPI/pkg/repository/cache"
	"github.com/dafuqqqyunglean/todoRestAPI/pkg/repository/sql"
)

type TodoListService interface {
	Create(ctx context.Context, userId int, list todo.TodoList) (int, error)
	GetAll(ctx context.Context, userId int) ([]todo.TodoList, error)
	GetById(ctx context.Context, userId, listId int) (todo.TodoList, error)
	Delete(ctx context.Context, userId, listId int) error
	Update(ctx context.Context, userId, listId int, input todo.UpdateListInput) error
}

type ImplTodoList struct {
	repo  sql.TodoListRepository
	cache cache.RedisCache
}

func NewTodoListService(repo sql.TodoListRepository, cache cache.RedisCache) *ImplTodoList {
	return &ImplTodoList{
		repo:  repo,
		cache: cache,
	}
}

func (s *ImplTodoList) Create(ctx context.Context, userId int, list todo.TodoList) (int, error) {
	return s.repo.Create(ctx, userId, list)
}

func (s *ImplTodoList) GetAll(ctx context.Context, userId int) ([]todo.TodoList, error) {
	return s.repo.GetAll(ctx, userId)
}

func (s *ImplTodoList) GetById(ctx context.Context, userId, listId int) (todo.TodoList, error) {
	list, err := s.cache.GetList(ctx, userId, listId)
	if err == nil {
		return list, nil
	}

	list, err = s.repo.GetById(ctx, userId, listId)
	if err != nil {
		return list, err
	}

	s.cache.SetList(ctx, userId, listId, list)

	return list, nil
}

func (s *ImplTodoList) Delete(ctx context.Context, userId, listId int) error {
	err := s.repo.Delete(ctx, userId, listId)
	if err != nil {
		return err
	}
	s.cache.Delete(ctx, userId, listId)

	return nil
}

func (s *ImplTodoList) Update(ctx context.Context, userId, listId int, input todo.UpdateListInput) error {
	if err := input.Validate(); err != nil {
		return err
	}

	err := s.repo.Update(ctx, userId, listId, input)
	if err != nil {
		return err
	}

	s.cache.Delete(ctx, userId, listId)
	return nil
}
