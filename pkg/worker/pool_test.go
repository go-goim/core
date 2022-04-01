package worker

import (
	"context"
	"log"
	"testing"
	"time"
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
	if got.Status() != TaskStatusRunning {
		t.Errorf("SubmitOrEnqueue() = %v, want %v", got.Status(), TaskStatusRunning)
		return
	}
	got.Wait()
	if cnt != concurrence {
		t.Errorf("cnt = %v, want %v", cnt, concurrence)
	}
	if got := p.Submit(context.Background(), tf, concurrence); got.Status() != TaskStatusRunning {
		t.Errorf("SubmitOrEnqueue() = %v, want %v", got, TaskStatusRunning)
		return
	}
	if got := p.Submit(context.Background(), tf, concurrence); got.Status() != TaskStatusEnqueue {
		t.Errorf("SubmitOrEnqueue() = %v, want %v", got.Status(), TaskStatusEnqueue)
		return
	}
	if got := p.Submit(context.Background(), tf, 6); got.Status() != TaskStatusTooManyWorker {
		t.Errorf("SubmitOrEnqueue() = %v, want %v", got.Status(), TaskStatusTooManyWorker)
		return
	}
	if got := p.Submit(context.Background(), tf, concurrence); got.Status() != TaskStatusQueueFull {
		t.Errorf("SubmitOrEnqueue() = %v, want %v", got.Status(), TaskStatusQueueFull)
		return
	}
	p.Stop()
	if got := p.Submit(context.Background(), tf, 1); got.Status() != TaskStatusPoolClosed {
		t.Errorf("SubmitOrEnqueue() = %v, want %v", got.Status(), TaskStatusPoolClosed)
		return
	}
}
