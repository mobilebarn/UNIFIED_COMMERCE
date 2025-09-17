# RETAIL OS GIT-BASED RAILWAY DEPLOYMENT
# This script creates individual git repos for each service and deploys them
# This bypasses the file size limitation by using Git instead of direct uploads

Write-Host "RETAIL OS GIT-BASED RAILWAY DEPLOYMENT" -ForegroundColor Green
Write-Host "=====================================" -ForegroundColor Green
Write-Host ""

# Function to create a minimal service repo and deploy
function Deploy-ServiceViaGit {
    param(
        [string]$ServiceName,
        [string]$ServiceDir,
        [string]$RailwayServiceName,
        [string]$Port
    )
    
    Write-Host "🚀 Deploying $ServiceName via Git..." -ForegroundColor Cyan
    
    # Create a temporary deployment directory
    $TempDir = "temp_deploy_$ServiceName"
    New-Item -ItemType Directory -Path $TempDir -Force | Out-Null
    
    try {
        # Copy only essential files to temp directory
        Write-Host "  📁 Creating minimal deployment package..." -ForegroundColor White
        
        # Copy service source code
        Copy-Item -Path "$ServiceDir\*" -Destination $TempDir -Recurse -Force
        
        # Copy shared dependencies if they exist
        if (Test-Path "shared") {
            New-Item -ItemType Directory -Path "$TempDir\shared" -Force | Out-Null
            Copy-Item -Path "shared\*" -Destination "$TempDir\shared" -Recurse -Force
        }
        
        # Copy root go.mod and go.sum
        if (Test-Path "go.mod") {
            Copy-Item -Path "go.mod" -Destination $TempDir -Force
        }
        if (Test-Path "go.sum") {
            Copy-Item -Path "go.sum" -Destination $TempDir -Force
        }
        
        # Navigate to temp directory
        Push-Location $TempDir
        
        # Initialize git repo
        Write-Host "  🔧 Initializing Git repository..." -ForegroundColor White
        git init -q
        git add .
        git commit -m "Initial commit for $ServiceName deployment" -q
        
        # Deploy to Railway
        Write-Host "  🚀 Deploying to Railway..." -ForegroundColor White
        railway link --project retail-os-platform --service $RailwayServiceName
        railway up --detach
        
        Write-Host "  ✅ $ServiceName deployed successfully!" -ForegroundColor Green
    }
    catch {
        Write-Host "  ❌ Failed to deploy $ServiceName" -ForegroundColor Red
        Write-Host "  Error: $_" -ForegroundColor Red
    }
    finally {
        # Clean up
        Pop-Location
        Remove-Item -Path $TempDir -Recurse -Force -ErrorAction SilentlyContinue
    }
    
    Write-Host ""
}

# Deploy each service
Write-Host "🎯 Starting Git-based deployment..." -ForegroundColor Yellow
Write-Host ""

Deploy-ServiceViaGit -ServiceName "Product Service" -ServiceDir "services\product-catalog" -RailwayServiceName "retail-os-product" -Port "8003"
Deploy-ServiceViaGit -ServiceName "Payment Service" -ServiceDir "services\payment" -RailwayServiceName "retail-os-payment" -Port "8005"
Deploy-ServiceViaGit -ServiceName "Analytics Service" -ServiceDir "services\analytics" -RailwayServiceName "retail-os-analytics" -Port "8001"
Deploy-ServiceViaGit -ServiceName "Inventory Service" -ServiceDir "services\inventory" -RailwayServiceName "retail-os-inventory" -Port "8002"
Deploy-ServiceViaGit -ServiceName "Identity Service" -ServiceDir "services\identity" -RailwayServiceName "retail-os-identity" -Port "8000"
Deploy-ServiceViaGit -ServiceName "Cart Service" -ServiceDir "services\cart" -RailwayServiceName "retail-os-cart" -Port "8080"

Write-Host "🎉 GIT-BASED DEPLOYMENT COMPLETE!" -ForegroundColor Green -BackgroundColor Black
Write-Host ""