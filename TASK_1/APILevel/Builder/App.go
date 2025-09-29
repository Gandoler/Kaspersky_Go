package Tempates

type App struct {
	queue  *infrastructure.MemoryQueue
	state  *infrastructure.MemoryStateStore
	pool   *usecase.WorkerPool
	server *infrastructure.HTTPServer
}
