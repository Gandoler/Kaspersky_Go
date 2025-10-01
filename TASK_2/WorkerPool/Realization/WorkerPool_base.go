package Realization

type WorkerPool struct {
}

func NewWorkerPool(numberOfWorkers int) *WorkerPool {
	return &WorkerPool{}
}

func (wp *WorkerPool) Submit(task func()) {

}

func (wp *WorkerPool) SubmitWait(task func()) {

}

func (wp *WorkerPool) Stop() {

}

func (wp *WorkerPool) StopWait() {

}
