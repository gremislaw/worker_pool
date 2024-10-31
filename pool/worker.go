package pool

/*
	 Воркер
		id: номер воркера
		quit: канал для завершения горутины, обрабатывающей job'ы этим воркером
*/
type worker struct {
	id   int
	quit chan struct{}
}

// Создание нового воркера
func CreateWorker(id int) *worker {
	return &worker{id, make(chan struct{})}
}

// Получение номера воркера
func (w *worker) GetId() int {
	return w.id
}

// Завершить обработку job'ов для этого воркера
func (w *worker) Kill() {
	close(w.quit)
}
