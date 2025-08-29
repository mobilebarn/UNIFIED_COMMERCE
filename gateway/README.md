# GraphQL Federation Gateway

This is the GraphQL Federation Gateway that unifies all microservices into a single API endpoint.

## Architecture

The gateway uses Apollo Federation to combine schemas from multiple microservices:
- Identity Service
- Merchant Account Service
- Product Catalog Service
- Inventory Service
- Order Service
- Cart & Checkout Service
- Payments Service
- Promotions Service

## Services

Each service exposes a GraphQL schema that can be federated:
- Services implement the `@key` directive to define their primary keys
- Services can extend types from other services using the `extend` keyword
- The gateway automatically stitches schemas together

## Development

```bash
# Install dependencies
npm install

# Start development server
npm run dev

# Build for production
npm run build
```