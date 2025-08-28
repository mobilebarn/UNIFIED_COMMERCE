# Promotions Service API Testing Guide

## Testing the Promotions Service

Once the Promotions Service is running, you can test it using these API calls.

### Prerequisites
- Promotions Service running on port 8007
- PostgreSQL database initialized
- Infrastructure services running via Docker Compose
- Valid JWT token for authenticated endpoints

### API Endpoints

#### 1. Health Check
```bash
curl -X GET http://localhost:8007/health
```

Expected Response:
```json
{
  "service": "promotions",
  "status": "healthy",
  "time": "2024-01-01T12:00:00Z",
  "checks": {
    "postgres": "healthy"
  }
}
```

#### 2. Create Promotion (Authenticated)
```bash
curl -X POST http://localhost:8007/api/v1/promotions \
  -H "Authorization: Bearer YOUR_JWT_TOKEN_HERE" \
  -H "Content-Type: application/json" \
  -d '{
    "merchant_id": "MERCHANT_UUID_HERE",
    "name": "Summer Sale",
    "description": "20% off all summer items",
    "type": "discount",
    "start_date": "2024-06-01T00:00:00Z",
    "end_date": "2024-08-31T23:59:59Z",
    "usage_limit": 1000,
    "priority": 1,
    "applies_to": {
      "all_products": true
    },
    "target": {
      "type": "order",
      "value": 20,
      "value_type": "percentage"
    },
    "allocation": {
      "method": "across"
    }
  }'
```

#### 3. Get Promotion (Authenticated)
```bash
curl -X GET http://localhost:8007/api/v1/promotions/PROMOTION_ID_HERE \
  -H "Authorization: Bearer YOUR_JWT_TOKEN_HERE"
```

#### 4. Get Promotions by Merchant (Authenticated)
```bash
curl -X GET http://localhost:8007/api/v1/promotions/merchant/MERCHANT_ID_HERE \
  -H "Authorization: Bearer YOUR_JWT_TOKEN_HERE"
```

#### 5. Update Promotion (Authenticated)
```bash
curl -X PUT http://localhost:8007/api/v1/promotions/PROMOTION_ID_HERE \
  -H "Authorization: Bearer YOUR_JWT_TOKEN_HERE" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Updated Summer Sale",
    "description": "25% off all summer items"
  }'
```

#### 6. Delete Promotion (Authenticated)
```bash
curl -X DELETE http://localhost:8007/api/v1/promotions/PROMOTION_ID_HERE \
  -H "Authorization: Bearer YOUR_JWT_TOKEN_HERE"
```

#### 7. Create Discount Code (Authenticated)
```bash
curl -X POST http://localhost:8007/api/v1/discount-codes \
  -H "Authorization: Bearer YOUR_JWT_TOKEN_HERE" \
  -H "Content-Type: application/json" \
  -d '{
    "promotion_id": "PROMOTION_ID_HERE",
    "code": "SUMMER25",
    "usage_limit": 100,
    "customer_uses": 1,
    "expires_at": "2024-08-31T23:59:59Z"
  }'
```

#### 8. Get Discount Code by ID (Authenticated)
```bash
curl -X GET http://localhost:8007/api/v1/discount-codes/DISCOUNT_CODE_ID_HERE \
  -H "Authorization: Bearer YOUR_JWT_TOKEN_HERE"
```

#### 9. Get Discount Code by Code (Authenticated)
```bash
curl -X GET http://localhost:8007/api/v1/discount-codes/code/SUMMER25 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN_HERE"
```

#### 10. Validate Discount Code (Public)
```bash
curl -X POST http://localhost:8007/api/v1/public/discount-codes/validate \
  -H "Content-Type: application/json" \
  -d '{
    "code": "SUMMER25",
    "order_amount": 100.00,
    "customer_id": "CUSTOMER_UUID_HERE"
  }'
```

#### 11. Create Gift Card (Authenticated)
```bash
curl -X POST http://localhost:8007/api/v1/gift-cards \
  -H "Authorization: Bearer YOUR_JWT_TOKEN_HERE" \
  -H "Content-Type: application/json" \
  -d '{
    "merchant_id": "MERCHANT_UUID_HERE",
    "balance": 50.00,
    "currency": "USD",
    "expires_at": "2025-12-31T23:59:59Z"
  }'
```

#### 12. Get Gift Card by Code (Authenticated)
```bash
curl -X GET http://localhost:8007/api/v1/gift-cards/code/GCXXXXXXXXXX \
  -H "Authorization: Bearer YOUR_JWT_TOKEN_HERE"
```

