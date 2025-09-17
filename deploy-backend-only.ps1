# RETAIL OS BACKEND DEPLOYMENT TO RAILWAY
# Frontend apps already deployed to Vercel - only deploying backend services

Write-Host "🚀 DEPLOYING RETAIL OS BACKEND SERVICES TO RAILWAY" -ForegroundColor Green -BackgroundColor Black
Write-Host "====================================================" -ForegroundColor Green

Write-Host "`n📋 Services to deploy:" -ForegroundColor Yellow
Write-Host "✅ 9 Backend microservices + 1 GraphQL Gateway" -ForegroundColor White
Write-Host "❌ Frontend apps (already on Vercel)" -ForegroundColor Gray

# 1. IDENTITY SERVICE
Write-Host "`n1️⃣ Deploying Identity Service..." -ForegroundColor Cyan
cd services/identity
railway service retail-os-identity
railway variables --set SERVICE_NAME=identity
railway variables --set ENVIRONMENT=production
railway variables --set JWT_SECRET=prod-identity-jwt-2024
railway up --detach
cd ../..

# 2. MERCHANT SERVICE  
Write-Host "`n2️⃣ Deploying Merchant Service..." -ForegroundColor Cyan
cd services/merchant
railway service retail-os-merchant
railway variables --set SERVICE_NAME=merchant
railway variables --set ENVIRONMENT=production
railway variables --set JWT_SECRET=prod-merchant-jwt-2024
railway up --detach
cd ../..

# 3. PRODUCT SERVICE
Write-Host "`n3️⃣ Deploying Product Service..." -ForegroundColor Cyan
cd services/product
railway service retail-os-product
railway variables --set SERVICE_NAME=product
railway variables --set ENVIRONMENT=production
railway variables --set JWT_SECRET=prod-product-jwt-2024
railway up --detach
cd ../..

# 4. INVENTORY SERVICE
Write-Host "`n4️⃣ Deploying Inventory Service..." -ForegroundColor Cyan
cd services/inventory
railway service retail-os-inventory
railway variables --set SERVICE_NAME=inventory
railway variables --set ENVIRONMENT=production
railway variables --set JWT_SECRET=prod-inventory-jwt-2024
railway up --detach
cd ../..

# 5. ORDER SERVICE
Write-Host "`n5️⃣ Deploying Order Service..." -ForegroundColor Cyan
cd services/order
railway service retail-os-order
railway variables --set SERVICE_NAME=order
railway variables --set ENVIRONMENT=production
railway variables --set JWT_SECRET=prod-order-jwt-2024
railway up --detach
cd ../..

# 6. PAYMENT SERVICE (PORT 8005)
Write-Host "`n6️⃣ Deploying Payment Service (Port 8005)..." -ForegroundColor Cyan
cd services/payment
railway service retail-os-payment
railway variables --set SERVICE_NAME=payment
railway variables --set SERVICE_PORT=8005
railway variables --set ENVIRONMENT=production
railway variables --set JWT_SECRET=prod-payment-jwt-2024
railway up --detach
cd ../..

# 7. CART SERVICE
Write-Host "`n7️⃣ Deploying Cart Service..." -ForegroundColor Cyan
cd services/cart
railway service retail-os-cart
railway variables --set SERVICE_NAME=cart
railway variables --set ENVIRONMENT=production
railway variables --set JWT_SECRET=prod-cart-jwt-2024
railway up --detach
cd ../..

# 8. PROMOTIONS SERVICE
Write-Host "`n8️⃣ Deploying Promotions Service..." -ForegroundColor Cyan
cd services/promotions
railway service retail-os-promotions
railway variables --set SERVICE_NAME=promotions
railway variables --set ENVIRONMENT=production
railway variables --set JWT_SECRET=prod-promotions-jwt-2024
railway up --detach
cd ../..

# 9. ANALYTICS SERVICE
Write-Host "`n9️⃣ Deploying Analytics Service..." -ForegroundColor Cyan
cd services/analytics
railway service retail-os-analytics
railway variables --set SERVICE_NAME=analytics
railway variables --set ENVIRONMENT=production
railway variables --set JWT_SECRET=prod-analytics-jwt-2024
railway up --detach
cd ../..

# 10. GRAPHQL GATEWAY
Write-Host "`n🔟 Deploying GraphQL Gateway..." -ForegroundColor Cyan
cd gateway
railway service retail-os-gateway
railway variables --set NODE_ENV=production
railway variables --set IDENTITY_SERVICE_URL=https://retail-os-identity.railway.app
railway variables --set MERCHANT_SERVICE_URL=https://retail-os-merchant.railway.app
railway variables --set PRODUCT_SERVICE_URL=https://retail-os-product.railway.app
railway variables --set INVENTORY_SERVICE_URL=https://retail-os-inventory.railway.app
railway variables --set ORDER_SERVICE_URL=https://retail-os-order.railway.app
railway variables --set PAYMENT_SERVICE_URL=https://retail-os-payment.railway.app
railway variables --set CART_SERVICE_URL=https://retail-os-cart.railway.app
railway variables --set PROMOTIONS_SERVICE_URL=https://retail-os-promotions.railway.app
railway variables --set ANALYTICS_SERVICE_URL=https://retail-os-analytics.railway.app
railway up --detach
cd ..

Write-Host "`n🎉 RETAIL OS BACKEND DEPLOYMENT COMPLETE!" -ForegroundColor Green -BackgroundColor Black
Write-Host "=============================================" -ForegroundColor Green

Write-Host "`n📊 Deployment Summary:" -ForegroundColor Yellow
Write-Host "✅ Backend Services: Deployed to Railway" -ForegroundColor Green
Write-Host "✅ GraphQL Gateway: https://retail-os-gateway.railway.app" -ForegroundColor Green
Write-Host "✅ Frontend Apps: Already on Vercel from last night" -ForegroundColor Green

Write-Host "`n🔗 Next Steps:" -ForegroundColor Yellow
Write-Host "1. Update Vercel frontend apps to use new Railway Gateway URL" -ForegroundColor White
Write-Host "2. Test end-to-end functionality" -ForegroundColor White
Write-Host "3. Your Retail OS platform is ready!" -ForegroundColor White