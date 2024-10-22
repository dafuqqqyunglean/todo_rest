package main

import (
	"context"
	"fmt"
	"github.com/dafuqqqyunglean/todoRestAPI/config"
	"github.com/dafuqqqyunglean/todoRestAPI/pkg/api"
	"github.com/dafuqqqyunglean/todoRestAPI/pkg/api/middlewares"
	"github.com/dafuqqqyunglean/todoRestAPI/pkg/api/utility"
	"github.com/dafuqqqyunglean/todoRestAPI/pkg/repository"
	"github.com/dafuqqqyunglean/todoRestAPI/pkg/service"
	"go.uber.org/zap"
	"log"
	"os"
	"os/signal"
	"syscall"
)

// @title Todo App API
// @version 1.0
// @description API Server for TodoList Application

// @host localhost:8000
// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	mainCtx := context.Background()
	ctx, cancel := context.WithCancel(mainCtx)
	defer cancel()

	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("failed to read settings: %s", err.Error())
	}

	mongoCollection, err := repository.NewMongoDB(ctx, cfg.Mongo)
	if err != nil {
		log.Fatalf("Could not create MongoDB logger: %v", err)
	}
	mongoCore, err := utility.NewMongoDBCore(mongoCollection)
	if err != nil {
		log.Fatalf("Could not create MongoDB core: %v", err)
	}

	prdLogger := zap.New(mongoCore)
	defer prdLogger.Sync()
	logger := prdLogger.Sugar()
	logger.Info("This is a test log message")

	fmt.Println(logger.Level())

	postgres := repository.NewPostgresDB(cfg.Postgres, logger)
	redisClient := repository.NewRedisDB(ctx, cfg.Redis, logger)
	services := service.NewService(ctx, postgres, redisClient)

	userAuthMiddleware := middlewares.NewUserAuthMiddleware(services.AuthService)

	srv := api.NewServer(userAuthMiddleware)
	srv.HandleAuth(services.AuthService)
	srv.HandleLists(ctx, services.ListService)
	srv.HandleItems(ctx, services.ItemService)

	go func() {
		if err := srv.Run(); err != nil {
			logger.Fatalf("error occured while running http server: %s", err.Error())
		}
		logger.Info("server started")
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)
	<-sigChan

	if err = srv.Shutdown(ctx); err != nil {
		logger.Errorf(err.Error())
		return
	}
}
