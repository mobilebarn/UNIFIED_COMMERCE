package graphql

import (
	"unified-commerce/services/product-catalog/service"
	"unified-commerce/services/shared/logger"
)

// Resolver provides resolution context for GraphQL resolvers.
// It holds dependencies like the IdentityService, making them available to resolver functions.
type Resolver struct {
	ProductService *service.ProductService
	Logger         *logger.Logger
}

// NewResolver creates a new Resolver instance.
func NewResolver(productService *service.ProductService, logger *logger.Logger) *Resolver {
	return &Resolver{
		ProductService: productService,
		Logger:         logger,
	}
}
