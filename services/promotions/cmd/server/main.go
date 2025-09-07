package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"unified-commerce/services/promotions/graphql"
	"unified-commerce/services/promotions/handlers"
	"unified-commerce/services/promotions/repository"
	"unified-commerce/services/promotions/service"
	"unified-commerce/services/shared/config"
	"unified-commerce/services/shared/database"
	"unified-commerce/services/shared/logger"
	"unified-commerce/services/shared/middleware"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig("promotions")
	if err != nil {
		panic("Failed to load configuration: " + err.Error())
	}

	// Initialize logger
	loggerConfig := logger.DefaultConfig("promotions")
	loggerConfig.Level = cfg.LogLevel
	log := logger.NewLogger(loggerConfig)

	// Connect to PostgreSQL
	postgresConfig := database.NewPostgresConfigFromEnv(
		cfg.DatabaseHost,
		cfg.DatabasePort,
		cfg.DatabaseUser,
		cfg.DatabasePassword,
		cfg.DatabaseName,
	)
	postgresDB, err := database.NewPostgresConnection(postgresConfig)
	if err != nil {
		log.WithError(err).Fatal("Failed to connect to PostgreSQL")
	}
	defer postgresDB.Close()

	// Run database migrations
	if err := runMigrations(postgresDB.DB); err != nil {
		log.WithError(err).Fatal("Failed to run database migrations")
	}

	// Initialize repository
	promotionsRepo := repository.NewPromotionsRepository(postgresDB.DB, log)

	// Initialize service
	promotionsService := service.NewPromotionsService(promotionsRepo, log)

	// Initialize handler
	promotionsHandler := handlers.NewPromotionsHandler(promotionsService, log)

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
		postgresStatus := "healthy"
		if err := postgresDB.Health(context.Background()); err != nil {
			postgresStatus = "unhealthy"
		}

		health := map[string]interface{}{
			"service": "promotions",
			"status":  "healthy",
			"time":    time.Now(),
			"checks": map[string]string{
				"postgres": postgresStatus,
			},
		}
		c.JSON(http.StatusOK, health)
	})

	// Register routes
	promotionsHandler.RegisterRoutes(router)

	// Add GraphQL endpoints
	graphqlHandler := graphql.NewGraphQLHandler(promotionsService, log)
	playgroundHandler := graphql.NewPlaygroundHandler()

	router.Any("/graphql", gin.WrapH(graphqlHandler))
	router.GET("/graphql/playground", gin.WrapH(playgroundHandler))

	// Start server
	srv := &http.Server{
		Addr:         ":" + cfg.ServicePort,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.WithField("port", cfg.ServicePort).Info("Starting Promotions Service")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.WithError(err).Fatal("Failed to start server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down Promotions Service...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.WithError(err).Fatal("Server forced to shutdown")
	}

	log.Info("Promotions Service stopped")
}

// runMigrations runs database migrations
func runMigrations(db *gorm.DB) error {
	// Add all promotions models here for migration
	return db.AutoMigrate()
}
