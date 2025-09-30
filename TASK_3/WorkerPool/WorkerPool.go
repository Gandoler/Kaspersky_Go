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

		p.wgWorkers.Wait()

		close(p.queue)
	})

	p.wgWorkers.Wait()
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
