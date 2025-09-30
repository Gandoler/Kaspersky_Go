package App

import (
	"errors"
	"log/slog"
	"Kaspersky_Go/TASK_3/WorkerPool"
)

type Config struct {
	Workers   int
	QueueSize int
	afterHook func()
}

func NewPool(config Config, log *slog.Logger) (WorkerPool.Pool, error) {
	if config.Workers < 1 {
		return nil, errors.New("workers must be greater than zero")
	}
	if config.QueueSize < 1 {
		return nil, errors.New("queue_size must be greater than zero")
	}

	p := &WorkerPool.workerPool{
		queue:     make(chan func(), config.QueueSize),
		afterHook: config.afterHook,
	}

	p.wgWorkers.Add(config.Workers)
	for i := 0; i < config.Workers; i++ {
		go p.worker()
	}
}
