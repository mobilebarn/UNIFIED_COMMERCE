# Payment Service API Testing Guide

## Overview

This guide provides comprehensive API testing instructions for the Payment Service in our Unified Commerce Platform. The service handles payment processing, refunds, settlements, and gateway management.

## Prerequisites

1. **Infrastructure Services Running:**
   ```bash
   docker-compose up -d postgres redis
   ```

2. **Environment Setup:**
   - Copy `.env.example` to `.env` in the payment service directory
   - Update database connection strings as needed
   - Ensure the Payment Service is running on port 8007

3. **Dependencies:**
   - Identity Service running on port 8001 (for authentication)
   - Order Service running on port 8005 (for order integration)

## Authentication Flow

All protected endpoints require JWT authentication. Here's the complete authentication flow:

### 1. User Registration & Login

```bash
# Register a new user (if needed)
curl -X POST http://localhost:8001/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "merchant@example.com",
    "username": "merchantuser",
    "password": "SecurePass123!",
    "first_name": "John",
    "last_name": "Merchant",
    "phone": "+1234567890"
  }'

# Login to get JWT token
curl -X POST http://localhost:8001/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "merchant@example.com",
    "password": "SecurePass123!"
  }'
```

**Save the `access_token` from the login response for authenticated requests.**

## Service Testing Workflows

### Payment Management

#### Health Check
```bash
curl -X GET http://localhost:8007/health
```

#### Create a New Payment
```bash
curl -X POST http://localhost:8007/api/v1/payments \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "order_id": "ORDER_UUID_HERE",
    "merchant_id": "MERCHANT_UUID_HERE",
    "payment_method_id": "PAYMENT_METHOD_UUID_HERE",
    "amount": 99.99,
    "currency": "USD",
    "description": "Order payment"
  }'
```

#### Get Payment by ID
```bash
curl -X GET http://localhost:8007/api/v1/payments/PAYMENT_ID_HERE \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

#### Get Payment by Order ID
```bash
curl -X GET http://localhost:8007/api/v1/payments/order/ORDER_ID_HERE \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

#### Get Payments by Merchant
```bash
curl -X GET "http://localhost:8007/api/v1/payments/merchant/MERCHANT_ID_HERE?page=1&per_page=10&status=captured" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

#### Process Payment
```bash
curl -X POST http://localhost:8007/api/v1/payments/PAYMENT_ID_HERE/process \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

#### Cancel Payment
```bash
curl -X POST http://localhost:8007/api/v1/payments/PAYMENT_ID_HERE/cancel \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### Payment Method Management

#### Create Payment Method
```bash
curl -X POST http://localhost:8007/api/v1/payment-methods \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "customer_id": "CUSTOMER_UUID_HERE",
    "type": "credit_card",
    "provider": "stripe",
    "token": "pm_example123",
    "last4": "4242",
    "expiry_month": 12,
    "expiry_year": 2025,
    "brand": "Visa",
    "name": "John Customer",
    "is_default": true
  }'
```

#### Get Payment Method by ID
```bash
curl -X GET http://localhost:8007/api/v1/payment-methods/PAYMENT_METHOD_ID_HERE \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

#### Get Payment Methods by Customer
```bash
curl -X GET http://localhost:8007/api/v1/payment-methods/customer/CUSTOMER_ID_HERE \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

#### Update Payment Method
```bash
curl -X PUT http://localhost:8007/api/v1/payment-methods/PAYMENT_METHOD_ID_HERE \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "is_default": true
  }'
```

#### Delete Payment Method
```bash
curl -X DELETE http://localhost:8007/api/v1/payment-methods/PAYMENT_METHOD_ID_HERE \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### Gateway Management

#### Create Payment Gateway
```bash
curl -X POST http://localhost:8007/api/v1/gateways \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Stripe Gateway",
    "provider": "stripe",
    "is_enabled": true,
    "is_sandbox": true,
    "credentials": {
      "secret_key": "sk_test_example123",
      "publishable_key": "pk_test_example456"
    },
    "settings": {
      "capture_method": "automatic",
      "statement_descriptor": "UNIFIED COMMERCE"
    }
  }'
```

#### Get Gateway by ID
```bash
curl -X GET http://localhost:8007/api/v1/gateways/GATEWAY_ID_HERE \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

#### Get Gateway by Name
```bash
curl -X GET http://localhost:8007/api/v1/gateways/name/Stripe Gateway \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

#### Get All Gateways
```bash
curl -X GET http://localhost:8007/api/v1/gateways \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

#### Update Gateway
```bash
curl -X PUT http://localhost:8007/api/v1/gateways/GATEWAY_ID_HERE \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "is_enabled": false
  }'
```

### Refund Management

#### Create Refund
```bash
curl -X POST http://localhost:8007/api/v1/refunds \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "payment_id": "PAYMENT_ID_HERE",
    "amount": 25.00,
    "reason": "customer_request"
  }'
```

#### Get Refund by ID
```bash
curl -X GET http://localhost:8007/api/v1/refunds/REFUND_ID_HERE \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

#### Get Refunds by Payment
```bash
curl -X GET http://localhost:8007/api/v1/refunds/payment/PAYMENT_ID_HERE \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

#### Process Refund
```bash
curl -X POST http://localhost:8007/api/v1/refunds/REFUND_ID_HERE/process \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### Event Management

