package main

import (
	"log"
	"notification-service/internal/db"
	"notification-service/internal/handlers"
	"notification-service/internal/rabbitmq"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize MongoDB
	if err := db.ConnectMongoDB(); err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Initialize RabbitMQ
	rabbitmq.InitRabbitMQ()

	// Start the Gin server
	router := gin.Default()

	// Routes
	router.POST("/notifications", handlers.CreateNotificationHandler)
	router.GET("/notifications/:userID", handlers.GetNotificationsHandler)
	router.PUT("/notifications/:userID/mark-read", handlers.MarkAsReadHandler)

	log.Println("Notification Service running on :8081")
	if err := router.Run(":8081"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
