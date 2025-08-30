package graphql

import (
	"fmt"
	"net/http"
	"unified-commerce/services/order/service"
	"unified-commerce/services/shared/logger"
)

// NewGraphQLHandler creates a simple GraphQL HTTP handler for the order service
func NewGraphQLHandler(orderService *service.OrderService, logger *logger.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if r.Method == "GET" {
			// Return GraphQL schema for introspection
			schema := `{
				"data": {
					"__schema": {
						"types": [
							{
								"name": "Order",
								"kind": "OBJECT",
								"fields": [
									{"name": "id", "type": {"name": "ID"}},
									{"name": "orderNumber", "type": {"name": "String"}},
									{"name": "customerId", "type": {"name": "ID"}},
									{"name": "merchantId", "type": {"name": "ID"}},
									{"name": "status", "type": {"name": "OrderStatus"}},
									{"name": "totalPrice", "type": {"name": "Float"}}
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
						"sdl": "extend type User @key(fields: \"id\") { id: ID! @external } type Order @key(fields: \"id\") { id: ID! orderNumber: String! customerId: ID merchantId: ID! status: OrderStatus! totalPrice: Float! }"
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
  <title>Order Service GraphQL Playground</title>
  <style>
    body { margin: 0; padding: 20px; font-family: Arial, sans-serif; }
    h1 { color: #333; }
  </style>
</head>
<body>
  <h1>Order Service GraphQL Playground</h1>
  <p>GraphQL endpoint available at <code>/graphql</code></p>
  <p>This is a development placeholder. Full GraphQL implementation pending.</p>
</body>
</html>`
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, playground)
	})
}
