package main

import (
	"log"
	"time"

	rabbit "github.com/rabbitmq/amqp091-go"
)

func main() {
	conn, err := rabbit.Dial("amqp://guest:guest@localhost:5672")
	if err != nil {
		log.Fatalln("Cannot connect to rabbitmq. ", err)
	}

	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalln("cannot open channel", err)
	}

	defer ch.Close()

	_, err = ch.QueueDeclare(
		"first",
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		log.Fatalln("Cannot create queue", err)
	}

	messages, err := ch.Consume(
		"first",
		"num1",
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		log.Fatalln("Cannot create consumer", err)
	}

	done := make(chan struct{})

	go func() {
		for message := range messages {
			log.Println(string(message.Body))
			time.Sleep(10 * time.Second)
			message.Ack(false)
		}
	}()

	<-done
}
