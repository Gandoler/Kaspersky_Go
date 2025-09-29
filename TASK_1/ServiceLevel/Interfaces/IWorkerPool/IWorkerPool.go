package IWorkerPool

import (
	"Kaspersky_Go/ModeLevel/Structures"
	"context"
)

type Queue interface {
	Enqueue(job Structures.Job) error
	Dequeue() (Structures.Job, bool)
	Close()
}

type StateStore interface {
	Set(id string, st Structures.JobStatus)
	Get(id string) (Structures.JobStatus, bool)
}

type JobProcessor interface {
	Process(ctx context.Context, job Structures.Job) Structures.JobStatus
}
