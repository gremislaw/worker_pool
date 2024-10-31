package pool

import (
	"fmt"
	"os"
	"time"
)

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

func (w *Worker) Start(jobs chan string, results chan string, file *os.File) {
	file.WriteString(fmt.Sprintf("Воркер %d добавлен.\n", w.id))
	for {
		select {
		case j, ok := <-jobs:
			if !ok {
				return
			}
			file.WriteString(fmt.Sprintf("Воркер %d обрабатывает строку: %s.\n", w.id, j))
			time.Sleep(time.Second * 5)
			file.WriteString(fmt.Sprintf("Воркер %d обработал строку %s.\n", w.id, j))
			results <- j + "обработана"
		case <-w.quit:
			file.WriteString(fmt.Sprintf("Воркер %d удален.\n", w.id))
			return
		}
	}
}
