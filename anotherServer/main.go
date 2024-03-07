package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	rabbit "github.com/rabbitmq/amqp091-go"
)

func main() {

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGTERM)

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

	if err = ch.ExchangeDeclare("logs", "fanout", false, false, false, false, nil); err != nil {
		log.Fatalln("Cannot create exchange", err)
	}

	_, err = ch.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil,
	)

	if err = ch.QueueBind("", "", "logs", false, nil); err != nil {
		log.Fatalln("Cannot bind queue", err)
	}

	if err != nil {
		log.Fatalln("Cannot create queue", err)
	}

	messages, err := ch.Consume(
		"",
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		log.Fatalln("Cannot create consumer", err)
	}

	go func() {
		for message := range messages {
			log.Println(string(message.Body))
		}
	}()

	<-done

}
