package cache

import (
	"context"
	"errors"
	"time"
)

var (
	globalCache Cache = NewMemoryCache()
)

var (
	ErrCacheMiss = errors.New("cache miss")
	ErrCacheFull = errors.New("cache full")
)

// SetGlobalCache sets the global cache.
func SetGlobalCache(c Cache) {
	globalCache = c
}

// GetGlobalCache returns the global cache.
func GetGlobalCache() Cache {
	return globalCache
}

// Get is wrapper for global cache.Get.
func Get(ctx context.Context, key string) ([]byte, error) {
	return globalCache.Get(ctx, key)
}

// Set is wrapper for global cache.Set.
func Set(ctx context.Context, key string, value []byte, expire time.Duration) error {
	return globalCache.Set(ctx, key, value, expire)
}

// Delete is wrapper for global cache.Delete.
func Delete(ctx context.Context, key string) error {
	return globalCache.Delete(ctx, key)
}

// Close is wrapper for global cache.Close.
func Close(ctx context.Context) error {
	return globalCache.Close(ctx)
}
