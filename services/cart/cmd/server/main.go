package main

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"unified-commerce/services/cart/graphql"
	"unified-commerce/services/cart/handlers"
	"unified-commerce/services/cart/models"
	"unified-commerce/services/cart/repository"
	"unified-commerce/services/cart/service"
	sharedService "unified-commerce/services/shared/service"
)

func main() {
	// Create base service with PostgreSQL and Redis
	baseService, err := sharedService.NewBaseService(sharedService.ServiceOptions{
		Name:        "cart",
		UsePostgres: true,
		UseMongoDB:  false,
		UseRedis:    true, // Enable Redis for cart session management
	})
	if err != nil {
		panic("Failed to create base service: " + err.Error())
	}

	// Run database migrations
	if err := runMigrations(baseService.PostgresDB.DB); err != nil {
		baseService.Logger.WithError(err).Fatal("Failed to run database migrations")
	}

	// Setup routes
	setupRoutes(baseService.Router, baseService)

	// Start background tasks
	go startBackgroundTasks(baseService)

	baseService.Logger.Info("Cart Service started successfully")

	// Start the service
	if err := baseService.Start(); err != nil {
		baseService.Logger.WithError(err).Fatal("Failed to start service")
	}
}

// setupRoutes configures the HTTP routes for the cart service
func setupRoutes(router *gin.Engine, baseService *sharedService.BaseService) {
	// Initialize repositories
	cartRepo := repository.NewCartRepository(baseService.PostgresDB.DB, baseService.Logger)

	// Initialize services
	cartService := service.NewCartService(cartRepo, baseService.Logger)

	// Initialize handlers
	cartHandler := handlers.NewCartHandler(cartService, baseService.Logger)

	// Register routes
	cartHandler.RegisterRoutes(router)

	// Add GraphQL endpoints
	graphqlHandler := graphql.NewGraphQLHandler(cartService, baseService.Logger)
	playgroundHandler := graphql.NewPlaygroundHandler()

	router.Any("/graphql", gin.WrapH(graphqlHandler))
	
	// Only expose playground in non-production environments
	if baseService.Config.Environment != "production" {
		router.GET("/graphql/playground", gin.WrapH(playgroundHandler))
	}
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
func startBackgroundTasks(baseService *sharedService.BaseService) {
	// Initialize repositories and services for background tasks
	cartRepo := repository.NewCartRepository(baseService.PostgresDB.DB, baseService.Logger)
	cartService := service.NewCartService(cartRepo, baseService.Logger)

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
			if err := cartService.MarkAbandonedCarts(context.Background(), abandonmentThreshold); err != nil {
				baseService.Logger.WithError(err).Error("Failed to mark abandoned carts")
			} else {
				baseService.Logger.Info("Abandoned cart processing completed")
			}
		case <-cleanupTicker.C:
			// Cleanup expired carts
			if err := cartService.CleanupExpiredCarts(context.Background()); err != nil {
				baseService.Logger.WithError(err).Error("Failed to cleanup expired carts")
			} else {
				baseService.Logger.Info("Expired cart cleanup completed")
			}
		}
	}
}
