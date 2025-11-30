package tests

import (
	"sync/atomic"
	"testing"

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
