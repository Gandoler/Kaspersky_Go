package WokerPool_test

import (
	"Kaspersky_Go/APILevel/Adapters"
	"Kaspersky_Go/ModeLevel/Structures"
	"Kaspersky_Go/ServiceLevel/UseCases/Processors"
	"Kaspersky_Go/ServiceLevel/UseCases/WokerPool"
	"context"
	"log/slog"
	"os"
	"testing"
	"time"
)

func TestWorkerPool_ProcessesJobToDone(t *testing.T) {
	q := Adapters.NewMemoryQueue(2)
	st := Adapters.NewMemoryStateStore()
	p := Processors.NewRetryProcessor()
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	wp := WokerPool.NewWorkerPool(1, q, st, p, logger)

	// enqueue a job
	_ = q.Enqueue(Structures.Job{ID: "job1", MaxRetries: 2})

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	go wp.Start(ctx)

	// wait until state becomes done or timeout
	deadline := time.Now().Add(3 * time.Second)
	for time.Now().Before(deadline) {
		s, ok := st.Get("job1")
		if ok && (s.State == Structures.StateDone || s.State == Structures.StateQueued || s.State == Structures.StateRunning) {
			if s.State == Structures.StateDone {
				return
			}
		}
		time.Sleep(50 * time.Millisecond)
	}
	s, _ := st.Get("job1")
	t.Fatalf("job did not reach done, last state: %v", s.State)
}
