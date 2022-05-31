package retry

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-goim/goim/pkg/log"
	"github.com/go-goim/goim/pkg/mq"
)

func Retry(f func() error, opts ...Option) error {
	o := newRetryOptions(opts...)
	if o.async {
		return retryAsync(f, o)
	}
	return retry(f, o)
}

func RetryWithQueue(f func() error, producer mq.Producer, topic string, data interface{}, opts ...Option) error { // nolint: golint
	if opts == nil {
		opts = make([]Option, 0)
	}

	b, err := json.Marshal(data)
	if err != nil {
		return err
	}

	opts = append(opts, WithPutQueueIfFail(func() error {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		msg := mq.NewMessage(topic, b)
		result, err := producer.SendSync(ctx, msg)
		if err != nil {
			return err
		}
		if result.Status != 0 {
			return fmt.Errorf("send message to queue failed, status: %d", result.Status)
		}

		return nil
	}))

	return Retry(f, opts...)
}

func retryAsync(f func() error, o *retryOptions) error {
	go func() {
		if err := retry(f, o); err != nil {
			log.Error("retry failed", "err", err)
		}
	}()

	return retry(f, o)
}

func retry(f func() error, o *retryOptions) error {
	var err error
	for i := 0; i < o.maxRetries; i++ {
		err = f()
		if err == nil {
			return nil
		}

		if o.putQueueIfFail {
			return o.putQueueFunc()
		}

		time.Sleep(o.backoff)
		o.backoff *= 2
		if o.backoff > o.maxBackoff {
			o.backoff = o.maxBackoff
		}
	}

	return err
}
