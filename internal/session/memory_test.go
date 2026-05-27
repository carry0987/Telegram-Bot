package session

import (
	"context"
	"testing"
	"time"
)

func TestMemoryStoreIncrementAndGet(t *testing.T) {
	store := NewMemoryStore()

	count, err := store.Increment(context.Background(), 42, "hits", time.Hour)
	if err != nil {
		t.Fatalf("increment: %v", err)
	}
	if count != 1 {
		t.Fatalf("expected count 1, got %d", count)
	}

	if err := store.Set(context.Background(), 42, "last_command", "/start", time.Hour); err != nil {
		t.Fatalf("set: %v", err)
	}

	value, err := store.Get(context.Background(), 42, "last_command")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if value != "/start" {
		t.Fatalf("expected /start, got %q", value)
	}
}

func TestMemoryStoreExpiresEntries(t *testing.T) {
	store := NewMemoryStore()

	if err := store.Set(context.Background(), 42, "last_command", "/start", time.Nanosecond); err != nil {
		t.Fatalf("set: %v", err)
	}

	time.Sleep(2 * time.Millisecond)

	_, err := store.Get(context.Background(), 42, "last_command")
	if err != ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}
