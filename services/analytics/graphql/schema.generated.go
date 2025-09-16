package graphql

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
	"github.com/vektah/gqlparser/v2/ast"
)

// region    ************************** generated!.gotpl **************************

// NewExecutableSchema creates an ExecutableSchema from the ResolverRoot interface.
func NewExecutableSchema(cfg Config) graphql.ExecutableSchema {
	return &executableSchema{
		resolvers:  cfg.Resolver,
		directives: cfg.Directives,
		complexity: cfg.Complexity,
	}
}

type Config struct {
	Resolver   ResolverRoot
	Directives DirectiveRoot
	Complexity ComplexityRoot
}

type ResolverRoot interface {
	CustomerBehavior() CustomerBehaviorResolver
	ProductRecommendation() ProductRecommendationResolver
	Query() QueryResolver
	Mutation() MutationResolver
}

type DirectiveRoot struct {
}

type ComplexityRoot struct {
	CustomerBehavior struct {
		Action             func(childComplexity int) int
		CreatedAt          func(childComplexity int) int
		CustomerID         func(childComplexity int) int
		EntityID           func(childComplexity int) int
		EntityType         func(childComplexity int) int
		ExpiresAt          func(childComplexity int) int
		ID                 func(childComplexity int) int
		Product            func(childComplexity int) int
		ProductID          func(childComplexity int) int
		RecommendationType func(childComplexity int) int
		Score              func(childComplexity int) int
		SessionID          func(childComplexity int) int
	}

	ProductRecommendation struct {
		Action             func(childComplexity int) int
		CreatedAt          func(childComplexity int) int
		CustomerID         func(childComplexity int) int
		EntityID           func(childComplexity int) int
		EntityType         func(childComplexity int) int
		ExpiresAt          func(childComplexity int) int
		ID                 func(childComplexity int) int
		Product            func(childComplexity int) int
		ProductID          func(childComplexity int) int
		RecommendationType func(childComplexity int) int
		Score              func(childComplexity int) int
		SessionID          func(childComplexity int) int
	}

	Query struct {
		CustomerBehaviors      func(childComplexity int, customerID string, limit *int) int
		ProductRecommendations func(childComplexity int, customerID string, limit *int) int
	}

	Mutation struct {
		GenerateProductRecommendations func(childComplexity int, input GenerateRecommendationsInput) int
		TrackCustomerBehavior          func(childComplexity int, input TrackCustomerBehaviorInput) int
	}
}

type executableSchema struct {
	resolvers  ResolverRoot
	directives DirectiveRoot
	complexity ComplexityRoot
}

func (e *executableSchema) Schema() *ast.Schema {
	return parsedSchema
}

func (e *executableSchema) Complexity(typeName, field string, childComplexity int, rawArgs map[string]interface{}) (int, bool) {
	// This is a placeholder implementation
	return 1, true
}

func (e *executableSchema) Exec(ctx context.Context) graphql.ResponseHandler {
	// This is a placeholder implementation
	return nil
}

// endregion ************************** generated!.gotpl **************************

// region    ************************** directives.gotpl **************************

// endregion ************************** directives.gotpl **************************

// region    ************************** models.gotpl **************************

// endregion ************************** models.gotpl **************************

// region    ************************** schema.gotpl **************************

// endregion ************************** schema.gotpl **************************

// region    ************************** types.gotpl **************************

func MarshalAny(v interface{}) graphql.Marshaler {
	return graphql.MarshalAny(v)
}

func UnmarshalAny(v interface{}) (interface{}, error) {
	return v, nil
}

// endregion ************************** types.gotpl **************************

// region    ************************** generated!.gotpl **************************

// endregion ************************** generated!.gotpl **************************
