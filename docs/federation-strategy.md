# Federation Standardization Plan

## Current Status
- ✅ merchant-account service: Working with custom handler + federation support
- ✅ cart service: Working with custom handler (fixed SDL)
- ✅ order service: Working with federation support (resolved compilation errors)
- ✅ payment service: Working with federation support (resolved compilation errors)
- ✅ inventory service: Working with federation support
- ⏳ product-catalog service: Not currently running
- ⏳ promotions service: Not currently running
- ⏳ identity service: Not currently running

## Strategy: Use Merchant-Account Pattern

Instead of fighting with complex gqlgen federation generation, standardize all services to use the merchant-account service pattern:

1. **Custom GraphQL Handler** - Simple, reliable
2. **Federation SDL** - Hardcoded but accurate
3. **Standard Service Structure** - Consistent across all services

## Implementation Steps

1. ✅ Copy merchant-account handler pattern to failing services
2. ✅ Update SDL strings with correct schema definitions
3. ✅ Remove dependency on generated federation files
4. ⏳ Restart all services with consistent architecture
5. ⏳ Connect all services to GraphQL Federation Gateway
6. ⏳ Test unified schema composition
7. ⏳ Verify cross-service queries work correctly

This approach is:
- ✅ Faster to implement
- ✅ More maintainable
- ✅ Avoids gqlgen version conflicts
- ✅ Proven to work (merchant-account is working)

## 📊 CURRENT STATUS: PARTIALLY CONNECTED TO FEDERATION GATEWAY

Currently, only 3 of 8 microservices are successfully connected to the Apollo GraphQL Federation Gateway:

### ✅ Connected Services:
- Order Service (8003)
- Payment Service (8004)
- Inventory Service (8005)

### ❌ Not Connected Services:
- Identity Service (8001)
- Cart Service (8002)
- Product Catalog Service (8006)
- Promotions Service (8007)
- Merchant Account Service (8008)

The GraphQL Federation Gateway is running on `http://localhost:4000/graphql` and can successfully introspect the 3 connected services.

## 🔧 Issues to Resolve

### Port Conflicts
Multiple services are failing to start due to port binding errors:
- "listen tcp :8005: bind: Only one usage of each socket address (protocol/network address/port) is normally permitted."
- "listen tcp :8003: bind: Only one usage of each socket address (protocol/network address/port)."

### Missing Services
Five services are not responding and need to be started:
- Identity Service (8001)
- Cart Service (8002)
- Product Catalog Service (8006)
- Promotions Service (8007)
- Merchant Account Service (8008)

## 🎯 Next Steps

1. **Resolve Port Conflicts**
   - Identify and stop duplicate service instances
   - Ensure each service runs on its designated port only

2. **Start Missing Services**
   - Start Identity Service (8001)
   - Start Cart Service (8002)
   - Start Product Catalog Service (8006)
   - Start Promotions Service (8007)
   - Start Merchant Account Service (8008)

3. **Verify GraphQL Federation**
   - Confirm all running services are properly federated
   - Test cross-service queries
   - Update gateway configuration to include newly started services

4. **Update Documentation**
   - Update this document when services are connected
   - Update gateway configuration as services come online