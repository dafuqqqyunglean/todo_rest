package config

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
)

type Config struct {
	Postgres PostgresConfig
	Redis    RedisConfig
	Mongo    MongoConfig
}

type PostgresConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

type RedisConfig struct {
	Address  string
	Password string
	DB       int
}

type MongoConfig struct {
	DB       string
	MongoURL string
}

func NewConfig() (Config, error) {
	err := initConfig()
	if err != nil {
		return Config{}, err
	}

	return Config{
		Postgres: PostgresConfig{
			Host:     viper.GetString("postgres.host"),
			Port:     viper.GetString("postgres.port"),
			Username: viper.GetString("postgres.username"),
			DBName:   viper.GetString("postgres.dbname"),
			SSLMode:  viper.GetString("postgres.sslmode"),
			Password: os.Getenv("POSTGRES_PASSWORD"),
		},
		Redis: RedisConfig{
			Address:  fmt.Sprintf("%s:%s", viper.GetString("redis.host"), viper.GetString("redis.port")),
			Password: viper.GetString("redis.password"),
			DB:       viper.GetInt("redis.db"),
		},
		Mongo: MongoConfig{
			DB:       viper.GetString("mongo.dbname"),
			MongoURL: fmt.Sprintf("mongodb://%s:%s", viper.GetString("mongo.host"), viper.GetString("mongo.port")),
		},
	}, nil
}

func initConfig() error {
	viper.AddConfigPath("config")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
