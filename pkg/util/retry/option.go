package retry

import (
	"time"
)

type retryOptions struct {
	maxRetries int
	backoff    time.Duration
	maxBackoff time.Duration
	async      bool

	putQueueIfFail bool
	putQueueFunc   func() error
}

func newRetryOptions(opts ...Option) *retryOptions {
	r := &retryOptions{}
	for _, o := range opts {
		o(r)
	}

	return r
}

type Option func(*retryOptions)

func WithMaxRetries(maxRetries int) Option {
	return func(r *retryOptions) {
		r.maxRetries = maxRetries
	}
}

func WithBackoff(backoff time.Duration) Option {
	return func(r *retryOptions) {
		r.backoff = backoff
	}
}

func WithMaxBackoff(maxBackoff time.Duration) Option {
	return func(r *retryOptions) {
		r.maxBackoff = maxBackoff
	}
}
func WithAsync(async bool) Option {
	return func(r *retryOptions) {
		r.async = async
	}
}

func WithPutQueueIfFail(putQueueFunc func() error) Option {
	return func(r *retryOptions) {
		r.putQueueIfFail = true
		r.putQueueFunc = putQueueFunc
	}
}
