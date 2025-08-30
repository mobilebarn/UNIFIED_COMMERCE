package graphql

import (
	"fmt"
	"net/http"
	"unified-commerce/services/inventory/service"
	"unified-commerce/services/shared/logger"
)

// NewGraphQLHandler creates a simple GraphQL HTTP handler for the inventory service
func NewGraphQLHandler(inventoryService *service.InventoryService, logger *logger.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if r.Method == "GET" {
			// Return GraphQL schema for introspection
			schema := `{
				"data": {
					"__schema": {
						"types": [
							{
								"name": "InventoryItem",
								"kind": "OBJECT",
								"fields": [
									{"name": "id", "type": {"name": "ID"}},
									{"name": "productId", "type": {"name": "ID"}},
									{"name": "locationId", "type": {"name": "ID"}},
									{"name": "quantityAvailable", "type": {"name": "Int"}},
									{"name": "sku", "type": {"name": "String"}}
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
						"sdl": "extend type User @key(fields: \"id\") { id: ID! @external } type InventoryItem @key(fields: \"id\") { id: ID! productId: ID! locationId: ID! quantityAvailable: Int! sku: String! }"
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
  <title>Inventory Service GraphQL Playground</title>
  <style>
    body { margin: 0; padding: 20px; font-family: Arial, sans-serif; }
    h1 { color: #333; }
  </style>
</head>
<body>
  <h1>Inventory Service GraphQL Playground</h1>
  <p>GraphQL endpoint available at <code>/graphql</code></p>
  <p>This is a development placeholder. Full GraphQL implementation pending.</p>
</body>
</html>`
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, playground)
	})
}