#### 13. Redeem Gift Card (Authenticated)
```bash
curl -X POST http://localhost:8007/api/v1/gift-cards/GIFT_CARD_ID_HERE/redeem \
  -H "Authorization: Bearer YOUR_JWT_TOKEN_HERE" \
  -H "Content-Type: application/json" \
  -d '{
    "order_id": "ORDER_UUID_HERE",
    "amount": 25.00
  }'
```

#### 14. Create Loyalty Program (Authenticated)
```bash
curl -X POST http://localhost:8007/api/v1/loyalty/programs \
  -H "Authorization: Bearer YOUR_JWT_TOKEN_HERE" \
  -H "Content-Type: application/json" \
  -d '{
    "merchant_id": "MERCHANT_UUID_HERE",
    "name": "Premium Rewards",
    "description": "Earn points on every purchase",
    "point_value": 1.00,
    "reward_ratio": 100.00,
    "settings": {
      "earn_on_purchase": true,
      "minimum_purchase_amount": 10.00,
      "redemption_enabled": true,
      "redemption_rate": 1.00,
      "minimum_redemption_points": 100
    }
  }'
```

#### 15. Create Loyalty Member (Authenticated)
```bash
curl -X POST http://localhost:8007/api/v1/loyalty/members \
  -H "Authorization: Bearer YOUR_JWT_TOKEN_HERE" \
  -H "Content-Type: application/json" \
  -d '{
    "loyalty_program_id": "LOYALTY_PROGRAM_ID_HERE",
    "customer_id": "CUSTOMER_UUID_HERE"
  }'
```

#### 16. Earn Loyalty Points (Authenticated)
```bash
curl -X POST http://localhost:8007/api/v1/loyalty/members/MEMBER_ID_HERE/earn \
  -H "Authorization: Bearer YOUR_JWT_TOKEN_HERE" \
  -H "Content-Type: application/json" \
  -d '{
    "points": 50,
    "description": "Purchase reward",
    "reference_id": "ORDER_UUID_HERE"
  }'
```

#### 17. Redeem Loyalty Points (Authenticated)
```bash
curl -X POST http://localhost:8007/api/v1/loyalty/members/MEMBER_ID_HERE/redeem \
  -H "Authorization: Bearer YOUR_JWT_TOKEN_HERE" \
  -H "Content-Type: application/json" \
  -d '{
    "points": 100,
    "description": "Redeemed for discount",
    "reference_id": "ORDER_UUID_HERE"
  }'
```

### Error Responses

#### Validation Error (422)
```json
{
  "success": false,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Validation failed",
    "details": {
      "validation": "Field validation error details"
    }
  }
}
```

#### Unauthorized (401)
```json
{
  "success": false,
  "error": {
    "code": "UNAUTHORIZED",
    "message": "Invalid or missing authentication token"
  }
}
```

#### Not Found (404)
```json
{
  "success": false,
  "error": {
    "code": "NOT_FOUND",
    "message": "Resource not found"
  }
}
```

#### Internal Server Error (500)
```json
{
  "success": false,
  "error": {
    "code": "INTERNAL_SERVER_ERROR",
    "message": "An unexpected error occurred"
  }
}
```

### Database Verification

You can verify the data is being stored correctly by connecting to PostgreSQL:

```bash
# Connect to PostgreSQL container
docker exec -it unified-commerce-postgres psql -U promotions_user -d promotions_service

# Check promotions table
SELECT id, name, type, status, start_date, end_date FROM promotions;

# Check discount codes table
SELECT id, promotion_id, code, status, usage_limit, used_count FROM discount_codes;

# Check gift cards table
SELECT id, code, balance, initial_balance, status FROM gift_cards;

# Check loyalty programs table
SELECT id, name, status, point_value, reward_ratio FROM loyalty_programs;

# Check loyalty members table
SELECT id, loyalty_program_id, customer_id, points, status FROM loyalty_members;
```

### Testing Workflow

1. **Start Infrastructure**: `docker-compose up -d postgres redis`
2. **Start Promotions Service**: `go run services/promotions/cmd/server/main.go`
3. **Test Health Check**: Verify service is running
4. **Create Promotion**: Create a new promotion
5. **Test Promotion Management**: Test CRUD operations for promotions
6. **Create Discount Codes**: Create discount codes for promotions
7. **Test Discount Validation**: Validate discount codes
8. **Test Gift Cards**: Create and redeem gift cards
9. **Test Loyalty Programs**: Create loyalty programs and members
10. **Test Loyalty Points**: Earn and redeem loyalty points

This testing guide demonstrates the complete promotions management capabilities of our unified commerce platform's Promotions Service.