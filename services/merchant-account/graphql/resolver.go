package graphql

import (
	"unified-commerce/services/merchant-account/service"
	"unified-commerce/services/shared/logger"
)

// Resolver provides resolution context for GraphQL resolvers.
// It holds dependencies like the IdentityService, making them available to resolver functions.
type Resolver struct {
	MerchantService *service.MerchantService
	Logger          *logger.Logger
}

// NewResolver creates a new Resolver instance.
func NewResolver(merchantService *service.MerchantService, logger *logger.Logger) *Resolver {
	return &Resolver{
		MerchantService: merchantService,
		Logger:          logger,
	}
}
