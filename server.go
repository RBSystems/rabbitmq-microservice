package main

import (
	"log"
	"os"
	"strings"
	"sync"

	"github.com/byuoitav/av-api/dbo"
	"github.com/byuoitav/rabbitmq-microservice/rabbitmq"
)

func main() {
	var wg sync.WaitGroup

	hostname := os.Getenv("PI_HOSTNAME")
	values := strings.Split(strings.TrimSpace(hostname), "-")
	devices, err := dbo.GetDevicesByBuildingAndRoomAndRole(values[0], values[1], "EventRouter")
	if err != nil {
		log.Fatal(err.Error())
	}

	addresses := []string{}
	for _, device := range devices {
		addresses = append(addresses, device.Address)
	}

	wg.Add(2)
	go rabbitmq.Recieve(rabbitmq.LocalAPI)
	go rabbitmq.Recieve(rabbitmq.External)
	go rabbitmq.Recieve(rabbitmq.TransmitAPI)
	wg.Wait()
}
