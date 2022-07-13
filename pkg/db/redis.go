package db

import (
	"context"

	redisv8 "github.com/go-redis/redis/v8"

	"github.com/go-goim/core/pkg/db/redis"
)

type redisCtxKey struct{}

// GetRedisFromCtx try to get redisv8.Client from context, if not found then return defaultRedisClient
func GetRedisFromCtx(ctx context.Context) *redisv8.Client {
	if ctx == nil {
		return redis.GetRedis().WithContext(context.Background())
	}

	v := ctx.Value(redisCtxKey{})
	if v == nil {
		return redis.GetRedis().WithContext(ctx)
	}

	// double check
	cli, ok := v.(*redisv8.Client)
	if !ok {
		// maybe set by others
		return redis.GetRedis().WithContext(ctx)
	}

	return cli
}

// CtxWithRedis return new context.Context contain value with redisv8.Client
func CtxWithRedis(ctx context.Context, redis *redisv8.Client) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}

	return context.WithValue(ctx, redisCtxKey{}, redis)
}

// Pipeline get redisv8.Client from ctx and run Pipeline Operation with ctx
func Pipeline(ctx context.Context, fn func(pipeline redisv8.Pipeliner) error) error {
	if ctx == nil {
		ctx = context.Background()
	}

	cli := GetRedisFromCtx(ctx)
	_, err := cli.Pipelined(ctx, fn)
	return err
}

// TxPipeline get redisv8.Client from ctx and run Transaction Pipeline Operation with ctx
//  It's same as Pipeline, but wraps queued commands with MULTI/EXEC.
// How to use:
// var (
//		get *redisv8.StringCmd
//	)
//	err := TxPipeline(ctx, func(pipeline redisv8.Pipeliner) error {
//		pipeline.Set(ctx, "key", "value", 0)
//		get = pipeline.Get(ctx, "key")
//		return nil
//	})
//	if err != nil {
//		return err
//	}
//
//	fmt.Println(get.Val())
//	// Output: value
func TxPipeline(ctx context.Context, fn func(pipeline redisv8.Pipeliner) error) error {
	if ctx == nil {
		ctx = context.Background()
	}

	cli := GetRedisFromCtx(ctx)
	_, err := cli.TxPipelined(ctx, fn)
	return err
}
