package Processors

import (
	"Kaspersky_Go/ModeLevel/Structures"
	"context"
	"log/slog"
	"math/rand"
	"time"
)

type RetryProcessor struct{}

func NewRetryProcessor() *RetryProcessor {
	return &RetryProcessor{}
}

func (p *RetryProcessor) Process(ctx context.Context, job Structures.Job, log *slog.Logger) Structures.JobStatus {
	// Логируем начало обработки
	log.Info("Starting job processing",
		"jobID", job.ID,
		"payload", job.Payload,
		"attempts", job.Attempts,
		"maxRetries", job.MaxRetries,
	)

	// Симуляция работы
	time.Sleep(time.Duration(100+rand.Intn(50)) * time.Millisecond)

	// Проверка на случайную ошибку (20% шанс)
	if rand.Intn(100) < 20 {
		job.Attempts++
		log.Warn("Job failed, will check retries",
			"jobID", job.ID,
			"attempts", job.Attempts,
		)

		if job.Attempts <= job.MaxRetries {
			delay := backoff(job.Attempts)
			log.Info("Retrying job with backoff",
				"jobID", job.ID,
				"nextAttemptDelay", delay,
				"attempts", job.Attempts,
			)
			time.Sleep(delay)
			return Structures.JobStatus{
				State:    Structures.StateQueued,
				Attempts: job.Attempts,
				Error:    "retrying",
			}
		} else {
			log.Error("Max retries reached, job failed permanently",
				"jobID", job.ID,
				"attempts", job.Attempts,
			)
		}
	}

	// Логируем успешное завершение
	log.Info("Job processed successfully",
		"jobID", job.ID,
		"attempts", job.Attempts,
	)
	return Structures.JobStatus{
		State:    Structures.StateDone,
		Attempts: job.Attempts,
	}
}

func backoff(attempt int) time.Duration {
	base := time.Duration(100*(1<<uint(attempt-1))) * time.Millisecond
	jitter := time.Duration(rand.Intn(101)) * time.Millisecond
	return base + jitter
}
