package Processors_test

import (
	"Kaspersky_Go/ModeLevel/Structures"
	"Kaspersky_Go/ServiceLevel/UseCases/Processors"
	"context"
	"log/slog"
	"math/rand"
	"os"
	"testing"
	"time"
)

func TestRetryProcessor_Success(t *testing.T) {
	// Seed to reduce flakiness: make random failure unlikely
	rand.Seed(1)
	p := Processors.NewRetryProcessor()
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	st := p.Process(ctx, Structures.Job{ID: "ok", MaxRetries: 0}, logger)
	if st.State != Structures.StateDone {
		t.Fatalf("got %v, want %v", st.State, Structures.StateDone)
	}
}

func TestRetryProcessor_RetryThenQueue(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	p := Processors.NewRetryProcessor()
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// With MaxRetries>0 there's a chance to return queued; run multiple times
	queuedSeen := false
	for i := 0; i < 20 && !queuedSeen; i++ {
		st := p.Process(ctx, Structures.Job{ID: "id", MaxRetries: 1}, logger)
		if st.State == Structures.StateQueued {
			queuedSeen = true
		}
	}
	if !queuedSeen {
		t.Fatalf("expected at least one queued status across attempts")
	}
}
