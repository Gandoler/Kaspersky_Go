package Adapters

import (
	"Kaspersky_Go/ModeLevel/Structures"
	"errors"
)

type MemoryQueue struct {
	ch chan Structures.Job
}

func NewMemoryQueue(size int) *MemoryQueue {
	return &MemoryQueue{ch: make(chan Structures.Job, size)}
}

func (mq *MemoryQueue) Enqueue(job Structures.Job) error {
	select {
	case mq.ch <- job:
		return nil
	default:
		return errors.New("Queue is full")
	}
}

func (mq *MemoryQueue) Dequeue() (Structures.Job, bool) {
	job, ok := <-mq.ch
	return job, ok
}

func (mq *MemoryQueue) Close() {
	close(mq.ch)
}
