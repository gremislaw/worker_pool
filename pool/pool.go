package pool

import (
	"fmt"
	"os"
	"time"
)

var MaxBuffSize = 100_000

type WorkerPool struct {
	Workers map[int]*Worker
	Jobs chan string
	Results chan string
	File *os.File
}

func CreateWorkerPool(file *os.File) *WorkerPool{
	workers := make(map[int]*Worker)
	jobs := make(chan string, MaxBuffSize)
	results := make(chan string, MaxBuffSize)
	return &WorkerPool{workers, jobs, results, file}
}

func (wp *WorkerPool) AddWorker(id int) {
	if _, ok := wp.Workers[id]; ok {
		fmt.Println("such worker is exist")
		return
	}
	w := CreateWorker(id)
	wp.Workers[id] = w
	go wp.StartWorker(id)
}

func (wp *WorkerPool) StartWorker(id int) {
	wp.File.WriteString(fmt.Sprintf("Воркер %d добавлен.\n", id))
	for {
		select {
		case j, ok := <-wp.Jobs:
			if !ok {
				return
			}
			wp.File.WriteString(fmt.Sprintf("Воркер %d обрабатывает строку: %s.\n", id, j))
			time.Sleep(time.Second * 5)
			wp.File.WriteString(fmt.Sprintf("Воркер %d обработал строку %s.\n", id, j))
			wp.Results <- j + "обработана"
		case <-wp.Workers[id].quit:
			wp.File.WriteString(fmt.Sprintf("Воркер %d удален.\n", id))
			return
		}
	}
}

func (wp *WorkerPool) AddJob(s string) {
	wp.Jobs <- s
}

func (wp *WorkerPool) DeleteWorker(id int) {
	if _, ok := wp.Workers[id]; !ok {
		fmt.Println("\nno such worker")
		return
	}
	w := wp.Workers[id]
	w.Kill()
	delete(wp.Workers, id)
}