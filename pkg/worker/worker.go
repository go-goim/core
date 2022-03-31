package worker

// worker run tasks
type worker struct {
	ws   *workerSet
	task TaskFunc
	err  error
}

func newWorker(ws *workerSet, t TaskFunc) *worker {
	return &worker{
		ws:   ws,
		task: t,
	}
}

func (w *worker) run() {
	defer w.ws.done()

	var ec = make(chan error, 1)
	defer close(ec)
	go func() {
		ec <- w.task()
	}()

	select {
	case e := <-ec:
		w.err = e
	case <-w.ws.ctx.Done():
	}
}
