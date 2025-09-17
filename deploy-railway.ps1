# Retail OS Railway Deployment Script for Windows PowerShell
# This script deploys the entire Retail OS platform to Railway

Write-Host "🚀 Starting Retail OS Platform Deployment to Railway" -ForegroundColor Green
Write-Host "==================================================" -ForegroundColor Green

# Check if Railway CLI is authenticated
try {
    railway whoami | Out-Null
    Write-Host "✅ Railway CLI authenticated" -ForegroundColor Green
} catch {
    Write-Host "❌ You need to login to Railway first" -ForegroundColor Red
    Write-Host "Please run: railway login" -ForegroundColor Yellow
    exit 1
}

# Function to deploy a service
function Deploy-Service {
    param(
        [string]$ServiceName,
        [string]$ServicePath,
        [string]$ProjectName
    )
    
    Write-Host "📤 Deploying $ServiceName service..." -ForegroundColor Blue
    Set-Location $ServicePath
    
    # Create a new Railway service
    railway init --name $ProjectName
    
    # Deploy the service
    railway up --detach
    
    Set-Location $PSScriptRoot
    Write-Host "✅ $ServiceName service deployed" -ForegroundColor Green
}

# Create main project and set up databases
Write-Host "📦 Creating Railway project..." -ForegroundColor Blue
railway init --name "retail-os-platform"

Write-Host "🗄️  Setting up databases..." -ForegroundColor Blue
railway add postgresql
railway add redis

# Wait for databases to be ready
Write-Host "⏳ Waiting for databases to initialize..." -ForegroundColor Yellow
Start-Sleep -Seconds 30

# Deploy backend services
Write-Host "🔧 Deploying backend microservices..." -ForegroundColor Blue

$services = @(
    @{Name="Identity"; Path="services/identity"; Project="retail-os-identity"},
    @{Name="Merchant"; Path="services/merchant"; Project="retail-os-merchant"},
    @{Name="Product"; Path="services/product"; Project="retail-os-product"},
    @{Name="Inventory"; Path="services/inventory"; Project="retail-os-inventory"},
    @{Name="Order"; Path="services/order"; Project="retail-os-order"},
    @{Name="Payment"; Path="services/payment"; Project="retail-os-payment"},
    @{Name="Cart"; Path="services/cart"; Project="retail-os-cart"},
    @{Name="Promotions"; Path="services/promotions"; Project="retail-os-promotions"},
    @{Name="Analytics"; Path="services/analytics"; Project="retail-os-analytics"}
)

foreach ($service in $services) {
    Deploy-Service -ServiceName $service.Name -ServicePath $service.Path -ProjectName $service.Project
}

# Deploy GraphQL Gateway
Write-Host "🌐 Deploying GraphQL Federation Gateway..." -ForegroundColor Blue
Deploy-Service -ServiceName "Gateway" -ServicePath "gateway" -ProjectName "retail-os-gateway"

# Deploy frontend applications
Write-Host "🖥️  Deploying frontend applications..." -ForegroundColor Blue
Deploy-Service -ServiceName "Storefront" -ServicePath "apps/storefront" -ProjectName "retail-os-storefront"
Deploy-Service -ServiceName "Admin Panel" -ServicePath "apps/admin" -ProjectName "retail-os-admin"

Write-Host "🎉 Retail OS Platform deployment to Railway completed!" -ForegroundColor Green
Write-Host "==================================================" -ForegroundColor Green
Write-Host "Your services are being deployed. You can check their status with:" -ForegroundColor Yellow
Write-Host "railway status" -ForegroundColor Cyan
Write-Host ""
Write-Host "🌍 Your applications will be available at:" -ForegroundColor Yellow
Write-Host "- Storefront: https://retail-os-storefront.railway.app" -ForegroundColor Cyan
Write-Host "- Admin Panel: https://retail-os-admin.railway.app" -ForegroundColor Cyan
Write-Host "- GraphQL Gateway: https://retail-os-gateway.railway.app" -ForegroundColor Cyan
Write-Host ""
Write-Host "📊 Monitor your deployments in the Railway dashboard:" -ForegroundColor Yellow
Write-Host "https://railway.app/dashboard" -ForegroundColor Cyan