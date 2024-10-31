package pool

import (
	"fmt"
	"os"
	"time"
	"sync"
	"errors"
)

var MaxBuffSize = 100_000

type WorkerPool struct {
	muFile *sync.Mutex
	muMap *sync.Mutex
	workers map[int]*worker
	Jobs chan string
	Results chan string
	File *os.File
}

func CreateWorkerPool(file *os.File) *WorkerPool{
	muFile := new(sync.Mutex)
	muMap := new(sync.Mutex)
	workers := make(map[int]*worker)
	jobs := make(chan string, MaxBuffSize)
	results := make(chan string, MaxBuffSize)
	return &WorkerPool{muFile, muMap, workers, jobs, results, file}
}

func (wp *WorkerPool) AddWorker(id int) error  {
	if _, ok := wp.workers[id]; ok {
		fmt.Println("\nтакой воркер уже существует")
		return errors.New("такой воркер уже существует")
	}
	wp.SetWorker(id)
	wp.safeWrite(fmt.Sprintf("Воркер %d добавлен.\n", id))
	go wp.StartWorker(id)
	return nil
}

func (wp *WorkerPool) StartWorker(id int) {
	w := wp.GetWorker(id)
	if w == nil {
		fmt.Printf("\nворкер с ID %d не найден\n", id)
		return
}
	for {
		select {
		case j, ok := <-wp.Jobs:
			if !ok {
				return
			}
			wp.safeWrite(fmt.Sprintf("Воркер %d обрабатывает строку: %s.\n", id, j))
			time.Sleep(time.Second * 3)
			wp.safeWrite(fmt.Sprintf("Воркер %d обработал строку %s.\n", id, j))
			wp.Results <- j + "обработана"
		case <-w.quit:
			wp.safeWrite(fmt.Sprintf("Воркер %d удален.\n", id))
			return
		}
	}
}

func (wp *WorkerPool) AddJob(s string) {
	wp.Jobs <- s
}

func (wp *WorkerPool) DeleteWorker(id int) error {
	if _, ok := wp.workers[id]; !ok {
		fmt.Println("\nтакого воркера нет")
		return errors.New("такого воркера нет")
	}
	w := wp.GetWorker(id)
	w.Kill()
	delete(wp.workers, id)
	return nil
}

func (wp *WorkerPool) GetWorker(id int) *worker {
	wp.muMap.Lock()
	w := wp.workers[id]
	wp.muMap.Unlock()
	return w
}

func (wp *WorkerPool) SetWorker(id int) {
	w := CreateWorker(id)
	wp.muMap.Lock()
	wp.workers[id] = w
	wp.muMap.Unlock()
}

func (wp *WorkerPool) safeWrite(s string) {
	wp.muFile.Lock()
	wp.File.WriteString(s)
	wp.muFile.Unlock()
}

