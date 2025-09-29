package IWorkerPool

import (
	"Kaspersky_Go/ModeLevel/Structures"
	"context"
)

type JobProcessor interface {
	Process(ctx context.Context, job Structures.Job) Structures.JobStatus
}
