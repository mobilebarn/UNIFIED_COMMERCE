# Retail OS Railway Step-by-Step Deployment
# Run these commands one by one in your terminal

Write-Host "🚀 Retail OS Railway CLI Deployment" -ForegroundColor Green
Write-Host "Follow these commands step by step:" -ForegroundColor Yellow

Write-Host "`n1️⃣ First, create the Identity Service:" -ForegroundColor Cyan
Write-Host "railway add --service" -ForegroundColor White
Write-Host "   → Enter service name: retail-os-identity" -ForegroundColor Gray
Write-Host "   → Press Enter to skip variables" -ForegroundColor Gray

Write-Host "`n2️⃣ Deploy Identity Service:" -ForegroundColor Cyan
Write-Host "cd services/identity" -ForegroundColor White
Write-Host "railway service retail-os-identity" -ForegroundColor White
Write-Host "railway variables set SERVICE_NAME=identity" -ForegroundColor White
Write-Host "railway variables set ENVIRONMENT=production" -ForegroundColor White
Write-Host "railway variables set JWT_SECRET=prod-secure-jwt-secret-identity-2024" -ForegroundColor White
Write-Host "railway up --detach" -ForegroundColor White
Write-Host "cd ../.." -ForegroundColor White

Write-Host "`n3️⃣ Create Merchant Service:" -ForegroundColor Cyan
Write-Host "railway add --service" -ForegroundColor White
Write-Host "   → Enter service name: retail-os-merchant" -ForegroundColor Gray
Write-Host "   → Press Enter to skip variables" -ForegroundColor Gray

Write-Host "`n4️⃣ Deploy Merchant Service:" -ForegroundColor Cyan
Write-Host "cd services/merchant" -ForegroundColor White
Write-Host "railway service retail-os-merchant" -ForegroundColor White
Write-Host "railway variables set SERVICE_NAME=merchant" -ForegroundColor White
Write-Host "railway variables set ENVIRONMENT=production" -ForegroundColor White
Write-Host "railway variables set JWT_SECRET=prod-secure-jwt-secret-merchant-2024" -ForegroundColor White
Write-Host "railway up --detach" -ForegroundColor White
Write-Host "cd ../.." -ForegroundColor White

Write-Host "`n📋 Continue this pattern for all 12 services:" -ForegroundColor Yellow
Write-Host "   - product (retail-os-product)" -ForegroundColor Gray
Write-Host "   - inventory (retail-os-inventory)" -ForegroundColor Gray
Write-Host "   - order (retail-os-order)" -ForegroundColor Gray
Write-Host "   - payment (retail-os-payment) - USE PORT 8005!" -ForegroundColor Red
Write-Host "   - cart (retail-os-cart)" -ForegroundColor Gray
Write-Host "   - promotions (retail-os-promotions)" -ForegroundColor Gray
Write-Host "   - analytics (retail-os-analytics)" -ForegroundColor Gray
Write-Host "   - gateway (retail-os-gateway)" -ForegroundColor Gray
Write-Host "   - storefront (retail-os-storefront)" -ForegroundColor Gray
Write-Host "   - admin (retail-os-admin)" -ForegroundColor Gray

Write-Host "`n🎯 Current Status: Ready to create retail-os-identity service!" -ForegroundColor Green
Write-Host "Execute the first command above ⬆️" -ForegroundColor Yellow