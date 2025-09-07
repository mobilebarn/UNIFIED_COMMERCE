package graphql

import (
	"unified-commerce/services/payment/service"
	"unified-commerce/services/shared/logger"
)

// Resolver provides resolution context for GraphQL resolvers.
// It holds dependencies like the PaymentService, making them available to resolver functions.
type Resolver struct {
	PaymentService *service.PaymentService
	Logger         *logger.Logger
}

// NewResolver creates a new Resolver instance.
func NewResolver(PaymentService *service.PaymentService, logger *logger.Logger) *Resolver {
	return &Resolver{
		PaymentService: PaymentService,
		Logger:         logger,
	}
}
