package worker

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPool_SubmitOrEnqueue(t *testing.T) {
	p := NewPool(5, 1)
	var (
		cnt         int
		concurrence = 5
	)

	tf := func() error {
		time.Sleep(time.Second)
		log.Println("hello world")
		cnt++
		return nil
	}

	got := p.Submit(context.Background(), tf, concurrence)
	assert.Equal(t, TaskStatusRunning, got.Status())
	got.Wait()
	assert.Equal(t, concurrence, cnt)

	got = p.Submit(context.Background(), tf, concurrence)
	assert.Equal(t, TaskStatusRunning, got.Status())

	got = p.Submit(context.Background(), tf, concurrence)
	assert.Equal(t, TaskStatusEnqueue, got.Status())

	got = p.Submit(context.Background(), tf, 6)
	assert.Equal(t, TaskStatusTooManyWorker, got.Status())

	got = p.Submit(context.Background(), tf, concurrence)
	assert.Equal(t, TaskStatusQueueFull, got.Status())

	_ = p.Shutdown(context.TODO())
	got = p.Submit(context.Background(), tf, 1)
	assert.Equal(t, TaskStatusPoolClosed, got.Status())
}
