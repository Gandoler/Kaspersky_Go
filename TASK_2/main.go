package main

import (
	"TASK_2/WorkerPool/Realization"
	"fmt"
	"time"
)

func main() {
	wp := Realization.NewWorkerPool(3)

	for i := 1; i <= 5; i++ {
		num := i // for thread safety
		wp.Submit(func() {
			fmt.Printf("Start task %d\n", num)
			time.Sleep(1 * time.Second)
			fmt.Printf("Done task %d\n", num)
		})
	}
	fmt.Println("All tasks submitted\n\n")

	wp.SubmitWait(func() {
		fmt.Println("SubmitWait task started")
		time.Sleep(2 * time.Second)
		fmt.Println("SubmitWait task finished")
	})

	fmt.Println("All tasks submitted\n\n")

	wp.StopWait()
	fmt.Println("WorkerPool stopped")
}
