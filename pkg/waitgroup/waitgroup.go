package waitgroup

import (
	"context"
	"sync"

	"go.uber.org/atomic"

	"github.com/go-goim/core/pkg/goroutine"
	"github.com/go-goim/core/pkg/graceful"
)

type WaitGroup struct {
	wg   *sync.WaitGroup
	ch   chan struct{}
	done *atomic.Bool
}

func NewWaitGroup(size int) *WaitGroup {
	wg := &WaitGroup{
		wg:   new(sync.WaitGroup),
		ch:   make(chan struct{}, size),
		done: atomic.NewBool(false),
	}

	graceful.Register(wg.waitWhenShutdown)
	return wg
}

func (wg *WaitGroup) Add(f func()) {
	wg.wg.Add(1)
	_ = goroutine.Submit(func() {
		defer wg.Done()

		wg.ch <- struct{}{}
		f()
	})
}

func (wg *WaitGroup) Wait() {
	wg.wg.Wait()
	wg.done.Store(true)
}

func (wg *WaitGroup) waitWhenShutdown(ctx context.Context) error {
	if wg.done.Load() {
		return nil
	}

	done := make(chan struct{})
	go func() {
		wg.wg.Wait()
		done <- struct{}{}
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-done:
		return nil
	}
}

func (wg *WaitGroup) Done() {
	<-wg.ch
	wg.wg.Done()
}
