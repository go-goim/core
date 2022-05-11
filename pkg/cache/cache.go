package cache

import (
	"errors"
	"time"
)

var (
	globalCache Cache = newMemoryCache()
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
func Get(key string) ([]byte, error) {
	return globalCache.Get(key)
}

// Set is wrapper for global cache.Set.
func Set(key string, value []byte, expire time.Duration) error {
	return globalCache.Set(key, value, expire)
}

// Delete is wrapper for global cache.Delete.
func Delete(key string) error {
	return globalCache.Delete(key)
}

// Close is wrapper for global cache.Close.
func Close() error {
	return globalCache.Close()
}