#### Get Payment Events
```bash
curl -X GET http://localhost:8007/api/v1/payment-events/payment/PAYMENT_ID_HERE \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### Settlement Management

#### Create Settlement
```bash
curl -X POST http://localhost:8007/api/v1/settlements \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "gateway_id": "GATEWAY_ID_HERE",
    "reference": "SETTLEMENT_REF_001",
    "amount": 1000.00,
    "currency": "USD"
  }'
```

#### Get Settlement by ID
```bash
curl -X GET http://localhost:8007/api/v1/settlements/SETTLEMENT_ID_HERE \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

#### Get Settlements by Gateway
```bash
curl -X GET "http://localhost:8007/api/v1/settlements/gateway/GATEWAY_ID_HERE?page=1&per_page=10" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

#### Update Settlement
```bash
curl -X PUT http://localhost:8007/api/v1/settlements/SETTLEMENT_ID_HERE \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "status": "completed",
    "deposited_at": "2024-01-15T10:30:00Z"
  }'
```

## Complete Payment Flow Test

Here's a complete workflow that demonstrates the full payment processing experience:

### 1. Setup Phase
```bash
# 1. Login to get JWT token
USER_TOKEN=$(curl -s -X POST http://localhost:8001/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"merchant@example.com","password":"SecurePass123!"}' | \
  jq -r '.data.access_token')

# 2. Create a payment gateway
GATEWAY_ID=$(curl -s -X POST http://localhost:8007/api/v1/gateways \
  -H "Authorization: Bearer $USER_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test Gateway",
    "provider": "test",
    "is_enabled": true,
    "is_sandbox": true
  }' | \
  jq -r '.data.id')

# 3. Create a payment method
PAYMENT_METHOD_ID=$(curl -s -X POST http://localhost:8007/api/v1/payment-methods \
  -H "Authorization: Bearer $USER_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "type": "credit_card",
    "provider": "test",
    "token": "pm_test123",
    "last4": "4242",
    "brand": "Visa",
    "name": "John Customer"
  }' | \
  jq -r '.data.id')
```

### 2. Payment Processing Phase
```bash
# 4. Create a payment
PAYMENT_ID=$(curl -s -X POST http://localhost:8007/api/v1/payments \
  -H "Authorization: Bearer $USER_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "order_id": "ORDER_UUID_HERE",
    "merchant_id": "MERCHANT_UUID_HERE",
    "customer_id": "CUSTOMER_UUID_HERE",
    "payment_method_id": "'$PAYMENT_METHOD_ID'",
    "amount": 99.99,
    "currency": "USD",
    "description": "Test payment"
  }' | \
  jq -r '.data.id')

# 5. Process the payment
curl -X POST http://localhost:8007/api/v1/payments/$PAYMENT_ID/process \
  -H "Authorization: Bearer $USER_TOKEN"

# 6. View payment details
curl -X GET http://localhost:8007/api/v1/payments/$PAYMENT_ID \
  -H "Authorization: Bearer $USER_TOKEN"
```

### 3. Refund Phase
```bash
# 7. Create a refund
REFUND_ID=$(curl -s -X POST http://localhost:8007/api/v1/refunds \
  -H "Authorization: Bearer $USER_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "payment_id": "'$PAYMENT_ID'",
    "amount": 25.00,
    "reason": "customer_request"
  }' | \
  jq -r '.data.id')

# 8. Process the refund
curl -X POST http://localhost:8007/api/v1/refunds/$REFUND_ID/process \
  -H "Authorization: Bearer $USER_TOKEN"

# 9. View refund details
curl -X GET http://localhost:8007/api/v1/refunds/$REFUND_ID \
  -H "Authorization: Bearer $USER_TOKEN"
```

## Error Handling

All services use consistent error response format:

```json
{
  "success": false,
  "error": {
    "code": "ERROR_CODE",
    "message": "Human readable error message",
    "details": {
      "field": "specific error details"
    }
  }
}
```

Common HTTP status codes:
- `200` - Success
- `201` - Created
- `400` - Bad Request
- `401` - Unauthorized
- `403` - Forbidden
- `404` - Not Found
- `409` - Conflict
- `422` - Validation Error
- `500` - Internal Server Error

## Performance Testing

Use tools like Apache Bench or Artillery for load testing:

```bash
# Basic load test for payment creation
ab -n 1000 -c 10 "http://localhost:8007/api/v1/payments"

# Payment retrieval
ab -n 500 -c 5 "http://localhost:8007/api/v1/payments/PAYMENT_ID_HERE"
```

## Database Verification

You can verify data persistence by connecting to the database:

```bash
# PostgreSQL services
docker exec -it unified-commerce-postgres psql -U payment_user -d payment_service
```

## Next Steps

With the Payment Service tested and validated, you can:

1. **Implement Promotions Service** - Discount codes, sales, and loyalty programs
2. **Create GraphQL Gateway** - Unified API endpoint
3. **Build Frontend Applications** - React admin panel and Next.js storefront

The Payment Service is designed to be independently deployable and scalable, following microservices best practices for the unified commerce platform.