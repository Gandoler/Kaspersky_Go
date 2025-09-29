package IWorkerPool

import (
	"Kaspersky_Go/ModeLevel/Structures"
	"context"
	"log/slog"
)

type JobProcessor interface {
	Process(ctx context.Context, job Structures.Job, log *slog.Logger) Structures.JobStatus
}
