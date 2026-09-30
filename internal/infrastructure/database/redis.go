package database

import (
	"context"
	"fmt"
	"log"

	"github.com/irfantrue/web-application-backend-go/internal/config"
	"github.com/redis/go-redis/v9"
)

// NewRedisClient creates and initializes a new Redis client using the provided configuration.
// It also pings the Redis server to ensure a successful connection before returning the client.
func NewRedisClient(cfg *config.Config) (*redis.Client, error) {
	addr := fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port)

	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: cfg.Redis.Password,
		DB:       0,
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %v", err)
	}

	log.Println("Successfully connected to the Redis database.")
	return client, nil
}
