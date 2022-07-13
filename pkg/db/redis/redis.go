package redis

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"

	"github.com/go-goim/core/pkg/graceful"
)

var (
	defaultRedisClient *redis.Client
)

func GetRedis() *redis.Client {
	return defaultRedisClient
}

func InitRedis(opts ...Option) error {
	var err error
	defaultRedisClient, err = NewRedis(opts...)
	if err != nil {
		return err
	}

	graceful.Register(func(ctx context.Context) error {
		return Close()
	})

	return nil
}

func Close() error {
	if defaultRedisClient != nil {
		return defaultRedisClient.Close()
	}

	return nil
}

func NewRedis(opts ...Option) (*redis.Client, error) {
	o := &options{}
	for _, opt := range opts {
		opt(o)
	}

	cli := redis.NewClient(&redis.Options{
		Addr:         o.addr,
		Password:     o.password,
		DialTimeout:  o.dialTimeout,
		PoolSize:     o.maxConns,
		IdleTimeout:  o.idleTimeout,
		MinIdleConns: o.minIdleConns,
	})

	// add open tracing for rdb

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if err := cli.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return cli, nil
}
