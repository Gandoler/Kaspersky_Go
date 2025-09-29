package main

import (
	"Kaspersky_Go/APILevel/Builder"
	"fmt"
	"os"
	"strconv"
)

func main() {
	QueueSizeEnv := os.Getenv("QUEUE_SIZE")
	WorkersEnv := os.Getenv("WORKERS")

	size, err := EnvTOInt(QueueSizeEnv, "QUEUE_SIZE")
	if err != nil {
		panic(err)
	}
	poolSize, err := EnvTOInt(WorkersEnv, "WORKERS")
	if err != nil {
		panic(err)
	}

	app := Builder.NewAppBuilder().
		WithQueueSize(size).
		WithWorkerCount(poolSize).
		WithServerAddr(":8080").
		Build()

	app.Start()

}

func EnvTOInt(str string, typeOf string) (int, error) {

	tmp := 0
	if str == "" {
		var err error
		tmp, err = strconv.Atoi(str)
		if err != nil {
			return 0, fmt.Errorf("ошибка при конвертации %s", typeOf)
		}
	}
	return tmp, nil
}
