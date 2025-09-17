# Docker Containerization Guide for Retail OS

## Overview
This document provides instructions for containerizing and deploying Retail OS using Docker. The platform consists of 8 microservices, a GraphQL Federation Gateway, and two frontend applications (storefront and admin panel).

## Architecture
The platform is composed of the following Docker services:

### Infrastructure Services (docker-compose.yml)
- PostgreSQL (5432)
- MongoDB (27017)
- Redis (6379)
- Elasticsearch (9200/9300)
- Kibana (5601)
- Logstash (5000/5044/9600)
- Kafka (9092)
- Zookeeper (2181)
- Prometheus (9090)
- Grafana (3001)

### Application Services (docker-compose.services.yml)
- Identity Service (8001)
- Merchant Account Service (8008)
- Product Catalog Service (8006)
- Inventory Service (8005)
- Order Service (8003)
- Cart Service (8002)
- Payment Service (8004)
- Promotions Service (8007)
- GraphQL Federation Gateway (4000)
- Storefront Application (3000)
- Admin Panel Application (3001)

## Prerequisites
- Docker Engine 20.10+
- Docker Compose 1.29+
- At least 8GB RAM available
- 20GB free disk space

## Deployment Steps

### 1. Clone the Repository
```bash
git clone <repository-url>
cd retail-os
```

### 2. Environment Configuration
Ensure all services have proper `.env` files with required configuration variables.

### 3. Start Infrastructure Services
```bash
docker-compose up -d
```

This will start all infrastructure services including databases, message queues, and monitoring tools.

### 4. Start Application Services
```bash
docker-compose -f docker-compose.services.yml up -d
```

This will build and start all application services including microservices, gateway, and frontend applications.

### 5. Monitor Services
```bash
# View logs for all services
docker-compose -f docker-compose.services.yml logs -f

# View logs for specific service
docker-compose -f docker-compose.services.yml logs -f <service-name>
```

## Service-Specific Dockerfiles

### Microservices (Go Applications)
All microservices use a multi-stage build approach:
1. Builder stage with Go dependencies
2. Production stage with minimal Alpine image
3. Non-root user for security
4. Health checks for monitoring

### Frontend Applications

#### Storefront (Next.js)
- Multi-stage build with Node.js builder
- Production stage with optimized build
- Non-root user for security
- Health checks for monitoring

#### Admin Panel (React/Vite)
- Multi-stage build with Node.js builder
- Production stage with Nginx server
- Static file serving
- Health checks for monitoring

### GraphQL Gateway (Node.js)
- Multi-stage build with Node.js
- Production dependencies only
- Non-root user for security
- Health checks for monitoring

## Environment Variables

### Infrastructure Services
Environment variables are configured in the docker-compose.yml file.

### Application Services
Each service requires specific environment variables:
- PORT: Service port
- DATABASE_URL: Database connection string
- REDIS_URL: Redis connection string
- KAFKA_BROKERS: Kafka broker addresses
- JWT_SECRET: JWT secret key for authentication
- ENVIRONMENT: Environment (development/production)

## Health Checks
All services include health checks:
- Infrastructure services: Built-in health endpoints
- Application services: Custom health endpoints
- Frontend services: HTTP endpoint checks

## Scaling Services
Individual services can be scaled using Docker Compose:
```bash
docker-compose -f docker-compose.services.yml up -d --scale <service-name>=3
```

## Monitoring and Logging
- Prometheus collects metrics from all services
- Grafana provides dashboards for monitoring
- ELK stack (Elasticsearch, Logstash, Kibana) handles logging
- Docker logs for real-time service logs

## Security Considerations
- All services run as non-root users
- Secrets are managed through environment variables
- Network isolation through Docker networks
- Health checks for service availability

## Backup and Recovery
- PostgreSQL data is persisted in named volumes
- MongoDB data is persisted in named volumes
- Redis data is persisted in named volumes
- Regular backups should be implemented for production

## Troubleshooting

### Common Issues
1. Port conflicts: Ensure ports are not in use by other applications
2. Memory issues: Allocate sufficient RAM to Docker
3. Build failures: Check Dockerfile syntax and dependencies
4. Connection issues: Verify network configuration and service dependencies

### Debugging Commands
```bash
# List running containers
docker-compose -f docker-compose.services.yml ps

# View container logs
docker-compose -f docker-compose.services.yml logs <service-name>

# Execute commands in container
docker-compose -f docker-compose.services.yml exec <service-name> sh

# Restart specific service
docker-compose -f docker-compose.services.yml restart <service-name>
```

## Production Considerations
1. Use production-ready database configurations
2. Implement proper secret management (HashiCorp Vault, AWS Secrets Manager)
3. Configure load balancing and reverse proxies
4. Set up automated backups
5. Implement monitoring alerts
6. Use container orchestration (Kubernetes) for production deployments
7. Configure SSL/TLS for secure communication
8. Implement proper logging and log rotation

## Cleanup
To stop and remove all services:
```bash
docker-compose -f docker-compose.services.yml down
docker-compose down
```

To remove volumes and data:
```bash
docker-compose -f docker-compose.services.yml down -v
docker-compose down -v
```