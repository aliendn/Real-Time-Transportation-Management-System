package utility

import (
	"log"
	"time"
)

func processNotificationWithRetry(notification models.Notification, retries int) {
	for i := 0; i < retries; i++ {
		err := processNotification(notification)
		if err == nil {
			return
		}
		log.Printf("Retry %d: Failed to process notification: %v", i+1, err)
		time.Sleep(time.Duration(i+1) * time.Second) // Exponential backoff
	}
	log.Printf("Failed to process notification after %d retries", retries)
}
