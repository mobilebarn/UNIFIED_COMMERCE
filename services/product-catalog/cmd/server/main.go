package main

import (
	"strings"

	"github.com/gin-gonic/gin"

	"unified-commerce/services/product-catalog/graphql"
	"unified-commerce/services/product-catalog/handlers"
	"unified-commerce/services/product-catalog/repository"
	productService "unified-commerce/services/product-catalog/service"
	sharedService "unified-commerce/services/shared/service"
)

func main() {
	// Create base service with MongoDB and Redis
	baseService, err := sharedService.NewBaseService(sharedService.ServiceOptions{
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
func setupRoutes(router *gin.Engine, baseService *sharedService.BaseService) {
	// Check if MongoDB is available
	if baseService.MongoDB == nil {
		baseService.Logger.Warn("MongoDB not available, setting up service with mock data for debugging")
		
		// Add a simple health endpoint for when MongoDB is unavailable
		router.GET("/status", func(c *gin.Context) {
			c.JSON(200, map[string]interface{}{
				"service": "product-catalog",
				"status":  "limited - MongoDB unavailable",
				"message": "Service is running but database is not available",
			})
		})
		
		// Setup GraphQL endpoint with mock data for debugging
		router.POST("/graphql", func(c *gin.Context) {
			var request struct {
				Query string `json:"query"`
			}
			
			if err := c.ShouldBindJSON(&request); err != nil {
				c.JSON(400, map[string]interface{}{"error": "Invalid request"})
				return
			}
			
			// Handle federation discovery query
			if strings.Contains(request.Query, "_service") {
				c.JSON(200, map[string]interface{}{
					"data": map[string]interface{}{
						"_service": map[string]interface{}{
							"sdl": `
								type Product @key(fields: "id") {
									id: ID!
									title: String!
									description: String!
									priceRange: PriceRange!
									variants: [ProductVariant!]!
									category: Category
									images: [String!]!
								}
								
								type PriceRange {
									minVariantPrice: Float!
									maxVariantPrice: Float!
								}
								
								type ProductVariant {
									id: ID!
									price: Float!
									inventoryQuantity: Int!
								}
								
								type Category {
									id: ID!
									name: String!
								}
								
								type Query {
									products: [Product!]!
									product(id: ID!): Product
								}
							`,
						},
					},
				})
				return
			}
			
			// Handle products query with mock data
			if strings.Contains(request.Query, "products") {
				c.JSON(200, map[string]interface{}{
					"data": map[string]interface{}{
						"products": []map[string]interface{}{
							{
								"id":          "1",
								"title":       "Sample Product",
								"description": "This is a sample product while MongoDB is unavailable",
								"priceRange": map[string]interface{}{
									"minVariantPrice": 29.99,
									"maxVariantPrice": 29.99,
								},
								"variants": []map[string]interface{}{
									{
										"id":                "1",
										"price":             29.99,
										"inventoryQuantity": 10,
									},
								},
								"category": map[string]interface{}{
									"id":   "1",
									"name": "Sample Category",
								},
								"images": []string{"https://via.placeholder.com/300"},
							},
						},
					},
				})
				return
			}
			
			// Default empty response
			c.JSON(200, map[string]interface{}{"data": nil})
		})
		return
	}

	// Initialize repositories with MongoDB
	repo := repository.NewRepository(baseService.MongoDB)

	// Initialize services
	productServiceInstance := productService.NewProductService(repo, baseService.Logger)

	// Initialize handlers
	productHandler := handlers.NewProductHandler(productServiceInstance, baseService.Logger)

	// Register routes
	productHandler.RegisterRoutes(router)

	// Add GraphQL endpoints
	graphqlHandler := graphql.NewGraphQLHandler(productServiceInstance, baseService.Logger)
	playgroundHandler := graphql.NewPlaygroundHandler()

	router.Any("/graphql", gin.WrapH(graphqlHandler))
	
	// Only expose playground in non-production environments
	if baseService.Config.Environment != "production" {
		router.GET("/graphql/playground", gin.WrapH(playgroundHandler))
	}
}
