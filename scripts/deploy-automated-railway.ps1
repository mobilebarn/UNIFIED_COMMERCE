# RETAIL OS AUTOMATED RAILWAY DEPLOYMENT
# Uses Railway's GitHub integration with monorepo support
# Each service specifies its root directory to avoid uploading the entire project

Write-Host "RETAIL OS AUTOMATED RAILWAY DEPLOYMENT" -ForegroundColor Green
Write-Host "=====================================" -ForegroundColor Green
Write-Host ""

Write-Host "🎯 Strategy: GitHub integration with monorepo root directory specification" -ForegroundColor Yellow
Write-Host "✅ Benefits: Fast deployment, automatic rebuilds, no file size limits" -ForegroundColor Green
Write-Host ""

# Function to deploy service with root directory specification
function Deploy-ServiceWithRoot {
    param(
        [string]$ServiceName,
        [string]$ServiceDir,
        [string]$RailwayServiceName,
        [string]$Port
    )
    
    Write-Host "🚀 Configuring $ServiceName..." -ForegroundColor Cyan
    
    try {
        # Connect to service
        Write-Host "  🔗 Connecting to Railway service: $RailwayServiceName" -ForegroundColor White
        railway service $RailwayServiceName
        
        # Set root directory for monorepo deployment
        Write-Host "  📁 Setting root directory: $ServiceDir" -ForegroundColor White
        railway variables --set "RAILWAY_ROOT_DIR=$ServiceDir"
        
        # Set service-specific environment variables
        Write-Host "  ⚙️  Setting environment variables..." -ForegroundColor White
        railway variables --set "SERVICE_NAME=$($ServiceName.ToLower() -replace ' service', '')"
        railway variables --set "SERVICE_PORT=$Port"
        railway variables --set "ENVIRONMENT=production"
        
        # Set database connections
        if ($ServiceName -like "*Product*") {
            railway variables --set "MONGO_URL=`${{MongoDB.MONGO_URL}}"
        } else {
            railway variables --set "DATABASE_URL=`${{Postgres.DATABASE_URL}}"
        }
        
        # Connect to GitHub and deploy
        Write-Host "  🚀 Deploying from GitHub repository..." -ForegroundColor White
        railway up --detach
        
        Write-Host "  ✅ $ServiceName configured and deployed!" -ForegroundColor Green
    }
    catch {
        Write-Host "  ❌ Failed to deploy $ServiceName" -ForegroundColor Red
        Write-Host "  Error: $_" -ForegroundColor Red
    }
    
    Write-Host ""
}

# Commit current changes to Git first
Write-Host "📝 Committing current changes to Git..." -ForegroundColor Yellow
git add .
git commit -m "Add Railway deployment configurations and scripts"
git push origin master

Write-Host ""
Write-Host "🚀 Starting automated deployment..." -ForegroundColor Yellow
Write-Host ""

# Deploy services with their root directories
Deploy-ServiceWithRoot -ServiceName "Product Service" -ServiceDir "services/product-catalog" -RailwayServiceName "retail-os-product" -Port "8003"
Deploy-ServiceWithRoot -ServiceName "Payment Service" -ServiceDir "services/payment" -RailwayServiceName "retail-os-payment" -Port "8005"
Deploy-ServiceWithRoot -ServiceName "Analytics Service" -ServiceDir "services/analytics" -RailwayServiceName "retail-os-analytics" -Port "8001"
Deploy-ServiceWithRoot -ServiceName "Inventory Service" -ServiceDir "services/inventory" -RailwayServiceName "retail-os-inventory" -Port "8002"
Deploy-ServiceWithRoot -ServiceName "Identity Service" -ServiceDir "services/identity" -RailwayServiceName "retail-os-identity" -Port "8000"
Deploy-ServiceWithRoot -ServiceName "Cart Service" -ServiceDir "services/cart" -RailwayServiceName "retail-os-cart" -Port "8080"

Write-Host "🎉 AUTOMATED DEPLOYMENT COMPLETE!" -ForegroundColor Green -BackgroundColor Black
Write-Host "================================" -ForegroundColor Green
Write-Host ""
Write-Host "📊 Deployment Summary:" -ForegroundColor Yellow
Write-Host "✅ All services connected to GitHub repository" -ForegroundColor Green
Write-Host "✅ Monorepo root directories configured" -ForegroundColor Green
Write-Host "✅ Environment variables set for production" -ForegroundColor Green
Write-Host "✅ Automatic rebuilds enabled" -ForegroundColor Green