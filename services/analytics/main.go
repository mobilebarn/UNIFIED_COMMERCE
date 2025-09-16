package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize Gin router
	router := gin.Default()

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		health := map[string]interface{}{
			"service": "analytics",
			"status":  "healthy",
		}
		c.JSON(http.StatusOK, health)
	})

	// Start server
	fmt.Println("Starting Analytics Service on port 8080")
	log.Fatal(router.Run(":8080"))
}
