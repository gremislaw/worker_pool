package main

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"worker_pool/pool"
	"worker_pool/core"
	"os"
)

var maxBuffSize = 100_000

func TestCreateWorker(t *testing.T) {
	assert := assert.New(t)
	id := 5
	w := pool.CreateWorker(id)
	if assert.NotNil(w) {
		assert.Equal(w.GetId(), id)
	}
}

func TestWorker(t *testing.T) {
	assert := assert.New(t)
	id := 5
	w := pool.CreateWorker(id)
	workers := make(map[int]*pool.Worker)
	jobs := make(chan string, maxBuffSize)
	results := make(chan string, maxBuffSize)
	file, err := os.Create("test.txt")
	if err != nil {
		panic(err)
	}
	if assert.NotNil(w) {
		w.Start(jobs, results, file)
		core.AddWorker(workers, id, jobs, results, file)
		core.AddJob("vk developer", jobs)
		assert.Equal(len(jobs), 1)
		w.Kill()
	}
}