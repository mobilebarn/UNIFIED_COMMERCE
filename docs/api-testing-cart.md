# Cart & Checkout Service API Testing Guide

## Overview

This guide provides comprehensive API testing instructions for the Cart & Checkout Service in our Unified Commerce Platform. The service handles shopping cart management, checkout flows, and order creation.

## Prerequisites

1. **Infrastructure Services Running:**
   ```bash
   docker-compose up -d postgres redis
   ```

2. **Environment Setup:**
   - Copy `.env.example` to `.env` in the cart service directory
   - Update database connection strings as needed
   - Ensure the Cart & Checkout Service is running on port 8006

3. **Dependencies:**
   - Identity Service running on port 8001 (for authentication)
   - Product Catalog Service running on port 8003 (for product data)
   - Inventory Service running on port 8004 (for stock management)

## Authentication Flow

All protected endpoints require JWT authentication. Here's the complete authentication flow:

### 1. User Registration & Login

```bash
# Register a new user (if needed)
curl -X POST http://localhost:8001/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "customer@example.com",
    "username": "customeruser",
    "password": "SecurePass123!",
    "first_name": "John",
    "last_name": "Customer",
    "phone": "+1234567890"
  }'

# Login to get JWT token
curl -X POST http://localhost:8001/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "customer@example.com",
    "password": "SecurePass123!"
  }'
```

**Save the `access_token` from the login response for authenticated requests.**

## Service Testing Workflows

### Cart Management

#### Health Check
```bash
curl -X GET http://localhost:8006/health
```

#### Create a New Cart
```bash
# For authenticated users
curl -X POST http://localhost:8006/api/v1/carts \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "merchant_id": "MERCHANT_UUID_HERE",
    "currency": "USD"
  }'

# For guest users
curl -X POST http://localhost:8006/api/v1/public/carts \
  -H "Content-Type: application/json" \
  -d '{
    "session_id": "SESSION_ID_HERE",
    "merchant_id": "MERCHANT_UUID_HERE",
    "currency": "USD"
  }'
```

#### Get Cart by ID
```bash
curl -X GET http://localhost:8006/api/v1/carts/CART_ID_HERE \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

#### Get Cart by Session (Guest)
```bash
curl -X GET http://localhost:8006/api/v1/public/carts/session/SESSION_ID_HERE
```

#### Get Customer Cart
```bash
curl -X GET "http://localhost:8006/api/v1/carts/customer/CUSTOMER_ID_HERE?merchant_id=MERCHANT_ID_HERE" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

#### Update Cart Information
```bash
curl -X PUT http://localhost:8006/api/v1/carts/CART_ID_HERE \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "customer_email": "customer@example.com",
    "customer_phone": "+1234567890",
    "customer_first_name": "John",
    "customer_last_name": "Customer",
    "billing_address": {
      "first_name": "John",
      "last_name": "Customer",
      "address1": "123 Main St",
      "city": "San Francisco",
      "province": "CA",
      "country": "US",
      "zip": "94105"
    },
    "shipping_address": {
      "first_name": "John",
      "last_name": "Customer",
      "address1": "123 Main St",
      "city": "San Francisco",
      "province": "CA",
      "country": "US",
      "zip": "94105"
    }
  }'
```

### Line Item Management

#### Add Line Item to Cart
```bash
curl -X POST http://localhost:8006/api/v1/carts/CART_ID_HERE/line-items \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "product_id": "PRODUCT_UUID_HERE",
    "sku": "PRODUCT_SKU_HERE",
    "name": "Product Name",
    "quantity": 2,
    "price": 29.99
  }'
```

#### Update Line Item
```bash
curl -X PUT http://localhost:8006/api/v1/carts/line-items/LINE_ITEM_ID_HERE \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "quantity": 3
  }'
```

#### Remove Line Item
```bash
curl -X DELETE http://localhost:8006/api/v1/carts/line-items/LINE_ITEM_ID_HERE \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### Checkout Management

#### Create Checkout from Cart
```bash
curl -X POST http://localhost:8006/api/v1/checkouts \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "cart_id": "CART_ID_HERE"
  }'
```

#### Get Checkout by ID
```bash
curl -X GET http://localhost:8006/api/v1/checkouts/CHECKOUT_ID_HERE \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

#### Get Checkout by Token
```bash
curl -X GET http://localhost:8006/api/v1/checkouts/token/CHECKOUT_TOKEN_HERE
```

#### Update Customer Information in Checkout
```bash
curl -X PUT http://localhost:8006/api/v1/checkouts/CHECKOUT_ID_HERE/customer-info \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "customer": {
      "email": "customer@example.com",
      "phone": "+1234567890",
      "first_name": "John",
      "last_name": "Customer"
    }
  }'
```

#### Update Shipping Address in Checkout
```bash
curl -X PUT http://localhost:8006/api/v1/checkouts/CHECKOUT_ID_HERE/shipping-address \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "address": {
      "first_name": "John",
      "last_name": "Customer",
      "address1": "123 Main St",
      "city": "San Francisco",
      "province": "CA",
      "country": "US",
      "zip": "94105"
    }
  }'
```

#### Add Shipping Line to Checkout
```bash
curl -X POST http://localhost:8006/api/v1/checkouts/CHECKOUT_ID_HERE/shipping-lines \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "shipping_line": {
      "title": "Standard Shipping",
      "code": "STD",
      "price": 5.99
    }
  }'
```

