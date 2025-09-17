# Retail OS Service Deployment Commands
# Execute these commands manually in Railway CLI

Write-Host "🚀 Retail OS Railway CLI Deployment Commands" -ForegroundColor Green
Write-Host "Execute these commands one by one:" -ForegroundColor Yellow

$commands = @"

# 1. Create and Deploy Identity Service
cd services/identity
railway service create retail-os-identity
railway variables set SERVICE_NAME=identity
railway variables set ENVIRONMENT=production
railway variables set JWT_SECRET=prod-secure-jwt-secret-identity-2024
railway variables set JWT_EXPIRATION=86400
railway variables set LOG_LEVEL=info
railway up --detach
cd ../..

# 2. Create and Deploy Merchant Service  
cd services/merchant
railway service create retail-os-merchant
railway variables set SERVICE_NAME=merchant
railway variables set ENVIRONMENT=production
railway variables set JWT_SECRET=prod-secure-jwt-secret-merchant-2024
railway variables set LOG_LEVEL=info
railway up --detach
cd ../..

# 3. Create and Deploy Product Service
cd services/product
railway service create retail-os-product
railway variables set SERVICE_NAME=product
railway variables set ENVIRONMENT=production
railway variables set JWT_SECRET=prod-secure-jwt-secret-product-2024
railway variables set LOG_LEVEL=info
railway up --detach
cd ../..

# 4. Create and Deploy Inventory Service
cd services/inventory
railway service create retail-os-inventory
railway variables set SERVICE_NAME=inventory
railway variables set ENVIRONMENT=production
railway variables set JWT_SECRET=prod-secure-jwt-secret-inventory-2024
railway variables set LOG_LEVEL=info
railway up --detach
cd ../..

# 5. Create and Deploy Order Service
cd services/order
railway service create retail-os-order
railway variables set SERVICE_NAME=order
railway variables set ENVIRONMENT=production
railway variables set JWT_SECRET=prod-secure-jwt-secret-order-2024
railway variables set LOG_LEVEL=info
railway up --detach
cd ../..

# 6. Create and Deploy Payment Service (Port 8005)
cd services/payment
railway service create retail-os-payment
railway variables set SERVICE_NAME=payment
railway variables set SERVICE_PORT=8005
railway variables set ENVIRONMENT=production
railway variables set JWT_SECRET=prod-secure-jwt-secret-payment-2024
railway variables set LOG_LEVEL=info
railway up --detach
cd ../..

# 7. Create and Deploy Cart Service
cd services/cart
railway service create retail-os-cart
railway variables set SERVICE_NAME=cart
railway variables set ENVIRONMENT=production
railway variables set JWT_SECRET=prod-secure-jwt-secret-cart-2024
railway variables set LOG_LEVEL=info
railway up --detach
cd ../..

# 8. Create and Deploy Promotions Service
cd services/promotions
railway service create retail-os-promotions
railway variables set SERVICE_NAME=promotions
railway variables set ENVIRONMENT=production
railway variables set JWT_SECRET=prod-secure-jwt-secret-promotions-2024
railway variables set LOG_LEVEL=info
railway up --detach
cd ../..

# 9. Create and Deploy Analytics Service
cd services/analytics
railway service create retail-os-analytics
railway variables set SERVICE_NAME=analytics
railway variables set ENVIRONMENT=production
railway variables set JWT_SECRET=prod-secure-jwt-secret-analytics-2024
railway variables set LOG_LEVEL=info
railway up --detach
cd ../..

# 10. Create and Deploy GraphQL Gateway
cd gateway
railway service create retail-os-gateway
railway variables set NODE_ENV=production
railway up --detach
cd ..

# 11. Create and Deploy Storefront
cd apps/storefront
railway service create retail-os-storefront
railway variables set NODE_ENV=production
railway variables set NEXT_PUBLIC_API_URL=https://retail-os-gateway.railway.app
railway up --detach
cd ../..

# 12. Create and Deploy Admin Panel
cd apps/admin
railway service create retail-os-admin
railway variables set NODE_ENV=production
railway variables set REACT_APP_API_URL=https://retail-os-gateway.railway.app
railway up --detach
cd ../..

"@

Write-Host $commands -ForegroundColor Cyan

Write-Host "`n📋 After deployment, update Gateway with service URLs:" -ForegroundColor Yellow
Write-Host "railway variables set IDENTITY_SERVICE_URL=https://retail-os-identity.railway.app" -ForegroundColor Green
Write-Host "railway variables set MERCHANT_SERVICE_URL=https://retail-os-merchant.railway.app" -ForegroundColor Green
Write-Host "railway variables set PRODUCT_SERVICE_URL=https://retail-os-product.railway.app" -ForegroundColor Green
Write-Host "railway variables set INVENTORY_SERVICE_URL=https://retail-os-inventory.railway.app" -ForegroundColor Green
Write-Host "railway variables set ORDER_SERVICE_URL=https://retail-os-order.railway.app" -ForegroundColor Green
Write-Host "railway variables set PAYMENT_SERVICE_URL=https://retail-os-payment.railway.app" -ForegroundColor Green
Write-Host "railway variables set CART_SERVICE_URL=https://retail-os-cart.railway.app" -ForegroundColor Green
Write-Host "railway variables set PROMOTIONS_SERVICE_URL=https://retail-os-promotions.railway.app" -ForegroundColor Green
Write-Host "railway variables set ANALYTICS_SERVICE_URL=https://retail-os-analytics.railway.app" -ForegroundColor Green