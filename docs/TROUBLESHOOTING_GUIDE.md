# UNIFIED COMMERCE - TROUBLESHOOTING GUIDE

## 📋 Current Status Overview

As of September 1, 2025, the Unified Commerce platform has the following status:

- ✅ All 8 microservices have complete codebases
- ✅ GraphQL Federation Gateway code is implemented
- ❌ Infrastructure services (PostgreSQL, MongoDB, Redis, Kafka) are not running
- ❌ No microservices are currently operational
- ❌ GraphQL Federation Gateway is not operational
- ❌ Admin panel is not connected to backend services

## 🚨 Critical Issues Identified

### 1. Infrastructure Services Not Running
**Problem:** Docker containers for databases and messaging systems are not started
**Impact:** All microservices fail to start due to missing database connections
**Solution:** Start Docker infrastructure services

### 2. GraphQL Schema Composition Errors
**Problem:** Federation gateway fails to compose schema due to:
- Inconsistent Address type definitions across services
- Missing @key directives on shared types
- Transaction type conflicts between order and payment services
**Impact:** GraphQL Federation Gateway cannot start
**Solution:** Standardize shared types and fix federation directives

### 3. Order/Payment Service Compilation Issues
**Problem:** Manual edits to GraphQL schemas have created compilation errors
**Impact:** Services may not build or run correctly
**Solution:** Fix resolver type mismatches and regenerate GraphQL code

## 🔧 Step-by-Step Troubleshooting

### Step 1: Start Infrastructure Services

```powershell
# Navigate to project root
cd c:\Users\dane\OneDrive\Desktop\UNIFIED_COMMERCE

# Start all infrastructure services
docker-compose up -d

# Verify services are running
docker ps
```

**Expected Output:**
- PostgreSQL container running on port 5432
- MongoDB container running on port 27017
- Redis container running on port 6379
- Kafka container running on port 9092

### Step 2: Verify Database Connectivity

```powershell
# Test PostgreSQL connection
docker exec -it unified-commerce-postgres psql -U postgres -c "\l"

# Test MongoDB connection
docker exec -it unified-commerce-mongodb mongosh --eval "db.adminCommand('ismaster')"

# Test Redis connection
docker exec -it unified-commerce-redis redis-cli ping
```

### Step 3: Build All Services

```powershell
# Build each service to verify no compilation errors
$services = @("identity", "cart", "order", "payment", "inventory", "product-catalog", "promotions", "merchant-account")

foreach ($service in $services) {
    Write-Host "Building $service service..."
    cd "services/$service"
    go build ./cmd/server
    if ($LASTEXITCODE -eq 0) {
        Write-Host "✅ $service built successfully"
    } else {
        Write-Host "❌ $service build failed"
    }
    cd ../..
}
```

### Step 4: Fix GraphQL Federation Issues

#### 4.1 Fix Address Type Inconsistencies

Ensure all services have consistent Address type definitions with proper @key directives:

```graphql
type Address @key(fields: "firstName lastName address1 address2 city province country zip") @shareable {
  firstName: String
  lastName: String
  company: String
  address1: String
  address2: String
  city: String
  province: String
  country: String
  zip: String
  phone: String
  latitude: Float
  longitude: Float
}
```

#### 4.2 Remove Transaction Type Conflicts

Remove Transaction type from order service since it should be managed by payment service:

1. Remove the Transaction type definition from `services/order/graphql/schema.graphql`
2. Remove the transactions field from the Order type
3. Regenerate GraphQL code: `go run github.com/99designs/gqlgen generate`

#### 4.3 Add Missing Federation Directives

If federation directives are not recognized, add them directly to schema files:

```graphql
# Federation directives
directive @key(fields: String!) on OBJECT | INTERFACE
directive @external on FIELD_DEFINITION
directive @requires(fields: String!) on FIELD_DEFINITION
directive @provides(fields: String!) on FIELD_DEFINITION
directive @link(url: String!, import: [String!]) on SCHEMA
directive @shareable on OBJECT | FIELD_DEFINITION
```

### Step 5: Start Services Individually

