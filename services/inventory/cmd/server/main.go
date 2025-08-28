package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"unified-commerce/services/inventory/handlers"
	"unified-commerce/services/inventory/models"
	"unified-commerce/services/inventory/repository"
	"unified-commerce/services/inventory/service"
	"unified-commerce/services/shared/config"
	"unified-commerce/services/shared/database"
	"unified-commerce/services/shared/logger"
	"unified-commerce/services/shared/middleware"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize logger
	log := logger.New(cfg.Environment)

	// Connect to PostgreSQL
	db, err := connectPostgreSQL(cfg.Database.PostgreSQL.URL)
	if err != nil {
		log.WithError(err).Fatal("Failed to connect to PostgreSQL")
	}

	// Run database migrations
	if err := runMigrations(db); err != nil {
		log.WithError(err).Fatal("Failed to run database migrations")
	}

	// Initialize repositories
	inventoryRepo := repository.NewInventoryRepository(db, log)

	// Initialize services
	inventoryService := service.NewInventoryService(inventoryRepo, log)

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
				"postgres": checkPostgreSQL(db),
			},
		}
		c.JSON(http.StatusOK, health)
	})

	// Register routes
	inventoryHandler.RegisterRoutes(router)

	// Start server
	srv := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(cfg.Server.IdleTimeout) * time.Second,
	}

	// Start background tasks
	go startBackgroundTasks(inventoryService, log)

	// Start server in a goroutine
	go func() {
		log.WithField("port", cfg.Server.Port).Info("Starting Inventory Service")
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

// connectPostgreSQL establishes connection to PostgreSQL
func connectPostgreSQL(databaseURL string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{
		Logger: database.NewGormLogger(),
	})
	if err != nil {
		return nil, err
	}

	// Get underlying sql.DB to configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// Configure connection pool
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db, nil
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
func checkPostgreSQL(db *gorm.DB) string {
	sqlDB, err := db.DB()
	if err != nil {
		return "unhealthy"
	}

	if err := sqlDB.Ping(); err != nil {
		return "unhealthy"
	}
	return "healthy"
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
