package graphql

import (
	"unified-commerce/services/promotions/service"
	"unified-commerce/services/shared/logger"
)

// Resolver provides resolution context for GraphQL resolvers.
// It holds dependencies like the IdentityService, making them available to resolver functions.
type Resolver struct {
	PromotionsService *service.PromotionsService
	Logger            *logger.Logger
}

// NewResolver creates a new Resolver instance.
func NewResolver(promotionsService *service.PromotionsService, logger *logger.Logger) *Resolver {
	return &Resolver{
		PromotionsService: promotionsService,
		Logger:            logger,
	}
}
