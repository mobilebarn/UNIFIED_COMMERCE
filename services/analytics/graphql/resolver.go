package graphql

import (
	"context"
	"retail-os/services/analytics/internal/service"
	"retail-os/services/shared/logger"
)

// Resolver provides resolution context for GraphQL resolvers.
// It holds dependencies like the AnalyticsService, making them available to resolver functions.
type Resolver struct {
	AnalyticsService *service.AnalyticsService
	Logger           *logger.Logger
}

// NewResolver creates a new Resolver instance.
func NewResolver(analyticsService *service.AnalyticsService, logger *logger.Logger) *Resolver {
	return &Resolver{
		AnalyticsService: analyticsService,
		Logger:           logger,
	}
}

// ResolverRoot interface
type ResolverRoot interface {
	CustomerBehavior() CustomerBehaviorResolver
	ProductRecommendation() ProductRecommendationResolver
	Query() QueryResolver
	Mutation() MutationResolver
}

// CustomerBehaviorResolver interface
type CustomerBehaviorResolver interface {
	// Add any field resolvers for CustomerBehavior here
}

// ProductRecommendationResolver interface
type ProductRecommendationResolver interface {
	// Add any field resolvers for ProductRecommendation here
}

// QueryResolver interface
type QueryResolver interface {
	CustomerBehaviors(ctx context.Context, customerID string, limit *int) ([]*CustomerBehavior, error)
	ProductRecommendations(ctx context.Context, customerID string, limit *int) ([]*ProductRecommendation, error)
}

// MutationResolver interface
type MutationResolver interface {
	TrackCustomerBehavior(ctx context.Context, input TrackCustomerBehaviorInput) (*CustomerBehavior, error)
	GenerateProductRecommendations(ctx context.Context, input GenerateRecommendationsInput) ([]*ProductRecommendation, error)
}
