# Unified Commerce Platform - API Testing Guide

## Overview

This guide provides comprehensive API testing instructions for the implemented microservices in our Unified Commerce Platform. All services follow RESTful API design principles with consistent error handling and response formats.

## Prerequisites

1. **Infrastructure Services Running:**
   ```bash
   docker-compose up -d postgres mongodb redis elasticsearch kafka
   ```

2. **Environment Setup:**
   - Copy `.env.example` to `.env` for each service
   - Update database connection strings as needed
   - Ensure all services are running on their designated ports

3. **Service Ports:**
   - Identity Service: `8001`
   - Merchant Account Service: `8002`
   - Product Catalog Service: `8003`
   - Inventory Service: `8004`

## Authentication Flow

All protected endpoints require JWT authentication. Here's the complete authentication flow:

### 1. User Registration & Login

```bash
# Register a new user
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

### Identity Service (Port 8001)

#### Health Check
```bash
curl -X GET http://localhost:8001/health
```

#### User Management
```bash
# Get user profile
curl -X GET http://localhost:8001/api/v1/users/me \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"

# Change password
curl -X POST http://localhost:8001/api/v1/users/change-password \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "old_password": "SecurePass123!",
    "new_password": "NewSecurePass456!"
  }'

# Admin: List all users (requires admin role)
curl -X GET http://localhost:8001/api/v1/admin/users \
  -H "Authorization: Bearer ADMIN_JWT_TOKEN"
```

### Merchant Account Service (Port 8002)

#### Create Merchant Account
```bash
curl -X POST http://localhost:8002/api/v1/merchants \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "business_name": "Tech Store Inc",
    "business_type": "retail",
    "industry": "technology",
    "description": "Leading technology retailer",
    "website": "https://techstore.example.com",
    "address": {
      "street1": "123 Business Ave",
      "city": "San Francisco",
      "state": "CA",
      "postal_code": "94105",
      "country": "US"
    },
    "contact": {
      "email": "info@techstore.example.com",
      "phone": "+1-555-0123"
    }
  }'
```

#### Subscription Management
```bash
# Get subscription plans
curl -X GET http://localhost:8002/api/v1/subscription-plans

# Subscribe to a plan
curl -X POST http://localhost:8002/api/v1/merchants/{merchant_id}/subscribe \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "plan_id": "PLAN_UUID",
    "billing_cycle": "monthly",
    "payment_method_id": "pm_example123"
  }'

# Get subscription details
curl -X GET http://localhost:8002/api/v1/merchants/{merchant_id}/subscription \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

#### Team Management
```bash
# Add team member
curl -X POST http://localhost:8002/api/v1/merchants/{merchant_id}/team \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "employee@techstore.com",
    "role": "manager",
    "permissions": ["inventory:read", "orders:read", "products:write"]
  }'

# List team members
curl -X GET http://localhost:8002/api/v1/merchants/{merchant_id}/team \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### Product Catalog Service (Port 8003)

#### Category Management
```bash
# Create category
curl -X POST http://localhost:8003/api/v1/categories \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Electronics",
    "description": "Electronic devices and accessories",
    "slug": "electronics",
    "is_active": true,
    "metadata": {
      "seo_title": "Electronics - Best Deals",
      "seo_description": "Shop the latest electronics"
    }
  }'

# List categories (public)
curl -X GET http://localhost:8003/api/v1/public/categories

# Get category with products
curl -X GET http://localhost:8003/api/v1/public/categories/{category_id}/products?page=1&limit=10
```

#### Product Management
```bash
# Create product
curl -X POST http://localhost:8003/api/v1/products \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "iPhone 15 Pro",
    "description": "Latest iPhone with advanced features",
    "sku": "IPHONE15PRO",
    "category_ids": ["CATEGORY_UUID"],
    "base_price": 999.00,
    "compare_at_price": 1199.00,
    "cost_per_item": 750.00,
    "track_inventory": true,
    "continue_selling_when_out_of_stock": false,
    "attributes": {
      "brand": "Apple",
      "model": "iPhone 15 Pro",
      "storage": "128GB",
      "color": "Natural Titanium"
    },
    "variants": [
      {
        "name": "128GB Natural Titanium",
        "sku": "IPHONE15PRO-128-NT",
        "price": 999.00,
        "attributes": {
          "storage": "128GB",
          "color": "Natural Titanium"
        }
      }
    ],
    "seo": {
      "title": "iPhone 15 Pro - 128GB Natural Titanium",
      "description": "Experience the latest iPhone technology"
    },
    "status": "active"
  }'