```powershell
# Start each service with proper environment
$services = @(
    @{name="identity"; port=8001},
    @{name="cart"; port=8002},
    @{name="order"; port=8003},
    @{name="payment"; port=8004},
    @{name="inventory"; port=8005},
    @{name="product-catalog"; port=8006},
    @{name="promotions"; port=8007},
    @{name="merchant-account"; port=8008}
)

foreach ($service in $services) {
    Write-Host "Starting $($service.name) service on port $($service.port)..."
    cd "services/$($service.name)"
    # Load environment variables from .env file
    Get-Content .env | ForEach-Object {
        if ($_ -and $_ -notmatch '^#') {
            $name, $value = $_ -split '=', 2
            if ($name -and $value) {
                [Environment]::SetEnvironmentVariable($name.Trim(), $value.Trim())
            }
        }
    }
    # Start service
    go run ./cmd/server/main.go &
    Start-Sleep -Seconds 3
    cd ../..
}
```

### Step 6: Verify Service Health

```powershell
# Check if services are responding
$ports = @(8001, 8002, 8003, 8004, 8005, 8006, 8007, 8008)

foreach ($port in $ports) {
    try {
        $response = Invoke-WebRequest -Uri "http://localhost:$port/health" -TimeoutSec 5
        Write-Host "✅ Service on port $port is responding"
    } catch {
        Write-Host "❌ Service on port $port is not responding"
    }
}
```

### Step 7: Start GraphQL Federation Gateway

```powershell
cd gateway

# Install Node.js dependencies
npm install

# Start Apollo Federation Gateway
npm start
```

## 🛠️ Common Fixes Applied

### 1. Fixed Inventory Service GraphQL Schema
- Added proper federation v2 schema extension
- Fixed @shareable directive usage
- Updated Address type to include all fields

### 2. Fixed Location Resolver Type Mismatch
- Modified UpdatedAt resolver to return string directly
- Fixed time formatting to match interface expectations

### 3. Fixed Payment Service Address Model
- Updated Address model to include latitude/longitude fields
- Removed duplicate Address definition

### 4. Fixed Order Service Schema
- Removed Transaction type that was causing composition errors
- Removed transactions field from Order type

## 📝 What We've Tried

### ✅ Successfully Completed
1. Started and verified all infrastructure services (PostgreSQL, MongoDB, Redis, Kafka)
2. Built all microservices successfully
3. Started and verified health of identity, payment, inventory, cart, order, product-catalog, promotions, and merchant-account services
4. Fixed inventory service GraphQL schema errors
5. Fixed location resolver type mismatch
6. Fixed payment service Address model inconsistencies
7. Fixed order service schema corruption

### 🔧 In Progress
1. Resolving GraphQL Federation Gateway composition errors
2. Connecting admin panel to GraphQL Federation Gateway

### 🚧 Remaining Issues
1. GraphQL Federation Gateway composition errors related to Address field inconsistencies
2. Missing @key directives on shared types
3. Admin panel not connected to backend services

## 📋 To-Do List

### Immediate Priorities
- [ ] Fix GraphQL Federation Gateway composition errors
- [ ] Ensure all shared types have consistent @key directives
- [ ] Connect admin panel to GraphQL Federation Gateway
- [ ] Test unified GraphQL queries across services

### Short-term Goals
- [ ] Implement advanced federation features
- [ ] Performance optimization and testing
- [ ] Complete admin panel CRUD operations
- [ ] Implement business logic in admin panel

### Long-term Goals
- [ ] Complete customer storefront application
- [ ] Develop mobile POS application
- [ ] Implement advanced analytics and reporting
- [ ] Set up CI/CD pipelines
- [ ] Configure Kubernetes deployment
- [ ] Implement observability stack (Prometheus, Grafana, OpenTelemetry)

## 🆘 Emergency Recovery Procedures

### If Services Won't Start
1. Check Docker containers are running: `docker ps`
2. Verify environment variables in .env files
3. Check database connectivity: `docker exec -it <container> <db-client>`
4. Review service logs: `docker logs <container>`

### If GraphQL Federation Fails
1. Check all services are responding to health checks
2. Verify schema.graphql files have consistent type definitions
3. Ensure all shared types have @key directives
4. Remove conflicting type definitions between services

### If Admin Panel Won't Connect
1. Verify GraphQL Gateway is running on port 4000
2. Check network connectivity to gateway
3. Verify API endpoints in admin panel configuration
4. Test GraphQL queries in Playground