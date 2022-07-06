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
	ErrKeyType   = errors.New("invalid key type")
)

// SetGlobalCache sets the global cache.
func SetGlobalCache(c Cache) {
	globalCache = c
}

// GetGlobalCache returns the global cache.
func GetGlobalCache() Cache {
	return globalCache
}

/*
 * string
 */

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

/*
 * set
 */

// IsInSet is wrapper for global cache.IsInSet.
func IsInSet(ctx context.Context, key string, member string) (bool, error) {
	return globalCache.IsInSet(ctx, key, member)
}

// AddToSet is wrapper for global cache.AddToSet.
func AddToSet(ctx context.Context, key string, member string) error {
	return globalCache.AddToSet(ctx, key, member)
}

// DeleteFromSet is wrapper for global cache.DeleteFromSet.
func DeleteFromSet(ctx context.Context, key string, member string) error {
	return globalCache.DeleteFromSet(ctx, key, member)
}

/*
 * hashmap
 */

// GetFromHash is wrapper for global cache.GetFromHash.
func GetFromHash(ctx context.Context, key string, field string) ([]byte, error) {
	return globalCache.GetFromHash(ctx, key, field)
}

// SetToHash is wrapper for global cache.SetToHash.
func SetToHash(ctx context.Context, key string, field string, value []byte) error {
	return globalCache.SetToHash(ctx, key, field, value)
}

// DeleteFromHash is wrapper for global cache.DeleteFromHash.
func DeleteFromHash(ctx context.Context, key string, field string) error {
	return globalCache.DeleteFromHash(ctx, key, field)
}

// Close is wrapper for global cache.Close.
func Close(ctx context.Context) error {
	return globalCache.Close(ctx)
}
