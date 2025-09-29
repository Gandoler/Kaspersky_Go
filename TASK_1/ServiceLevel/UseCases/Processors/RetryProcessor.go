package Processors

import (
	"Kaspersky_Go/ModeLevel/Structures"
	"context"
	"fmt"
	"math/rand"
	"time"
)

type RetryProcessor struct{}

func NewRetryProcessor() *RetryProcessor {
	return &RetryProcessor{}
}

func (p *RetryProcessor) Process(ctx context.Context, job Structures.Job) Structures.JobStatus {

	time.Sleep(time.Duration(100+rand.Intn(50)) * time.Millisecond)
	fmt.Printf("process running  ID%d:\tpaylod:%s", job.ID, job.Payload)
	if rand.Intn(100) < 20 {
		job.Attempts++
		if job.Attempts <= job.MaxRetries {

			delay := backoff(job.Attempts)
			time.Sleep(delay)
			return Structures.JobStatus{
				State:    Structures.StateQueued,
				Attempts: job.Attempts,
				Error:    "retrying",
			}
		}
	}

	return Structures.JobStatus{State: Structures.StateDone, Attempts: job.Attempts}
}

func backoff(attempt int) time.Duration {
	base := time.Duration(100*(1<<uint(attempt-1))) * time.Millisecond
	jitter := time.Duration(rand.Intn(101)) * time.Millisecond
	return base + jitter
}
