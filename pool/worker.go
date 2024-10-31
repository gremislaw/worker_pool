package pool

type worker struct {
	id   int
	quit chan struct{}
}

func CreateWorker(id int) *worker {
	return &worker{id, make(chan struct{})}
}

func (w *worker) GetId() int {
	return w.id
}

func (w *worker) Kill() {
	close(w.quit)
}