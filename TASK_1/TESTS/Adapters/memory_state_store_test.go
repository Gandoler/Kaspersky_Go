package Adapters_test

import (
	"Kaspersky_Go/APILevel/Adapters"
	"Kaspersky_Go/ModeLevel/Structures"
	"sync"
	"testing"
)

func TestMemoryStateStore_SetGet(t *testing.T) {
	st := Adapters.NewMemoryStateStore()
	want := Structures.JobStatus{State: Structures.StateQueued, Attempts: 1}
	st.Set("job1", want)

	got, ok := st.Get("job1")
	if !ok {
		t.Fatalf("expected ok=true")
	}
	if got != want {
		t.Fatalf("got %+v, want %+v", got, want)
	}
}

func TestMemoryStateStore_ConcurrentAccess(t *testing.T) {
	st := Adapters.NewMemoryStateStore()
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			st.Set("k", Structures.JobStatus{State: Structures.StateRunning, Attempts: i})
			_, _ = st.Get("k")
		}(i)
	}
	wg.Wait()
}
