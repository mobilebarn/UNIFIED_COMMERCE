package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strings"
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
			log.WithError(err).Error("Failed to parse JSON request")
			c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "Invalid JSON"})
			return
		}

		query, ok := request["query"].(string)
		if !ok {
			log.Error("No query field in request")
			c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "Query required"})
			return
		}

		// Debug: log the incoming query and full request
		log.WithFields(map[string]interface{}{
			"query": query,
			"fullRequest": request,
		}).Info("Received GraphQL query")

		// Check for federation service discovery queries (more flexible matching)
		if strings.Contains(query, "_service") && strings.Contains(query, "sdl") {
			log.Info("Responding to federation service discovery query")
			response := map[string]interface{}{
				"data": map[string]interface{}{
					"_service": map[string]interface{}{
						"sdl": "extend type Query { analytics: String }",
					},
				},
			}
			log.WithFields(map[string]interface{}{"response": response}).Info("Sending federation SDL response")
			c.JSON(http.StatusOK, response)
			return
		}

		// Handle introspection queries
		if strings.Contains(query, "__schema") || strings.Contains(query, "__type") {
			log.Info("Responding to GraphQL introspection query")
			c.JSON(http.StatusOK, map[string]interface{}{
				"data": map[string]interface{}{
					"__schema": map[string]interface{}{
						"queryType": map[string]interface{}{
							"name": "Query",
						},
					},
				},
			})
			return
		}

		// Default response for other queries
		log.WithFields(map[string]interface{}{"query": query}).Info("Responding to regular GraphQL query")
		c.JSON(http.StatusOK, map[string]interface{}{
			"data": map[string]interface{}{
				"analytics": "Analytics service is operational",
			},
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
