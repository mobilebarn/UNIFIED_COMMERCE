# Force Deploy Script for Storefront

# This script will force deploy the storefront to Vercel
# Run this if GitHub integration is not working

Write-Host \"🚀 Starting manual Vercel deployment...\" -ForegroundColor Green

# Navigate to storefront directory
Set-Location \"C:\\Users\\dane\\OneDrive\\Desktop\\UNIFIED_COMMERCE\\storefront\"

# Remove any existing Vercel configuration
if (Test-Path \".vercel\") {
    Remove-Item -Recurse -Force \".vercel\"
    Write-Host \"✅ Removed old Vercel config\" -ForegroundColor Yellow
}

# Install dependencies
Write-Host \"📦 Installing dependencies...\" -ForegroundColor Cyan
npm install

# Build the project
Write-Host \"🔨 Building project...\" -ForegroundColor Cyan
npm run build

# Deploy to Vercel with force flag
Write-Host \"🚀 Deploying to Vercel...\" -ForegroundColor Green
vercel --prod --yes --force

Write-Host \"✅ Deployment completed!\" -ForegroundColor Green
Write-Host \"🔗 Check your Vercel dashboard for the new URL\" -ForegroundColor Yellow

# Pause to show results
Read-Host \"Press Enter to continue...\"