# Search products (public)
curl -X GET "http://localhost:8003/api/v1/public/products/search?q=iPhone&page=1&limit=10"

# Get product details (public)
curl -X GET http://localhost:8003/api/v1/public/products/{product_id}

# Update product
curl -X PUT http://localhost:8003/api/v1/products/{product_id} \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "description": "Updated description with new features",
    "base_price": 949.00
  }'
```

### Inventory Service (Port 8004)

#### Location Management
```bash
# Create inventory location
curl -X POST http://localhost:8004/api/v1/locations \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "merchant_id": "MERCHANT_UUID",
    "name": "Main Warehouse",
    "type": "warehouse",
    "code": "WH001",
    "description": "Primary fulfillment center",
    "address": {
      "street1": "456 Warehouse Blvd",
      "city": "Oakland",
      "state": "CA",
      "postal_code": "94607",
      "country": "US"
    },
    "settings": {
      "allow_negative_stock": false,
      "low_stock_threshold": 10,
      "auto_reorder_enabled": true,
      "reorder_quantity": 100
    }
  }'

# List locations
curl -X GET "http://localhost:8004/api/v1/locations?merchant_id=MERCHANT_UUID" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"

# Get location summary
curl -X GET http://localhost:8004/api/v1/locations/{location_id}/summary \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

#### Inventory Management
```bash
# Create inventory item
curl -X POST http://localhost:8004/api/v1/inventory \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "location_id": "LOCATION_UUID",
    "product_id": "PRODUCT_UUID",
    "sku": "IPHONE15PRO-128-NT",
    "quantity": 50,
    "cost": 750.00,
    "retail_price": 999.00,
    "low_stock_threshold": 5,
    "bin": "A1-B2-C3"
  }'

# Check stock availability
curl -X GET "http://localhost:8004/api/v1/inventory/check-availability?sku=IPHONE15PRO-128-NT&location_id=LOCATION_UUID&quantity=5" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"

# Adjust inventory
curl -X POST http://localhost:8004/api/v1/inventory/{inventory_id}/adjust \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "quantity": 45,
    "reason": "damage",
    "notes": "5 units damaged in transit",
    "reference": "DAMAGE_REPORT_001"
  }'

# Get low stock items
curl -X GET "http://localhost:8004/api/v1/inventory/low-stock?location_id=LOCATION_UUID" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"

# Get stock movements
curl -X GET "http://localhost:8004/api/v1/stock-movements?location_id=LOCATION_UUID&date_from=2024-01-01" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

#### Stock Reservations
```bash
# Create stock reservation
curl -X POST http://localhost:8004/api/v1/reservations \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "inventory_item_id": "INVENTORY_ITEM_UUID",
    "quantity": 2,
    "type": "order",
    "reference": "ORDER_123",
    "expires_at": "2024-12-31T23:59:59Z",
    "notes": "Reserved for pending order"
  }'

# Get reservations by reference
curl -X GET "http://localhost:8004/api/v1/reservations?reference=ORDER_123" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"

# Fulfill reservation
curl -X POST http://localhost:8004/api/v1/reservations/{reservation_id}/fulfill \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "actual_quantity": 2
  }'

# Cancel reservation
curl -X POST http://localhost:8004/api/v1/reservations/{reservation_id}/cancel \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## Complete E-commerce Workflow Test

Here's a complete workflow that demonstrates the integration between all services:

