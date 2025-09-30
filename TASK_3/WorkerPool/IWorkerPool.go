package WorkerPool

type Pool interface {
	Submit(task func()) error
	Stop() error
}

type Config struct {
	Workers   int
	QueueSize int
	AfterTask func()
}
