package main

import (
	"Kaspersky_Go/APILevel/Adapters"
	"fmt"
	"os"
	"strconv"
)

func main() {
	QUEUE_SIZE_env := os.Getenv("QUEUE_SIZE")
	WORKERS_env     := os.Getenv("WORKERS")

	size, err := EnvTOInt(QUEUE_SIZE_env, "QUEUE_SIZE")
	if err != nil {
		panic(err)
	}
	worker, err := EnvTOInt(WORKERS_env, "WORKERS")
	if err != nil {
		panic(err)
	}


	queue := Adapters.NewMemoryQueue(size)
	state := Adapters.NewMemoryStateStore()
	pool  :=

}


func EnvTOInt(str string, typeOf string) (int ,error) {

	tmp := 0
	if str != "" {
		var err error
		tmp, err = strconv.Atoi(str)
		if err != nil {
			return 0, fmt.Errorf("ошибка при конвертации %s", typeOf)
		}
	}
	return tmp, nil
}