Write-Host "DEPLOYING RETAIL OS PLATFORM NOW" -ForegroundColor Red -BackgroundColor Yellow

# Open Railway Dashboard
Start-Process "https://railway.com/project/1f76e4aa-b36a-4670-a445-1bfde6bba0a3"

Write-Host "Railway dashboard opened. Create these services manually:" -ForegroundColor Yellow
Write-Host "1. retail-os-identity" -ForegroundColor Cyan
Write-Host "2. retail-os-merchant" -ForegroundColor Cyan
Write-Host "3. retail-os-product" -ForegroundColor Cyan
Write-Host "4. retail-os-inventory" -ForegroundColor Cyan
Write-Host "5. retail-os-order" -ForegroundColor Cyan
Write-Host "6. retail-os-payment" -ForegroundColor Cyan
Write-Host "7. retail-os-cart" -ForegroundColor Cyan
Write-Host "8. retail-os-promotions" -ForegroundColor Cyan
Write-Host "9. retail-os-analytics" -ForegroundColor Cyan
Write-Host "10. retail-os-gateway" -ForegroundColor Cyan
Write-Host "11. retail-os-storefront" -ForegroundColor Cyan
Write-Host "12. retail-os-admin" -ForegroundColor Cyan

Write-Host "`nThen run: ./railway-auto-deploy.ps1" -ForegroundColor Green