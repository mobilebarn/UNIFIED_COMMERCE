package graphql

import (
	"unified-commerce/services/identity/service"
	"unified-commerce/services/shared/logger"
)

// Resolver provides resolution context for GraphQL resolvers.
// It holds dependencies like the IdentityService, making them available to resolver functions.
type Resolver struct {
	IdentityService *service.IdentityService
	Logger          *logger.Logger
}

// NewResolver creates a new Resolver instance.
func NewResolver(identityService *service.IdentityService, logger *logger.Logger) *Resolver {
	return &Resolver{
		IdentityService: identityService,
		Logger:          logger,
	}
}
