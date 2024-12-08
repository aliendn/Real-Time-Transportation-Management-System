package rabbitmq

import (
	"encoding/json"
	"log"
	"notification-service/internal/models"

	"github.com/streadway/amqp"
)

var conn *amqp.Connection
var channel *amqp.Channel

func InitRabbitMQ() {
	var err error
	conn, err = amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}

	channel, err = conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}

	log.Println("RabbitMQ initialized")
}

func ListenForNotifications() {
	q, err := channel.QueueDeclare(
		"notification_queue",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}

	msgs, err := channel.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	go func() {
		for msg := range msgs {
			var notification models.Notification
			if err := json.Unmarshal(msg.Body, &notification); err != nil {
				log.Printf("Failed to decode notification: %v", err)
				continue
			}
			log.Printf("Received notification: %v", notification)
		}
	}()
}
