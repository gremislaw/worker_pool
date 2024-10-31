package main

import (
	"os"
	"worker_pool/cmd"
	"worker_pool/pool"
)

func main() {
	file, err := os.Create("bin/out.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	wp := pool.CreateWorkerPool(file)
	cmd.Run(wp)
	close(wp.Jobs)
}
