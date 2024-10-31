package main

import (
	"os"
	"worker_pool/pool"
	"worker_pool/core"
)

func main() {
	workers := make(map[int]*pool.Worker)
	jobs := make(chan string, core.MaxBuffSize)
	results := make(chan string,core.MaxBuffSize)
	file, err := os.Create("bin/out.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	core.Run(workers, jobs, results, file)
	close(jobs)
}
