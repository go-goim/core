package worker

type TaskFunc func() error

type task struct {
	tf          TaskFunc   // task function
	concurrence int        // concurrence of task
	ws          *workerSet // assign value after task distribute to worker.
	status      TaskStatus // store task status.
}

// TaskResult is a manager of submitted task.
type TaskResult interface {
	// get error if task failed.
	Err() error
	// wait for task done.
	Wait()
	// get task status.
	Status() TaskStatus
	// kill task.
	Kill()
}

func newTask(tf TaskFunc, concurrence int) *task {
	return &task{
		tf:          tf,
		concurrence: concurrence,
	}
}

// Err returns the first error that occurred in the workerSet.
func (t *task) Err() error {
	// check t.ws if nil return nil.
	if t.ws == nil {
		return nil
	}

	return t.ws.err()
}

// Wait for task done.
// Please make sure task is done or running before call this function.
func (t *task) Wait() {
	// check t.ws if nil.
	if t.ws == nil {
		return
	}

	t.ws.wait()
}

// Status returns task status.
func (t *task) Status() TaskStatus {
	return t.status
}

// Kill task.
func (t *task) Kill() {
	// check t.ws if nil.
	if t.ws == nil {
		return
	}

	t.ws.stopAll()
}

func (t *task) assignWorkerSet(ws *workerSet) {
	t.ws = ws
}

func (t *task) updateStatus(status TaskStatus) {
	t.status = status
}

// TaskStatus is the status of task.
// TaskStatus implements TaskResult interface.
type TaskStatus int

const (
	// TaskStatusEnqueue means task is enqueued.
	TaskStatusEnqueue TaskStatus = iota
	// TaskStatusRunning means task is running.
	TaskStatusRunning
	// TaskStatusDone means task is done.
	TaskStatusDone
	// TaskStatusKilled means task is killed.
	TaskStatusKilled
	// TaskStatusError means task error.
	TaskStatusError
	// TaskStatusQueueFull means task queue is full.
	TaskStatusQueueFull
	// TaskStatusTooManyWorker means task concurrence is greater then worker pool max worker count.
	TaskStatusTooManyWorker
	// TaskStatusPoolClosed means task run worker pool closed.
	TaskStatusPoolClosed
)

func (t TaskStatus) String() string {
	switch t {
	case TaskStatusEnqueue:
		return "TaskStatusEnqueue"
	case TaskStatusRunning:
		return "TaskStatusRunning"
	case TaskStatusDone:
		return "TaskStatusDone"
	case TaskStatusKilled:
		return "TaskStatusKilled"
	case TaskStatusError:
		return "TaskStatusError"
	case TaskStatusQueueFull:
		return "TaskStatusQueueFull"
	case TaskStatusTooManyWorker:
		return "TaskStatusTooManyWorker"
	case TaskStatusPoolClosed:
		return "TaskStatusPoolClosed"
	default:
		return "TaskStatusUnknown"
	}
}

func (t TaskStatus) Error() string {
	return t.String()
}

func (t TaskStatus) Err() error {
	switch t {
	case TaskStatusError, TaskStatusQueueFull, TaskStatusTooManyWorker, TaskStatusPoolClosed, TaskStatusKilled:
		return t
	}

	return nil
}

func (t TaskStatus) Wait() {
	// do nothing.
}

func (t TaskStatus) Status() TaskStatus {
	return t
}

func (t TaskStatus) Kill() {
	// do nothing.
}
