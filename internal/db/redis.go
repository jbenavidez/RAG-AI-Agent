package db

import (
	"context"
	"fmt"
	"os"

	redis "github.com/redis/go-redis/v9"
)

func NewRedisClient() (*redis.Client, error) {
	ctx := context.Background()
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		redisURL = "redis://localhost:6379"
	}

	options, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, fmt.Errorf("invalid REDIS_URL: %w", err)
	}

	rdb := redis.NewClient(options)
	// test connn
	ping, err := rdb.Ping(ctx).Result()
	if err != nil {
		_ = rdb.Close()
		return nil, fmt.Errorf("could not connect to redis: %w", err)
	}
	fmt.Printf("Redis connected successfully: %s\n", ping)
	return rdb, nil
}
