# Analytics Service

The Analytics Service is a microservice responsible for collecting, processing, and analyzing customer behavior data, generating product recommendations, and providing business intelligence for the Retail OS Platform.

## Features

- Customer behavior tracking (views, clicks, purchases, etc.)
- Product recommendation engine
- Customer segmentation
- Business metrics tracking
- Analytics reporting

## Architecture

The service follows a clean architecture pattern with the following components:

- **Handlers**: HTTP request handlers
- **Service**: Business logic layer
- **Repository**: Data access layer
- **Models**: Data structures

## Technologies

- Go 1.21
- Gin Framework for HTTP handling
- MongoDB for analytics data storage
- PostgreSQL for relational data
- GraphQL for API interface

## API Endpoints

### REST Endpoints

- `POST /api/v1/analytics/behavior` - Track customer behavior
- `GET /api/v1/analytics/behavior/:customerID` - Get customer behaviors
- `POST /api/v1/analytics/recommendations/generate/:customerID` - Generate product recommendations
- `GET /api/v1/analytics/recommendations/:customerID` - Get product recommendations
- `POST /api/v1/analytics/segments` - Create customer segment
- `GET /api/v1/analytics/segments/:id` - Get customer segment
- `POST /api/v1/analytics/metrics` - Track business metric
- `GET /api/v1/analytics/metrics/:name` - Get business metrics

### GraphQL Schema

The service exposes a GraphQL API for integration with the platform's GraphQL Federation Gateway.

## Deployment

The service can be deployed using Docker:

```bash
docker build -t analytics-service .
docker run -p 8009:8009 analytics-service
```

## Configuration

The service is configured using environment variables. See `.env` for available configuration options.

## Monitoring

The service exposes health check endpoints and integrates with the platform's observability stack.