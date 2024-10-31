package main

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"worker_pool/pool"
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
	if assert.NotNil(w) {
		assert.Equal(w.GetId(), 5)
		w.Kill()
	}
}

func TestWorkerPool(t *testing.T) {
	assert := assert.New(t)
	id := 5
	file, err := os.Create("bin/out.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	wp := pool.CreateWorkerPool(file)
	if assert.NotNil(wp) {
		wp.StartWorker(id)
		wp.AddWorker(id)
		wp.AddJob("vk developer")
		assert.Equal(len(wp.Jobs), 1)
		wp.Workers[id].Kill()
	}
}