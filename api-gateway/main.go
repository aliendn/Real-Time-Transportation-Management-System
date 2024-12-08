package api_gateway

import "net/http"

import (
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// Route to route-service
	router.GET("/routes/:id", proxyToService("http://route-service:8080"))
	router.POST("/routes", proxyToService("http://route-service:8080"))

	// Route to fleet-service
	router.GET("/fleet/:id", proxyToService("http://fleet-service:8080"))
	router.POST("/fleet", proxyToService("http://fleet-service:8080"))

	// Route to notification-service
	router.GET("/notifications", proxyToService("http://notification-service:8080"))
	router.POST("/notifications", proxyToService("http://notification-service:8080"))

	router.Run(":8080")
}

func proxyToService(target string) gin.HandlerFunc {
	return func(c *gin.Context) {
		targetURL := target + c.Request.URL.Path
		http.Redirect(c.Writer, c.Request, targetURL, http.StatusTemporaryRedirect)
	}
}
