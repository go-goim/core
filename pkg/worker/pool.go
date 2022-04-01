package worker

import (
	"context"
	"log"
	"sync"
	"time"

	"go.uber.org/atomic"
)

// Pool is a buffered worker pool
type Pool struct {
	// TODO: taskQueue should be a linked list, so that we can get the task from the head of the list and put it back to the head.
	// If we use a channel as taskQueue, we can't get the task from the head of the list and put it back to the head.
	// But make sure that before change it to linked list, we should have the ability run the task in min(taskQueue length, concurrence) goroutines.
	taskQueue         chan *task
	enqueuedTaskCount atomic.Int32 // count of unhandled tasks
	bufferSize        int          // size of taskQueue buffer, means can count of bufferSize task can wait to be handled
	maxWorker         int          // count of how many worker run in concurrence
	workerSets        []*workerSet
	lock              *sync.Mutex
	stopFlag          atomic.Bool
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

	if p.enqueueTask(t, true) {
		return t
	}

	return TaskStatusQueueFull
}

func (p *Pool) Stop() {
	p.stopFlag.Store(true)
	// stop queue daemon
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
	var ticker = time.NewTicker(time.Second)
	for {
		select {
		case t, ok := <-p.taskQueue:
			if !ok {
				// channel closed
				return
			}

			p.lock.Lock()
			if p.tryRunTask(context.Background(), t) {
				p.enqueuedTaskCount.Sub(1)
				goto unlock
			}

			// if enqueueTask return false, means channel is closed.
			if !p.enqueueTask(t, false) {
				// channel is closed
				goto unlock
			}
			// sleep little while if try to run task failed
			time.Sleep(time.Millisecond * 20)

		unlock:
			p.lock.Unlock()
		case <-ticker.C:
			// check if there has any worker place left
			// TODO: check if there is any workerSet is idle and remove it
			// TODO: try to run enqueued tasks even if there is no enough worker to run.
			log.Printf("current running worker num: %d", p.curRunningWorkerNum())
		}
	}

	// never reach here
}
