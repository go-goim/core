package worker

import (
	"github.com/go-goim/core/pkg/goroutine"
)

// worker run tasks
type worker struct {
	ws    *workerSet
	err   error
	state int // 0: idle, 1: working, 2: done
}

func newWorker(ws *workerSet) *worker {
	return &worker{
		ws: ws,
	}
}

func (w *worker) isIdle() bool {
	return w.state == 0
}

func (w *worker) setRunning() {
	w.state = 1
}

func (w *worker) setDone() {
	w.state = 2
	w.ws.done()
}

func (w *worker) run() {
	var ec = make(chan error, 1)
	_ = goroutine.Submit(func() {
		ec <- w.ws.task.tf()
		close(ec)
	})

	select {
	case e := <-ec:
		w.err = e
	case <-w.ws.ctx.Done():
	}

	w.setDone()
}
