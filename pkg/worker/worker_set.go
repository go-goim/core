package worker

import (
	"container/list"
	"context"
	"log"
	"sync"

	"go.uber.org/atomic"

	"github.com/go-goim/goim/pkg/errors"
)

// workerSet represent a group of task handle workers
type workerSet struct {
	task            *task
	finishWorkers   atomic.Int32
	runningWorker   atomic.Int32
	runnableWorkers atomic.Int32
	// list of workers w0 -> w1 -> w2 -> w3 -> w4
	// start from w0, w0 will be used firstly then move w0 to end of list
	workers *list.List
	// lock for workers
	lock   *sync.Mutex
	ctx    context.Context
	cancel context.CancelFunc
	wg     *sync.WaitGroup
}

func newWorkerSet(ctx context.Context, t *task) *workerSet {
	if ctx == nil {
		ctx = context.Background()
	}

	ws := &workerSet{
		task:    t,
		workers: list.New(),
		lock:    new(sync.Mutex),
		wg:      new(sync.WaitGroup),
	}
	t.assignWorkerSet(ws)
	ws.ctx, ws.cancel = context.WithCancel(ctx)

	for i := 0; i < t.concurrence; i++ {
		ws.workers.PushBack(newWorker(ws))
	}

	return ws
}

func (ws *workerSet) run(concurrence int) {
	if concurrence <= 0 {
		return
	}

	if ws.isDone() {
		return
	}

	ws.task.updateStatus(TaskStatusRunning)
	for i := 0; i < concurrence; i++ {
		if ws.runOne() {
			ws.runnableWorkers.Inc()
		}
	}
}

func (ws *workerSet) runOne() bool {
	ws.lock.Lock()
	defer ws.lock.Unlock()

	worker := ws.workers.Front().Value.(*worker)
	if !worker.isIdle() {
		// means all workers are running or finished
		return false
	}

	ws.addOne()
	ws.workers.MoveToBack(ws.workers.Front())

	worker.setRunning()
	go worker.run()
	return true
}

func (ws *workerSet) stopAll() {
	ws.cancel()
	ws.task.updateStatus(TaskStatusKilled)
}

// err returns the first error that occurred in the workerSet.
func (ws *workerSet) err() error {
	var err = make(errors.ErrorSet, 0)
	// range over all workers and return all errors.
	for e := ws.workers.Front(); e != nil; e = e.Next() {
		worker := e.Value.(*worker)
		if worker.err != nil {
			err = append(err, worker.err)
		}
	}

	return err.Err()
}

func (ws *workerSet) curRunningWorkerNum() int {
	return int(ws.runningWorker.Load())
}

func (ws *workerSet) needMoreWorker() int {
	// if runnable workers is less than concurrence, need more workers
	return ws.task.concurrence - int(ws.runnableWorkers.Load())
}

func (ws *workerSet) isDone() bool {
	return ws.finishWorkers.Load() == int32(ws.task.concurrence)
}

// done called when worker stop.
func (ws *workerSet) done() {
	ws.addRunningWorker(-1)
	ws.wg.Add(-1)

	if ws.isDone() {
		log.Println("all workers are done")
		ws.task.updateStatus(TaskStatusDone)
		return
	}

	ws.runOne()
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
		ws.finishWorkers.Add(1)
	}
}
