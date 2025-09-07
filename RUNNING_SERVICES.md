# Running Services Summary

## Backend Services

All 8 microservices are running and connected to the GraphQL Federation Gateway:

1. **Identity Service** - Port 8001
   - Health Check: http://localhost:8001/health
   - GraphQL Endpoint: http://localhost:8001/graphql

2. **Cart Service** - Port 8002
   - Health Check: http://localhost:8002/health
   - GraphQL Endpoint: http://localhost:8002/graphql

3. **Order Service** - Port 8003
   - Health Check: http://localhost:8003/health
   - GraphQL Endpoint: http://localhost:8003/graphql

4. **Payment Service** - Port 8004
   - Health Check: http://localhost:8004/health
   - GraphQL Endpoint: http://localhost:8004/graphql

5. **Inventory Service** - Port 8005
   - Health Check: http://localhost:8005/health
   - GraphQL Endpoint: http://localhost:8005/graphql

6. **Product Catalog Service** - Port 8006
   - Health Check: http://localhost:8006/health
   - GraphQL Endpoint: http://localhost:8006/graphql

7. **Promotions Service** - Port 8007
   - Health Check: http://localhost:8007/health
   - GraphQL Endpoint: http://localhost:8007/graphql

8. **Merchant Account Service** - Port 8008
   - Health Check: http://localhost:8008/health
   - GraphQL Endpoint: http://localhost:8008/graphql

## GraphQL Federation Gateway

- **Gateway URL**: http://localhost:4000/graphql
- **Health Check**: http://localhost:4000/health
- **GraphQL Playground**: http://localhost:4000/graphql

All 8 services are successfully federated and accessible through the single GraphQL endpoint.

## Frontend Applications

1. **Admin Panel** (React)
   - URL: http://localhost:3004/
   - UI Complete: Yes
   - Connected to GraphQL Federation Gateway: Yes
   - Real data integration: Yes (partial, transitioning from mock data)

2. **Storefront** (Next.js)
   - URL: http://localhost:3002/
   - Connected to GraphQL Federation Gateway: Yes
   - Real product data: Yes

## Infrastructure Services

All infrastructure services are running in Docker containers:

- **PostgreSQL**: Port 5432
- **MongoDB**: Port 27017
- **Redis**: Port 6379
- **Kafka**: Port 9092
- **Zookeeper**: Port 2181
- **Prometheus**: Port 9090
- **Grafana**: Port 3001
- **Elasticsearch**: Port 9200
- **Kibana**: Port 5601
- **Logstash**: Port 5044

## Access Points Summary

| Service | URL | Port |
|---------|-----|------|
| GraphQL Gateway | http://localhost:4000/graphql | 4000 |
| GraphQL Health | http://localhost:4000/health | 4000 |
| Admin Panel | http://localhost:3004/ | 3004 |
| Storefront | http://localhost:3002/ | 3002 |
| Identity Service | http://localhost:8001/ | 8001 |
| Cart Service | http://localhost:8002/ | 8002 |
| Order Service | http://localhost:8003/ | 8003 |
| Payment Service | http://localhost:8004/ | 8004 |
| Inventory Service | http://localhost:8005/ | 8005 |
| Product Catalog | http://localhost:8006/ | 8006 |
| Promotions | http://localhost:8007/ | 8007 |
| Merchant Account | http://localhost:8008/ | 8008 |
| Grafana | http://localhost:3001/ | 3001 |

## Testing GraphQL Federation

You can test the unified GraphQL API with a query like:

```graphql
query {
  products {
    id
    title
    price
  }
}
```

This query will be resolved across multiple services through the federation gateway.