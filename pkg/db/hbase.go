package db

import (
	"context"

	"github.com/go-goim/core/pkg/db/hbase"
)

type hbaseCtxKey struct{}

// GetHBaseFromCtx try to get gohbase.Client from context, if not found then return defaultHBaseClient
func GetHBaseFromCtx(ctx context.Context) hbase.Client {
	if ctx == nil {
		return hbase.GetClient().Context(context.Background())
	}

	v := ctx.Value(hbaseCtxKey{})
	if v == nil {
		return hbase.GetClient().Context(ctx)
	}

	// double check
	cli, ok := v.(hbase.Client)
	if !ok {
		// maybe set by others
		return hbase.NewClient().Context(ctx)
	}

	return cli
}

// CtxWithHBase return new context.Context contain value with gohbase.Client
func CtxWithHBase(ctx context.Context, hbase hbase.Client) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}

	return context.WithValue(ctx, hbaseCtxKey{}, hbase.Context(ctx))
}
