package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"unified-commerce/services/promotions/graphql"
	"unified-commerce/services/promotions/handlers"
	"unified-commerce/services/promotions/repository"
	promotionsService "unified-commerce/services/promotions/service"
	"unified-commerce/services/promotions/models"
	sharedService "unified-commerce/services/shared/service"
)

func main() {
	// Create base service with PostgreSQL and Redis
	baseService, err := sharedService.NewBaseService(sharedService.ServiceOptions{
		Name:        "promotions",
		UsePostgres: true,
		UseMongoDB:  false,
		UseRedis:    true, // Enable Redis for promotion caching
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

	baseService.Logger.Info("Promotions Service started successfully")

	// Start the service
	if err := baseService.Start(); err != nil {
		baseService.Logger.WithError(err).Fatal("Failed to start service")
	}
}

// setupRoutes configures the HTTP routes for the promotions service
func setupRoutes(router *gin.Engine, baseService *sharedService.BaseService) {
	// Initialize repository
	promotionsRepo := repository.NewPromotionsRepository(baseService.PostgresDB.DB, baseService.Logger)

	// Initialize service
	promotionsServiceInstance := promotionsService.NewPromotionsService(promotionsRepo, baseService.Logger)

	// Initialize handler
	promotionsHandler := handlers.NewPromotionsHandler(promotionsServiceInstance, baseService.Logger)

	// Register routes
	promotionsHandler.RegisterRoutes(router)

	// Add GraphQL endpoints
	graphqlHandler := graphql.NewGraphQLHandler(promotionsServiceInstance, baseService.Logger)
	playgroundHandler := graphql.NewPlaygroundHandler()

	router.Any("/graphql", gin.WrapH(graphqlHandler))
	
	// Only expose playground in non-production environments
	if baseService.Config.Environment != "production" {
		router.GET("/graphql/playground", gin.WrapH(playgroundHandler))
	}
}

// runMigrations runs database migrations
func runMigrations(db *gorm.DB) error {
	// Add all promotions models here for migration
	return db.AutoMigrate(
		&models.Promotion{},
		&models.DiscountCode{},
		&models.CodeUsage{},
		&models.GiftCard{},
		&models.GiftCardTransaction{},
		&models.LoyaltyProgram{},
		&models.LoyaltyMember{},
		&models.LoyaltyTier{},
		&models.LoyaltyActivity{},
	)
}
