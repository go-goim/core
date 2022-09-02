package goroutine

import (
	"context"
	"time"

	"github.com/panjf2000/ants/v2"

	"github.com/go-goim/core/pkg/graceful"
)

var (
	_defaultPool *ants.Pool
)

func init() {
	var err error
	_defaultPool, err = ants.NewPool(10_000)
	if err != nil {
		panic(err)
	}

	graceful.Register(func(ctx context.Context) error {
		var timeout = time.Second * 5
		ddl, ok := ctx.Deadline()
		if ok {
			timeout = ddl.Sub(time.Now())
		}
		return _defaultPool.ReleaseTimeout(timeout)
	})
}

// Submit new task to pool
// It will block if goroutine is up to max
func Submit(f func()) error {
	return _defaultPool.Submit(f)
}
