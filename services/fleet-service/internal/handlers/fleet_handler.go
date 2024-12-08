package handlers

import (
	"context"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"services/fleet-service/internal/models"
	"strconv"
	"sync"
	"time"
)

func GetFleetHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	vehicle, err := models.FetchVehicleByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Vehicle not found"})
		return
	}

	c.JSON(http.StatusOK, vehicle)
}

func CreateFleetHandler(c *gin.Context) {
	var vehicle models.Vehicle
	if err := c.ShouldBindJSON(&vehicle); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := models.SaveVehicle(&vehicle)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save vehicle"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Vehicle created successfully"})
}

func AssignFleetHandler(c *gin.Context) {
	// Sample tasks
	tasks := []string{"Task1", "Task2", "Task3"}
	var wg sync.WaitGroup

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for _, task := range tasks {
		wg.Add(1)
		go func(task string) {
			defer wg.Done()
			// Simulate task assignment
			select {
			case <-time.After(1 * time.Second):
				log.Printf("Assigned %s to a vehicle", task)
			case <-ctx.Done():
				log.Printf("Context deadline exceeded for %s", task)
			}
		}(task)
	}

	wg.Wait()
	c.JSON(http.StatusOK, gin.H{"message": "All tasks assigned"})
}
