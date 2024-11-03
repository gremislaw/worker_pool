package pool

import (
	"fmt"
	"os"
	"time"
)

/*
	 ВоркерПул
		Jobs: канал для передачи job'ов воркерам
		Results: канал для получения результатов обработки job'ов
		File: файл, в который будет записываться информация о происходящем в ВоркерПуле
*/
type WorkerPool struct {
	cntWorkers    int
	Jobs          chan string
	Results       chan string
	deleteWorkers chan struct{}
	File          *os.File
}

// Максимальный размер каналов
var MaxBuffSize = 100_000

// Cоздание нового ВоркерПула
func CreateWorkerPool(file *os.File) *WorkerPool {
	workers := 0
	jobs := make(chan string, MaxBuffSize)
	results := make(chan string, MaxBuffSize)
	deleteWorkers := make(chan struct{}, MaxBuffSize)
	return &WorkerPool{workers, jobs, results, deleteWorkers, file}
}

// Добавление нового воркера в ВоркерПул и сразу запуская воркер
func (wp *WorkerPool) AddWorkers(cnt int) {
	for i := 1; i <= cnt; i++ {
		wp.cntWorkers++
		wp.Write(fmt.Sprintf("Воркер %d добавлен.\n", i))
		go wp.startWorker(i) // Запуск воркера
	}
}

// Запуск воркера и обработка приходящих job'ов
func (wp *WorkerPool) startWorker(id int) {
	for {
		select {
		case j, ok := <-wp.Jobs: // тут приходят job'ы
			if !ok {
				return
			}
			wp.Write(fmt.Sprintf("Воркер %d обрабатывает строку: %s.\n", id, j))
			time.Sleep(time.Second * 3)
			wp.Write(fmt.Sprintf("Воркер %d обработал строку %s.\n", id, j))
			wp.Results <- j + "обработана"
		case <-wp.deleteWorkers: // Конец работы воркера
			wp.Write(fmt.Sprintf("Воркер %d удален.\n", id))
			wp.cntWorkers--
			return
		}
	}
}

// Добавление cnt новых job'ов в ВоркерПул со строкой data
func (wp *WorkerPool) AddJobs(cnt int, data string) {
	wp.Write(fmt.Sprintf("%d джобов со строкой %s добавлено.\n", cnt, data))
	for i := 1; i <= cnt; i++ {
		wp.Jobs <- data
	}
}

// Завершение работы воркера и его удаление из ВоркерПула
func (wp *WorkerPool) DeleteWorkers(cnt int) {
	for i := 1; i <= cnt; i++ {
		wp.deleteWorkers <- struct{}{}
	}
}

// Установка определенного количества воркеров
func (wp *WorkerPool) SetWorkers(cnt int) {
	wCnt := wp.GetWorkersCnt()
	if cnt > wp.cntWorkers {
		wp.AddWorkers(cnt - wCnt)
	} else {
		wp.DeleteWorkers(wCnt - cnt)
	}
}

// Получение количества воркеров в пуле
func (wp *WorkerPool) GetWorkersCnt() int {
	return wp.cntWorkers
}

// Получение количества воркеров для удаления с пула
func (wp *WorkerPool) GetWorkersCntForDelete() int {
	return len(wp.deleteWorkers)
}

// Запись в файл
func (wp *WorkerPool) Write(s string) {
	wp.File.WriteString(s)
}
