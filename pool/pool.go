package pool

import (
	"errors"
	"fmt"
	"os"
	"sync"
	"time"
)

// Для обработки ошибок
const (
	NoSuchWorkerError       = "такого воркера нет"
	WorkerExistsError       = "такой воркер существует"
	CannotCreateWorkerError = "не удалось создать воркер"
)

/*
	 ВоркерПул
		workers: потокобезопасная мапа, хранящая воркеров по их номеру
		Jobs: канал для передачи job'ов воркерам
		Results: канал для получения результатов обработки job'ов
		File: файл, в который будет записываться информация о происходящем в ВоркерПуле
*/
type WorkerPool struct {
	workers *sync.Map
	Jobs    chan string
	Results chan string
	File    *os.File
}

// Максимальный размер каналов
var MaxBuffSize = 100_000

// Cоздание нового ВоркерПула
func CreateWorkerPool(file *os.File) *WorkerPool {
	workers := new(sync.Map)
	jobs := make(chan string, MaxBuffSize)
	results := make(chan string, MaxBuffSize)
	return &WorkerPool{workers, jobs, results, file}
}

// Добавление нового воркера в ВоркерПул и сразу запуская воркер
func (wp *WorkerPool) AddWorker(id int) error {
	if _, ok := wp.workers.Load(id); ok {
		fmt.Println(WorkerExistsError)
		return errors.New(WorkerExistsError)
	}
	w := CreateWorker(id)
	if w == nil { // Проверка на валидость созданного воркера
		fmt.Println(CannotCreateWorkerError)
		return errors.New(CannotCreateWorkerError)
	}
	wp.SetWorker(w) // Добавление воркера в ВоркерПул
	wp.Write(fmt.Sprintf("Воркер %d добавлен.\n", id))
	go wp.startWorker(w) // Запуск воркера
	return nil
}

// Запуск воркера и обработка приходящих job'ов
func (wp *WorkerPool) startWorker(w *worker) {
	for {
		select {
		case j, ok := <-wp.Jobs: // тут приходят job'ы
			if !ok {
				return
			}
			wp.Write(fmt.Sprintf("Воркер %d обрабатывает строку: %s.\n", w.id, j))
			time.Sleep(time.Second * 3)
			wp.Write(fmt.Sprintf("Воркер %d обработал строку %s.\n", w.id, j))
			wp.Results <- j + "обработана"
		case <-w.quit: // Конец работы воркера
			wp.Write(fmt.Sprintf("Воркер %d удален.\n", w.id))
			return
		}
	}
}

// Добавление нового job'а в ВоркерПул
func (wp *WorkerPool) AddJob(s string) {
	wp.Jobs <- s
}

// Завершение работы воркера и его удаление из ВоркерПула
func (wp *WorkerPool) DeleteWorker(id int) error {
	if _, ok := wp.workers.Load(id); !ok {
		fmt.Println(NoSuchWorkerError)
		return errors.New(NoSuchWorkerError)
	}
	w := wp.GetWorker(id)
	w.Kill()
	wp.workers.Delete(id)
	return nil
}

// Получение воркера из ВоркерПула по номеру
func (wp *WorkerPool) GetWorker(id int) *worker {
	w, ok := wp.workers.Load(id)
	if !ok {
		return nil
	}
	return w.(*worker)
}

// Добавление воркера в Воркерпул
func (wp *WorkerPool) SetWorker(w *worker) {
	wp.workers.Store(w.id, w)
}

// Запись в файл
func (wp *WorkerPool) Write(s string) {
	wp.File.WriteString(s)
}
