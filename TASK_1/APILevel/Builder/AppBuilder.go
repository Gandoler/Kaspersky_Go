package Builder

import (
	"Kaspersky_Go/APILevel/Adapters"
	"Kaspersky_Go/APILevel/HTTPServer"
	"Kaspersky_Go/ServiceLevel/UseCases/Processors"
	"Kaspersky_Go/ServiceLevel/UseCases/WokerPool"
)

type AppBuilder struct {
	queueSize   int
	workerCount int
	serverAddr  string
}

func NewAppBuilder() *AppBuilder {
	return &AppBuilder{
		queueSize:   64,
		workerCount: 4,
		serverAddr:  ":8080",
	}
}
func (b *AppBuilder) WithQueueSize(size int) *AppBuilder {
	b.queueSize = size
	return b
}

func (b *AppBuilder) WithWorkerCount(count int) *AppBuilder {
	b.workerCount = count
	return b
}

func (b *AppBuilder) WithServerAddr(addr string) *AppBuilder {
	b.serverAddr = addr
	return b
}

func (b *AppBuilder) Build() *App {
	queue := Adapters.NewMemoryQueue(b.queueSize)
	state := Adapters.NewMemoryStateStore()
	processor := Processors.NewRetryProcessor()
	pool := WokerPool.NewWorkerPool(b.workerCount, queue, state, processor)
	server := HTTPServer.NewHTTPServer(b.serverAddr, queue, state)

	return &App{
		queue:  queue,
		state:  state,
		pool:   pool,
		server: server,
	}
}
