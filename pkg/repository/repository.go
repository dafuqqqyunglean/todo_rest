package repository

import (
	"context"
	"fmt"
	"github.com/dafuqqqyunglean/todoRestAPI/config"
	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"log"
)

func NewPostgresDB(cfg config.PostgresConfig, logger *zap.SugaredLogger) *sqlx.DB {
	db, err := sqlx.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.DBName, cfg.Password, cfg.SSLMode))
	if err != nil {
		logger.Fatalf("failed to connect to postgres: %v", err)
	}

	err = db.Ping()
	if err != nil {
		logger.Fatalf("failed to ping database: %v", err)
	}

	return db
}

func NewRedisDB(ctx context.Context, cfg config.RedisConfig, logger *zap.SugaredLogger) *redis.Client {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.Address,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		logger.Fatalf("failed to ping Redis: %v", err)
	}

	return redisClient
}

func NewMongoDB(ctx context.Context, cfg config.MongoConfig) (*mongo.Collection, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoURL))
	if err != nil {
		log.Fatalf("failed to connect to MongoDB: %v", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("failed to ping MongoDB: %v", err)
	}

	mongoDB := client.Database(cfg.DB).Collection("logs")

	return mongoDB, nil
}
