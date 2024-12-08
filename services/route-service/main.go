package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"route-service/internal/db"
	"route-service/internal/handlers"
)

func main() {
	// Initialize DB
	err := db.ConnectPostgres()
	if err != nil {
		log.Fatalf("Failed to connect to Postgres: %v", err)
	}

	router := gin.Default()

	// Routes
	router.GET("/routes/:id", handlers.GetRoute)
	router.POST("/routes", handlers.CreateRoute)

	log.Println("Route Service is running on :8080")
	http.ListenAndServe(":8080", router)
}
