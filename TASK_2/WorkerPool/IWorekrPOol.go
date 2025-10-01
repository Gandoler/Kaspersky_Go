package WorkerPool

type Pool interface {
	Submit(task func())
	SubmitWait(task func())
	Stop()
	StopWait()
}
