package worker

import (
	"context"
	"log"
	"sync"
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

	wg := new(sync.WaitGroup)
	if got := p.SubmitOrEnqueue(context.Background(), tf, concurrence, wg); got != SubmitResultStarted {
		t.Errorf("SubmitOrEnqueue() = %v, want %v", got, SubmitResultStarted)
		return
	}
	wg.Wait()
	if cnt != concurrence {
		t.Errorf("cnt = %v, want %v", cnt, concurrence)
	}
	if got := p.SubmitOrEnqueue(context.Background(), tf, concurrence, nil); got != SubmitResultStarted {
		t.Errorf("SubmitOrEnqueue() = %v, want %v", got, SubmitResultStarted)
		return
	}
	if got := p.SubmitOrEnqueue(context.Background(), tf, concurrence, nil); got != SubmitResultEnqueue {
		t.Errorf("SubmitOrEnqueue() = %v, want %v", got, SubmitResultEnqueue)
		return
	}
	if got := p.SubmitOrEnqueue(context.Background(), tf, 6, nil); got != SubmitResultOutOfSize {
		t.Errorf("SubmitOrEnqueue() = %v, want %v", got, SubmitResultOutOfSize)
		return
	}
	if got := p.SubmitOrEnqueue(context.Background(), tf, concurrence, nil); got != SubmitResultBufferFull {
		t.Errorf("SubmitOrEnqueue() = %v, want %v", got, SubmitResultBufferFull)
		return
	}
	p.Stop()
	if got := p.SubmitOrEnqueue(context.Background(), tf, 1, nil); got != SubmitResultClosed {
		t.Errorf("SubmitOrEnqueue() = %v, want %v", got, SubmitResultClosed)
		return
	}
}
