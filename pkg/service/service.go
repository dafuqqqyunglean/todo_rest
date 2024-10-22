package service

import (
	"context"
	"github.com/dafuqqqyunglean/todoRestAPI/pkg/repository/cache"
	"github.com/dafuqqqyunglean/todoRestAPI/pkg/repository/sql"
	"github.com/dafuqqqyunglean/todoRestAPI/pkg/service/auth"
	"github.com/dafuqqqyunglean/todoRestAPI/pkg/service/item"
	"github.com/dafuqqqyunglean/todoRestAPI/pkg/service/list"
	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	"time"
)

const (
	cacheKey = "todo_item:%d:%d"
	ttl      = time.Minute * 10
)

type Service struct {
	AuthService *auth.ImplAuthorizationService
	ListService *list.ImplTodoList
	ItemService *item.ImplTodoItem
}

func NewService(ctx context.Context, postgres *sqlx.DB, redis *redis.Client) *Service {
	authService := auth.NewAuthorizationService(sql.NewAuthorizationPostgres(postgres), ctx)
	todoLists := list.NewTodoListService(sql.NewTodoListPostgres(postgres), cache.NewRedisCache(redis, cacheKey, ttl))
	todoItems := item.NewTodoItemService(sql.NewTodoItemPostgres(postgres), todoLists, cache.NewRedisCache(redis, cacheKey, ttl))
	return &Service{
		AuthService: authService,
		ListService: todoLists,
		ItemService: todoItems,
	}
}
