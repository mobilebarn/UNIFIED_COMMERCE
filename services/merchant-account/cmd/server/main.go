package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"

	"unified-commerce/services/merchant-account/graphql"
	"unified-commerce/services/merchant-account/handlers"
	"unified-commerce/services/merchant-account/models"
	"unified-commerce/services/merchant-account/repository"
	merchantService "unified-commerce/services/merchant-account/service"
	"unified-commerce/services/shared/service"
)

func main() {
	// Create base service with PostgreSQL and Redis
	baseService, err := service.NewBaseService(service.ServiceOptions{
		Name:        "merchant-account",
		UsePostgres: true,
		UseRedis:    false,
		UseMongoDB:  false,
	})
	if err != nil {
		log.Fatalf("Failed to create base service: %v", err)
	}

	// Setup custom routes
	setupRoutes(baseService.Router, baseService)

	// Run database migrations
	repo := repository.NewRepository(baseService.PostgresDB)
	if err := repo.Migrate(); err != nil {
		baseService.Logger.WithError(err).Fatal("Failed to run database migrations")
	}

	// Seed initial data
	if err := seedInitialData(repo); err != nil {
		baseService.Logger.WithError(err).Warn("Failed to seed initial data")
	}

	baseService.Logger.Info("Merchant Account service started successfully")

	// Start the service
	if err := baseService.Start(); err != nil {
		baseService.Logger.WithError(err).Fatal("Failed to start service")
	}
}

// setupRoutes configures the HTTP routes for the merchant account service
func setupRoutes(router *gin.Engine, baseService *service.BaseService) {
	// Create repository
	repo := repository.NewRepository(baseService.PostgresDB)

	// Create service
	merchantService := merchantService.NewMerchantService(repo, baseService.Logger)

	// Create handlers
	handler := handlers.NewMerchantHandler(merchantService, baseService.Logger)
	graphqlHandler := graphql.NewGraphQLHandler(merchantService, baseService.Logger)

	// Register REST routes
	handler.RegisterRoutes(router)

	// Register GraphQL endpoint
	router.POST("/graphql", gin.WrapH(graphqlHandler))

}

// seedInitialData seeds the database with initial subscription plans
func seedInitialData(repo *repository.Repository) error {
	ctx := context.Background()

	// Create default subscription plans
	defaultPlans := []models.Plan{
		{
			Name:         "starter",
			DisplayName:  "Starter Plan",
			Description:  "Perfect for small businesses just getting started",
			PlanType:     "basic",
			MonthlyPrice: 29.00,
			AnnualPrice:  290.00,
			Currency:     "USD",
			TrialDays:    14,
			IsActive:     true,
			Features: map[string]interface{}{
				"products":         1000,
				"orders_per_month": 500,
				"storage_gb":       10,
				"staff_accounts":   2,
				"support":          "email",
				"analytics":        "basic",
				"api_calls":        10000,
				"payment_gateways": []string{"stripe", "paypal"},
			},
			Limits: map[string]int{
				"products":        1000,
				"orders_monthly":  500,
				"storage_gb":      10,
				"staff_accounts":  2,
				"api_calls_daily": 1000,
				"webhooks":        5,
			},
		},
		{
			Name:         "professional",
			DisplayName:  "Professional Plan",
			Description:  "For growing businesses with advanced needs",
			PlanType:     "professional",
			MonthlyPrice: 79.00,
			AnnualPrice:  790.00,
			Currency:     "USD",
			TrialDays:    14,
			IsActive:     true,
			Features: map[string]interface{}{
				"products":           10000,
				"orders_per_month":   2000,
				"storage_gb":         100,
				"staff_accounts":     10,
				"support":            "priority",
				"analytics":          "advanced",
				"api_calls":          100000,
				"payment_gateways":   []string{"stripe", "paypal", "square", "authorize_net"},
				"multi_location":     true,
				"advanced_reporting": true,
				"bulk_operations":    true,
			},
			Limits: map[string]int{
				"products":        10000,
				"orders_monthly":  2000,
				"storage_gb":      100,
				"staff_accounts":  10,
				"api_calls_daily": 10000,
				"webhooks":        20,
				"locations":       5,
			},
		},
		{
			Name:         "enterprise",
			DisplayName:  "Enterprise Plan",
			Description:  "For large businesses with complex requirements",
			PlanType:     "enterprise",
			MonthlyPrice: 299.00,
			AnnualPrice:  2990.00,
			Currency:     "USD",
			TrialDays:    30,
			IsActive:     true,
			Features: map[string]interface{}{
				"products":            "unlimited",
				"orders_per_month":    "unlimited",
				"storage_gb":          "unlimited",
				"staff_accounts":      "unlimited",
				"support":             "dedicated",
				"analytics":           "enterprise",
				"api_calls":           "unlimited",
				"payment_gateways":    []string{"stripe", "paypal", "square", "authorize_net", "custom"},
				"multi_location":      true,
				"advanced_reporting":  true,
				"bulk_operations":     true,
				"white_label":         true,
				"custom_integrations": true,
				"sla":                 "99.9%",
			},
			Limits: map[string]int{
				"products":        -1, // -1 means unlimited
				"orders_monthly":  -1,
				"storage_gb":      -1,
				"staff_accounts":  -1,
				"api_calls_daily": -1,
				"webhooks":        -1,
				"locations":       -1,
			},
		},
	}

	for _, plan := range defaultPlans {
		// Check if plan already exists
		_, err := repo.Plan.GetByName(ctx, plan.Name)
		if err != nil {
			// Plan doesn't exist, create it
			if err := repo.Plan.Create(ctx, &plan); err != nil {
				return err
			}
		}
	}

	return nil
}
