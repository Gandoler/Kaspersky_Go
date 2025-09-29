package WokerPool

import (
	"Kaspersky_Go/ModeLevel/Structures"
	"Kaspersky_Go/ServiceLevel/Interfaces/IWorkerPool"
	"context"
)

type WorkerPool struct {
	workers   int
	queue     IWorkerPool.Queue
	state     IWorkerPool.StateStore
	processor IWorkerPool.JobProcessor
}

func NewWorkerPool(workers int, queue IWorkerPool.Queue, state IWorkerPool.StateStore,
	processor IWorkerPool.JobProcessor) *WorkerPool {

	return &WorkerPool{
		workers:   workers,
		queue:     queue,
		state:     state,
		processor: processor,
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

					status := wp.processor.Process(ctx, job)
					wp.state.Set(job.ID, status)

					if status.State == Structures.StateQueued {
						_ = wp.queue.Enqueue(job)
					}
				}

			}
		}(i)
	}
}
