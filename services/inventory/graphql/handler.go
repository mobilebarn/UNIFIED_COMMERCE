package graphql

import (
	"context"
	"fmt"
	"net/http"

	"unified-commerce/services/inventory/service"
	"unified-commerce/services/shared/logger"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
)

// NewGraphQLHandler creates a new GraphQL HTTP handler
func NewGraphQLHandler(inventoryService *service.InventoryService, log *logger.Logger) http.Handler {
	// Create a simple executable schema
	schema := NewExecutableSchema(Config{
		Resolvers: NewResolver(inventoryService, log),
	})

	// Create the GraphQL server
	srv := handler.NewDefaultServer(schema)

	// Add recovery handler
	srv.SetRecoverFunc(func(ctx context.Context, err interface{}) error {
		log.WithField("panic", err).Error("GraphQL panic recovered")
		return fmt.Errorf("internal server error")
	})

	return srv
}

// NewPlaygroundHandler creates a new GraphQL playground handler
func NewPlaygroundHandler() http.Handler {
	return playground.Handler("GraphQL Playground", "/graphql")
}
