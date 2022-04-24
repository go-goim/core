package graceful

import (
	"context"
	"sync"
	"time"

	"github.com/yusank/goim/pkg/errors"
	"github.com/yusank/goim/pkg/log"
)

var (
	// DefaultTimeout is the default timeout for graceful shutdown.
	DefaultTimeout = 30 * time.Second

	// need a function set to store all the graceful shutdown functions
	// so that we can call them all at once when the server is shutdown.
	gracefulShutdownFuncs []gracefulShutdownFunc
)

// gracefulShutdownFunc is a function that is called when the server is shutdown.
// gracefulShutdownFunc need pass context to limit the time of function execution.
type gracefulShutdownFunc func(ctx context.Context) error

// Shutdown gracefully shuts down the server.
//
// This function will block until the shutdown is complete.
// It is recommended to pass a context with timeout to limit the time of server shutdown.
func Shutdown(ctx context.Context) error {
	log.Info("graceful shutdown...")
	var cancel context.CancelFunc
	if ctx == nil {
		ctx = context.Background()
		ctx, cancel = context.WithTimeout(ctx, DefaultTimeout)
		defer cancel()
	}

	var (
		done = make(chan struct{})
		errs = make(errors.ErrorSet, 0)
		wg   sync.WaitGroup
	)

	for _, f := range gracefulShutdownFuncs {
		wg.Add(1)
		go func(f gracefulShutdownFunc) {
			defer wg.Done()
			if err := f(ctx); err != nil {
				errs = append(errs, err)
			}
		}(f)
	}

	go func() {
		wg.Wait()
		done <- struct{}{}
	}()

	select {
	case <-ctx.Done(): // timeout
		return ctx.Err()
	case <-done: // shutdown complete
		return errs.Err()
	}

}

// Register registers a function to be called when the server is shutdown.
func Register(f gracefulShutdownFunc) {
	gracefulShutdownFuncs = append(gracefulShutdownFuncs, f)
}
