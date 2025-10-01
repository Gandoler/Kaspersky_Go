package tests

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"TASK_2/WorkerPool/Realization"
)

func TestNewWorkerPool_InvalidWorkers(t *testing.T) {
	if Realization.NewWorkerPool(0) != nil {
		t.Fatalf("expected nil for 0 workers")
	}
}

func TestSubmit_ExecutesAllTasks(t *testing.T) {
	wp := Realization.NewWorkerPool(3)
	if wp == nil {
		t.Fatal("nil worker pool")
	}
	var executed int64
	const total = 20
	for i := 0; i < total; i++ {
		wp.Submit(func() {
			atomic.AddInt64(&executed, 1)
		})
	}
	wp.StopWait()
	if executed != total {
		t.Fatalf("expected %d executed, got %d", total, executed)
	}
}

func TestSubmit_NilTaskIgnored(t *testing.T) {
	wp := Realization.NewWorkerPool(1)
	if wp == nil {
		t.Fatal("nil worker pool")
	}
	// Should not panic
	wp.Submit(nil)
	wp.SubmitWait(nil)
	wp.StopWait()
}

func TestSubmitWait_BlocksUntilTaskDone(t *testing.T) {
	wp := Realization.NewWorkerPool(1)
	if wp == nil {
		t.Fatal("nil worker pool")
	}
	start := time.Now()
	const delay = 150 * time.Millisecond
	wp.SubmitWait(func() { time.Sleep(delay) })
	elapsed := time.Since(start)
	if elapsed < delay {
		t.Fatalf("SubmitWait returned too early: %v < %v", elapsed, delay)
	}
	wp.StopWait()
}

func TestQueueCapacity_DropsWhenFull(t *testing.T) {
	wp := Realization.NewWorkerPool(1)
	if wp == nil {
		t.Fatal("nil worker pool")
	}
	// Block the single worker so queued tasks accumulate up to capacity (64)
	blockCh := make(chan struct{})
	wp.Submit(func() { <-blockCh })

	var executed int64
	// Enqueue up to the channel capacity with quick tasks
	for i := 0; i < 64; i++ {
		wp.Submit(func() { atomic.AddInt64(&executed, 1) })
	}
	// This one should be rejected due to full queue
	wp.Submit(func() { atomic.AddInt64(&executed, 1) })

	// Unblock worker so it can drain the queue
	close(blockCh)
	wp.StopWait()

	if executed > 64 {
		t.Fatalf("expected at most 64 tasks executed, got %d", executed)
	}
}

func TestStop_PreventsNewSubmissions(t *testing.T) {
	wp := Realization.NewWorkerPool(2)
	if wp == nil {
		t.Fatal("nil worker pool")
	}
	var executed int64
	wp.Stop()
	// Any further submissions should be ignored
	wp.Submit(func() { atomic.AddInt64(&executed, 1) })
	wp.SubmitWait(func() { atomic.AddInt64(&executed, 1) })
	wp.StopWait()
	if executed != 0 {
		t.Fatalf("expected no tasks after Stop, got %d", executed)
	}
}

func TestStopWait_WaitsForAllQueuedTasks(t *testing.T) {
	wp := Realization.NewWorkerPool(3)
	if wp == nil {
		t.Fatal("nil worker pool")
	}
	var executed int64
	const total = 50
	for i := 0; i < total; i++ {
		wp.Submit(func() {
			time.Sleep(5 * time.Millisecond)
			atomic.AddInt64(&executed, 1)
		})
	}
	wp.StopWait()
	if executed != total {
		t.Fatalf("expected %d executed, got %d", total, executed)
	}
}

func TestStop_Idempotent(t *testing.T) {
	wp := Realization.NewWorkerPool(2)
	if wp == nil {
		t.Fatal("nil worker pool")
	}
	var executed int64
	wp.Submit(func() { atomic.AddInt64(&executed, 1) })
	wp.Stop()
	// Second Stop should be safe
	wp.Stop()
	wp.StopWait()
	if executed != 1 {
		t.Fatalf("expected 1 executed before stop, got %d", executed)
	}
}

func TestConcurrent_SubmitAndSubmitWait(t *testing.T) {
	wp := Realization.NewWorkerPool(5)
	if wp == nil {
		t.Fatal("nil worker pool")
	}
	var executed int64
	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			wp.Submit(func() { atomic.AddInt64(&executed, 1) })
		}()
	}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			wp.SubmitWait(func() { atomic.AddInt64(&executed, 1) })
		}()
	}
	wg.Wait()
	wp.StopWait()
	if executed != 30 {
		t.Fatalf("expected 30 executed, got %d", executed)
	}
}
