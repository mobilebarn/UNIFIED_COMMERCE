#!/usr/bin/env pwsh

Write-Host "Creating service environment files..." -ForegroundColor Yellow

$services = @(
    @{Name="identity"; Port="8001"; Database="identity_service"},
    @{Name="cart"; Port="8002"; Database="cart_service"},
    @{Name="order"; Port="8003"; Database="order_service"},
    @{Name="payment"; Port="8004"; Database="payment_service"},
    @{Name="inventory"; Port="8005"; Database="inventory_service"},
    @{Name="product-catalog"; Port="8006"; Database="product_catalog"},
    @{Name="promotions"; Port="8007"; Database="promotions_service"},
    @{Name="merchant-account"; Port="8008"; Database="merchant_account_service"}
)

foreach ($service in $services) {
    $serviceName = $service.Name
    $servicePort = $service.Port
    $databaseName = $service.Database
    
    $envContent = @"
# $serviceName Service Configuration

# Service Configuration
SERVICE_NAME=$serviceName
SERVICE_PORT=$servicePort
ENVIRONMENT=development

# Database Configuration
DATABASE_HOST=localhost
DATABASE_PORT=5432
DATABASE_USER=postgres
DATABASE_PASSWORD=postgres
DATABASE_NAME=$databaseName

# MongoDB Configuration (for product-catalog service)
MONGO_URL=mongodb://mongodb:mongodb@localhost:27017
MONGO_DATABASE=product_catalog
MONGO_USER=mongodb
MONGO_PASSWORD=mongodb

# Redis Configuration
REDIS_URL=redis://localhost:6379
REDIS_PASSWORD=redis
REDIS_DB=0

# JWT Configuration
JWT_SECRET=your-super-secure-jwt-secret-key-for-development
JWT_EXPIRATION=86400

# Kafka Configuration
KAFKA_BROKERS=localhost:9092

# External Services
ELASTICSEARCH_URL=http://localhost:9200
JAEGER_ENDPOINT=http://localhost:14268/api/traces

# GraphQL Gateway
GATEWAY_URL=http://localhost:4000

# Logging
LOG_LEVEL=debug
"@

    $envFile = "services\$serviceName\.env"
    $envContent | Out-File -FilePath $envFile -Encoding UTF8
    Write-Host "Created $envFile" -ForegroundColor Green
}

Write-Host "All service environment files created successfully!" -ForegroundColor Green
