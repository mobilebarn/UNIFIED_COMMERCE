package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"retail-os/services/shared/config"
	"retail-os/services/shared/database"
	"retail-os/services/shared/logger"
	"retail-os/services/shared/middleware"

	"./graphql"
	"./internal/analytics"
	"./internal/repository"
	"./internal/service"
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

	// Connect to MongoDB for analytics data
	mongoConfig := database.NewMongoConfigFromEnv(
		cfg.MongoURL,
		cfg.MongoDatabase,
		cfg.MongoUser,
		cfg.MongoPassword,
	)
	mongoDB, err := database.NewMongoConnection(mongoConfig)
	if err != nil {
		log.WithError(err).Fatal("Failed to connect to MongoDB")
	}
	defer mongoDB.Close(context.Background())

	// Connect to PostgreSQL for relational data
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

	// Initialize repositories
	repo := repository.NewRepository(mongoDB, postgresDB)

	// Initialize services
	analyticsService := service.NewAnalyticsService(repo, log)

	// Initialize handlers
	analyticsHandler := analytics.NewHandler(analyticsService, log)

	// Initialize GraphQL resolver
	graphQLResolver := graphql.NewResolver(analyticsService, log)

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
		mongoStatus := "healthy"
		if err := mongoDB.Health(context.Background()); err != nil {
			mongoStatus = "unhealthy"
		}

		postgresStatus := "healthy"
		if err := postgresDB.Health(context.Background()); err != nil {
			postgresStatus = "unhealthy"
		}

		health := map[string]interface{}{
			"service": "analytics",
			"status":  "healthy",
			"time":    time.Now(),
			"checks": map[string]string{
				"mongodb":  mongoStatus,
				"postgres": postgresStatus,
			},
		}
		c.JSON(http.StatusOK, health)
	})

	// Register routes
	analyticsHandler.RegisterRoutes(router)

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

	// Start GraphQL server in a separate goroutine
	go func() {
		graphQLServer := graphql.NewServer(graphQLResolver)

		// Create a new mux for GraphQL
		mux := http.NewServeMux()
		mux.Handle("/graphql", graphQLServer.Handler())
		mux.Handle("/", graphQLServer.Playground())

		log.Info("Starting GraphQL server on port 8081")
		if err := http.ListenAndServe(":8081", mux); err != nil {
			log.WithError(err).Fatal("Failed to start GraphQL server")
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
