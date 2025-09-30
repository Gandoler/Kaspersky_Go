package WorkerPool

import (
	"errors"
	"log/slog"
	"sync"
)

var (
	ErrStopped   = errors.New("pool stopped")
	ErrQueueFull = errors.New("queue full")
)

type WorkerPool struct {
	queue     chan func()
	stopOnce  sync.Once
	stopped   bool
	stateMu   sync.Mutex
	wgTasks   sync.WaitGroup
	wgWorkers sync.WaitGroup
	afterHook func()
	logger    *slog.Logger
}

func (p *WorkerPool) Submit(task func()) error {
	if task == nil {
		return nil
	}

	p.stateMu.Lock()
	if p.stopped {
		p.stateMu.Unlock()
		return ErrStopped
	}

	if len(p.queue) == cap(p.queue) {
		p.stateMu.Unlock()
		return ErrQueueFull
	}
	p.wgTasks.Add(1)
	p.stateMu.Unlock()

	p.queue <- task
	return nil
}

func (p *WorkerPool) Stop() error {
	p.stopOnce.Do(func() {
		p.stateMu.Lock()
		p.stopped = true
		p.stateMu.Unlock()

		close(p.queue)
		p.wgWorkers.Wait()
	})

	return nil
}

func (p *WorkerPool) worker() {
	defer p.wgWorkers.Done()
	for task := range p.queue {
		p.Start(task)
		if p.afterHook != nil {
			func() {
				defer func() { _ = recover() }()
				p.afterHook()
			}()
		}
		p.wgTasks.Done()
	}
}

func (p *WorkerPool) Start(task func()) {
	defer func() { _ = recover() }()
	task()
}

func New(config Config, logger *slog.Logger) (*WorkerPool, error) {
	if config.Workers < 1 {
		return nil, errors.New("workers must be greater than zero")
	}
	if config.QueueSize < 1 {
		return nil, errors.New("queue size must be greater than zero")
	}

	pool := &WorkerPool{
		queue:     make(chan func(), config.QueueSize),
		afterHook: config.AfterTask,
		logger:    logger,
	}

	pool.wgWorkers.Add(config.Workers)
	for i := 0; i < config.Workers; i++ {
		go pool.worker()
	}
	return pool, nil
}
