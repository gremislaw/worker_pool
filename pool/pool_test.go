package pool

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
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

func TestWorkerPool_AddWorkers(t *testing.T) {
	assert := assert.New(t)
	cnt := 10
	var file *os.File = nil
	wp := CreateWorkerPool(file)
	wp.AddWorkers(cnt)
	assert.Equal(int(wp.cntWorkers.Load()), cnt)
}

func TestWorkerPool_AddWorkers_Many(t *testing.T) {
	assert := assert.New(t)
	cnt := 100000
	var file *os.File = nil
	wp := CreateWorkerPool(file)
	wp.AddWorkers(cnt)
	assert.Equal(int(wp.cntWorkers.Load()), cnt)
}

func TestWorkerPool_AddWorkers_ZeroCnt(t *testing.T) {
	assert := assert.New(t)
	cnt := 0
	var file *os.File = nil
	wp := CreateWorkerPool(file)
	wp.AddWorkers(cnt)
	assert.Equal(int(wp.cntWorkers.Load()), cnt)
}

func TestWorkerPool_AddWorkers_NegativeCnt(t *testing.T) {
	assert := assert.New(t)
	cnt := -5
	var file *os.File = nil
	wp := CreateWorkerPool(file)
	assert.NotPanics(func() { wp.AddWorkers(cnt) })
}

func TestWorkerPool_AddJob(t *testing.T) {
	assert := assert.New(t)
	var file *os.File = nil
	wp := CreateWorkerPool(file)
	wp.AddJobs(1, "vk developer")
	wp.AddJobs(1, "intern backend")
	assert.Equal(len(wp.Jobs), 2)
}

func TestWorkerPool_DeleteWorkers(t *testing.T) {
	assert := assert.New(t)
	var file *os.File = nil
	n := 10
	wp := CreateWorkerPool(file)
	wp.DeleteWorkers(n)
	assert.Equal(wp.GetWorkersCntForDelete(), n)
}

func TestWorkerPool_DeleteWorkers_Many(t *testing.T) {
	assert := assert.New(t)
	var file *os.File = nil
	n := 100000
	wp := CreateWorkerPool(file)
	wp.DeleteWorkers(n)
	assert.Equal(wp.GetWorkersCntForDelete(), n)
}

func TestWorkerPool_DeleteWorker_ZeroCnt(t *testing.T) {
	assert := assert.New(t)
	var file *os.File = nil
	cnt := 0
	wp := CreateWorkerPool(file)
	wp.DeleteWorkers(cnt)
	assert.Equal(wp.GetWorkersCntForDelete(), cnt)
}

func TestWorkerPool_DeleteWorker_NegativeCnt(t *testing.T) {
	assert := assert.New(t)
	var file *os.File = nil
	cnt := -10
	wp := CreateWorkerPool(file)
	assert.NotPanics(func() { wp.DeleteWorkers(cnt) })
}

func TestWorkerPool_WorkerDoJobs(t *testing.T) {
	assert := assert.New(t)
	n := 4
	var file *os.File = nil
	wp := CreateWorkerPool(file)
	if assert.NotNil(wp) {
		wp.AddWorkers(n)
		wp.AddJobs(1, "vk developer")
		wp.AddJobs(1, "intern backend")
		<-wp.Results
		<-wp.Results
		assert.Equal(len(wp.Jobs), 0)
		wp.DeleteWorkers(n)
	}
}

func TestWorkerPool_ManyWorkersDoJobs(t *testing.T) {
	assert := assert.New(t)
	var file *os.File = nil
	n := 100000
	wp := CreateWorkerPool(file)
	if assert.NotNil(wp) {
		wp.AddWorkers(n)
		wp.AddWorkers(n)
		wp.AddJobs(1, "vk developer")
		wp.AddJobs(1, "intern backend")
		<-wp.Results
		<-wp.Results
		assert.Equal(len(wp.Jobs), 0)
	}
}

func TestWorkerPool_ManyWorkersDoManyJobs(t *testing.T) {
	assert := assert.New(t)
	var file *os.File = nil
	n := 100000
	wp := CreateWorkerPool(file)
	if assert.NotNil(wp) {
		wp.AddWorkers(n)
		wp.AddJobs(n, "vk developer")
		for i := 1; i < 100000; i++ {
			<-wp.Results
		}
		assert.Equal(len(wp.Jobs), 0)
	}
}

func TestWorkerPool_ManyWorkersDoManyJobs_DynamicAddWorkers(t *testing.T) {
	assert := assert.New(t)
	var file *os.File = nil
	wp := CreateWorkerPool(file)
	if assert.NotNil(wp) {
		wp.AddWorkers(10000)
		wp.AddJobs(100000, "vk developer")
		assert.Greater(len(wp.Jobs), 10000)
		for i := 1; i < 10000; i++ {
			<-wp.Results
		}
		wp.AddWorkers(90000)
		for i := 10000; i < 100000; i++ {
			<-wp.Results
		}
		assert.Equal(len(wp.Jobs), 0)
	}
}

func TestWorkerPool_ManyWorkersDoManyJobs_DynamicDeleteWorkers(t *testing.T) {
	assert := assert.New(t)
	var file *os.File = nil
	wp := CreateWorkerPool(file)
	if assert.NotNil(wp) {
		wp.AddWorkers(60000)
		wp.AddJobs(100000, "vk developer")
		assert.Greater(len(wp.Jobs), 10000)
		for i := 1; i < 10000; i++ {
			<-wp.Results
		}
		wp.DeleteWorkers(10000)
		for i := 10000; i < 100000; i++ {
			<-wp.Results
		}
		assert.Equal(len(wp.Jobs), 0)
	}
}

func TestWorkerPool_ManyWorkersDoManyJobs_DynamicAddDeleteWorkers(t *testing.T) {
	assert := assert.New(t)
	var file *os.File = nil
	wp := CreateWorkerPool(file)
	if assert.NotNil(wp) {
		wp.AddWorkers(10000)
		wp.AddJobs(100000, "vk developer")
		assert.Greater(len(wp.Jobs), 10000)
		for i := 1; i < 10000; i++ {
			<-wp.Results
		}
		wp.DeleteWorkers(5000)
		for i := 1; i < 10000; i++ {
			<-wp.Results
		}
		wp.AddWorkers(80000)
		for i := 20000; i < 100000; i++ {
			<-wp.Results
		}
		assert.Equal(len(wp.Jobs), 0)
	}
}
