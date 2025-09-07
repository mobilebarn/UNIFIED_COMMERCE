# UNIFIED COMMERCE - OPERATIONAL STARTUP GUIDE

## 🚀 Quick Start Guide (Based on Testing - September 1, 2025)

This guide provides step-by-step instructions to get the Unified Commerce platform operational, based on actual testing of the current system state.

---

## ⚠️ **BEFORE YOU START**

**Current Reality Check:**
- ✅ All code is written and builds successfully
- ✅ Docker infrastructure is now started and running
- ⏳ GraphQL Federation Gateway composition errors being resolved
- ⏳ Admin panel connection in progress
- ⏳ Services starting individually

**Estimated Time to Full Operation:** 6-8 hours for experienced developers

---

## 📋 **PREREQUISITES**

### Required Software
- [x] Docker Desktop (for Windows)
- [x] Go 1.21+ (installed)
- [x] Node.js 18+ (installed)
- [x] PostgreSQL client tools
- [x] Git

### Environment Setup
```powershell
# Navigate to project root
cd c:\Users\dane\OneDrive\Desktop\UNIFIED_COMMERCE

# Verify required tools
docker --version
go version
node --version
```

---

## 🔧 **STEP 1: START INFRASTRUCTURE (15 minutes)**

### Start Docker Services
```powershell
# Start all required infrastructure services
docker-compose up -d

# Verify services are running
docker ps

# Expected containers:
# - PostgreSQL (port 5432)
# - MongoDB (port 27017) 
# - Redis (port 6379)
# - Kafka (port 9092)
```

### Verify Database Connectivity
```powershell
# Test PostgreSQL connection
docker exec -it unified-commerce-postgres psql -U postgres -c "\l"

# Test MongoDB connection
docker exec -it unified-commerce-mongodb mongosh --eval "db.adminCommand('ismaster')"

# Test Redis connection
docker exec -it unified-commerce-redis redis-cli ping
```

**Troubleshooting:**
- If containers fail to start, check `docker-compose.yml` exists
- Ensure ports 5432, 27017, 6379, 9092 are not in use
- On Windows, verify Docker Desktop is running

---

## 🏗️ **STEP 2: BUILD ALL SERVICES (30 minutes)**

### Test Service Builds
```powershell
# Build each service to verify no compilation errors
$services = @("identity", "cart", "inventory", "product-catalog", "promotions", "merchant-account")

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

### Fix Known Issues
```powershell
# Order and Payment services may need manual review
# These were identified as having compilation issues

# Build Order service
cd services/order
go build ./cmd/server
cd ../..

# Build Payment service  
cd services/payment
go build ./cmd/server
cd ../..
```

**If builds fail:**
1. Check for missing dependencies: `go mod tidy`
2. Review recent manual edits in resolver files
3. Check import statements and package declarations

---

## 🚀 **STEP 3: START MICROSERVICES (2-3 hours)**

### Option A: Start Services Individually (Recommended for Testing)

```powershell
# Create environment variables for each service
$env:DB_HOST = "localhost"
$env:DB_PORT = "5432"
$env:DB_USER = "postgres"
$env:DB_PASSWORD = "postgres"
$env:REDIS_HOST = "localhost"
$env:REDIS_PORT = "6379"
$env:KAFKA_BROKERS = "localhost:9092"

# Start Identity Service (Port 8001)
$env:SERVICE_PORT = "8001"
$env:DB_NAME = "identity_db"
cd services/identity
Start-Process powershell -ArgumentList "-NoExit", "-Command", "go run ./cmd/server/main.go"
cd ../..

# Start Cart Service (Port 8002)
$env:SERVICE_PORT = "8002"
$env:DB_NAME = "cart_db"
cd services/cart
Start-Process powershell -ArgumentList "-NoExit", "-Command", "go run ./cmd/server/main.go"
cd ../..

# Continue for all services...
# Inventory (8005), Product-Catalog (8006), Promotions (8007), Merchant-Account (8008)
```

### Option B: Use PowerShell Script (If Available)
```powershell
# Check if startup script exists
if (Test-Path "scripts/start-services.ps1") {
    .\scripts\start-services.ps1
} else {
    Write-Host "Manual startup required - use Option A above"
}
```

### Verify Service Health
```powershell
# Check if services are responding
$ports = @(8001, 8002, 8005, 8006, 8007, 8008)

foreach ($port in $ports) {
    try {
        $response = Invoke-WebRequest -Uri "http://localhost:$port/health" -TimeoutSec 5
        Write-Host "✅ Service on port $port is responding"
    } catch {
        Write-Host "❌ Service on port $port is not responding"
    }
}
```

---

## 🌐 **STEP 4: FIX GRAPHQL FEDERATION ISSUES (2-4 hours)**

### Current Federation Issues
The GraphQL Federation Gateway is currently failing to compose the schema due to:

1. **Address Type Inconsistencies** - Different services have different Address type definitions
2. **Missing @key Directives** - Shared types need proper federation keys
3. **Transaction Type Conflicts** - Order and Payment services both define Transaction types

### Fix Steps

#### 4.1 Standardize Address Types
Ensure all services have consistent Address type definitions:

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
Remove the Transaction type from the order service since it should be managed by the payment service.

#### 4.3 Add Missing Federation Directives
If federation directives are not recognized, add them directly to schema files.

---

## 🌐 **STEP 5: START GRAPHQL FEDERATION GATEWAY (30 minutes)**

### Install Dependencies and Start Gateway
```powershell
cd gateway

# Install Node.js dependencies
npm install

# Start Apollo Federation Gateway
npm start

# Gateway should start on port 4000
# Check GraphQL Playground: http://localhost:4000/graphql
```

### Test Federation
```powershell
# Test individual service schemas
curl -X POST http://localhost:4000/graphql \
  -H "Content-Type: application/json" \
  -d '{"query":"{ __schema { types { name } } }"}'
```

---

## 🖥️ **STEP 6: CONNECT ADMIN PANEL (2-3 hours)**

### Update Admin Panel Configuration
1. Point admin panel to GraphQL Federation Gateway (http://localhost:4000/graphql)
2. Replace mock data with real GraphQL queries
3. Implement authentication flow

### Test Admin Panel
1. Verify login/logout functionality
2. Test CRUD operations with real data
3. Validate cross-service queries work correctly

---

## 📋 **TROUBLESHOOTING CHECKLIST**

### Infrastructure Issues
- [x] Docker containers running
- [x] Database connectivity verified
- [x] Redis connectivity verified
- [x] Kafka connectivity verified

### Service Issues
- [x] All services building successfully
- [x] Environment variables configured
- [ ] All services responding to health checks
- [ ] Services communicating properly

### Federation Issues
- [x] Federation directives implemented
- [x] Shared types defined
- [ ] Schema composition successful
- [ ] Cross-service queries working

### Frontend Issues
- [x] Admin panel UI complete
- [x] Authentication UI implemented
- [ ] Backend connection established
- [ ] Real data flow implemented

---

## 📞 **SUPPORT RESOURCES**

### Documentation
- [Troubleshooting Guide](./TROUBLESHOOTING_GUIDE.md) - Detailed issue resolution
- [Architecture Documentation](./architecture.md) - System design overview
- [Build Completion Plan](./BUILD_COMPLETION_PLAN.md) - Long-term roadmap

### Common Issues
1. **Services won't start**: Check Docker containers and environment variables
2. **Federation fails**: Review shared type definitions and @key directives
3. **Admin panel not connecting**: Verify GraphQL Gateway is running and accessible

### Emergency Contacts
For critical issues, contact the development team or refer to the troubleshooting guide.