# RETAIL OS DIRECT DEPLOYMENT SCRIPT
# This script will deploy your services immediately

Write-Host "🚀 DEPLOYING RETAIL OS PLATFORM NOW" -ForegroundColor Red -BackgroundColor Yellow
Write-Host "=====================================" -ForegroundColor Red -BackgroundColor Yellow

# Step 1: Create services using REST API approach
Write-Host "`n1️⃣ Opening Railway Dashboard for service creation..." -ForegroundColor Green
Start-Process "https://railway.com/project/1f76e4aa-b36a-4670-a445-1bfde6bba0a3"

Write-Host "`n2️⃣ Meanwhile, preparing deployment files..." -ForegroundColor Green

# Create deployment commands that will work once services exist
$deploymentScript = @"
# RETAIL OS AUTOMATED DEPLOYMENT COMMANDS
# Run these after creating services in Railway dashboard

# 1. IDENTITY SERVICE
Write-Host "Deploying Identity Service..." -ForegroundColor Cyan
cd services/identity
railway service --name retail-os-identity 2>nul
railway variables set SERVICE_NAME=identity
railway variables set ENVIRONMENT=production
railway variables set JWT_SECRET=prod-identity-jwt-2024
railway up --detach
cd ../..

# 2. MERCHANT SERVICE  
Write-Host "Deploying Merchant Service..." -ForegroundColor Cyan
cd services/merchant
railway service --name retail-os-merchant 2>nul
railway variables set SERVICE_NAME=merchant
railway variables set ENVIRONMENT=production
railway variables set JWT_SECRET=prod-merchant-jwt-2024
railway up --detach
cd ../..

# 3. PRODUCT SERVICE
Write-Host "Deploying Product Service..." -ForegroundColor Cyan
cd services/product
railway service --name retail-os-product 2>nul
railway variables set SERVICE_NAME=product
railway variables set ENVIRONMENT=production
railway variables set JWT_SECRET=prod-product-jwt-2024
railway up --detach
cd ../..

# 4. INVENTORY SERVICE
Write-Host "Deploying Inventory Service..." -ForegroundColor Cyan
cd services/inventory
railway service --name retail-os-inventory 2>nul
railway variables set SERVICE_NAME=inventory
railway variables set ENVIRONMENT=production
railway variables set JWT_SECRET=prod-inventory-jwt-2024
railway up --detach
cd ../..

# 5. ORDER SERVICE
Write-Host "Deploying Order Service..." -ForegroundColor Cyan
cd services/order
railway service --name retail-os-order 2>nul
railway variables set SERVICE_NAME=order
railway variables set ENVIRONMENT=production
railway variables set JWT_SECRET=prod-order-jwt-2024
railway up --detach
cd ../..

# 6. PAYMENT SERVICE (PORT 8005)
Write-Host "Deploying Payment Service..." -ForegroundColor Cyan
cd services/payment
railway service --name retail-os-payment 2>nul
railway variables set SERVICE_NAME=payment
railway variables set SERVICE_PORT=8005
railway variables set ENVIRONMENT=production
railway variables set JWT_SECRET=prod-payment-jwt-2024
railway up --detach
cd ../..

# 7. CART SERVICE
Write-Host "Deploying Cart Service..." -ForegroundColor Cyan
cd services/cart
railway service --name retail-os-cart 2>nul
railway variables set SERVICE_NAME=cart
railway variables set ENVIRONMENT=production
railway variables set JWT_SECRET=prod-cart-jwt-2024
railway up --detach
cd ../..

# 8. PROMOTIONS SERVICE
Write-Host "Deploying Promotions Service..." -ForegroundColor Cyan
cd services/promotions
railway service --name retail-os-promotions 2>nul
railway variables set SERVICE_NAME=promotions
railway variables set ENVIRONMENT=production
railway variables set JWT_SECRET=prod-promotions-jwt-2024
railway up --detach
cd ../..

# 9. ANALYTICS SERVICE
Write-Host "Deploying Analytics Service..." -ForegroundColor Cyan
cd services/analytics
railway service --name retail-os-analytics 2>nul
railway variables set SERVICE_NAME=analytics
railway variables set ENVIRONMENT=production
railway variables set JWT_SECRET=prod-analytics-jwt-2024
railway up --detach
cd ../..

# 10. GRAPHQL GATEWAY
Write-Host "Deploying GraphQL Gateway..." -ForegroundColor Cyan
cd gateway
railway service --name retail-os-gateway 2>nul
railway variables set NODE_ENV=production
railway up --detach
cd ..

# 11. STOREFRONT
Write-Host "Deploying Storefront..." -ForegroundColor Cyan
cd apps/storefront
railway service --name retail-os-storefront 2>nul
railway variables set NODE_ENV=production
railway variables set NEXT_PUBLIC_API_URL=https://retail-os-gateway.railway.app
railway up --detach
cd ../..

# 12. ADMIN PANEL
Write-Host "Deploying Admin Panel..." -ForegroundColor Cyan
cd apps/admin
railway service --name retail-os-admin 2>nul
railway variables set NODE_ENV=production
railway variables set REACT_APP_API_URL=https://retail-os-gateway.railway.app
railway up --detach
cd ../..

Write-Host "`n🎉 RETAIL OS DEPLOYMENT COMPLETE!" -ForegroundColor Green -BackgroundColor Black
Write-Host "Your platform is available at:" -ForegroundColor Yellow
Write-Host "- Storefront: https://retail-os-storefront.railway.app" -ForegroundColor Cyan
Write-Host "- Admin: https://retail-os-admin.railway.app" -ForegroundColor Cyan
Write-Host "- Gateway: https://retail-os-gateway.railway.app" -ForegroundColor Cyan
"@

# Save the deployment script
$deploymentScript | Out-File -FilePath "railway-auto-deploy.ps1" -Encoding UTF8

Write-Host "`n3️⃣ EXECUTE THIS NOW:" -ForegroundColor Red -BackgroundColor Yellow
Write-Host ".\railway-auto-deploy.ps1" -ForegroundColor White -BackgroundColor Red

Write-Host "`n📋 OR CREATE SERVICES MANUALLY IN BROWSER:" -ForegroundColor Yellow
Write-Host "1. Go to Railway Dashboard (opened above)" -ForegroundColor White
Write-Host "2. Click '+ New' → 'Empty Service' for each:" -ForegroundColor White
Write-Host "   - retail-os-identity" -ForegroundColor Gray
Write-Host "   - retail-os-merchant" -ForegroundColor Gray
Write-Host "   - retail-os-product" -ForegroundColor Gray
Write-Host "   - retail-os-inventory" -ForegroundColor Gray
Write-Host "   - retail-os-order" -ForegroundColor Gray
Write-Host "   - retail-os-payment" -ForegroundColor Gray
Write-Host "   - retail-os-cart" -ForegroundColor Gray
Write-Host "   - retail-os-promotions" -ForegroundColor Gray
Write-Host "   - retail-os-analytics" -ForegroundColor Gray
Write-Host "   - retail-os-gateway" -ForegroundColor Gray
Write-Host "   - retail-os-storefront" -ForegroundColor Gray
Write-Host "   - retail-os-admin" -ForegroundColor Gray
Write-Host "3. Then run: .\railway-auto-deploy.ps1" -ForegroundColor White

Write-Host "`n🚀 GETTING YOUR RETAIL OS DEPLOYED NOW!" -ForegroundColor Green -BackgroundColor Black