# Unified Commerce Platform

A next-generation commerce operating system that truly unifies online and offline retail operations through a modern microservices architecture.

## 📊 Project Status
**Current Phase:** Backend Complete ✅ | Frontend Development In Progress 🔄  
**Overall Completion:** ~65%  
**For detailed status:** See [`UNIFIED_IMPLEMENTATION_STATUS.md`](./UNIFIED_IMPLEMENTATION_STATUS.md)

## � Quick Overview

### What's Complete ✅
- **8 Microservices** with GraphQL Federation (Identity, Cart, Order, Payment, Inventory, Product Catalog, Promotions, Merchant Account)
- **Apollo Federation Gateway** on port 4000 with unified GraphQL schema
- **Authentication System** with JWT tokens and protected routes
- **React Admin Panel** with working login and dashboard
- **Infrastructure** with Docker, PostgreSQL, MongoDB, Redis, Kafka

### What's In Progress 🔄
- **Full Admin Panel** - Need CRUD operations for all entities
- **Customer Storefront** - Next.js structure exists, needs implementation
- **Mobile POS** - Directory structure exists, needs development

## 🚀 Quick Start

### Prerequisites
- Docker and Docker Compose
- Node.js 18+
- Go 1.21+
- PowerShell (Windows)

### 1. Start Infrastructure
```powershell
docker-compose up -d
```

### 2. Start All Backend Services
```powershell
.\scripts\start-services.ps1 -All
```

### 3. Start GraphQL Gateway
```powershell
cd gateway
npm install
npm start
```

### 4. Start Admin Panel
```powershell
cd admin-panel-new
npm install
npm run dev
```

### 5. Access Points
- **GraphQL Federation:** http://localhost:4000/graphql
- **GraphQL Playground:** http://localhost:4000/playground
- **Admin Panel:** http://localhost:3003
- **Individual Services:** http://localhost:8001-8008

## 🏗️ Architecture Overview

### System Design
```
Frontend Applications (3003)
     ↓ HTTP/GraphQL
GraphQL Federation Gateway (4000)
     ↓ GraphQL Federation
┌─────────────────────────────────────┐
│ Microservices (8001-8008)           │
│ ├─ Identity (Auth/Users)            │
│ ├─ Cart (Shopping Carts)            │
│ ├─ Order (Order Management)         │
│ ├─ Payment (Payment Processing)     │
│ ├─ Inventory (Stock Management)     │
│ ├─ Product Catalog (Products)       │
│ ├─ Promotions (Discounts/Coupons)   │
│ └─ Merchant Account (Business Mgmt) │
└─────────────────────────────────────┘
     ↓ Database Connections
┌─────────────────────────────────────┐
│ Databases                           │
│ ├─ PostgreSQL (Primary)             │
│ ├─ MongoDB (Product Catalog)        │
│ ├─ Redis (Session/Cache)            │
│ └─ Kafka (Message Queue)            │
└─────────────────────────────────────┘
```

### Key Features
- **GraphQL Federation** - Single endpoint exposing all microservice functionality
- **JWT Authentication** - Secure authentication with context forwarding
- **Database-per-Service** - Polyglot persistence with PostgreSQL and MongoDB
- **Event-Driven Architecture** - Kafka for inter-service communication
- **Cloud-Native** - Docker containerization ready for Kubernetes

## 📁 Project Structure

```
unified-commerce/
├── 🏗️ Backend (Complete)
│   ├── gateway/                # GraphQL Federation Gateway
│   ├── services/               # 8 Microservices
│   │   ├── identity/          # Authentication & users
│   │   ├── cart/              # Shopping cart management
│   │   ├── order/             # Order processing
│   │   ├── payment/           # Payment processing
│   │   ├── inventory/         # Stock management
│   │   ├── product-catalog/   # Product data (MongoDB)
│   │   ├── promotions/        # Discounts & campaigns
│   │   └── merchant-account/  # Business accounts
│   └── shared/                # Common utilities
├── 🖥️ Frontend (Partial)
│   ├── admin-panel-new/       # React Admin (Working auth)
│   ├── storefront/            # Next.js Store (Skeleton)
│   └── mobile-pos/            # Mobile App (Skeleton)
├── 🐳 Infrastructure
│   ├── docker-compose.yml     # Local development
│   ├── infrastructure/        # K8s, Terraform, monitoring
│   └── scripts/               # PowerShell automation
└── 📚 Documentation
    ├── README.md              # This file
    ├── PROJECT_SUMMARY.md     # High-level accomplishments
    ├── UNIFIED_IMPLEMENTATION_STATUS.md  # Detailed status
    └── architecture.md        # Technical architecture
```

## 🧪 API Testing

### Authentication Flow
```bash
# 1. Login to get JWT token
curl -X POST http://localhost:8001/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@example.com",
    "password": "Admin123!"
  }'

# 2. Use token for authenticated requests
curl -X GET http://localhost:8001/api/v1/users/profile \
  -H "Authorization: Bearer YOUR_JWT_TOKEN_HERE"
```

