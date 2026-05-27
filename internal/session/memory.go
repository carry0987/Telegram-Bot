package session

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"
)

type memoryEntry struct {
	value     string
	expiresAt time.Time
}

type MemoryStore struct {
	mu     sync.Mutex
	values map[string]memoryEntry
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{values: map[string]memoryEntry{}}
}

func (s *MemoryStore) Increment(_ context.Context, chatID int64, key string, ttl time.Duration) (int64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	cacheKey := sessionKey(chatID, key)
	entry, ok := s.values[cacheKey]
	if ok && entry.isExpired() {
		delete(s.values, cacheKey)
		ok = false
	}

	var current int64
	if ok {
		parsed, err := strconv.ParseInt(entry.value, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("parse counter: %w", err)
		}
		current = parsed
	}

	current++
	s.values[cacheKey] = memoryEntry{value: strconv.FormatInt(current, 10), expiresAt: expiresAt(ttl)}

	return current, nil
}

func (s *MemoryStore) Set(_ context.Context, chatID int64, key, value string, ttl time.Duration) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.values[sessionKey(chatID, key)] = memoryEntry{value: value, expiresAt: expiresAt(ttl)}

	return nil
}

func (s *MemoryStore) Get(_ context.Context, chatID int64, key string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	cacheKey := sessionKey(chatID, key)
	entry, ok := s.values[cacheKey]
	if !ok {
		return "", ErrNotFound
	}
	if entry.isExpired() {
		delete(s.values, cacheKey)
		return "", ErrNotFound
	}

	return entry.value, nil
}

func (s *MemoryStore) Backend() string {
	return "memory"
}

func (s *MemoryStore) Close() error {
	return nil
}

func (e memoryEntry) isExpired() bool {
	if e.expiresAt.IsZero() {
		return false
	}

	return time.Now().After(e.expiresAt)
}

func expiresAt(ttl time.Duration) time.Time {
	if ttl <= 0 {
		return time.Time{}
	}

	return time.Now().Add(ttl)
}

func sessionKey(chatID int64, key string) string {
	return fmt.Sprintf("%d:%s", chatID, key)
}
