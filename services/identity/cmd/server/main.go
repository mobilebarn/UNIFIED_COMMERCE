package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"

	"unified-commerce/services/identity/handlers"
	"unified-commerce/services/identity/models"
	"unified-commerce/services/identity/repository"
	identityService "unified-commerce/services/identity/service"
	"unified-commerce/services/shared/service"
)

func main() {
	// Create base service with PostgreSQL and Redis
	baseService, err := service.NewBaseService(service.ServiceOptions{
		Name:        "identity",
		UsePostgres: true,
		UseRedis:    true,
		UseMongoDB:  false,
		CustomRoutes: func(router *gin.Engine) {
			setupRoutes(router, baseService)
		},
	})
	if err != nil {
		log.Fatalf("Failed to create base service: %v", err)
	}

	// Run database migrations
	repo := repository.NewRepository(baseService.PostgresDB)
	if err := repo.Migrate(); err != nil {
		baseService.Logger.WithError(err).Fatal("Failed to run database migrations")
	}

	// Seed initial data
	if err := seedInitialData(repo); err != nil {
		baseService.Logger.WithError(err).Warn("Failed to seed initial data")
	}

	baseService.Logger.Info("Identity service started successfully")

	// Start the service
	if err := baseService.Start(); err != nil {
		baseService.Logger.WithError(err).Fatal("Failed to start service")
	}
}

// setupRoutes configures the HTTP routes for the identity service
func setupRoutes(router *gin.Engine, baseService *service.BaseService) {
	// Create repository
	repo := repository.NewRepository(baseService.PostgresDB)

	// Create service
	identityService := identityService.NewIdentityService(repo, baseService.Logger, baseService.Config.JWTSecret)

	// Create handler
	handler := handlers.NewIdentityHandler(identityService, baseService.Logger)

	// Register routes
	handler.RegisterRoutes(router)
}

// seedInitialData seeds the database with initial roles and permissions
func seedInitialData(repo *repository.Repository) error {
	ctx := context.Background()

	// Create default roles if they don't exist
	defaultRoles := []models.Role{
		{
			Name:        "super_admin",
			DisplayName: "Super Administrator",
			Description: "Full system access",
			IsActive:    true,
		},
		{
			Name:        "admin",
			DisplayName: "Administrator",
			Description: "Administrative access",
			IsActive:    true,
		},
		{
			Name:        "merchant",
			DisplayName: "Merchant",
			Description: "Merchant account owner",
			IsActive:    true,
		},
		{
			Name:        "staff",
			DisplayName: "Staff",
			Description: "Merchant staff member",
			IsActive:    true,
		},
		{
			Name:        "customer",
			DisplayName: "Customer",
			Description: "Customer account",
			IsActive:    true,
		},
	}

	for _, role := range defaultRoles {
		// Check if role already exists
		_, err := repo.Role.GetByName(ctx, role.Name)
		if err != nil {
			// Role doesn't exist, create it
			if err := repo.Role.Create(ctx, &role); err != nil {
				return err
			}
		}
	}

	return nil
}
