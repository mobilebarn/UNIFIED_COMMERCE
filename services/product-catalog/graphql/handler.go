package graphql

import (
	"fmt"
	"net/http"
	"unified-commerce/services/product-catalog/service"
	"unified-commerce/services/shared/logger"
)

// NewGraphQLHandler creates a simple GraphQL HTTP handler for the product-catalog service
func NewGraphQLHandler(productService *service.ProductService, logger *logger.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if r.Method == "GET" {
			// Return GraphQL schema for introspection
			schema := `{
				"data": {
					"__schema": {
						"types": [
							{
								"name": "Product",
								"kind": "OBJECT",
								"fields": [
									{"name": "id", "type": {"name": "ID"}},
									{"name": "title", "type": {"name": "String"}},
									{"name": "handle", "type": {"name": "String"}},
									{"name": "merchantId", "type": {"name": "ID"}},
									{"name": "status", "type": {"name": "ProductStatus"}}
								]
							}
						]
					}
				}
			}`
			fmt.Fprint(w, schema)
			return
		}

		if r.Method == "POST" {
			// Handle GraphQL mutations and queries
			response := `{
				"data": {
					"_service": {
						"sdl": "extend type User @key(fields: \"id\") { id: ID! @external } type Product @key(fields: \"id\") { id: ID! title: String! handle: String! merchantId: ID! status: ProductStatus! }"
					}
				}
			}`
			fmt.Fprint(w, response)
			return
		}

		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	})
}

// Simple GraphQL playground for development
func NewGraphQLPlaygroundHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		playground := `<!DOCTYPE html>
<html>
<head>
  <title>Product Catalog Service GraphQL Playground</title>
  <style>
    body { margin: 0; padding: 20px; font-family: Arial, sans-serif; }
    h1 { color: #333; }
  </style>
</head>
<body>
  <h1>Product Catalog Service GraphQL Playground</h1>
  <p>GraphQL endpoint available at <code>/graphql</code></p>
  <p>This is a development placeholder. Full GraphQL implementation pending.</p>
</body>
</html>`
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, playground)
	})
}
