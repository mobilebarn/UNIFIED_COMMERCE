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

	"unified-commerce/services/inventory/eventhandlers"
	"unified-commerce/services/inventory/graphql"
	"unified-commerce/services/inventory/handlers"
	"unified-commerce/services/inventory/models"
	"unified-commerce/services/inventory/repository"
	"unified-commerce/services/inventory/service"
	"unified-commerce/services/shared/config"
	"unified-commerce/services/shared/database"
	"unified-commerce/services/shared/logger"
	"unified-commerce/services/shared/middleware"
	"unified-commerce/shared/messaging"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig("inventory")
	if err != nil {
		panic("Failed to load configuration: " + err.Error())
	}

	// Initialize logger
	loggerConfig := logger.DefaultConfig("inventory")
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
	postgresDB, err := database.NewPostgresConnection(postgresConfig)
	if err != nil {
		log.WithError(err).Fatal("Failed to connect to PostgreSQL")
	}
	defer postgresDB.Close()

	// Run database migrations
	if err := runMigrations(postgresDB.DB); err != nil {
		log.WithError(err).Fatal("Failed to run database migrations")
	}

	// Initialize repositories
	inventoryRepo := repository.NewInventoryRepository(postgresDB.DB, log)

	// Initialize services
	inventoryService := service.NewInventoryService(inventoryRepo, log)

	// Initialize event handlers
	orderEventHandler := eventhandlers.NewOrderEventHandler(inventoryService, log)

	// Initialize Event Consumer
	consumerConfig := messaging.ConsumerConfig{
		Brokers:   cfg.KafkaBrokers,
		GroupID:   "inventory-service",
		Topics:    []string{"orders.placed"},
		UseDocker: messaging.DetectEnvironment(),
	}
	consumer, err := messaging.NewEventConsumer(consumerConfig)
	if err != nil {
		log.WithError(err).Fatal("Failed to create event consumer")
	}
	defer consumer.Close()

	// Start consuming messages in a separate goroutine
	go startEventConsumption(consumer, orderEventHandler, log)

	// Initialize handlers
	inventoryHandler := handlers.NewInventoryHandler(inventoryService, log)

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
			"service": "inventory",
			"status":  "healthy",
			"time":    time.Now(),
			"checks": map[string]string{
				"postgres": checkPostgreSQL(postgresDB),
			},
		}
		c.JSON(http.StatusOK, health)
	})

	// Register routes
	inventoryHandler.RegisterRoutes(router)

	// Add GraphQL endpoints
	graphqlHandler := graphql.NewGraphQLHandler(inventoryService, log)
	playgroundHandler := graphql.NewPlaygroundHandler()

	router.Any("/graphql", gin.WrapH(graphqlHandler))
	router.GET("/graphql/playground", gin.WrapH(playgroundHandler))

	// Start server
	srv := &http.Server{
		Addr:         ":8003", // Temporarily hardcode inventory service port
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start background tasks
	go startBackgroundTasks(inventoryService, log)

	// Start server in a goroutine
	go func() {
		log.WithFields(map[string]interface{}{"port": "8003"}).Info("Starting Inventory Service")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.WithError(err).Fatal("Failed to start server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down Inventory Service...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.WithError(err).Fatal("Server forced to shutdown")
	}

	log.Info("Inventory Service stopped")
}

// runMigrations runs database migrations
func runMigrations(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.Location{},
		&models.InventoryItem{},
		&models.StockMovement{},
		&models.StockReservation{},
		&models.StockTransfer{},
		&models.StockTransferItem{},
		&models.StockAlert{},
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

// startEventConsumption handles event consumption from the messaging system
func startEventConsumption(consumer messaging.EventConsumer, handler *eventhandlers.OrderEventHandler, log *logger.Logger) {
	if err := consumer.Subscribe([]string{"orders.placed"}); err != nil {
		log.WithError(err).Fatal("Failed to subscribe to topics")
	}

	for {
		msg, err := consumer.ReadMessage()
		if err != nil {
			log.WithError(err).Error("Failed to read message")
			time.Sleep(1 * time.Second)
			continue
		}

		// Convert to the format expected by the handler
		if err := handler.HandleOrderPlacedEvent(msg.Value); err != nil {
			log.WithError(err).Error("Failed to handle order placed event")
		} else {
			if err := consumer.CommitMessage(msg); err != nil {
				log.WithError(err).Error("Failed to commit message")
			}
		}
	}
}

// startBackgroundTasks starts background tasks for inventory management
func startBackgroundTasks(service *service.InventoryService, log *logger.Logger) {
	// Process expired reservations every 5 minutes
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			if err := service.ProcessExpiredReservations(ctx); err != nil {
				log.WithError(err).Error("Failed to process expired reservations")
			}
			cancel()
		}
	}
}
