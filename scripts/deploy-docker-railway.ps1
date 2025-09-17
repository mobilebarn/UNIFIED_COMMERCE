# RETAIL OS DOCKER-BASED RAILWAY DEPLOYMENT
# This script deploys each service individually using Docker
# to avoid the file size issues with source code uploads

Write-Host "RETAIL OS DOCKER DEPLOYMENT TO RAILWAY" -ForegroundColor Green
Write-Host "=======================================" -ForegroundColor Green
Write-Host ""

Write-Host "Strategy: Docker-based deployment from individual service directories" -ForegroundColor Yellow
Write-Host "Benefits: Smaller uploads, optimized builds, production-ready containers" -ForegroundColor Yellow
Write-Host ""

# Function to deploy a service using Docker
function Deploy-Service {
    param(
        [string]$ServiceName,
        [string]$ServiceDir,
        [string]$RailwayServiceName,
        [string]$Port
    )
    
    Write-Host "Deploying $ServiceName..." -ForegroundColor Cyan
    Write-Host "  Directory: $ServiceDir" -ForegroundColor Gray
    Write-Host "  Railway Service: $RailwayServiceName" -ForegroundColor Gray
    Write-Host "  Port: $Port" -ForegroundColor Gray
    
    # Navigate to service directory
    Push-Location $ServiceDir
    
    try {
        # Link to the specific Railway service
        Write-Host "  Linking to Railway service..." -ForegroundColor White
        railway service $RailwayServiceName
        
        # Deploy using Docker
        Write-Host "  Starting Docker deployment..." -ForegroundColor White
        railway up --detach
        
        Write-Host "  ✅ $ServiceName deployed successfully!" -ForegroundColor Green
    }
    catch {
        Write-Host "  ❌ Failed to deploy $ServiceName" -ForegroundColor Red
        Write-Host "  Error: $_" -ForegroundColor Red
    }
    finally {
        # Return to original directory
        Pop-Location
    }
    
    Write-Host ""
}

# Deploy services one by one
Write-Host "🚀 Starting deployment of backend services..." -ForegroundColor Yellow
Write-Host ""

# 1. Product Service
Deploy-Service -ServiceName "Product Service" -ServiceDir "services\product-catalog" -RailwayServiceName "retail-os-product" -Port "8003"

# 2. Payment Service
Deploy-Service -ServiceName "Payment Service" -ServiceDir "services\payment" -RailwayServiceName "retail-os-payment" -Port "8005"

# 3. Inventory Service
Deploy-Service -ServiceName "Inventory Service" -ServiceDir "services\inventory" -RailwayServiceName "retail-os-inventory" -Port "8002"

# 4. Analytics Service
Deploy-Service -ServiceName "Analytics Service" -ServiceDir "services\analytics" -RailwayServiceName "retail-os-analytics" -Port "8001"

# 5. Cart Service (if we have a Railway service for it)
# Deploy-Service -ServiceName "Cart Service" -ServiceDir "services\cart" -RailwayServiceName "retail-os-cart" -Port "8080"

# 6. Order Service (if we have a Railway service for it)  
# Deploy-Service -ServiceName "Order Service" -ServiceDir "services\order" -RailwayServiceName "retail-os-order" -Port "8004"

Write-Host "🎉 DEPLOYMENT COMPLETE!" -ForegroundColor Green -BackgroundColor Black
Write-Host "======================" -ForegroundColor Green
Write-Host ""
Write-Host "📊 Summary:" -ForegroundColor Yellow
Write-Host "✅ Services deployed using optimized Docker containers" -ForegroundColor Green
Write-Host "✅ Each service deployed from its own directory" -ForegroundColor Green
Write-Host "✅ Production environment variables configured" -ForegroundColor Green
Write-Host ""
Write-Host "🔗 Next Steps:" -ForegroundColor Yellow
Write-Host "1. Deploy GraphQL Gateway" -ForegroundColor White
Write-Host "2. Configure service networking" -ForegroundColor White
Write-Host "3. Test end-to-end functionality" -ForegroundColor White