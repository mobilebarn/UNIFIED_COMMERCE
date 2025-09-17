# Retail OS Railway CLI Service Creation Script
# Execute these commands one by one to create all services first

Write-Host "🚀 Creating Retail OS Services on Railway" -ForegroundColor Green
Write-Host "==========================================" -ForegroundColor Green

# First, create all empty services
Write-Host "📦 Creating empty services..." -ForegroundColor Yellow

# Create Identity Service
Write-Host "Creating Identity Service..." -ForegroundColor Cyan
railway add --service --variables "SERVICE_NAME=identity" --variables "ENVIRONMENT=production"

# Create Merchant Service
Write-Host "Creating Merchant Service..." -ForegroundColor Cyan
railway add --service --variables "SERVICE_NAME=merchant" --variables "ENVIRONMENT=production"

# Create Product Service
Write-Host "Creating Product Service..." -ForegroundColor Cyan
railway add --service --variables "SERVICE_NAME=product" --variables "ENVIRONMENT=production"

# Create Inventory Service
Write-Host "Creating Inventory Service..." -ForegroundColor Cyan
railway add --service --variables "SERVICE_NAME=inventory" --variables "ENVIRONMENT=production"

# Create Order Service
Write-Host "Creating Order Service..." -ForegroundColor Cyan
railway add --service --variables "SERVICE_NAME=order" --variables "ENVIRONMENT=production"

# Create Payment Service (Port 8005)
Write-Host "Creating Payment Service..." -ForegroundColor Cyan
railway add --service --variables "SERVICE_NAME=payment" --variables "SERVICE_PORT=8005" --variables "ENVIRONMENT=production"

# Create Cart Service
Write-Host "Creating Cart Service..." -ForegroundColor Cyan
railway add --service --variables "SERVICE_NAME=cart" --variables "ENVIRONMENT=production"

# Create Promotions Service
Write-Host "Creating Promotions Service..." -ForegroundColor Cyan
railway add --service --variables "SERVICE_NAME=promotions" --variables "ENVIRONMENT=production"

# Create Analytics Service
Write-Host "Creating Analytics Service..." -ForegroundColor Cyan
railway add --service --variables "SERVICE_NAME=analytics" --variables "ENVIRONMENT=production"

# Create GraphQL Gateway
Write-Host "Creating GraphQL Gateway..." -ForegroundColor Cyan
railway add --service --variables "NODE_ENV=production"

# Create Storefront
Write-Host "Creating Storefront..." -ForegroundColor Cyan
railway add --service --variables "NODE_ENV=production"

# Create Admin Panel
Write-Host "Creating Admin Panel..." -ForegroundColor Cyan
railway add --service --variables "NODE_ENV=production"

Write-Host "✅ All services created! Now we can deploy to them." -ForegroundColor Green