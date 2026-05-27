package session

import (
	"context"
	"errors"
	"strings"
	"time"
)

var ErrNotFound = errors.New("session value not found")

type Store interface {
	Increment(ctx context.Context, chatID int64, key string, ttl time.Duration) (int64, error)
	Set(ctx context.Context, chatID int64, key, value string, ttl time.Duration) error
	Get(ctx context.Context, chatID int64, key string) (string, error)
	Backend() string
	Close() error
}

func NewStore(redisURL string) (Store, error) {
	if strings.TrimSpace(redisURL) == "" {
		return NewMemoryStore(), nil
	}

	return NewRedisStore(redisURL)
}