#### Apply Discount Code
```bash
curl -X POST http://localhost:8006/api/v1/checkouts/CHECKOUT_ID_HERE/discounts \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "code": "SAVE10"
  }'
```

#### Remove Discount Code
```bash
curl -X DELETE http://localhost:8006/api/v1/checkouts/CHECKOUT_ID_HERE/discounts/SAVE10 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

#### Complete Checkout
```bash
curl -X POST http://localhost:8006/api/v1/checkouts/CHECKOUT_ID_HERE/complete \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "payment_method_id": "pm_example123"
  }'
```

### Configuration and Reference Data

#### Get Shipping Rates
```bash
curl -X GET http://localhost:8006/api/v1/config/shipping-rates \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

#### Get Payment Methods
```bash
curl -X GET http://localhost:8006/api/v1/config/payment-methods \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### Administrative Operations

#### Process Abandoned Carts
```bash
curl -X POST http://localhost:8006/api/v1/admin/abandonment/process \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

#### Get Abandoned Carts
```bash
curl -X GET "http://localhost:8006/api/v1/admin/abandonment/carts?page=1&per_page=10" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

#### Cleanup Expired Carts
```bash
curl -X POST http://localhost:8006/api/v1/admin/cleanup/expired \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## Complete Shopping Flow Test

Here's a complete workflow that demonstrates the full shopping experience:

### 1. Setup Phase
```bash
# 1. Login to get JWT token
USER_TOKEN=$(curl -s -X POST http://localhost:8001/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"customer@example.com","password":"SecurePass123!"}' | \
  jq -r '.data.access_token')

# 2. Create a cart
CART_ID=$(curl -s -X POST http://localhost:8006/api/v1/carts \
  -H "Authorization: Bearer $USER_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"merchant_id":"MERCHANT_UUID_HERE","currency":"USD"}' | \
  jq -r '.data.id')
```

### 2. Shopping Phase
```bash
# 3. Add items to cart
curl -X POST http://localhost:8006/api/v1/carts/$CART_ID/line-items \
  -H "Authorization: Bearer $USER_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "product_id": "PRODUCT_UUID_HERE",
    "sku": "PRODUCT_SKU_HERE",
    "name": "Product Name",
    "quantity": 2,
    "price": 29.99
  }'

# 4. View cart
curl -X GET http://localhost:8006/api/v1/carts/$CART_ID \
  -H "Authorization: Bearer $USER_TOKEN"
```

### 3. Checkout Phase
```bash
# 5. Create checkout
CHECKOUT_ID=$(curl -s -X POST http://localhost:8006/api/v1/checkouts \
  -H "Authorization: Bearer $USER_TOKEN" \
  -H "Content-Type: application/json" \
  -d "{\"cart_id\":\"$CART_ID\"}" | \
  jq -r '.data.id')

# 6. Add customer info
curl -X PUT http://localhost:8006/api/v1/checkouts/$CHECKOUT_ID/customer-info \
  -H "Authorization: Bearer $USER_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "customer": {
      "email": "customer@example.com",
      "phone": "+1234567890",
      "first_name": "John",
      "last_name": "Customer"
    }
  }'

# 7. Add shipping address
curl -X PUT http://localhost:8006/api/v1/checkouts/$CHECKOUT_ID/shipping-address \
  -H "Authorization: Bearer $USER_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "address": {
      "first_name": "John",
      "last_name": "Customer",
      "address1": "123 Main St",
      "city": "San Francisco",
      "province": "CA",
      "country": "US",
      "zip": "94105"
    }
  }'

# 8. Add shipping method
curl -X POST http://localhost:8006/api/v1/checkouts/$CHECKOUT_ID/shipping-lines \
  -H "Authorization: Bearer $USER_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "shipping_line": {
      "title": "Standard Shipping",
      "code": "STD",
      "price": 5.99
    }
  }'

# 9. Complete checkout
ORDER_ID=$(curl -s -X POST http://localhost:8006/api/v1/checkouts/$CHECKOUT_ID/complete \
  -H "Authorization: Bearer $USER_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"payment_method_id":"pm_example123"}' | \
  jq -r '.data.order_id')
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
- `410` - Gone (expired cart)
- `500` - Internal Server Error

## Performance Testing

Use tools like Apache Bench or Artillery for load testing:

```bash
# Basic load test for cart creation
ab -n 1000 -c 10 "http://localhost:8006/api/v1/public/carts"

# Cart retrieval
ab -n 500 -c 5 "http://localhost:8006/api/v1/carts/CART_ID_HERE"
```

## Database Verification

You can verify data persistence by connecting to the database:

```bash
# PostgreSQL services
docker exec -it unified-commerce-postgres psql -U cart_user -d cart_checkout_service
```

## Next Steps

With the Cart & Checkout Service tested and validated, you can:

1. **Build Payments Service** - Complete payment processing integrations
2. **Implement Promotions Service** - Discount codes, sales, and loyalty programs
3. **Create GraphQL Gateway** - Unified API endpoint
4. **Build Frontend Applications** - React admin panel and Next.js storefront

The Cart & Checkout Service is designed to be independently deployable and scalable, following microservices best practices for the unified commerce platform.