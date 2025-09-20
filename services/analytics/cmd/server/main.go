package main

import (
	"context"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"unified-commerce/services/shared/service"
)

func main() {
	// Create base service with PostgreSQL, MongoDB and Redis
	baseService, err := service.NewBaseService(service.ServiceOptions{
		Name:        "analytics",
		UsePostgres: true,
		UseMongoDB:  true,
		UseRedis:    true, // Enable Redis for analytics caching
	})
	if err != nil {
		panic("Failed to create base service: " + err.Error())
	}

	// Setup routes
	setupRoutes(baseService.Router, baseService)

	baseService.Logger.Info("Analytics Service started successfully")

	// Start the service
	if err := baseService.Start(); err != nil {
		baseService.Logger.WithError(err).Fatal("Failed to start service")
	}
}

// setupRoutes configures the HTTP routes for the analytics service
func setupRoutes(router *gin.Engine, baseService *service.BaseService) {
	// Basic analytics endpoint
	router.GET("/analytics", func(c *gin.Context) {
		c.JSON(200, map[string]interface{}{
			"message": "Analytics service is running",
			"version": "1.0.0",
		})
	})

	// GraphQL endpoint for federation
	router.POST("/graphql", func(c *gin.Context) {
		var request map[string]interface{}
		if err := c.ShouldBindJSON(&request); err != nil {
			baseService.Logger.WithError(err).Error("Failed to parse JSON request")
			c.JSON(400, map[string]interface{}{"error": "Invalid JSON"})
			return
		}

		query, ok := request["query"].(string)
		if !ok {
			baseService.Logger.Error("No query field in request")
			c.JSON(400, map[string]interface{}{"error": "Query required"})
			return
		}

		// Check for federation service discovery queries
		if strings.Contains(query, "_service") && strings.Contains(query, "sdl") {
			baseService.Logger.Info("Responding to federation service discovery query")
			
			// Analytics SDL for GraphQL Federation
			analyticsSDL := `
				directive @key(fields: String!) on OBJECT | INTERFACE
				directive @extends on OBJECT | INTERFACE  
				directive @external on OBJECT | FIELD_DEFINITION

				extend type User @key(fields: "id") {
					id: ID! @external
					behaviors: [CustomerBehavior!]!
				}

				type CustomerBehavior {
					id: ID!
					customerId: ID!
					action: String!
					timestamp: String!
				}

				extend type Query {
					customerBehaviors(customerId: ID!): [CustomerBehavior!]!
				}
			`
			
			response := map[string]interface{}{
				"data": map[string]interface{}{
					"_service": map[string]interface{}{
						"sdl": analyticsSDL,
					},
				},
			}
			c.JSON(200, response)
			return
		}

		// Handle introspection queries
		if strings.Contains(query, "__schema") || strings.Contains(query, "__type") {
			baseService.Logger.Info("Responding to GraphQL introspection query")
			c.JSON(200, map[string]interface{}{
				"data": map[string]interface{}{
					"__schema": map[string]interface{}{
						"queryType": map[string]interface{}{
							"name": "Query",
						},
					},
				},
			})
			return
		}

		// Sample data for customer behaviors
		if strings.Contains(query, "customerBehaviors") {
			c.JSON(200, map[string]interface{}{
				"data": map[string]interface{}{
					"customerBehaviors": []map[string]interface{}{
						{
							"id":         "behavior-1",
							"customerId": "customer-123",
							"action":     "VIEW",
							"timestamp":  time.Now().Format(time.RFC3339),
						},
					},
				},
			})
			return
		}

		// Default response
		c.JSON(200, map[string]interface{}{
			"data": map[string]interface{}{
				"analytics": "Analytics service is operational",
			},
		})
	})
}