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

	"unified-commerce/services/order/graphql"
	"unified-commerce/services/order/handlers"
	"unified-commerce/services/order/models"
	"unified-commerce/services/order/repository"
	"unified-commerce/services/order/service"
	"unified-commerce/services/shared/config"
	"unified-commerce/services/shared/database"
	"unified-commerce/services/shared/logger"
	"unified-commerce/services/shared/middleware"
	"unified-commerce/shared/messaging"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig("order")
	if err != nil {
		panic("Failed to load configuration: " + err.Error())
	}

	// Initialize logger
	loggerConfig := logger.DefaultConfig("order")
	loggerConfig.Level = cfg.LogLevel
	log := logger.NewLogger(loggerConfig)

	// Connect to PostgreSQL using shared database utility
	postgresConfig := database.NewPostgresConfigFromEnv(
		cfg.DatabaseHost,
		cfg.DatabasePort,
		cfg.DatabaseUser,
		cfg.DatabasePassword,
		cfg.DatabaseName,
	)
	// Set the full DATABASE_URL if available (for cloud providers like Render)
	postgresConfig.DatabaseURL = cfg.DatabaseURL
	postgresDB, err := database.NewPostgresConnection(postgresConfig)
	if err != nil {
		log.WithError(err).Fatal("Failed to connect to PostgreSQL")
	}
	defer postgresDB.Close()

	// Initialize Event Producer
	producerConfig := messaging.ProducerConfig{
		Brokers:   cfg.KafkaBrokers,
		UseDocker: messaging.DetectEnvironment(),
	}
	producer, err := messaging.NewEventProducer(producerConfig)
	if err != nil {
		log.WithError(err).Warn("Failed to create event producer, using no-op producer for graceful degradation")
		// Create a no-op producer to allow service to continue without Kafka
		producer = messaging.NewNoOpProducer()
	}
	defer producer.Close()

	// Run database migrations
	if err := runMigrations(postgresDB.DB); err != nil {
		log.WithError(err).Fatal("Failed to run database migrations")
	}

	// Initialize repositories
	orderRepo := repository.NewOrderRepository(postgresDB.DB, log)

	// Initialize services
	orderService := service.NewOrderService(orderRepo, log, producer)

	// Initialize handlers
	orderHandler := handlers.NewOrderHandler(orderService, log)

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
	router.Use(middleware.RateLimit())

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		health := map[string]interface{}{
			"service": "order",
			"status":  "healthy",
			"time":    time.Now(),
			"checks": map[string]string{
				"postgres": checkPostgreSQL(postgresDB),
			},
		}
		c.JSON(http.StatusOK, health)
	})

	// Register routes
	orderHandler.RegisterRoutes(router)

	// Add GraphQL endpoints
	graphqlHandler := graphql.NewGraphQLHandler(orderService, log)
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

	// Start background tasks
	go startBackgroundTasks(orderService, log)

	// Start server in a goroutine
	go func() {
		log.WithFields(map[string]interface{}{"port": cfg.ServicePort}).Info("Starting Order Service")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.WithError(err).Fatal("Failed to start server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down Order Service...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.WithError(err).Fatal("Server forced to shutdown")
	}

	log.Info("Order Service stopped")
}

// runMigrations runs database migrations
func runMigrations(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.Order{},
		&models.OrderLineItem{},
		&models.Fulfillment{},
		&models.FulfillmentLineItem{},
		&models.Transaction{},
		&models.Return{},
		&models.ReturnLineItem{},
		&models.OrderEvent{},
	)
}

// checkPostgreSQL checks PostgreSQL connection health
func checkPostgreSQL(db *database.PostgresDB) string {
	sqlDB, err := db.DB.DB()
	if err != nil {
		return "unhealthy"
	}

	if err := sqlDB.Ping(); err != nil {
		return "unhealthy"
	}
	return "healthy"
}

// startBackgroundTasks starts background tasks for order management
func startBackgroundTasks(service *service.OrderService, log *logger.Logger) {
	// Order processing tasks could be added here
	// For example: auto-confirm orders after payment, send notifications, etc.

	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// Background task placeholder
			log.Debug("Running background order processing tasks")
		}
	}
}
