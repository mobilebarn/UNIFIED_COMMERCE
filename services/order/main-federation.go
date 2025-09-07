package main

import (
	"net/http"
	"time"

	"unified-commerce/services/order/graphql"
	"unified-commerce/services/order/service"
	"unified-commerce/services/shared/logger"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize logger
	loggerConfig := logger.DefaultConfig("order")
	log := logger.NewLogger(loggerConfig)

	// Create a mock order service for federation testing
	orderService := &service.OrderService{}

	// Initialize Gin router
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	// Add middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		health := map[string]interface{}{
			"service": "order",
			"status":  "healthy",
			"time":    time.Now(),
		}
		c.JSON(http.StatusOK, health)
	})

	// Add GraphQL endpoints
	graphqlHandler := graphql.NewGraphQLHandler(orderService, log)
	playgroundHandler := graphql.NewPlaygroundHandler()

	router.Any("/graphql", gin.WrapH(graphqlHandler))
	router.GET("/graphql/playground", gin.WrapH(playgroundHandler))

	// Start server
	log.WithFields(map[string]interface{}{"port": "8003"}).Info("Starting Order Service (Federation-only)")
	if err := http.ListenAndServe(":8003", router); err != nil {
		log.WithError(err).Fatal("Failed to start server")
	}
}
