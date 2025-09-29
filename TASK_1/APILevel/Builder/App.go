package Builder

import (
	"Kaspersky_Go/APILevel/Adapters"
	"Kaspersky_Go/APILevel/HTTPServer"
	"Kaspersky_Go/ServiceLevel/UseCases/WokerPool"
	"context"
	"log"
	"os/signal"
	"syscall"
	"time"
)

type App struct {
	queue  *Adapters.MemoryQueue
	state  *Adapters.MemoryStateStore
	pool   *WokerPool.WorkerPool
	server *HTTPServer.HTTPServer
}

func (a *App) Start() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go a.pool.Start(ctx)

	errCh := make(chan error, 1)
	go func() {
		if err := a.server.Start(); err != nil {
			errCh <- err
		}
	}()

	// Ожидание сигналов или ошибок сервера
	select {
	case <-ctx.Done():
		log.Println("shutdown signal received")
	case err := <-errCh:
		log.Fatalf("server error: %v", err)
	}

	// Graceful shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := a.server.Stop(shutdownCtx); err != nil {
		log.Printf("server shutdown error: %v", err)
	}

	a.queue.Close()
	log.Println("graceful shutdown complete")
}
