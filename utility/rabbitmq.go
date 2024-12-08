package utility

import (
	"context"
	"encoding/json"
	"log"
	"sync"
)

import (
	"github.com/streadway/amqp"
	"notification-service/internal/db"
	"notification-service/internal/models"
)

const (
	workerPoolSize = 10
)

func StartWorkerPool(msgs <-chan amqp.Delivery) {
	var wg sync.WaitGroup
	workerQueue := make(chan models.Notification, workerPoolSize)

	// Start worker Goroutines
	for i := 0; i < workerPoolSize; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for notification := range workerQueue {
				processNotification(notification)
			}
			log.Printf("Worker %d stopped", id)
		}(i)
	}

	// Push messages to the worker queue
	go func() {
		for msg := range msgs {
			var notification models.Notification
			if err := json.Unmarshal(msg.Body, &notification); err != nil {
				log.Printf("Failed to unmarshal message: %v", err)
				continue
			}
			workerQueue <- notification
		}
		close(workerQueue)
	}()

	wg.Wait()
}

func processNotification(notification models.Notification) {
	// Process the notification and store it in MongoDB
	_, err := db.NotificationCollection.InsertOne(context.Background(), notification)
	if err != nil {
		log.Printf("Failed to store notification: %v", err)
	}
}
