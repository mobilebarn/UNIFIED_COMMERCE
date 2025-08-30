package graphql

import (
	"unified-commerce/services/cart/service"
	"unified-commerce/services/shared/logger"
)

// Resolver is the root resolver
type Resolver struct {
	cartService *service.CartService
	logger      *logger.Logger
}

// NewResolver creates a new resolver
func NewResolver(cartService *service.CartService, logger *logger.Logger) *Resolver {
	return &Resolver{
		cartService: cartService,
		logger:      logger,
	}
}
