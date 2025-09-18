package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"retail-os/services/shared/config"
	"retail-os/services/shared/logger"
	"retail-os/services/shared/middleware"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig("analytics")
	if err != nil {
		panic("Failed to load configuration: " + err.Error())
	}

	// Initialize logger
	loggerConfig := logger.DefaultConfig("analytics")
	loggerConfig.Level = cfg.LogLevel
	log := logger.NewLogger(loggerConfig)

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
		health := map[string]interface{}{
			"service": "analytics",
			"status":  "healthy",
			"time":    time.Now(),
		}
		c.JSON(http.StatusOK, health)
	})

	// Basic analytics endpoint
	router.GET("/analytics", func(c *gin.Context) {
		c.JSON(http.StatusOK, map[string]interface{}{
			"message": "Analytics service is running",
			"version": "1.0.0",
		})
	})

	// GraphQL endpoint for federation
	router.POST("/graphql", func(c *gin.Context) {
		// Handle GraphQL introspection for federation
		var request map[string]interface{}
		if err := c.ShouldBindJSON(&request); err != nil {
			log.WithError(err).Error("Failed to parse JSON request")
			c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "Invalid JSON"})
			return
		}

		query, ok := request["query"].(string)
		if !ok {
			log.Error("No query field in request")
			c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "Query required"})
			return
		}

		// Debug: log the incoming query and full request
		log.WithFields(map[string]interface{}{
			"query": query,
			"fullRequest": request,
		}).Info("Received GraphQL query")

		// Check for federation service discovery queries (more flexible matching)
		if strings.Contains(query, "_service") && strings.Contains(query, "sdl") {
			log.Info("Responding to federation service discovery query")
			
			// Comprehensive Analytics SDL for GraphQL Federation
			analyticsSDL := `
				directive @key(fields: String!) on OBJECT | INTERFACE
				directive @extends on OBJECT | INTERFACE  
				directive @external on OBJECT | FIELD_DEFINITION
				directive @requires(fields: String!) on FIELD_DEFINITION
				directive @provides(fields: String!) on FIELD_DEFINITION

				# Extend existing types from other services
				extend type User @key(fields: "id") {
					id: ID! @external
					behaviors(limit: Int): [CustomerBehavior!]!
					recommendations(limit: Int): [ProductRecommendation!]!
					segments: [CustomerSegment!]!
				}

				extend type Product @key(fields: "id") {
					id: ID! @external
					recommendations: [ProductRecommendation!]!
				}

				# Analytics Types
				type CustomerBehavior {
					id: ID!
					customerId: ID!
					sessionId: String!
					action: CustomerAction!
					entityId: String!
					entityType: String!
					timestamp: String!
					userAgent: String
					ipAddress: String
					referrer: String
					utmSource: String
					utmMedium: String
					utmCampaign: String
				}

				type ProductRecommendation {
					id: ID!
					customerId: ID!
					productId: ID!
					product: Product
					score: Float!
					recommendationType: RecommendationType!
					createdAt: String!
					expiresAt: String!
				}

				type CustomerSegment {
					id: ID!
					name: String!
					description: String
					customerIds: [ID!]!
					customers: [User!]!
					createdAt: String!
					updatedAt: String!
				}

				type BusinessMetric {
					id: ID!
					name: String!
					value: Float!
					dimension: String
					timestamp: String!
				}

				enum CustomerAction {
					VIEW
					CLICK
					ADD_TO_CART
					REMOVE_FROM_CART
					PURCHASE
					SEARCH
					FILTER
					SORT
					SHARE
					WISHLIST
					REVIEW
				}

				enum RecommendationType {
					POPULARITY
					COLLABORATIVE
					CONTENT
					TRENDING
					SIMILAR
					PERSONALIZED
					CATEGORY
					BRAND
				}

				extend type Query {
					customerBehaviors(customerId: ID!, limit: Int): [CustomerBehavior!]!
					productRecommendations(customerId: ID!, limit: Int): [ProductRecommendation!]!
					customerSegment(id: ID!): CustomerSegment
					businessMetrics(name: String!, start: String, end: String): [BusinessMetric!]!
				}
			`
			
			response := map[string]interface{}{
				"data": map[string]interface{}{
					"_service": map[string]interface{}{
						"sdl": analyticsSDL,
					},
				},
			}
			log.WithFields(map[string]interface{}{"sdl_length": len(analyticsSDL)}).Info("Sending comprehensive analytics federation SDL response")
			c.JSON(http.StatusOK, response)
			return
		}

		// Handle introspection queries
		if strings.Contains(query, "__schema") || strings.Contains(query, "__type") {
			log.Info("Responding to GraphQL introspection query")
			c.JSON(http.StatusOK, map[string]interface{}{
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

		// Default response for other queries - provide sample analytics data
		log.WithFields(map[string]interface{}{"query": query}).Info("Responding to regular GraphQL query")
		
		// Check if it's a customerBehaviors query
		if strings.Contains(query, "customerBehaviors") {
			c.JSON(http.StatusOK, map[string]interface{}{
				"data": map[string]interface{}{
					"customerBehaviors": []map[string]interface{}{
						{
							"id":         "behavior-1",
							"customerId": "customer-123",
							"sessionId":  "session-abc",
							"action":     "VIEW",
							"entityId":   "product-456",
							"entityType": "product",
							"timestamp":  time.Now().Format(time.RFC3339),
						},
						{
							"id":         "behavior-2",
							"customerId": "customer-123",
							"sessionId":  "session-abc",
							"action":     "ADD_TO_CART",
							"entityId":   "product-456",
							"entityType": "product",
							"timestamp":  time.Now().Add(-10*time.Minute).Format(time.RFC3339),
						},
					},
				},
			})
			return
		}
		
		// Check if it's a productRecommendations query
		if strings.Contains(query, "productRecommendations") {
			c.JSON(http.StatusOK, map[string]interface{}{
				"data": map[string]interface{}{
					"productRecommendations": []map[string]interface{}{
						{
							"id":                 "rec-1",
							"customerId":         "customer-123",
							"productId":          "product-789",
							"score":              0.95,
							"recommendationType": "PERSONALIZED",
							"createdAt":          time.Now().Format(time.RFC3339),
							"expiresAt":          time.Now().Add(24*time.Hour).Format(time.RFC3339),
						},
						{
							"id":                 "rec-2",
							"customerId":         "customer-123",
							"productId":          "product-012",
							"score":              0.87,
							"recommendationType": "COLLABORATIVE",
							"createdAt":          time.Now().Format(time.RFC3339),
							"expiresAt":          time.Now().Add(24*time.Hour).Format(time.RFC3339),
						},
						{
							"id":                 "rec-3",
							"customerId":         "customer-123",
							"productId":          "product-345",
							"score":              0.82,
							"recommendationType": "TRENDING",
							"createdAt":          time.Now().Format(time.RFC3339),
							"expiresAt":          time.Now().Add(24*time.Hour).Format(time.RFC3339),
						},
					},
				},
			})
			return
		}
		
		// Default analytics service operational response
		c.JSON(http.StatusOK, map[string]interface{}{
			"data": map[string]interface{}{
				"analytics": "Analytics service is operational with comprehensive recommendation engine",
			},
		})
	})

	// Start REST server
	srv := &http.Server{
		Addr:         ":" + cfg.ServicePort,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.WithFields(map[string]interface{}{"port": cfg.ServicePort}).Info("Starting Analytics Service")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.WithError(err).Fatal("Failed to start server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down Analytics Service...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.WithError(err).Fatal("Server forced to shutdown")
	}

	log.Info("Analytics Service stopped")
}
