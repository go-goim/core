package worker

import (
	"container/list"
	"context"
	"sync"
	"time"

	"go.uber.org/atomic"

	"github.com/go-goim/core/pkg/util"
)

// Pool is a buffered worker pool
type Pool struct {
	taskList   *list.List // double linked list
	poolSize   int        // size of max task in list
	maxWorker  int        // count of how many worker run in concurrence
	workerSets *list.List // list of worker sets
	lock       *sync.Mutex
	stopChan   chan struct{}
	stopFlag   atomic.Bool
}

const (
	defaultWorkerSize = 100
	defaultPoolSize   = 20 // assume one task need run 5 worker concurrence
)

func NewPool(maxWorker, poolSize int) *Pool {
	p := &Pool{
		poolSize:   defaultPoolSize,
		maxWorker:  defaultWorkerSize,
		lock:       new(sync.Mutex),
		taskList:   list.New(),
		workerSets: list.New(),
		stopChan:   make(chan struct{}, 1),
		stopFlag:   atomic.Bool{},
	}

	if maxWorker > 0 {
		p.maxWorker = maxWorker
	}

	if poolSize >= 0 {
		p.poolSize = poolSize
	}

	go p.consumeQueue()
	return p
}

func (p *Pool) Submit(ctx context.Context, tf TaskFunc, concurrence int) TaskResult {
	if p.stopFlag.Load() {
		return TaskStatusPoolClosed
	}

	if concurrence > p.maxWorker {
		return TaskStatusTooManyWorker
	}

	// check if there has any worker place left
	p.lock.Lock()
	defer p.lock.Unlock()

	t := newTask(tf, concurrence)
	if p.tryRunTask(ctx, t) {
		return t
	}

	if p.enqueueTask(t) {
		return t
	}

	return TaskStatusQueueFull
}

func (p *Pool) Shutdown(ctx context.Context) error {
	p.stopChan <- struct{}{}
	p.stopFlag.Store(true)
	// stop all workers
	for e := p.workerSets.Front(); e != nil; e = e.Next() {
		ws := e.Value.(*workerSet)
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			ws.stopAll()
			ws.wait()
		}
	}

	return nil
}

// tryRunTask try to put task into workerSet and run it.Return false if capacity not enough.
// Make sure get p.Lock before call this func
func (p *Pool) tryRunTask(ctx context.Context, t *task) bool {
	// calculate how many worker can run in concurrence
	availableWorkerCount := p.maxWorker - p.curRunningWorkerNum()
	if availableWorkerCount <= 0 {
		return false
	}

	ws := newWorkerSet(ctx, t)
	ws.run(util.Min(availableWorkerCount, t.concurrence))
	p.workerSets.PushBack(ws)
	return true
}

// curRunningWorkerNum make sure lock mutex before call this func
func (p *Pool) curRunningWorkerNum() int {
	var (
		cnt int
	)
	// range worker set list
	for e := p.workerSets.Front(); e != nil; e = e.Next() {
		ws := e.Value.(*workerSet)
		cnt += ws.curRunningWorkerNum()
	}
	return cnt
}

func (p *Pool) enqueueTask(t *task) bool {
	// double check to avoid got panic: write to closed channel
	if p.stopFlag.Load() {
		return false
	}

	// check if queue is full
	if p.taskList.Len() >= p.poolSize {
		return false
	}

	p.taskList.PushBack(t)
	return true
}

func (p *Pool) checkWorkerNum() {
	// make sure lock mutex before call this func
	var (
		lastEmptyWorkerSet *list.Element
	)

	for e := p.workerSets.Front(); e != nil; e = e.Next() {
		if lastEmptyWorkerSet != nil {
			p.workerSets.Remove(lastEmptyWorkerSet)
			lastEmptyWorkerSet = nil
		}

		ws := e.Value.(*workerSet)
		if ws.isDone() {
			lastEmptyWorkerSet = e
			continue
		}

		if cnt := ws.needMoreWorker(); cnt > 0 {
			ws.run(util.Min(p.maxWorker-p.curRunningWorkerNum(), cnt))
		}
	}

	if lastEmptyWorkerSet != nil {
		p.workerSets.Remove(lastEmptyWorkerSet)
	}
}

func (p *Pool) consumeQueue() {
	var ticker = time.NewTicker(time.Millisecond * 20)
	for {
		if p.stopFlag.Load() {
			return
		}

		select {
		// check stop chan
		case <-p.stopChan:
			ticker.Stop()
			return
		default:
		}

		// check worker set first
		p.lock.Lock()
		p.checkWorkerNum()

		// try to run task from task list
		e := p.taskList.Front()
		if e != nil {
			if p.tryRunTask(context.Background(), e.Value.(*task)) {
				p.taskList.Remove(e)
			}
		}
		// release lock first
		p.lock.Unlock()
		// check again
		if e == nil {
			// no task to run, wait for next tick
			<-ticker.C
		}
	}
}
