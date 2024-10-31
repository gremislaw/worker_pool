package core

import (
	"fmt"
	"worker_pool/pool"
	"os"
)

var MaxBuffSize = 100_000

func AddJob(s string, jobs chan string) {
	jobs <- s
}

func AddWorker(workers map[int]*pool.Worker, id int, jobs, results chan string, file *os.File) {
	if _, ok := workers[id]; ok {
		fmt.Println("such worker is exist")
		return
	}
	w := pool.CreateWorker(id)
	workers[id] = w
	go w.Start(jobs, results, file)
}

func DeleteWorker(workers map[int]*pool.Worker, id int) {
	if _, ok := workers[id]; !ok {
		fmt.Println("\nno such worker")
		return
	}
	w := workers[id]
	w.Kill()
	delete(workers, id)
}

func Run(workers map[int]*pool.Worker, jobs, results chan string, file *os.File) {
	format, n := "", 0
	fmt.Scan(&format)
	for ; ; fmt.Scan(&format) {
		switch format {
		case "add_job":
			fmt.Scan(&format)
			AddJob(format, jobs)
		case "add_worker":
			fmt.Scan(&n)
			AddWorker(workers, n, jobs, results, file)
		case "delete_worker":
			fmt.Scan(&n)
			DeleteWorker(workers, n)
		case "stop":
			return
		default:
			fmt.Println("invalid format.\nAvailable format: add_job <data>, add_worker <id>, delete_worker <id>, stop")
		}
	}
}