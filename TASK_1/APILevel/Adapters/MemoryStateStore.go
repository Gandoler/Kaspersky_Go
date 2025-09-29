package Adapters

import (
	"Kaspersky_Go/ModeLevel/Structures"
	"sync"
)

type MemoryStateStore struct {
	mu    sync.RWMutex
	state map[string]Structures.JobStatus
}

func NewMemoryStateStore() *MemoryStateStore {
	return &MemoryStateStore{state: make(map[string]Structures.JobStatus)}

}

func (ms *MemoryStateStore) Set(id string, st Structures.JobStatus) {
	ms.mu.Lock()
	ms.state[id] = st
	ms.mu.Unlock()
}

func (ms *MemoryStateStore) Get(id string) (Structures.JobStatus, bool) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	st, ok := ms.state[id]
	return st, ok
}
