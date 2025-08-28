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

	"unified-commerce/services/cart/handlers"
	"unified-commerce/services/cart/models"
	"unified-commerce/services/cart/repository"
	"unified-commerce/services/cart/service"
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
	cartRepo := repository.NewCartRepository(db, log)

	// Initialize services
	cartService := service.NewCartService(cartRepo, log)

	// Initialize handlers
	cartHandler := handlers.NewCartHandler(cartService, log)

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
			"service": "cart",
			"status":  "healthy",
			"time":    time.Now(),
			"checks": map[string]string{
				"postgres": checkPostgreSQL(db),
			},
		}
		c.JSON(http.StatusOK, health)
	})

	// Register routes
	cartHandler.RegisterRoutes(router)

	// Start server
	srv := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(cfg.Server.IdleTimeout) * time.Second,
	}

	// Start background tasks
	go startBackgroundTasks(cartService, log)

	// Start server in a goroutine
	go func() {
		log.WithField("port", cfg.Server.Port).Info("Starting Cart & Checkout Service")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.WithError(err).Fatal("Failed to start server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down Cart & Checkout Service...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.WithError(err).Fatal("Server forced to shutdown")
	}

	log.Info("Cart & Checkout Service stopped")
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
		&models.Cart{},
		&models.CartLineItem{},
		&models.CartTaxLine{},
		&models.CartShippingLine{},
		&models.CartDiscountApplication{},
		&models.CartLineItemDiscountAllocation{},
		&models.Checkout{},
		&models.CheckoutEvent{},
		&models.ShippingRate{},
		&models.PaymentMethod{},
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

// startBackgroundTasks starts background tasks for cart management
func startBackgroundTasks(service *service.CartService, log *logger.Logger) {
	// Process abandoned carts every hour
	abandonmentTicker := time.NewTicker(1 * time.Hour)
	defer abandonmentTicker.Stop()

	// Cleanup expired carts daily
	cleanupTicker := time.NewTicker(24 * time.Hour)
	defer cleanupTicker.Stop()

	for {
		select {
		case <-abandonmentTicker.C:
			// Mark carts as abandoned after 24 hours of inactivity
			abandonmentThreshold := 24 * time.Hour
			if err := service.MarkAbandonedCarts(context.Background(), abandonmentThreshold); err != nil {
				log.WithError(err).Error("Failed to mark abandoned carts")
			} else {
				log.Info("Abandoned cart processing completed")
			}
		case <-cleanupTicker.C:
			// Cleanup expired carts
			if err := service.CleanupExpiredCarts(context.Background()); err != nil {
				log.WithError(err).Error("Failed to cleanup expired carts")
			} else {
				log.Info("Expired cart cleanup completed")
			}
		}
	}
}
