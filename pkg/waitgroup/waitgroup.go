package waitgroup

import (
	"sync"

	"github.com/go-goim/core/pkg/goroutine"
)

type WaitGroup struct {
	wg *sync.WaitGroup
	ch chan struct{}
}

func NewWaitGroup(size int) *WaitGroup {
	return &WaitGroup{
		wg: new(sync.WaitGroup),
		ch: make(chan struct{}, size),
	}
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
}

func (wg *WaitGroup) Done() {
	<-wg.ch
	wg.wg.Done()
}
