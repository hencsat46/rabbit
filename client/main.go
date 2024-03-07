package main

import (
	"context"
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

	if err = ch.ExchangeDeclare("logs", "fanout", false, false, false, false, nil); err != nil {
		log.Fatalln("Cannot create exchange", err)
	}

	// _, err = ch.QueueDeclare(
	// 	"first",
	// 	false,
	// 	false,
	// 	false,
	// 	false,
	// 	nil,
	// )

	// if err != nil {
	// 	log.Fatalln("Cannot create queue", err)
	// }

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	body := "Дарова"

	err = ch.PublishWithContext(ctx,
		"logs",
		"",
		false,
		false,
		rabbit.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})

	if err != nil {
		log.Fatalln("cannot send message", err)
	}

}