### GraphQL Federation Testing
```graphql
# Single query spanning multiple services
query UnifiedQuery {
  user(id: "1") {
    id
    email
    firstName
    lastName
    cart {
      id
      items {
        quantity
        product {
          title
          price
        }
      }
    }
    orders {
      id
      status
      total
      payments {
        status
        amount
      }
    }
  }
  
  products(filter: { limit: 5 }) {
    id
    title
    status
    variants {
      sku
      price
      inventory {
        quantity
        location
      }
    }
  }
}
```

### Service Health Checks
```bash
# Check all service health
curl http://localhost:8001/health  # Identity
curl http://localhost:8002/health  # Cart
curl http://localhost:8003/health  # Order
curl http://localhost:8004/health  # Payment
curl http://localhost:8005/health  # Inventory
curl http://localhost:8006/health  # Product Catalog
curl http://localhost:8007/health  # Promotions
curl http://localhost:8008/health  # Merchant Account
```

## 🎯 Core Entities

### User & Authentication
- **Users** - Customer and merchant accounts
- **Roles** - Role-based access control
- **Sessions** - JWT token management
- **Audit Logs** - Authentication event tracking

### Commerce Core
- **Products** - Catalog with variants and categories
- **Inventory** - Multi-location stock management
- **Cart** - Shopping cart with user relationships
- **Orders** - Order lifecycle and fulfillment
- **Payments** - Payment processing and transactions

### Business Management
- **Merchants** - Business accounts and profiles
- **Stores** - Physical and online store locations
- **Promotions** - Discounts, campaigns, loyalty programs
- **Subscriptions** - Business plan management

## 🔧 Development

### Prerequisites
- Docker and Docker Compose installed
- Go 1.21+ for backend services
- Node.js 18+ for frontend applications
- PostgreSQL client (optional, for direct DB access)

### Environment Setup
1. **Clone repository**
2. **Copy environment files:**
   ```powershell
   .\scripts\create-env-files.ps1
   ```
3. **Start infrastructure:**
   ```powershell
   docker-compose up -d
   ```
4. **Initialize databases:**
   ```powershell
   .\scripts\init-databases.sql
   ```

### Building Services
```powershell
# Build all services
cd services/identity && go build ./...
cd services/cart && go build ./...
cd services/order && go build ./...
cd services/payment && go build ./...
cd services/inventory && go build ./...
cd services/product-catalog && go build ./...
cd services/promotions && go build ./...
cd services/merchant-account && go build ./...
```

### Testing
```powershell
# Run service tests
cd services/identity && go test ./...

# Run integration tests
cd services/identity && go test ./integration -v
```

## 📈 Performance & Scalability

### Current Capabilities
- **GraphQL Federation** - Optimized queries across services
- **Database Sharding** - Ready for horizontal scaling
- **Caching Strategy** - Redis for session and query caching
- **Event Processing** - Kafka for asynchronous operations

### Production Readiness
- **Containerization** - All services dockerized
- **Health Checks** - Comprehensive service monitoring
- **Logging** - Structured logging across services
- **Security** - JWT authentication with proper validation

## 🚀 Deployment

### Local Development
```powershell
# Complete local setup
.\scripts\start-services.ps1 -All
cd gateway && npm start
cd admin-panel-new && npm run dev
```

### Production (Kubernetes)
```bash
# Deploy to Kubernetes cluster
kubectl apply -f infrastructure/k8s/

# Install via Helm
helm install unified-commerce infrastructure/helm/
```

## 📚 Documentation

- **[Project Summary](./PROJECT_SUMMARY.md)** - High-level accomplishments and vision
- **[Implementation Status](./UNIFIED_IMPLEMENTATION_STATUS.md)** - Detailed current status and roadmap
- **[Architecture Guide](./architecture.md)** - Technical architecture and design decisions

## 🤝 Contributing

### Development Workflow
1. **Backend Changes** - Update Go services, test locally
2. **Frontend Changes** - Update React/Next.js apps
3. **Schema Changes** - Update GraphQL schemas and regenerate
4. **Database Changes** - Create migration scripts
5. **Documentation** - Update relevant documentation

### Code Quality
- **Go** - Follow standard Go practices, use `gofmt`
- **JavaScript/TypeScript** - ESLint and Prettier configured
- **GraphQL** - Schema-first development with gqlgen
- **Testing** - Unit tests for business logic, integration tests for APIs

## 🏆 Vision & Goals

### Competitive Advantages
- **Unified Commerce** - True omnichannel capabilities
- **Developer Experience** - GraphQL Federation with rich tooling
- **Scalability** - Microservices architecture ready for enterprise scale
- **Flexibility** - Headless-first design for unlimited customization

### Market Position
- **vs. Shopify** - More flexible, true unified commerce
- **vs. Square** - Better online capabilities, open architecture
- **vs. Others** - Modern tech stack, microservices, GraphQL Federation

---

*Last Updated: August 30, 2025*  
*Status: Backend Complete, Frontend Development In Progress*