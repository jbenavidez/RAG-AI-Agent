package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	keyExpireTime = 24 * time.Hour
)

type MemoryService struct {
	redisClient *redis.Client
}

type Chat struct {
	Question  string    `json:"question"`
	Answer    string    `json:"answer"`
	createdAt time.Time `json:"created_at"`
}

func NewMemoryService(rdb *redis.Client) *MemoryService {

	return &MemoryService{
		redisClient: rdb,
	}
}

func (m *MemoryService) Store(ctx context.Context, sessionID, question, answer string) error {
	// create redis key using current sessionID
	key := fmt.Sprintf("chat:%s:messages", sessionID)
	// create one chat turn
	turn := Chat{
		Question:  question,
		Answer:    answer,
		createdAt: time.Now(),
	}
	// convert chat into json
	data, err := json.Marshal(turn)
	if err != nil {
		return err
	}
	// Save the chat turn at the end of the Redis list.
	if err := m.redisClient.RPush(ctx, key, data).Err(); err != nil {
		return err
	}
	// keep chat history for 24 hours
	if err := m.redisClient.Expire(ctx, key, keyExpireTime).Err(); err != nil {
		return err
	}
	return nil

}

func (m *MemoryService) GetChatsHistory(ctx context.Context, sessionID string) ([]Chat, error) {
	key := fmt.Sprintf("chat:%s:messages", sessionID)
	values, err := m.redisClient.LRange(ctx, key, 0, -1).Result()
	if err != nil {
		return nil, err
	}
	chatHistory := make([]Chat, len(values))
	for i, v := range values {
		var chat Chat
		if err := json.Unmarshal([]byte(v), &chat); err != nil {
			return nil, err
		}
		chatHistory[i] = chat
	}

	return chatHistory, nil
}
