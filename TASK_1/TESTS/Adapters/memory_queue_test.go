package Adapters_test

import (
	"Kaspersky_Go/APILevel/Adapters"
	"Kaspersky_Go/ModeLevel/Structures"
	"testing"
)

func TestMemoryQueue_EnqueueDequeue(t *testing.T) {
	q := Adapters.NewMemoryQueue(2)
	job := Structures.Job{ID: "1"}

	if err := q.Enqueue(job); err != nil {
		t.Fatalf("enqueue error: %v", err)
	}

	got, ok := q.Dequeue()
	if !ok {
		t.Fatalf("expected ok=true from dequeue")
	}
	if got.ID != job.ID {
		t.Fatalf("got %s, want %s", got.ID, job.ID)
	}
}

func TestMemoryQueue_Full(t *testing.T) {
	q := Adapters.NewMemoryQueue(1)
	_ = q.Enqueue(Structures.Job{ID: "1"})
	if err := q.Enqueue(Structures.Job{ID: "2"}); err == nil {
		t.Fatalf("expected error when queue is full")
	}
}

func TestMemoryQueue_Close(t *testing.T) {
	q := Adapters.NewMemoryQueue(1)
	q.Close()
	// After close, a receive returns zero value with ok=false
	_, ok := q.Dequeue()
	if ok {
		t.Fatalf("expected ok=false after close")
	}
}
