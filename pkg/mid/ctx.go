package mid

import (
	"context"

	"github.com/gin-gonic/gin"
)

const (
	ctxKey = "ctx"
)

func SetContext(c *gin.Context, ctx context.Context) {
	c.Set(ctxKey, ctx)
}

func GetContext(c *gin.Context) context.Context {
	v, ok := c.Get(ctxKey)
	if !ok {
		return context.Background()
	}

	ctx, ok := v.(context.Context)
	if !ok {
		return context.Background()
	}

	return ctx
}
