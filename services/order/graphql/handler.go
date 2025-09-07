package graphql

import (
	"context"
	"fmt"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"unified-commerce/services/order/service"
	"unified-commerce/services/shared/logger"
)

// NewGraphQLHandler creates a new GraphQL HTTP handler
func NewGraphQLHandler(orderService *service.OrderService, logger *logger.Logger) http.Handler {
	// Create a simple executable schema
	schema := NewExecutableSchema(Config{
		Resolvers: NewResolver(orderService, logger),
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

