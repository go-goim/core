package worker

import (
	"context"
	"sync"

	"go.uber.org/atomic"
)

// workerSet represent a group of task handle workers
type workerSet struct {
	task          *task
	runningWorker atomic.Int32
	workers       []*worker
	ctx           context.Context
	cancel        context.CancelFunc
	wg            *sync.WaitGroup
}

func newWorkerSet(ctx context.Context, t *task) *workerSet {
	if t.wg == nil {
		t.wg = new(sync.WaitGroup)
	}

	if ctx == nil {
		ctx = context.Background()
	}

	ws := &workerSet{
		task:    t,
		workers: make([]*worker, t.concurrence),
		wg:      t.wg,
	}
	t.assignWorkerSet(ws)
	ws.ctx, ws.cancel = context.WithCancel(ctx)

	for i := 0; i < t.concurrence; i++ {
		ws.workers[i] = newWorker(ws)
	}

	return ws
}

func (ws *workerSet) run() {
	ws.task.updateStatus(TaskStatusRunning)
	for _, w := range ws.workers {
		ws.addOne()
		go w.run()
	}
}

func (ws *workerSet) stopAll() {
	ws.cancel()
	ws.task.updateStatus(TaskStatusKilled)
}

// err returns the first error that occurred in the workerSet.
func (ws *workerSet) err() error {
	for _, w := range ws.workers {
		if w.err != nil {
			return w.err
		}
	}

	return nil
}

func (ws *workerSet) getRunningWorker() int {
	return int(ws.runningWorker.Load())
}

// done called when worker stop.
func (ws *workerSet) done() {
	ws.addRunningWorker(-1)
	ws.wg.Done()
	if ws.getRunningWorker() == 0 {
		ws.task.updateStatus(TaskStatusDone)
	}
}

// addOne called when worker start running.
func (ws *workerSet) addOne() {
	ws.addRunningWorker(1)
	ws.wg.Add(1)
}

func (ws *workerSet) wait() {
	ws.wg.Wait()
}

func (ws *workerSet) addRunningWorker(delta int) {
	if delta > 0 {
		ws.runningWorker.Add(int32(delta))
	}

	if delta < 0 {
		ws.runningWorker.Sub(-int32(delta))
	}
}
