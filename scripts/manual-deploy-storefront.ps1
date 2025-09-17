# Retail OS Storefront - Manual Deploy Script

# Remove existing Vercel project link
if (Test-Path \".vercel\") {
    Remove-Item -Recurse -Force \".vercel\"
}

# Navigate to storefront directory
Set-Location \"C:\\Users\\dane\\OneDrive\\Desktop\\UNIFIED_COMMERCE\\storefront\"

# Install dependencies if needed
Write-Host \"Installing dependencies...\" -ForegroundColor Green
npm install

# Build the project
Write-Host \"Building project...\" -ForegroundColor Green
npm run build

# Deploy to new Vercel project
Write-Host \"Deploying to Vercel...\" -ForegroundColor Green
vercel --prod --yes --force

Write-Host \"Deployment complete!\" -ForegroundColor Green
Write-Host \"Check your Vercel dashboard for the new URL\" -ForegroundColor Yellow