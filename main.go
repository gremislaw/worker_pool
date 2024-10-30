package main

import (
	"fmt"
	"time"
	"os"
)

var maxBuffSize = 100_000
var curJobNum = 1

type Worker struct {
	id int
	quit chan struct{}
}

func (w *Worker) kill() {
	close(w.quit)
}

func (w *Worker) start(jobs chan string, file *os.File) {
	file.WriteString(fmt.Sprintf("worker %d arrived.\n", w.id))
	for {
		select {
		case j, ok := <-jobs:
			if !ok {
					return
			}
			file.WriteString(fmt.Sprintf("worker %d started job %s.\n", w.id, j))
			time.Sleep(time.Second * 5)
			file.WriteString(fmt.Sprintf("worker %d finished job %s.\n", w.id, j))
		case <-w.quit:
			file.WriteString(fmt.Sprintf("worker %d kicked.\n", w.id))
			return
		}
	}
}

func add_jobs(n int, jobs chan string) {
	for j := 0; j < n; j++ {
		jobs <- fmt.Sprintf("%d", curJobNum)
		curJobNum++
	}
}

func solve(workers map[int]*Worker, jobs chan string, file *os.File) {
	format, n := "", 0
	fmt.Scan(&format)
	for ;; fmt.Scan(&format) {
		switch format {
		case "add_job":
			fmt.Scan(&n)
			add_jobs(n, jobs)
		case "add_worker":
			fmt.Scan(&n)
			if _, ok := workers[n]; ok {
				fmt.Println("such worker is exist")
			}
			w := &Worker{n, make(chan struct{})}
			workers[n] = w
			go w.start(jobs, file)
		case "delete_worker":
			fmt.Scan(&n)
			if w, ok := workers[n]; ok {
				w.kill()
			} else {
				fmt.Println("no such worker")
			}
		case "stop":
			return
		default:
			fmt.Println("invalid format")
		}
	}
}

func main() {
	workers := make(map[int]*Worker)
	jobs := make(chan	string, maxBuffSize)
	file, err := os.Create("bin/out.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	solve(workers, jobs, file)
	close(jobs)
}