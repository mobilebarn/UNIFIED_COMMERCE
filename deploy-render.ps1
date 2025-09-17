# RETAIL OS RENDER DEPLOYMENT - SIMPLE & FAST
# Deploy all services to Render.com in 15 minutes

Write-Host "RETAIL OS RENDER DEPLOYMENT" -ForegroundColor Green
Write-Host "==========================" -ForegroundColor Green
Write-Host ""

Write-Host "Why Render?" -ForegroundColor Yellow
Write-Host "- Built for monorepos (perfect for our setup)" -ForegroundColor Green
Write-Host "- No file size limitations" -ForegroundColor Green
Write-Host "- Automatic Docker builds" -ForegroundColor Green
Write-Host "- Free tier available" -ForegroundColor Green
Write-Host "- Much simpler than Railway" -ForegroundColor Green
Write-Host ""

Write-Host "Quick Setup Steps:" -ForegroundColor Yellow
Write-Host "1. Go to render.com and sign up" -ForegroundColor White
Write-Host "2. Connect your GitHub repository" -ForegroundColor White
Write-Host "3. Create services with these settings:" -ForegroundColor White
Write-Host ""

# Service configurations for Render
$services = @(
    @{Name="product-service"; Dir="services/product-catalog"; Port="8003"; DB="MongoDB"},
    @{Name="payment-service"; Dir="services/payment"; Port="8005"; DB="PostgreSQL"},
    @{Name="analytics-service"; Dir="services/analytics"; Port="8001"; DB="PostgreSQL"},
    @{Name="inventory-service"; Dir="services/inventory"; Port="8002"; DB="PostgreSQL"},
    @{Name="identity-service"; Dir="services/identity"; Port="8000"; DB="PostgreSQL"},
    @{Name="cart-service"; Dir="services/cart"; Port="8080"; DB="PostgreSQL"},
    @{Name="order-service"; Dir="services/order"; Port="8004"; DB="PostgreSQL"},
    @{Name="merchant-service"; Dir="services/merchant-account"; Port="8006"; DB="PostgreSQL"},
    @{Name="promotions-service"; Dir="services/promotions"; Port="8007"; DB="PostgreSQL"},
    @{Name="graphql-gateway"; Dir="gateway"; Port="4000"; DB="N/A"}
)

foreach ($service in $services) {
    Write-Host "Service: $($service.Name)" -ForegroundColor Cyan
    Write-Host "  Root Directory: $($service.Dir)" -ForegroundColor White
    Write-Host "  Port: $($service.Port)" -ForegroundColor White
    Write-Host "  Database: $($service.DB)" -ForegroundColor White
    Write-Host ""
}

Write-Host "Total Time: ~15 minutes for all services" -ForegroundColor Green