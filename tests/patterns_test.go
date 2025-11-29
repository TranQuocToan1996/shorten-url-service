package tests

import (
	"sync/atomic"
	"testing"
	"time"

	"shorten/pkg/utils/patterns"
)

func TestWorkerPool_AllTasksComplete(t *testing.T) {
	var completed int32
	pool := patterns.NewWorkerPool(4)
	pool.Run()

	numTasks := 10
	for i := 0; i < numTasks; i++ {
		pool.Submit(func() error {
			atomic.AddInt32(&completed, 1)
			return nil
		})
	}
	pool.Close()

	if atomic.LoadInt32(&completed) != int32(numTasks) {
		t.Fatalf("not all tasks completed, got %d, want %d", completed, numTasks)
	}
}

func TestWorkerPool_ErrorsAndPanics(t *testing.T) {
	pool := patterns.NewWorkerPool(2)
	pool.Run()

	var completed int32
	pool.Submit(func() error {
		atomic.AddInt32(&completed, 1)
		return nil
	})
	pool.Submit(func() error {
		atomic.AddInt32(&completed, 1)
		return assertError()
	})
	pool.Submit(func() error {
		atomic.AddInt32(&completed, 1)
		panic("panic in task")
	})
	pool.Close()

	if atomic.LoadInt32(&completed) != 3 {
		t.Fatalf("expected 3 tasks completed, got %d", completed)
	}
}

func assertError() error {
	return &customErr{"dummy error"}
}

type customErr struct{ msg string }

func (e *customErr) Error() string { return e.msg }

func TestSemaphore_BoundedConcurrency(t *testing.T) {
	maxConcurrent := 3
	s := patterns.NewSemaphore(maxConcurrent)
	n := 10
	var running int32
	var maxRunning int32

	for i := 0; i < n; i++ {
		s.Submit(func() error {
			r := atomic.AddInt32(&running, 1)
			for j := 0; j < 5; j++ {
				time.Sleep(1 * time.Millisecond)
			}
			for {
				rx := atomic.LoadInt32(&maxRunning)
				if r > rx {
					atomic.CompareAndSwapInt32(&maxRunning, rx, r)
					break
				}
				break
			}
			atomic.AddInt32(&running, -1)
			return nil
		})
	}
	s.Wait()
	if maxRunning > int32(maxConcurrent) {
		t.Fatalf("concurrent tasks exceeded max: got %d, want <= %d", maxRunning, maxConcurrent)
	}
}

func TestSemaphore_TaskErrorsAndPanics(t *testing.T) {
	s := patterns.NewSemaphore(2)
	var completed int32
	s.Submit(func() error {
		atomic.AddInt32(&completed, 1)
		return assertError()
	})
	s.Submit(func() error {
		atomic.AddInt32(&completed, 1)
		panic("panic in semaphore task")
	})
	s.Submit(func() error {
		atomic.AddInt32(&completed, 1)
		return nil
	})
	s.Wait()
	if atomic.LoadInt32(&completed) != 3 {
		t.Fatalf("expected 3 tasks completed, got %d", completed)
	}
}
