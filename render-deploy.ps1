# Retail OS - Automated Render Deployment Script
# This script automates the entire Render deployment process

Write-Host "🚀 Retail OS - Automated Render Deployment" -ForegroundColor Cyan
Write-Host "================================================" -ForegroundColor Cyan

# Check if Render CLI is installed
$renderInstalled = Get-Command render -ErrorAction SilentlyContinue
if (-not $renderInstalled) {
    Write-Host "📦 Installing Render CLI..." -ForegroundColor Yellow
    
    # Install Render CLI
    if ($IsWindows -or $env:OS -eq "Windows_NT") {
        # Windows installation
        Invoke-WebRequest -Uri "https://github.com/render-oss/cli/releases/latest/download/render-windows-amd64.exe" -OutFile "$env:TEMP\render.exe"
        $renderPath = "$env:LOCALAPPDATA\Render"
        New-Item -ItemType Directory -Force -Path $renderPath | Out-Null
        Move-Item "$env:TEMP\render.exe" "$renderPath\render.exe" -Force
        $env:PATH += ";$renderPath"
        [Environment]::SetEnvironmentVariable("PATH", $env:PATH + ";$renderPath", [EnvironmentVariableTarget]::User)
    }
    
    Write-Host "✅ Render CLI installed successfully!" -ForegroundColor Green
}

# Login to Render
Write-Host "🔐 Authenticating with Render..." -ForegroundColor Yellow
Write-Host "Please visit the URL that opens and paste your API key when prompted." -ForegroundColor Cyan
render auth login

# Get user input for database URLs
Write-Host "📋 Please provide your database connection details:" -ForegroundColor Cyan
$postgresUrl = Read-Host "PostgreSQL URL (from Render dashboard)"
$redisUrl = Read-Host "Redis URL (from Render dashboard)"

# Validate URLs
if (-not $postgresUrl -or -not $redisUrl) {
    Write-Host "❌ Database URLs are required. Please run the script again with valid URLs." -ForegroundColor Red
    exit 1
}

# Service configurations
$services = @(
    @{
        name = "retail-os-identity"
        path = "services/identity"
        port = 8001
        description = "Identity & Authentication Service"
    },
    @{
        name = "retail-os-product-catalog"
        path = "services/product-catalog"
        port = 8008
        description = "Product Catalog Service"
    },
    @{
        name = "retail-os-inventory"
        path = "services/inventory"
        port = 8004
        description = "Inventory Management Service"
    },
    @{
        name = "retail-os-cart"
        path = "services/cart"
        port = 8002
        description = "Shopping Cart Service"
    },
    @{
        name = "retail-os-order"
        path = "services/order"
        port = 8005
        description = "Order Processing Service"
    },
    @{
        name = "retail-os-payment"
        path = "services/payment"
        port = 8006
        description = "Payment Processing Service"
    },
    @{
        name = "retail-os-promotions"
        path = "services/promotions"
        port = 8007
        description = "Promotions & Discounts Service"
    },
    @{
        name = "retail-os-merchant"
        path = "services/merchant-account"
        port = 8003
        description = "Merchant Account Service"
    },
    @{
        name = "retail-os-analytics"
        path = "services/analytics"
        port = 8009
        description = "Analytics & Reporting Service"
    },
    @{
        name = "retail-os-gateway"
        path = "gateway"
        port = 4000
        description = "GraphQL Federation Gateway"
    }
)

Write-Host "🚀 Starting automated deployment of $($services.Count) services..." -ForegroundColor Green

$deployedServices = @()

foreach ($service in $services) {
    Write-Host ""
    Write-Host "📦 Deploying: $($service.description)" -ForegroundColor Yellow
    Write-Host "   Service: $($service.name)" -ForegroundColor Cyan
    Write-Host "   Path: $($service.path)" -ForegroundColor Cyan
    Write-Host "   Port: $($service.port)" -ForegroundColor Cyan
    
    # Create render.yaml for the service
    $renderYaml = @"
services:
  - type: web
    name: $($service.name)
    runtime: go
    repo: https://github.com/mobilebarn/UNIFIED_COMMERCE
    rootDir: $($service.path)
    buildCommand: go build -o app ./cmd/server
    startCommand: ./app
    plan: free
    region: singapore
    envVars:
      - key: PORT
        value: $($service.port)
      - key: ENVIRONMENT
        value: production
      - key: LOG_LEVEL
        value: info
      - key: DATABASE_URL
        value: $postgresUrl
      - key: REDIS_URL
        value: $redisUrl
      - key: SERVICE_NAME
        value: $($service.name)
      - key: SERVICE_PORT
        value: $($service.port)
"@
    
    # Write render.yaml to service directory
    $yamlPath = "$($service.path)\render.yaml"
    $renderYaml | Out-File -FilePath $yamlPath -Encoding UTF8
    
    try {
        # Deploy the service
        Write-Host "   ⏳ Deploying..." -ForegroundColor Yellow
        $result = render deploy --yaml $yamlPath 2>&1
        
        if ($LASTEXITCODE -eq 0) {
            Write-Host "   ✅ $($service.name) deployed successfully!" -ForegroundColor Green
            $deployedServices += $service
        } else {
            Write-Host "   ❌ Failed to deploy $($service.name): $result" -ForegroundColor Red
        }
    } catch {
        Write-Host "   ❌ Error deploying $($service.name): $_" -ForegroundColor Red
    }
    
    # Small delay between deployments
    Start-Sleep -Seconds 5
}

Write-Host ""
Write-Host "🎯 Deployment Summary:" -ForegroundColor Cyan
Write-Host "=====================" -ForegroundColor Cyan
Write-Host "✅ Successfully deployed: $($deployedServices.Count)/$($services.Count) services" -ForegroundColor Green

if ($deployedServices.Count -gt 0) {
    Write-Host ""
    Write-Host "🌐 Deployed Services:" -ForegroundColor Green
    foreach ($service in $deployedServices) {
        Write-Host "   • $($service.description) - https://$($service.name).onrender.com" -ForegroundColor Cyan
    }
    
    Write-Host ""
    Write-Host "🔗 Next Steps:" -ForegroundColor Yellow
    Write-Host "1. Check your Render dashboard: https://dashboard.render.com" -ForegroundColor Cyan
    Write-Host "2. Monitor deployment logs for any issues" -ForegroundColor Cyan
    Write-Host "3. Update your Vercel frontend to use the new URLs" -ForegroundColor Cyan
    Write-Host "4. Test the GraphQL Gateway: https://retail-os-gateway.onrender.com/graphql" -ForegroundColor Cyan
}

Write-Host ""
Write-Host "🚀 Retail OS deployment automation complete!" -ForegroundColor Green