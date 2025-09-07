package graphql

import (
	"unified-commerce/services/cart/service"
	"unified-commerce/services/shared/logger"
)

// Resolver provides resolution context for GraphQL resolvers.
// It holds dependencies like the CartService, making them available to resolver functions.
type Resolver struct {
	CartService *service.CartService
	Logger      *logger.Logger
}

// NewResolver creates a new Resolver instance.
func NewResolver(cartService *service.CartService, logger *logger.Logger) *Resolver {
	return &Resolver{
		CartService: cartService,
		Logger:      logger,
	}
}
