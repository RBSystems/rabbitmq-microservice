package rabbitmq

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

func failOnError(err error, failed string, success string) {
	if err != nil {
		log.Fatalf("%s: %s", failed, err)
		panic(fmt.Sprintf("%s: %s", failed, err))
	} else if len(success) > 0 {
		log.Printf("%s", success)
	}
}

func publish(message string, channel string) {
	log.Printf("Publishing message: %s, on channel: %s", message, channel)
	conn, err := amqp.Dial("amqp://guest@localhost:5672")
	failOnError(err, "\tfailed to connect to RabbitMQ", "\tConnected to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "\tfailed to open an channel", "\tSuccessfully opened a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		channel, // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "\tfailed to declare a queue", "\tSuccessfully declared a queue")

	body := message
	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	failOnError(err, "\tfailed to publish a message", "\tSuccessfully published a message")
}

func recieve(channel string) {
	conn, err := amqp.Dial("amqp://guest@localhost:5672")
	failOnError(err, "\tfailed to connect to RabbitMQ", "\tConnected to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "\tfailed to open an channel", "\tSuccessfully opened the channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		channel, // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "\tfailed to declare a queue", "\tSuccessfully declared a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "\tFailed to register a consumer", "\tSuccessfully registered a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("\tRecieved a message: %s", d.Body)
		}
	}()

	log.Printf("[*] Waiting for messages (channel: %s). To exit press CTRL+C", channel)
	<-forever
}
