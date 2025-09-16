package graphql

import (
	"context"
	"log"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
)

// Server represents the GraphQL server
type Server struct {
	srv *handler.Server
}

// NewServer creates a new GraphQL server
func NewServer(resolver ResolverRoot) *Server {
	// Create a new server with the schema
	config := Config{Resolvers: resolver}
	executableSchema := NewExecutableSchema(config)
	srv := handler.NewDefaultServer(executableSchema)

	return &Server{
		srv: srv,
	}
}

// Handler returns the GraphQL handler
func (s *Server) Handler() http.Handler {
	return s.srv
}

// Playground returns the GraphQL playground handler
func (s *Server) Playground() http.Handler {
	return playground.Handler("GraphQL playground", "/graphql")
}

// Start starts the GraphQL server
func (s *Server) Start(port string) {
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", s.srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

// Resolver implementation
func (r *Resolver) CustomerBehavior() CustomerBehaviorResolver {
	return &customerBehaviorResolver{r}
}

func (r *Resolver) ProductRecommendation() ProductRecommendationResolver {
	return &productRecommendationResolver{r}
}

func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}

type customerBehaviorResolver struct{ *Resolver }
type productRecommendationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }

// Implement the resolver methods
func (r *queryResolver) CustomerBehaviors(ctx context.Context, customerID string, limit *int) ([]*CustomerBehavior, error) {
	// This will be implemented properly later
	return nil, nil
}

func (r *queryResolver) ProductRecommendations(ctx context.Context, customerID string, limit *int) ([]*ProductRecommendation, error) {
	// This will be implemented properly later
	return nil, nil
}

func (r *mutationResolver) TrackCustomerBehavior(ctx context.Context, input TrackCustomerBehaviorInput) (*CustomerBehavior, error) {
	// This will be implemented properly later
	return nil, nil
}

func (r *mutationResolver) GenerateProductRecommendations(ctx context.Context, input GenerateRecommendationsInput) ([]*ProductRecommendation, error) {
	// This will be implemented properly later
	return nil, nil
}