### 1. Setup Phase
```bash
# 1. Register and login user
USER_TOKEN=$(curl -s -X POST http://localhost:8001/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"merchant@example.com","password":"SecurePass123!"}' | \
  jq -r '.data.access_token')

# 2. Create merchant account
MERCHANT_ID=$(curl -s -X POST http://localhost:8002/api/v1/merchants \
  -H "Authorization: Bearer $USER_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"business_name":"Tech Store","business_type":"retail"}' | \
  jq -r '.data.id')

# 3. Create inventory location
LOCATION_ID=$(curl -s -X POST http://localhost:8004/api/v1/locations \
  -H "Authorization: Bearer $USER_TOKEN" \
  -H "Content-Type: application/json" \
  -d "{\"merchant_id\":\"$MERCHANT_ID\",\"name\":\"Main Store\",\"type\":\"store\",\"code\":\"STORE001\"}" | \
  jq -r '.data.id')
```

### 2. Product Setup
```bash
# 4. Create product category
CATEGORY_ID=$(curl -s -X POST http://localhost:8003/api/v1/categories \
  -H "Authorization: Bearer $USER_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"Smartphones","slug":"smartphones"}' | \
  jq -r '.data.id')

# 5. Create product
PRODUCT_ID=$(curl -s -X POST http://localhost:8003/api/v1/products \
  -H "Authorization: Bearer $USER_TOKEN" \
  -H "Content-Type: application/json" \
  -d "{\"name\":\"iPhone 15\",\"sku\":\"IPHONE15\",\"category_ids\":[\"$CATEGORY_ID\"],\"base_price\":799}" | \
  jq -r '.data.id')

# 6. Add inventory
INVENTORY_ID=$(curl -s -X POST http://localhost:8004/api/v1/inventory \
  -H "Authorization: Bearer $USER_TOKEN" \
  -H "Content-Type: application/json" \
  -d "{\"location_id\":\"$LOCATION_ID\",\"product_id\":\"$PRODUCT_ID\",\"sku\":\"IPHONE15\",\"quantity\":100,\"cost\":600,\"retail_price\":799}" | \
  jq -r '.data.id')
```

### 3. Order Simulation
```bash
# 7. Reserve stock for order
RESERVATION_ID=$(curl -s -X POST http://localhost:8004/api/v1/reservations \
  -H "Authorization: Bearer $USER_TOKEN" \
  -H "Content-Type: application/json" \
  -d "{\"inventory_item_id\":\"$INVENTORY_ID\",\"quantity\":2,\"type\":\"order\",\"reference\":\"ORDER_001\"}" | \
  jq -r '.data.id')

# 8. Fulfill the reservation (simulate order completion)
curl -X POST http://localhost:8004/api/v1/reservations/$RESERVATION_ID/fulfill \
  -H "Authorization: Bearer $USER_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"actual_quantity":2}'

# 9. Check updated inventory
curl -X GET http://localhost:8004/api/v1/inventory/$INVENTORY_ID \
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
# Basic load test for product search
ab -n 1000 -c 10 "http://localhost:8003/api/v1/public/products/search?q=iPhone"

# Inventory availability check
ab -n 500 -c 5 "http://localhost:8004/api/v1/inventory/check-availability?sku=IPHONE15&location_id=$LOCATION_ID&quantity=1"
```

## Database Verification

You can verify data persistence by connecting to the databases:

```bash
# PostgreSQL services
docker exec -it unified-commerce-postgres psql -U identity_user -d identity_service
docker exec -it unified-commerce-postgres psql -U merchant_user -d merchant_account_service
docker exec -it unified-commerce-postgres psql -U inventory_user -d inventory_service

# MongoDB (Product Catalog)
docker exec -it unified-commerce-mongodb mongosh product_catalog_service
```

## Next Steps

With these core services tested and validated, you can:

1. **Build Order Service** - Complete order lifecycle management
2. **Implement Cart & Checkout** - Shopping cart and payment processing
3. **Add Payment Processing** - Multiple gateway integrations
4. **Create GraphQL Gateway** - Unified API endpoint
5. **Build Frontend Applications** - React admin panel and Next.js storefront

Each service is designed to be independently deployable and scalable, following microservices best practices for the unified commerce platform.