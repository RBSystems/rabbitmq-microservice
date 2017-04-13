package main

import (
	"os"
	"sync"

	"github.com/byuoitav/rabbitmq-microservice/rabbitmq"
)

func main() {
	var wg sync.WaitGroup

	hostname := os.Getenv("PI_HOSTNAME")

	wg.Add(1)
	rabbitmq.Publish("hey", "test")

	go rabbitmq.Recieve("test")
	wg.Wait()
}
