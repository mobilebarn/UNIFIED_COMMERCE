# Retail OS Railway URL Update Script
# Run this after all services are deployed to update service discovery URLs

Write-Host "🔄 Updating Retail OS service URLs in Railway..." -ForegroundColor Green

# Service URLs (Update these with your actual Railway URLs)
$services = @{
    "IDENTITY_SERVICE_URL" = "https://retail-os-identity.railway.app"
    "MERCHANT_SERVICE_URL" = "https://retail-os-merchant.railway.app"
    "PRODUCT_SERVICE_URL" = "https://retail-os-product.railway.app"
    "INVENTORY_SERVICE_URL" = "https://retail-os-inventory.railway.app"
    "ORDER_SERVICE_URL" = "https://retail-os-order.railway.app"
    "PAYMENT_SERVICE_URL" = "https://retail-os-payment.railway.app"
    "CART_SERVICE_URL" = "https://retail-os-cart.railway.app"
    "PROMOTIONS_SERVICE_URL" = "https://retail-os-promotions.railway.app"
    "ANALYTICS_SERVICE_URL" = "https://retail-os-analytics.railway.app"
}

Write-Host "📋 Copy these URLs to your Gateway service environment variables:" -ForegroundColor Yellow

foreach ($service in $services.GetEnumerator()) {
    Write-Host "$($service.Key)=$($service.Value)" -ForegroundColor Cyan
}

Write-Host "`n🌐 Frontend URLs:" -ForegroundColor Yellow
Write-Host "NEXT_PUBLIC_API_URL=https://retail-os-gateway.railway.app" -ForegroundColor Cyan
Write-Host "REACT_APP_API_URL=https://retail-os-gateway.railway.app" -ForegroundColor Cyan

Write-Host "`n✅ Update complete! Don't forget to:" -ForegroundColor Green
Write-Host "1. Copy these URLs to Railway dashboard environment variables" -ForegroundColor White
Write-Host "2. Redeploy Gateway and Frontend services after updating URLs" -ForegroundColor White
Write-Host "3. Test all services at their new URLs" -ForegroundColor White