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

	"unified-commerce/services/cart/graphql"
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
	cfg, err := config.LoadConfig("cart")
	if err != nil {
		panic("Failed to load configuration: " + err.Error())
	}

	// Initialize logger
	loggerConfig := logger.DefaultConfig("cart")
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

	// Initialize repositories
	cartRepo := repository.NewCartRepository(postgresDB.DB, log)

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

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		postgresStatus := "healthy"
		if err := postgresDB.Health(context.Background()); err != nil {
			postgresStatus = "unhealthy"
		}

		health := map[string]interface{}{
			"service": "cart",
			"status":  "healthy",
			"time":    time.Now(),
			"checks": map[string]string{
				"postgres": postgresStatus,
			},
		}
		c.JSON(http.StatusOK, health)
	})

	// Register routes
	cartHandler.RegisterRoutes(router)

	// Add GraphQL endpoints
	graphqlHandler := graphql.NewGraphQLHandler(cartService, log)
	playgroundHandler := graphql.NewGraphQLPlaygroundHandler()

	router.Any("/graphql", gin.WrapH(graphqlHandler))
	router.GET("/graphql/playground", gin.WrapH(playgroundHandler))

	// Start server
	srv := &http.Server{
		Addr:         ":" + cfg.ServicePort,
		Handler:      router,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start background tasks
	go startBackgroundTasks(cartService, log)

	// Start server in a goroutine
	go func() {
		log.WithField("port", cfg.ServicePort).Info("Starting Cart & Checkout Service")
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
