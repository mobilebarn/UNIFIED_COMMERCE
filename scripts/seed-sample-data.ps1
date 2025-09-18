# Unified Commerce OS - Sample Data Seeding Script (PowerShell)
# This script runs the Node.js seeding script with proper environment setup

param(
    [string]$GatewayUrl = "https://unified-commerce-gateway.onrender.com/graphql",
    [switch]$Help
)

if ($Help) {
    Write-Host "Unified Commerce OS - Sample Data Seeding" -ForegroundColor Green
    Write-Host "=========================================" -ForegroundColor Green
    Write-Host ""
    Write-Host "This script populates your Unified Commerce OS with sample data including:"
    Write-Host "  • Product categories (Electronics, Fashion, Home & Kitchen, etc.)"
    Write-Host "  • Sample products with realistic details"
    Write-Host "  • Customer accounts for testing"
    Write-Host "  • Promotional campaigns and discount codes"
    Write-Host ""
    Write-Host "Usage:" -ForegroundColor Yellow
    Write-Host "  .\scripts\seed-sample-data.ps1                    # Use default gateway URL"
    Write-Host "  .\scripts\seed-sample-data.ps1 -GatewayUrl <url>  # Use custom gateway URL"
    Write-Host "  .\scripts\seed-sample-data.ps1 -Help              # Show this help"
    Write-Host ""
    Write-Host "Examples:" -ForegroundColor Yellow
    Write-Host "  .\scripts\seed-sample-data.ps1"
    Write-Host "  .\scripts\seed-sample-data.ps1 -GatewayUrl 'http://localhost:4000/graphql'"
    Write-Host ""
    return
}

Write-Host "🌱 Unified Commerce OS - Sample Data Seeding" -ForegroundColor Green
Write-Host "=============================================" -ForegroundColor Green
Write-Host ""

# Check if Node.js is available
try {
    $nodeVersion = node --version 2>$null
    if ($nodeVersion) {
        Write-Host "✅ Node.js found: $nodeVersion" -ForegroundColor Green
    } else {
        throw "Node.js not found"
    }
} catch {
    Write-Host "❌ Node.js is required but not found" -ForegroundColor Red
    Write-Host ""
    Write-Host "Please install Node.js from https://nodejs.org/" -ForegroundColor Yellow
    Write-Host "Then run this script again." -ForegroundColor Yellow
    exit 1
}

# Set environment variable for the seeding script
$env:GATEWAY_URL = $GatewayUrl

Write-Host "🎯 Gateway URL: $GatewayUrl" -ForegroundColor Cyan
Write-Host ""

# Change to the project root directory
$scriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$projectRoot = Split-Path -Parent $scriptDir
Set-Location $projectRoot

Write-Host "📂 Working directory: $(Get-Location)" -ForegroundColor Gray
Write-Host ""

try {
    # Run the Node.js seeding script
    Write-Host "🚀 Starting sample data seeding..." -ForegroundColor Yellow
    Write-Host ""
    
    node .\scripts\seed-sample-data.js
    
    if ($LASTEXITCODE -eq 0) {
        Write-Host ""
        Write-Host "🎉 Sample data seeding completed successfully!" -ForegroundColor Green
        Write-Host ""
        Write-Host "Your Unified Commerce OS is now populated with:" -ForegroundColor Cyan
        Write-Host "  • Product categories and sample products" -ForegroundColor White
        Write-Host "  • Customer accounts for testing" -ForegroundColor White
        Write-Host "  • Promotional campaigns and discount codes" -ForegroundColor White
        Write-Host ""
        Write-Host "Next steps:" -ForegroundColor Yellow
        Write-Host "  1. Visit your storefront to browse products" -ForegroundColor White
        Write-Host "  2. Use the admin panel to manage data" -ForegroundColor White
        Write-Host "  3. Test GraphQL queries through the gateway" -ForegroundColor White
        Write-Host ""
    } else {
        Write-Host ""
        Write-Host "⚠️  Seeding completed with some warnings." -ForegroundColor Yellow
        Write-Host "Some data might already exist or services might be unavailable." -ForegroundColor Yellow
        Write-Host "This is normal and expected." -ForegroundColor Yellow
        Write-Host ""
    }
} catch {
    Write-Host ""
    Write-Host "❌ Error during seeding: $($_.Exception.Message)" -ForegroundColor Red
    Write-Host ""
    Write-Host "Troubleshooting:" -ForegroundColor Yellow
    Write-Host "  1. Ensure the GraphQL Gateway is running" -ForegroundColor White
    Write-Host "  2. Check your network connection" -ForegroundColor White
    Write-Host "  3. Verify the gateway URL is correct" -ForegroundColor White
    Write-Host ""
    exit 1
}