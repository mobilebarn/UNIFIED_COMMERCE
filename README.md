# Retail OS

A complete, production-ready e-commerce platform built with modern microservices architecture, GraphQL federation, and cutting-edge frontend technologies.

## 🏗️ Architecture Overview

This platform consists of 8 independent microservices connected through a GraphQL Federation Gateway, providing a unified API for frontend applications:

### Backend Microservices

| Service | Port | Description |
|---------|------|-------------|
| Identity | 8001 | Authentication, authorization, and user management |
| Cart | 8002 | Shopping cart management and checkout workflows |
| Order | 8003 | Order processing and lifecycle management |
| Payment | 8004 | Payment processing and transaction management |
| Inventory | 8005 | Real-time inventory tracking across locations |
| Product Catalog | 8006 | Product information management with flexible schemas |
| Promotions | 8007 | Discount codes, sales, and promotional campaigns |
| Merchant Account | 8008 | Merchant profiles, subscriptions, and billing |

### GraphQL Federation Gateway

All services are unified through a GraphQL Federation Gateway running on port 4000, providing:
- Single endpoint for all API requests
- Cross-service relationships and entity resolution
- Real-time data fetching with powerful querying capabilities

### Frontend Applications

1. **Admin Panel** - React-based dashboard for business management (port 5173)
2. **Storefront** - Next.js e-commerce frontend (port 3000)

### Infrastructure

- PostgreSQL (Primary database)
- MongoDB (Flexible document storage)
- Redis (Caching and session management)
- Kafka (Event streaming and messaging)

## ✅ Current Status

All core components are fully implemented and operational:

- ✅ All 8 microservices built and running
- ✅ GraphQL Federation Gateway connecting all services
- ✅ Admin panel with real data integration
- ✅ Storefront with real product data
- ✅ Docker containerization for all services
- ✅ Comprehensive documentation

## 🚀 Getting Started

### Prerequisites

- Docker and Docker Compose
- Go 1.19+
- Node.js 16+
- npm or yarn

### Quick Start

1. **Start infrastructure services:**
   ```bash
   docker-compose up -d
   ```

2. **Start all microservices:**
   ```bash
   # In PowerShell
   .\start-all-services.ps1
   ```

3. **Start GraphQL Federation Gateway:**
   ```bash
   cd gateway
   npm start
   ```

4. **Start Admin Panel:**
   ```bash
   cd admin-panel-new
   npm run dev
   ```

5. **Start Storefront:**
   ```bash
   cd storefront
   npm run dev
   ```

### Access Points

- **GraphQL Playground:** http://localhost:4000/graphql
- **Admin Panel:** http://localhost:5173
- **Storefront:** http://localhost:3000
- **Health Check:** http://localhost:4000/health

## 📚 Documentation

- [Progress Summary](RETAIL_OS_PROGRESS_SUMMARY.md) - Complete status of implementation
- [GraphQL Federation Guide](docs/GRAPHQL_FEDERATION_GUIDE.md) - Detailed federation implementation
- [Troubleshooting Guide](docs/TROUBLESHOOTING_GUIDE.md) - Common issues and solutions
- [Startup Guide](docs/STARTUP_GUIDE.md) - How to start all services
- [Implementation Status](docs/UNIFIED_IMPLEMENTATION_STATUS.md) - Technical implementation details

## 🛠️ Development

### Project Structure

```
RETAIL_OS/
├── services/              # Go microservices
│   ├── identity/
│   ├── cart/
│   ├── order/
│   ├── payment/
│   ├── inventory/
│   ├── product-catalog/
│   ├── promotions/
│   └── merchant-account/
├── gateway/               # GraphQL Federation Gateway
├── admin-panel-new/       # React admin dashboard
├── storefront/            # Next.js storefront
├── infrastructure/        # Docker configurations
├── docs/                  # Documentation
└── scripts/               # Utility scripts
```

### Building Services

Each service can be built independently:
```bash
cd services/[service-name]
go build
```

### Testing GraphQL Federation

Test the unified GraphQL API:
```bash
curl -X POST \
  -H "Content-Type: application/json" \
  --data '{"query":"{ __schema { types { name } } }"}' \
  http://localhost:4000/graphql
```

## 🎯 Next Steps

With the core platform complete, we're focusing on:

1. **Enhancing Frontend Applications**
   - Full CRUD operations in admin panel
   - Complete shopping experience in storefront
   - Server-side rendering optimizations

2. **Kubernetes Deployment**
   - Helm charts for all services
   - Production-ready configurations
   - CI/CD pipeline implementation

3. **Observability**
   - Centralized logging
   - Metrics collection with Prometheus
   - Distributed tracing with OpenTelemetry

## 🤝 Contributing

This is a solo project developed as part of a comprehensive learning experience. The codebase follows enterprise standards and best practices.

## 📄 License

This project is for educational and demonstration purposes only.

## 📞 Support

For questions about the implementation or architecture, please refer to the documentation files in the [docs](docs/) directory.