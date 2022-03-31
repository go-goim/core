package worker

import (
	"context"
	"sync"
	"time"

	"go.uber.org/atomic"
)

// Pool is a buffered worker pool
type Pool struct {
	taskQueue         chan *task
	enqueuedTaskCount atomic.Int32 // count of unhandled tasks
	bufferSize        int          // size of taskQueue buffer, means can count of bufferSize task can wait to be handled
	maxWorker         int          // count of how many worker run in concurrence
	workerSets        []*workerSet
	lock              *sync.Mutex
	stop              chan struct{}
	stopFlag          atomic.Bool
}

type TaskFunc func() error

type task struct {
	wg          *sync.WaitGroup // store waitGroup passed from task submiter.
	tf          TaskFunc
	concurrence int
}

const (
	defaultWorkerSize = 100
	defaultQueueSize  = 20 // assume one task need run 5 worker concurrence
)

func NewPool(workerSize, queueSize int) *Pool {
	p := &Pool{
		enqueuedTaskCount: atomic.Int32{},
		bufferSize:        defaultQueueSize,
		maxWorker:         defaultWorkerSize,
		lock:              new(sync.Mutex),
		workerSets:        make([]*workerSet, 0),
		stop:              make(chan struct{}, 1),
		stopFlag:          atomic.Bool{},
	}

	if workerSize > 0 {
		p.maxWorker = workerSize
	}

	if queueSize >= 0 {
		p.bufferSize = queueSize
	}

	// check p.enqueue to find out why make this channel size with p.bufferSize+1.
	p.taskQueue = make(chan *task, p.bufferSize+1)
	go p.consumeQueue()
	return p
}

type SubmitResult int

const (
	SubmitResultStarted SubmitResult = iota
	SubmitResultEnqueue
	SubmitResultBufferFull
	SubmitResultOutOfSize
	SubmitResultClosed
)

func (p *Pool) SubmitOrEnqueue(ctx context.Context, tf TaskFunc, concurrence int, wg *sync.WaitGroup) SubmitResult {
	if p.stopFlag.Load() {
		return SubmitResultClosed
	}

	if concurrence > p.maxWorker {
		return SubmitResultOutOfSize
	}

	// check if there has any worker place left
	p.lock.Lock()
	defer p.lock.Unlock()
	t := &task{
		tf:          tf,
		concurrence: concurrence,
		wg:          wg,
	}

	if p.tryRunTask(ctx, t) {
		return SubmitResultStarted
	}

	if p.enqueueTask(t, true) {
		return SubmitResultEnqueue
	}

	return SubmitResultBufferFull
}

func (p *Pool) Stop() {
	p.stopFlag.Store(true)
	// stop queue daemon
	p.stop <- struct{}{}
	close(p.taskQueue)
	// stop all workers
	for _, ws := range p.workerSets {
		ws.stopAll()
		ws.wait()
	}
}

// tryRunTask try to put task into workerSet and run it.Return false if capacity not enough.
// Make sure get p.Lock before call this func
func (p *Pool) tryRunTask(ctx context.Context, t *task) bool {
	if p.curRunningWorkerNum()+t.concurrence <= p.maxWorker {
		ws := newWorkerSet(ctx, t)
		p.workerSets = append(p.workerSets, ws)
		ws.run()
		return true
	}

	return false
}

// curRunningWorkerNum make sure lock mutex before call this func
func (p *Pool) curRunningWorkerNum() int {
	var (
		cnt         int
		needRemoved = make(map[int]bool)
	)
	for i, w := range p.workerSets {
		rw := w.getRunningWorker()
		if rw == 0 {
			needRemoved[i] = true
		}

		cnt += rw
	}

	if len(needRemoved) == 0 {
		return cnt
	}

	// remove finished workerSet
	temp := make([]*workerSet, 0, len(p.workerSets)-len(needRemoved))
	for i := range p.workerSets {
		if needRemoved[i] {
			continue
		}

		temp = append(temp, p.workerSets[i])
	}

	p.workerSets = temp
	return cnt
}

func (p *Pool) enqueueTask(t *task, isNewTask bool) bool {
	// double check to avoid got panic: write to closed channel
	if p.stopFlag.Load() {
		return false
	}

	// Use atomic value instead of len(p.taskQueue).
	// Because taskQueue need be read by p.consumeQueue and try to run the task,
	// when try to run task fail and before put it back to taskQueue, there are len(p.taskQueue) + 1 tasks
	// need to be handled.So it may cause unpredictable problem if we use len(p.taskQueue) as total count of
	// enqueued tasks.
	if int(p.enqueuedTaskCount.Load()) >= p.bufferSize && isNewTask {
		return false
	}

	// if this is a put back old task action,then won't check channel capacity,
	// because channel length is p.bufferSize + 1, so put back task will not be blocked.
	p.taskQueue <- t
	if isNewTask {
		p.enqueuedTaskCount.Add(1)
	}
	return true
}

func (p *Pool) consumeQueue() {
	for {
		select {
		case t := <-p.taskQueue:
			p.lock.Lock()
			if p.tryRunTask(context.Background(), t) {
				p.enqueuedTaskCount.Sub(1)
				goto unlock
			}

			// put it back
			if !p.enqueueTask(t, false) {
				// channel is closed
				continue
			}
			// sleep little while if try to run task failed
			time.Sleep(time.Millisecond * 20)

		unlock:
			p.lock.Unlock()
		case <-p.stop:
			// taskQueue closed
			return
		}
	}
}

// workerSet represent a group of task handle workers
type workerSet struct {
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
		workers: make([]*worker, t.concurrence),
		wg:      t.wg,
	}
	ws.ctx, ws.cancel = context.WithCancel(ctx)

	for i := 0; i < t.concurrence; i++ {
		ws.workers[i] = newWorker(ws, t.tf)
	}

	return ws
}

func (ws *workerSet) run() {
	for _, w := range ws.workers {
		ws.addOne()
		go w.run()
	}
}

func (ws *workerSet) stopAll() {
	ws.cancel()
}

func (ws *workerSet) getRunningWorker() int {
	return int(ws.runningWorker.Load())
}

// done called when worker stop.
func (ws *workerSet) done() {
	ws.addRunningWorker(-1)
	ws.wg.Done()
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
