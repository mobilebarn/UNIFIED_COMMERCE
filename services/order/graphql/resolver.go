package graphql

import (
	"unified-commerce/services/order/service"
	"unified-commerce/services/shared/logger"
)

// Resolver provides resolution context for GraphQL resolvers.
// It holds dependencies like the IdentityService, making them available to resolver functions.
type Resolver struct {
	OrderService *service.OrderService
	Logger       *logger.Logger
}

// NewResolver creates a new Resolver instance.
func NewResolver(OrderService *service.OrderService, logger *logger.Logger) *Resolver {
	return &Resolver{
		OrderService: OrderService,
		Logger:       logger,
	}
}
