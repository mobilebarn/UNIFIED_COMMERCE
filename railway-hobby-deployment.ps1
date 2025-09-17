# RETAIL OS RAILWAY HOBBY PLAN DEPLOYMENT
# Optimized deployment script for Railway Hobby plan
# Takes advantage of higher limits and better performance

Write-Host "RETAIL OS RAILWAY HOBBY DEPLOYMENT" -ForegroundColor Green
Write-Host "==================================" -ForegroundColor Green
Write-Host ""

Write-Host "Railway Hobby Plan Benefits:" -ForegroundColor Yellow
Write-Host "✅ Higher resource limits (no file size issues)" -ForegroundColor Green
Write-Host "✅ Better build performance" -ForegroundColor Green
Write-Host "✅ Priority support" -ForegroundColor Green
Write-Host "✅ More concurrent deployments" -ForegroundColor Green
Write-Host ""

Write-Host "Let's make your $20 investment count!" -ForegroundColor Cyan
Write-Host ""

# First, let's try direct deployment from the main directory
# The Hobby plan should handle the larger file size
Write-Host "Step 1: Testing direct deployment with Hobby plan..." -ForegroundColor Yellow

# Check current Railway status
Write-Host "Checking Railway connection..." -ForegroundColor Cyan
railway status

Write-Host ""
Write-Host "Step 2: Let's try deploying the Product service first" -ForegroundColor Yellow
Write-Host "With Hobby plan, this should work without file size issues" -ForegroundColor Cyan

# Navigate to Product service
Set-Location "services\product-catalog"

Write-Host "Current directory: $(Get-Location)" -ForegroundColor Gray
Write-Host "Attempting deployment..." -ForegroundColor Cyan

# Try deployment - Hobby plan should handle this better
railway up --service retail-os-product --detach

Write-Host ""
Write-Host "If this works, we'll automate the rest!" -ForegroundColor Green
Write-Host "The Hobby plan should resolve our previous issues." -ForegroundColor Yellow