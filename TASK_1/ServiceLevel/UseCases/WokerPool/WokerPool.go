package WokerPool

import (
	"Kaspersky_Go/ModeLevel/Structures"
	"Kaspersky_Go/ServiceLevel/Interfaces/IAdapters"
	"Kaspersky_Go/ServiceLevel/Interfaces/IWorkerPool"
	"context"
	"log/slog"
)

type WorkerPool struct {
	workers   int
	queue     IAdapters.Queue
	state     IAdapters.StateStore
	processor IWorkerPool.JobProcessor
	logger    *slog.Logger
}

func NewWorkerPool(workers int, queue IAdapters.Queue, state IAdapters.StateStore,
	processor IWorkerPool.JobProcessor, log *slog.Logger) *WorkerPool {

	return &WorkerPool{
		workers:   workers,
		queue:     queue,
		state:     state,
		processor: processor,
		logger:    log,
	}
}

func (wp *WorkerPool) Start(ctx context.Context) {
	for i := 0; i < wp.workers; i++ {
		go func(workerId int) {
			for {
				select {
				case <-ctx.Done():
					return
				default:
					job, ok := wp.queue.Dequeue()
					if !ok {
						return
					}
					wp.state.Set(job.ID, Structures.JobStatus{
						State:    Structures.StateRunning,
						Attempts: job.Attempts,
					})

					status := wp.processor.Process(ctx, job, wp.logger)

					wp.state.Set(job.ID, status)

					if status.State == Structures.StateQueued {
						_ = wp.queue.Enqueue(job)
					}
				}

			}
		}(i)
	}
}
