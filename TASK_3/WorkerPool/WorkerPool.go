package WorkerPool

import (
	"log/slog"
	"sync"
)

type workerPool struct {
	queue     chan func()
	stopOnce  sync.Once
	stopped   bool
	stateMu   sync.Mutex
	wgTasks   sync.WaitGroup
	wgWorkers sync.WaitGroup
	afterHook func()
	logger    *slog.Logger
}
