# RAILWAY GO VERSION COMPATIBILITY FIX
# Fix Go version compatibility for Railway deployment

Write-Host "FIXING GO VERSION COMPATIBILITY FOR RAILWAY" -ForegroundColor Red
Write-Host "===========================================" -ForegroundColor Red
Write-Host ""

Write-Host "Railway Build Environment: Go 1.21.13" -ForegroundColor Yellow
Write-Host "Our Project: Downgrading from Go 1.23.0 → Go 1.21" -ForegroundColor Yellow
Write-Host ""

# Fix all service go.mod files
$services = @(
    "services/analytics",
    "services/cart", 
    "services/identity",
    "services/inventory",
    "services/merchant-account",
    "services/order",
    "services/payment",
    "services/product-catalog",
    "services/promotions"
)

Write-Host "Updating Go version in service go.mod files..." -ForegroundColor Cyan

foreach ($service in $services) {
    if (Test-Path "$service/go.mod") {
        Write-Host "  Fixing $service/go.mod" -ForegroundColor White
        
        # Read content
        $content = Get-Content "$service/go.mod"
        
        # Replace go version
        $newContent = $content -replace "go 1\.23\.0", "go 1.21"
        $newContent = $newContent -replace "go 1\.23", "go 1.21"
        
        # Write back
        $newContent | Set-Content "$service/go.mod"
        
        Write-Host "    ✅ Updated to Go 1.21" -ForegroundColor Green
    } else {
        Write-Host "    ❌ go.mod not found in $service" -ForegroundColor Red
    }
}

Write-Host ""
Write-Host "Committing and pushing fixes..." -ForegroundColor Cyan
git add .
git commit -m "Fix: Downgrade Go version to 1.21 for Railway compatibility"
git push origin master

Write-Host ""
Write-Host "✅ FIXES APPLIED!" -ForegroundColor Green
Write-Host ""
Write-Host "Next steps:" -ForegroundColor Yellow
Write-Host "1. Railway will automatically detect the new commit" -ForegroundColor White
Write-Host "2. Failed services should start rebuilding" -ForegroundColor White
Write-Host "3. Add RAILWAY_DOCKERFILE_PATH=Dockerfile to any remaining failed services" -ForegroundColor White
Write-Host ""
Write-Host "Monitor your Railway dashboard - services should be green in 5-10 minutes!" -ForegroundColor Green