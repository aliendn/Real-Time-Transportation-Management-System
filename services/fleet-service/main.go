package fleet_service

import (
	"log"
)

import (
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize database
	if err := db.ConnectPostgres(); err != nil {
		log.Fatalf("Failed to connect to Postgres: %v", err)
	}

	// Initialize RabbitMQ
	rabbitmq.InitRabbitMQ()

	// Start Gin server
	router := gin.Default()
	router.GET("/fleet/:id", handlers.GetFleetHandler)
	router.POST("/fleet", handlers.CreateFleetHandler)
	router.POST("/assign", handlers.AssignFleetHandler)

	log.Println("Fleet Service running on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
