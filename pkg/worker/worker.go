package worker

// worker run tasks
type worker struct {
	ws  *workerSet
	err error
}

func newWorker(ws *workerSet) *worker {
	return &worker{
		ws: ws,
	}
}

func (w *worker) run() {
	defer w.ws.done()

	var ec = make(chan error, 1)
	defer close(ec)
	go func() {
		ec <- w.ws.task.tf()
	}()

	select {
	case e := <-ec:
		w.err = e
	case <-w.ws.ctx.Done():
	}
}
