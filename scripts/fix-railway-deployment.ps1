# 🚀 Railway Deployment Fix Script
Write-Host "🔧 Fixing Railway Deployment - Adding Missing Dockerfiles" -ForegroundColor Cyan

# Check if we're in a git repository
if (-not (Test-Path ".git")) {
    Write-Host "❌ Not in a git repository. Initializing..." -ForegroundColor Red
    git init
    git remote add origin https://github.com/YOUR_USERNAME/UNIFIED_COMMERCE.git
}

# Add all new Dockerfiles
Write-Host "📦 Adding Dockerfiles to git..." -ForegroundColor Yellow
git add services/*/Dockerfile
git add gateway/Dockerfile

# Add other changed files
git add .

# Commit the changes
Write-Host "💾 Committing changes..." -ForegroundColor Yellow
git commit -m "🔧 Add Dockerfiles for Railway deployment

- Added Dockerfiles to all 8 microservices
- Added Dockerfile to GraphQL gateway
- Fixed Railway deployment issue: 'dockerfile does not exist'
- All services now use Go 1.21 for compatibility
- Gateway uses Node.js 18

This should resolve the build failures in Railway."

# Push to GitHub
Write-Host "🚀 Pushing to GitHub..." -ForegroundColor Green
git push origin main

Write-Host "✅ Done! Railway should automatically redeploy all services." -ForegroundColor Green
Write-Host ""
Write-Host "🔍 Next steps:" -ForegroundColor Cyan
Write-Host "1. Go to Railway dashboard" -ForegroundColor White
Write-Host "2. Watch services rebuild (should turn green)" -ForegroundColor White
Write-Host "3. Once all services are green, configure gateway environment variables" -ForegroundColor White
Write-Host "4. Use the RAILWAY-COMPLETE-SOLUTION.html tool" -ForegroundColor White

Write-Host ""
Write-Host "⏱️  Estimated time: 5-10 minutes for all services to rebuild" -ForegroundColor Yellow