package redisdb

import (
	"context"
	"errors"
	"os"

	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog/log"
)

// Database struct
type Database struct {
	Client *redis.Client
}

var (
	// ErrNil error message
	ErrNil = errors.New("no matching record found in redis database")
	// Ctx empty context
	Ctx = context.TODO()
)

// NewDatabase returns redisClient that can be shared across the app
func NewDatabase() *Database {
	log.Info().Msg("Initializing Redis client...")
	client := redis.NewClient(&redis.Options{
		Addr:     getEnv("REDIS_URL", "localhost:6379"),
		Password: getEnv("REDIS_PASSWORD", ""),
		DB:       0,
	})

	if err := client.Ping(Ctx).Err(); err != nil {
		panic("Could not ping Redis")
	}
	return &Database{
		Client: client,
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
