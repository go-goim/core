package cache

import (
	"context"
	"time"

	redisv8 "github.com/go-redis/redis/v8"
)

// redisCache is a wrapper around the redis client that implements the
// Cache interface.
type redisCache struct {
	client *redisv8.Client
}

var _ Cache = &redisCache{}

// newRedisCache creates a new redisCache instance.
func newRedisCache(cli *redisv8.Client) *redisCache { //nolint:deadcode,unused
	return &redisCache{
		client: cli,
	}
}

func (r *redisCache) Get(key string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	return r.client.Get(ctx, key).Bytes()
}

func (r *redisCache) Set(key string, value []byte, expire time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	return r.client.Set(ctx, key, value, expire).Err()
}

func (r *redisCache) Delete(key string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	return r.client.Del(ctx, key).Err()
}

func (r *redisCache) Close() error {
	return r.client.Close()
}
