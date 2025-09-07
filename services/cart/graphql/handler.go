package graphql

import (
	"context"
	"fmt"
	"net/http"

	"unified-commerce/services/cart/service"
	"unified-commerce/services/shared/logger"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
)

// NewGraphQLHandler creates a new GraphQL HTTP handler
func NewGraphQLHandler(cartService *service.CartService, logger *logger.Logger) http.Handler {
	// Create a simple executable schema
	schema := NewExecutableSchema(Config{
		Resolvers: NewResolver(cartService, logger),
	})

	// Create the GraphQL server
	srv := handler.NewDefaultServer(schema)

	// Add recovery handler
	srv.SetRecoverFunc(func(ctx context.Context, err interface{}) error {
		logger.WithField("panic", err).Error("GraphQL panic recovered")
		return fmt.Errorf("internal server error")
	})

	return srv
}

// NewPlaygroundHandler creates a new GraphQL playground handler
func NewPlaygroundHandler() http.Handler {
	return playground.Handler("GraphQL Playground", "/graphql")
}
