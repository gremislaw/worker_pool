package pool

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestWorker_Create(t *testing.T) {
	assert := assert.New(t)
	id := 5
	w := CreateWorker(id)
	assert.NotNil(w)
}

func TestWorker_Create_NegativeId(t *testing.T) {
	assert := assert.New(t)
	id := -5
	w := CreateWorker(id)
	assert.NotNil(w)
}

func TestWorker(t *testing.T) {
	assert := assert.New(t)
	id := 5
	w := CreateWorker(id)
	if assert.NotNil(w) {
		assert.Equal(w.GetId(), id)
	}
}

func TestWorkerPool_Create(t *testing.T) {
	assert := assert.New(t)
	file, err := os.Create("test.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	wp := CreateWorkerPool(file)
	assert.NotNil(wp)
}

func TestWorkerPool_CreateWithNilFile(t *testing.T) {
	assert := assert.New(t)
	var file *os.File = nil
	wp := CreateWorkerPool(file)
	assert.NotNil(wp)
}

func TestWorkerPool_AddWorker(t *testing.T) {
	assert := assert.New(t)
	id := -5
	var file *os.File = nil
	wp := CreateWorkerPool(file)
	wp.AddWorker(id)
	assert.NotNil(wp.GetWorker(id))
	wp.DeleteWorker(id)
}

func TestWorkerPool_AddWorker_Exists(t *testing.T) {
	assert := assert.New(t)
	id := -5
	var file *os.File = nil
	wp := CreateWorkerPool(file)
	err := wp.AddWorker(id)
	assert.Nil(err)
	err = wp.AddWorker(id)
	assert.NotNil(err)
	wp.DeleteWorker(id)
}

func TestWorkerPool_AddWorker_NegativeId(t *testing.T) {
	assert := assert.New(t)
	id := -5
	var file *os.File = nil
	wp := CreateWorkerPool(file)
	wp.AddWorker(id)
	assert.NotNil(wp.GetWorker(id))
	wp.DeleteWorker(id)
}

func TestWorkerPool_AddJob(t *testing.T) {
	assert := assert.New(t)
	var file *os.File = nil
	wp := CreateWorkerPool(file)
	wp.AddJob("vk developer")
	wp.AddJob("intern backend")
	assert.Equal(len(wp.Jobs), 2)
}

func TestWorkerPool_DeleteWorker(t *testing.T) {
	assert := assert.New(t)
	var file *os.File = nil
	id := 5
	wp := CreateWorkerPool(file)
	wp.AddWorker(id)
	wp.DeleteWorker(id)
	assert.Nil(wp.GetWorker(id))
}

func TestWorkerPool_DeleteWorker_NotExist(t *testing.T) {
	assert := assert.New(t)
	var file *os.File = nil
	id := 5
	wp := CreateWorkerPool(file)
	wp.AddWorker(id)
	err := wp.DeleteWorker(id)
	assert.Nil(err)
	err = wp.DeleteWorker(id)
	assert.Nil(wp.GetWorker(id))
	assert.NotNil(err)
}

func TestWorkerPool_WorkerDoJobs(t *testing.T) {
	assert := assert.New(t)
	id := -5
	var file *os.File = nil
	wp := CreateWorkerPool(file)
	if assert.NotNil(wp) {
		wp.AddWorker(id)
		wp.AddJob("vk developer")
		wp.AddJob("intern backend")
		<-wp.Results
		<-wp.Results
		assert.Equal(len(wp.Jobs), 0)
		wp.DeleteWorker(id)
	}
}

func TestWorkerPool_ManyWorkersDoJobs(t *testing.T) {
	assert := assert.New(t)
	var file *os.File = nil
	wp := CreateWorkerPool(file)
	if assert.NotNil(wp) {
		for i := 1; i < 100000; i++ {
			wp.AddWorker(i)
		}
		wp.AddJob("vk developer")
		wp.AddJob("intern backend")
		<-wp.Results
		<-wp.Results
		assert.Equal(len(wp.Jobs), 0)
		for i := 1; i < 100000; i++ {
			wp.DeleteWorker(i)
		}
	}
}

func TestWorkerPool_ManyWorkersDoManyJobs(t *testing.T) {
	assert := assert.New(t)
	var file *os.File = nil
	wp := CreateWorkerPool(file)
	if assert.NotNil(wp) {
		for i := 1; i < 100000; i++ {
			wp.AddWorker(i)
		}
		for i := 1; i < 100000; i++ {
			wp.AddJob("vk developer")
		}
		for i := 1; i < 100000; i++ {
			<-wp.Results
		}
		assert.Equal(len(wp.Jobs), 0)
		for i := 1; i < 100000; i++ {
			wp.DeleteWorker(i)
		}
	}
}

func TestWorkerPool_ManyWorkersDoManyJobs_DynamicAddWorkers(t *testing.T) {
	assert := assert.New(t)
	var file *os.File = nil
	wp := CreateWorkerPool(file)
	if assert.NotNil(wp) {
		for i := 1; i < 10000; i++ {
			wp.AddWorker(i)
		}
		for i := 1; i < 100000; i++ {
			wp.AddJob("vk developer")
		}
		assert.Greater(len(wp.Jobs), 10000)
		for i := 1; i < 10000; i++ {
			<-wp.Results
		}
		for i := 10000; i < 100000; i++ {
			wp.AddWorker(i)
		}
		for i := 10000; i < 100000; i++ {
			<-wp.Results
		}
		assert.Equal(len(wp.Jobs), 0)
		for i := 1; i < 100000; i++ {
			wp.DeleteWorker(i)
		}
	}
}

func TestWorkerPool_ManyWorkersDoManyJobs_DynamicDeleteWorkers(t *testing.T) {
	assert := assert.New(t)
	var file *os.File = nil
	wp := CreateWorkerPool(file)
	if assert.NotNil(wp) {
		for i := 1; i < 60000; i++ {
			wp.AddWorker(i)
		}
		for i := 1; i < 100000; i++ {
			wp.AddJob("vk developer")
		}
		assert.Greater(len(wp.Jobs), 10000)
		for i := 1; i < 10000; i++ {
			<-wp.Results
		}
		for i := 1; i < 10000; i++ {
			wp.DeleteWorker(i)
		}
		for i := 10000; i < 100000; i++ {
			<-wp.Results
		}
		assert.Equal(len(wp.Jobs), 0)
		for i := 10000; i < 60000; i++ {
			wp.DeleteWorker(i)
		}
	}
}

func TestWorkerPool_ManyWorkersDoManyJobs_DynamicAddDeleteWorkers(t *testing.T) {
	assert := assert.New(t)
	var file *os.File = nil
	wp := CreateWorkerPool(file)
	if assert.NotNil(wp) {
		for i := 1; i < 10000; i++ {
			wp.AddWorker(i)
		}
		for i := 1; i < 100000; i++ {
			wp.AddJob("vk developer")
		}
		assert.Greater(len(wp.Jobs), 10000)
		for i := 1; i < 10000; i++ {
			<-wp.Results
		}
		for i := 1; i < 5000; i++ {
			wp.DeleteWorker(i)
		}
		for i := 10000; i < 100000; i++ {
			wp.AddWorker(i)
		}
		for i := 10000; i < 100000; i++ {
			<-wp.Results
		}
		assert.Equal(len(wp.Jobs), 0)
		for i := 5000; i < 100000; i++ {
			wp.DeleteWorker(i)
		}
	}
}
