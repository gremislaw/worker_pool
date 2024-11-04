package pool

import (
	"fmt"
	"math/rand"
	"os"
	"sync/atomic"
	"time"
)

/*
	 ВоркерПул
	  cntWorkers: количество воркеров в пуле
	 	curId: автоинкремент для предотвращения коллизий
		Jobs: канал для передачи job'ов воркерам
		cntJobs: количество работ в пуле
		Results: канал для получения результатов обработки job'ов
		deleteWorkers: канал для удаления воркеров
		File: файл, в который будет записываться информация о происходящем в ВоркерПуле
*/
type WorkerPool struct {
	cntWorkers    *atomic.Int64
	curId         *atomic.Int64
	Jobs          chan string
	cntJobs       *atomic.Int64
	Results       chan string
	deleteWorkers chan struct{}
	File          *os.File
}

// Максимальный размер каналов
var MaxBuffSize = 100_000

// Cоздание нового ВоркерПула
func CreateWorkerPool(file *os.File) *WorkerPool {
	workers := new(atomic.Int64)
	curId := new(atomic.Int64)
	jobs := make(chan string, MaxBuffSize)
	cntJobs := new(atomic.Int64)
	results := make(chan string, MaxBuffSize)
	deleteWorkers := make(chan struct{}, MaxBuffSize)
	return &WorkerPool{workers, curId, jobs, cntJobs, results, deleteWorkers, file}
}

// Добавление нового воркера в ВоркерПул и сразу запуская воркер
func (wp *WorkerPool) AddWorkers(cnt int) {
	if cnt + int(wp.cntWorkers.Load()) > MaxBuffSize && cnt >= 0{
		cnt = MaxBuffSize - wp.GetWorkersCnt()
	}
	for i := 1; i <= cnt; i++ {
		wp.cntWorkers.Add(1)
		wp.curId.Add(1)
		wp.Write(fmt.Sprintf("Воркер %d добавлен.\n", int(wp.curId.Load())))
		go wp.startWorker(int(wp.curId.Load())) // Запуск воркера
	}
}

// Запуск воркера и обработка приходящих job'ов
func (wp *WorkerPool) startWorker(id int) {
	for {
		select {
		case j, ok := <-wp.Jobs: // тут приходят job'ы
			if !ok {
				wp.Write(fmt.Sprintf("Воркер %d удален.\n", id))
				wp.cntWorkers.Add(-1)
				return
			}
			wp.Write(fmt.Sprintf("Воркер %d обрабатывает строку: %s.\n", id, j))
			time.Sleep(time.Second + time.Millisecond * time.Duration(rand.Int31n(3000)))
			wp.Write(fmt.Sprintf("Воркер %d обработал строку %s.\n", id, j))
			wp.Results <- j + "обработана"
			wp.ClearOutChannel()
			wp.cntJobs.Add(-1)
		case <-wp.deleteWorkers: // Конец работы воркера
			wp.Write(fmt.Sprintf("Воркер %d удален.\n", id))
			wp.cntWorkers.Add(-1)
			return
		}
	}
}

// Добавление cnt новых job'ов в ВоркерПул со строкой data
func (wp *WorkerPool) AddJobs(cnt int, data string) {
	if cnt + int(wp.cntJobs.Load()) > MaxBuffSize && cnt >= 0{
		cnt = MaxBuffSize - wp.GetJobCnt()
	}
	wp.Write(fmt.Sprintf("%d джобов со строкой %s добавлено.\n", cnt, data))
	for i := 1; i <= cnt; i++ {
		wp.Jobs <- data
		wp.cntJobs.Add(1)
	}
}

// Завершение работы воркера и его удаление из ВоркерПула
func (wp *WorkerPool) DeleteWorkers(cnt int) {
	if cnt + wp.GetWorkersCntForDelete() > MaxBuffSize && cnt >= 0{
		cnt = MaxBuffSize - wp.GetWorkersCntForDelete() 
	}
	if cnt > wp.GetWorkersCnt() {
		cnt = wp.GetWorkersCnt()
	}
	for i := 1; i <= cnt; i++ {
		wp.deleteWorkers <- struct{}{}
	}
}

// Установка определенного количества воркеров
func (wp *WorkerPool) SetWorkers(cnt int) {
	wCnt := wp.GetWorkersCnt()
	if cnt > wp.GetWorkersCnt() {
		wp.AddWorkers(cnt - wCnt)
	} else {
		wp.DeleteWorkers(wCnt - cnt)
	}
}

// Получение количества воркеров в пуле
func (wp *WorkerPool) GetWorkersCnt() int {
	return int(wp.cntWorkers.Load())
}

// Получение количества работ в пуле
func (wp *WorkerPool) GetJobCnt() int {
	return int(wp.cntJobs.Load())
}

// Получение количества воркеров для удаления с пула
func (wp *WorkerPool) GetWorkersCntForDelete() int {
	return len(wp.deleteWorkers)
}

// Запись в файл
func (wp *WorkerPool) Write(s string) {
	wp.File.WriteString(s)
}

func (wp *WorkerPool) ClearOutChannel() {
	if len(wp.Results) == MaxBuffSize {
		for {
			select {
			case <-wp.Results:
				// Игнорируем элемент
			default:
				return
			}
		}
	}
}