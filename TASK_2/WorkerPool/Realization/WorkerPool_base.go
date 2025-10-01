package Realization

import (
	"sync"
)

type WorkerPool struct {
	queue     chan func()
	stopOnce  sync.Once
	stopped   bool
	stateMu   sync.Mutex
	wgTasks   sync.WaitGroup
	wgWorkers sync.WaitGroup
}

func NewWorkerPool(numberOfWorkers int) *WorkerPool {
	if numberOfWorkers < 1 {
		return nil
	}

	pool := &WorkerPool{
		queue: make(chan func(), 64), // default value
	}

	pool.wgWorkers.Add(numberOfWorkers)
	for i := 0; i < numberOfWorkers; i++ {
		go func(id int) {
			defer pool.wgWorkers.Done()
			for task := range pool.queue {
				task()
				pool.wgTasks.Done()
			}
		}(i)
	}
	return pool
}

func (wp *WorkerPool) Submit(task func()) {
	if task == nil {
		return
	}

	wp.stateMu.Lock()
	if wp.stopped {
		wp.stateMu.Unlock()
		return
	}

	if len(wp.queue) == cap(wp.queue) {
		wp.stateMu.Unlock()
		return
	}
	wp.wgTasks.Add(1)
	wp.stateMu.Unlock()

	wp.queue <- task
	return

}

func (wp *WorkerPool) SubmitWait(task func()) {
	if task == nil {
		return
	}

	wp.stateMu.Lock()
	if wp.stopped {
		wp.stateMu.Unlock()
		return
	}

	if len(wp.queue) == cap(wp.queue) {
		wp.stateMu.Unlock()
		return
	}
	wp.wgTasks.Add(1)
	wp.stateMu.Unlock()

	var wg sync.WaitGroup
	wg.Add(1)

	wp.queue <- func() {
		defer wg.Done()
		task()
	}

	wg.Wait()
	return
}

func (wp *WorkerPool) Stop() {
	wp.stopOnce.Do(func() {
		wp.stateMu.Lock()
		wp.stopped = true
		wp.stateMu.Unlock()

		close(wp.queue)
		wp.wgWorkers.Wait()
	})

	return
}

func (wp *WorkerPool) StopWait() {
	wp.Stop()
	wp.wgTasks.Wait()
	wp.wgWorkers.Wait()
}
