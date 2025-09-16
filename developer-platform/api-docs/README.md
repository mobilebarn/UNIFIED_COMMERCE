# Retail OS Developer Platform API Documentation

## Overview

Welcome to the Retail OS API documentation. This comprehensive guide provides everything developers need to integrate with our platform, including detailed information about our GraphQL API, SDKs, authentication methods, and best practices.

Retail OS provides a unified API for all e-commerce operations through GraphQL Federation, allowing developers to access all services through a single endpoint while maintaining the benefits of a microservices architecture.

## API Endpoints

### GraphQL Federation Gateway
- **Production**: `https://api.retail-os.com/graphql`
- **Development**: `http://localhost:4000/graphql`

### Health Check
- **Endpoint**: `/health`
- **Method**: GET
- **Description**: Returns the health status of the GraphQL gateway and all connected services

## Authentication

All API requests require authentication using JWT tokens. Tokens can be obtained through the authentication flow:

```graphql
mutation {
  login(input: {
    email: "user@example.com"
    password: "password123"
  }) {
    user {
      id
      email
      firstName
      lastName
    }
    accessToken
    refreshToken
    expiresIn
  }
}
```

Include the access token in the Authorization header for subsequent requests:

```
Authorization: Bearer <access_token>
```

## Core Services

Retail OS consists of 8 core microservices, all accessible through the GraphQL Federation Gateway:

1. **Identity Service** - User authentication and authorization
2. **Product Catalog Service** - Product, category, and inventory management
3. **Order Service** - Order processing and fulfillment
4. **Payment Service** - Payment processing and transactions
5. **Cart Service** - Shopping cart management
6. **Inventory Service** - Real-time inventory tracking
7. **Promotions Service** - Discounts, coupons, and promotional campaigns
8. **Merchant Account Service** - Merchant profile and account management

## Rate Limiting

To ensure fair usage and system stability, the API implements rate limiting:

- **Anonymous requests**: 100 requests per hour
- **Authenticated requests**: 1,000 requests per hour
- **Enterprise accounts**: Custom rate limits based on plan

Exceeding rate limits will result in a 429 (Too Many Requests) response.

## Error Handling

The API uses standard HTTP status codes and provides detailed error information in the response body:

```json
{
  "errors": [
    {
      "message": "User not found",
      "locations": [
        {
          "line": 2,
          "column": 3
        }
      ],
      "path": [
        "user"
      ],
      "extensions": {
        "code": "USER_NOT_FOUND"
      }
    }
  ]
}
```

## Versioning

The API uses semantic versioning. Breaking changes are introduced through new major versions, while backward-compatible changes are added to existing versions.

## SDKs

We provide official SDKs for the following programming languages:
- JavaScript/TypeScript
- Python
- Java
- Go
- PHP

## Support

For API-related questions and support:
- **Email**: developers@retail-os.com
- **Documentation**: https://docs.retail-os.com
- **Status Page**: https://status.retail-os.com

## Getting Started

1. Register for a developer account at https://developer.retail-os.com
2. Create an application to obtain your API credentials
3. Choose an authentication method (OAuth 2.0 or API keys)
4. Start integrating with our API using the documentation below

## Changelog

For information about recent changes to the API, see our [changelog](CHANGELOG.md).