package Builder

import (
	"Kaspersky_Go/APILevel/HTTPServer"
	"Kaspersky_Go/ServiceLevel/Interfaces/IAdapters"
	"Kaspersky_Go/ServiceLevel/Interfaces/IWorkerPool"
	"Kaspersky_Go/ServiceLevel/UseCases/WokerPool"
	"context"
	"log/slog"
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
	logger    *slog.Logger
}

func (a *App) Start() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go a.pool.Start(ctx)
	a.logger.Info("Worker pool started")

	errCh := make(chan error, 1)
	go func() {
		a.logger.Info("Starting HTTP server on port 8080...")
		if err := a.server.Start(); err != nil {
			errCh <- err
		}
	}()

	a.logger.Info("HTTP server goroutine launched")

	select {
	case <-ctx.Done():
		a.logger.Info("shutdown signal received")
	case err := <-errCh:
		a.logger.Error("server error", "err", err)
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := a.server.Stop(shutdownCtx); err != nil {
		a.logger.Error("server shutdown error", "err", err)
	}

	a.queue.Close()
	a.logger.Info("graceful shutdown complete")
}
