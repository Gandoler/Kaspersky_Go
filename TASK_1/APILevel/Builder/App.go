package Builder

import (
	"Kaspersky_Go/APILevel/HTTPServer"
	"Kaspersky_Go/ServiceLevel/Interfaces/IAdapters"
	"Kaspersky_Go/ServiceLevel/Interfaces/IWorkerPool"
	"Kaspersky_Go/ServiceLevel/UseCases/WokerPool"
	"context"
	"fmt"
	"log"
	"os/signal"
	"syscall"
	"time"
)

type App struct {
	queue     IAdapters.Queue
	state     IAdapters.StateStore
	pool      *WokerPool.WorkerPool
	server    *HTTPServer.HTTPServer
	processor IWorkerPool.JobProcessor
}

func (a *App) Start() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go a.pool.Start(ctx)
	fmt.Println("Worker pool started")

	errCh := make(chan error, 1)
	go func() {
		fmt.Println("Starting HTTP server on port 8080...")
		if err := a.server.Start(); err != nil {
			errCh <- err
		}
	}()

	fmt.Println("HTTP server goroutine launched")

	select {
	case <-ctx.Done():
		fmt.Println("shutdown signal received")
	case err := <-errCh:
		log.Fatalf("server error: %v", err)
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := a.server.Stop(shutdownCtx); err != nil {
		fmt.Printf("server shutdown error: %v", err)
	}

	a.queue.Close()
	fmt.Println("graceful shutdown complete")
}
