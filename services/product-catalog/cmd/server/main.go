package main

import (
	"github.com/gin-gonic/gin"

	"unified-commerce/services/product-catalog/graphql"
	"unified-commerce/services/product-catalog/handlers"
	"unified-commerce/services/product-catalog/repository"
	"unified-commerce/services/product-catalog/service"
	"unified-commerce/services/shared/service"
)

func main() {
	// Create base service with MongoDB and Redis
	baseService, err := service.NewBaseService(service.ServiceOptions{
		Name:        "product-catalog",
		UsePostgres: false,
		UseMongoDB:  true,
		UseRedis:    true, // Enable Redis for product caching
	})
	if err != nil {
		panic("Failed to create base service: " + err.Error())
	}

	// Setup routes
	setupRoutes(baseService.Router, baseService)

	baseService.Logger.Info("Product Catalog Service started successfully")

	// Start the service
	if err := baseService.Start(); err != nil {
		baseService.Logger.WithError(err).Fatal("Failed to start service")
	}
}

// setupRoutes configures the HTTP routes for the product catalog service
func setupRoutes(router *gin.Engine, baseService *service.BaseService) {
	// Initialize repositories
	repo := repository.NewRepository(baseService.MongoDB)

	// Initialize services
	productService := service.NewProductService(repo, baseService.Logger)

	// Initialize handlers
	productHandler := handlers.NewProductHandler(productService, baseService.Logger)

	// Register routes
	productHandler.RegisterRoutes(router)

	// Add GraphQL endpoints
	graphqlHandler := graphql.NewGraphQLHandler(productService, baseService.Logger)
	playgroundHandler := graphql.NewPlaygroundHandler()

	router.Any("/graphql", gin.WrapH(graphqlHandler))
	
	// Only expose playground in non-production environments
	if baseService.Config.Environment != "production" {
		router.GET("/graphql/playground", gin.WrapH(playgroundHandler))
	}
}
