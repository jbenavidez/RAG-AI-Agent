package services

import "github.com/redis/go-redis/v9"

type MemoryService struct {
	redisClient *redis.Client
}

func NewMemoryService(rdb *redis.Client) *MemoryService {

	return &MemoryService{
		redisClient: rdb,
	}
}
