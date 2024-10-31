package pool

type Worker struct {
	id   int
	quit chan struct{}
}

func CreateWorker(id int) *Worker {
	return &Worker{id, make(chan struct{})}
}

func (w *Worker) GetId() int {
	return w.id
}

func (w *Worker) Kill() {
	close(w.quit)
}