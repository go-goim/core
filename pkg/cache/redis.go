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

// NewRedisCache creates a new redisCache instance.
func NewRedisCache(cli *redisv8.Client) Cache { //nolint:deadcode,unused
	return &redisCache{
		client: cli,
	}
}

func (r *redisCache) Get(ctx context.Context, key string) ([]byte, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	return r.client.Get(ctx, key).Bytes()
}

func (r *redisCache) Set(ctx context.Context, key string, value []byte, expire time.Duration) error {
	if ctx == nil {
		ctx = context.Background()
	}

	return r.client.Set(ctx, key, value, expire).Err()
}

func (r *redisCache) Delete(ctx context.Context, key string) error {
	if ctx == nil {
		ctx = context.Background()
	}

	return r.client.Del(ctx, key).Err()
}

func (r *redisCache) Close(_ context.Context) error {
	return r.client.Close()
}
