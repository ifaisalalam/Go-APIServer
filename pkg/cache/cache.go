package cache

import (
	"context"
	"time"
)

// Cache interface.
type Cache interface {
	// Get value of a key from cache.
	Get(ctx context.Context, key string) ([]byte, error)
	// Set a key in cache with TTL.
	Set(ctx context.Context, key string, value []byte, ttl time.Duration) error
	// Delete a key in cache.
	Delete(ctx context.Context, key string) error
	// IsAlive performs a healthcheck on the cache.
	IsAlive(context.Context) (bool, error)
}
