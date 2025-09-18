package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"unified-commerce/services/shared/config"
	"unified-commerce/services/shared/logger"
	"unified-commerce/services/shared/middleware"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig("analytics")
	if err != nil {
		panic("Failed to load configuration: " + err.Error())
	}

	// Initialize logger
	loggerConfig := logger.DefaultConfig("analytics")
	loggerConfig.Level = cfg.LogLevel
	log := logger.NewLogger(loggerConfig)

	// Initialize Gin router
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.New()

	// Add middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middleware.CORS())
	router.Use(middleware.RequestID())

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		health := map[string]interface{}{
			"service": "analytics",
			"status":  "healthy",
			"time":    time.Now(),
		}
		c.JSON(http.StatusOK, health)
	})

	// Basic analytics endpoint
	router.GET("/analytics", func(c *gin.Context) {
		c.JSON(http.StatusOK, map[string]interface{}{
			"message": "Analytics service is running",
			"version": "1.0.0",
		})
	})

	// GraphQL endpoint for federation
	router.POST("/graphql", func(c *gin.Context) {
		// Handle GraphQL introspection for federation
		var request map[string]interface{}
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "Invalid JSON"})
			return
		}

		query, ok := request["query"].(string)
		if !ok {
			c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "Query required"})
			return
		}

		// Respond to federation service discovery
		if query == "query{_service{sdl}}" {
			c.JSON(http.StatusOK, map[string]interface{}{
				"data": map[string]interface{}{
					"_service": map[string]interface{}{
						"sdl": "# Analytics Service GraphQL Schema\ntype Query { _service: _Service }\ntype _Service { sdl: String }",
					},
				},
			})
			return
		}

		// Default empty response for other queries
		c.JSON(http.StatusOK, map[string]interface{}{
			"data": map[string]interface{}{},
		})
	})

	// Start REST server
	srv := &http.Server{
		Addr:         ":" + cfg.ServicePort,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.WithFields(map[string]interface{}{"port": cfg.ServicePort}).Info("Starting Analytics Service")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.WithError(err).Fatal("Failed to start server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down Analytics Service...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.WithError(err).Fatal("Server forced to shutdown")
	}

	log.Info("Analytics Service stopped")
}
