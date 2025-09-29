package Builder

import (
	"Kaspersky_Go/APILevel/Adapters"
	"Kaspersky_Go/APILevel/HTTPServer"
	"Kaspersky_Go/ServiceLevel/UseCases/Processors"
	"Kaspersky_Go/ServiceLevel/UseCases/WokerPool"
	"log/slog"
	"os"
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
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})

	logger := slog.New(handler)
	queue := Adapters.NewMemoryQueue(b.queueSize)
	state := Adapters.NewMemoryStateStore()
	processor := Processors.NewRetryProcessor()
	pool := WokerPool.NewWorkerPool(b.workerCount, queue, state, processor, logger)
	server := HTTPServer.NewHTTPServer(b.serverAddr, queue, state)

	return &App{
		queue:  queue,
		state:  state,
		pool:   pool,
		server: server,
		logger: logger,
	}
}
