package cmd

import (
	"fmt"
	"worker_pool/pool"
)

func Run(wp *pool.WorkerPool) {
	format, n := "", 0
	fmt.Scan(&format)
	for ; ; fmt.Scan(&format) {
		switch format {
		case "add_job":
			fmt.Scan(&format)
			wp.AddJob(format)
		case "add_worker":
			fmt.Scan(&n)
			wp.AddWorker(n)
		case "delete_worker":
			fmt.Scan(&n)
			wp.DeleteWorker(n)
		case "stop":
			return
		default:
			fmt.Println("invalid format.\nAvailable format: add_job <data>, add_worker <id>, delete_worker <id>, stop")
		}
	}
}