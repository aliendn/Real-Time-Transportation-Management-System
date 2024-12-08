package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"route-service/internal/models"
)

func GetRoute(c *gin.Context) {
	id := c.Param("id")
	route, err := models.FetchRouteByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Route not found"})
		return
	}

	c.JSON(http.StatusOK, route)
}

func CreateRoute(c *gin.Context) {
	var route models.Route
	if err := c.ShouldBindJSON(&route); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := models.SaveRoute(&route)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save route"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Route created successfully"})
}
