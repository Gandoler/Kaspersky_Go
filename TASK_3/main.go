package main

import (
	"TASK_3/WorkerPool"
	"fmt"
	"log/slog"
	"math/rand"
	"os"
	"time"
)

func main() {

	rand.Seed(time.Now().UnixNano())

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	p, err := WorkerPool.New(WorkerPool.Config{
		Workers:   3,
		QueueSize: 2,
		AfterTask: func() { fmt.Println("after task hook") },
	}, logger)
	if err != nil {
		panic(err)
	}

	for i := 0; i < 10; i++ {
		i := i
		err := p.Submit(func() {
			workMs := 200 + rand.Intn(400)
			time.Sleep(time.Duration(workMs) * time.Millisecond)
			fmt.Printf("task %d done in %dms\n", i, workMs)
		})
		if err != nil {
			fmt.Printf("submit task %d error: %v\n", i, err)
		}
	}

	if err := p.Stop(); err != nil {
		panic(err)
	}

	fmt.Println("pool stopped")

}
