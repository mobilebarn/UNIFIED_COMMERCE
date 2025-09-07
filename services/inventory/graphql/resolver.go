package graphql

import (
	"unified-commerce/services/inventory/service"
	"unified-commerce/services/shared/logger"
)

// Resolver provides resolution context for GraphQL resolvers.
// It holds dependencies like the IdentityService, making them available to resolver functions.
type Resolver struct {
	InventoryService *service.InventoryService
	Logger           *logger.Logger
}

// NewResolver creates a new Resolver instance.
func NewResolver(inventoryService *service.InventoryService, log *logger.Logger) *Resolver {
	return &Resolver{
		InventoryService: inventoryService,
		Logger:           log,
	}
}
