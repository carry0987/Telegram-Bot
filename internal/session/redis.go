package session

import (
	"context"
	"fmt"
	"time"

	redis "github.com/redis/go-redis/v9"
)

type RedisStore struct {
	client *redis.Client
}

func NewRedisStore(redisURL string) (*RedisStore, error) {
	options, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, fmt.Errorf("parse redis url: %w", err)
	}

	return &RedisStore{client: redis.NewClient(options)}, nil
}

func (s *RedisStore) Increment(ctx context.Context, chatID int64, key string, ttl time.Duration) (int64, error) {
	count, err := s.client.Incr(ctx, redisSessionKey(chatID, key)).Result()
	if err != nil {
		return 0, err
	}

	if ttl > 0 {
		if err := s.client.Expire(ctx, redisSessionKey(chatID, key), ttl).Err(); err != nil {
			return 0, err
		}
	}

	return count, nil
}

func (s *RedisStore) Set(ctx context.Context, chatID int64, key, value string, ttl time.Duration) error {
	return s.client.Set(ctx, redisSessionKey(chatID, key), value, ttl).Err()
}

func (s *RedisStore) Get(ctx context.Context, chatID int64, key string) (string, error) {
	value, err := s.client.Get(ctx, redisSessionKey(chatID, key)).Result()
	if err == redis.Nil {
		return "", ErrNotFound
	}

	return value, err
}

func (s *RedisStore) Backend() string {
	return "redis"
}

func (s *RedisStore) Close() error {
	return s.client.Close()
}

func redisSessionKey(chatID int64, key string) string {
	return fmt.Sprintf("telegram-bot:session:%d:%s", chatID, key)
}
