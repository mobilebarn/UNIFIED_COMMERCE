# Railway Service Connectivity Fix
# This script configures environment variables for Railway service internal connectivity

Write-Host "🔧 Railway Service Connectivity Fix" -ForegroundColor Cyan
Write-Host "=======================================" -ForegroundColor Cyan

# Railway service internal URL pattern: https://{service-name}-production.up.railway.app
# For GraphQL endpoints, add /graphql to the end

$services = @{
    "gateway" = @{
        "IDENTITY_SERVICE_URL" = "https://identity-service-production.up.railway.app/graphql"
        "CART_SERVICE_URL" = "https://cart-service-production.up.railway.app/graphql"
        "ORDER_SERVICE_URL" = "https://order-service-production.up.railway.app/graphql" 
        "PAYMENT_SERVICE_URL" = "https://payment-service-production.up.railway.app/graphql"
        "INVENTORY_SERVICE_URL" = "https://inventory-service-production.up.railway.app/graphql"
        "PRODUCT_SERVICE_URL" = "https://product-catalog-service-production.up.railway.app/graphql"
        "PROMOTIONS_SERVICE_URL" = "https://promotions-service-production.up.railway.app/graphql"
        "MERCHANT_SERVICE_URL" = "https://merchant-account-service-production.up.railway.app/graphql"
        "CORS_ORIGINS" = "https://storefront-eta-six.vercel.app,https://admin-panel-tau-eight.vercel.app"
        "NODE_ENV" = "production"
    }
}

Write-Host "📝 Environment Variables to Configure:" -ForegroundColor Yellow
Write-Host ""

foreach ($service in $services.Keys) {
    Write-Host "Service: $service" -ForegroundColor Green
    
    foreach ($envVar in $services[$service].Keys) {
        $value = $services[$service][$envVar]
        Write-Host "  $envVar = $value" -ForegroundColor Gray
    }
    Write-Host ""
}

Write-Host "🚀 Railway Dashboard Configuration Steps:" -ForegroundColor Cyan
Write-Host ""
Write-Host "1. Go to Railway Dashboard: https://railway.app/dashboard" -ForegroundColor White
Write-Host "2. Select your project" -ForegroundColor White
Write-Host "3. Click on 'GraphQL Gateway' service" -ForegroundColor White
Write-Host "4. Go to 'Variables' tab" -ForegroundColor White
Write-Host "5. Add each environment variable listed above" -ForegroundColor White
Write-Host "6. Deploy the gateway service" -ForegroundColor White
Write-Host ""

Write-Host "⚠️  IMPORTANT NOTES:" -ForegroundColor Yellow
Write-Host ""
Write-Host "• Replace 'production' with your actual Railway environment name if different" -ForegroundColor Red
Write-Host "• Update service names to match your actual Railway service names" -ForegroundColor Red
Write-Host "• Ensure all backend services are deployed and running before configuring gateway" -ForegroundColor Red
Write-Host "• Check Railway service names in your dashboard for exact URLs" -ForegroundColor Red

# Create a Railway environment variables template file
$envTemplate = @"
# Railway Environment Variables for GraphQL Gateway
# Copy these to Railway Dashboard > Service > Variables

# Backend Service URLs (replace with actual Railway URLs)
IDENTITY_SERVICE_URL=https://identity-service-production.up.railway.app/graphql
CART_SERVICE_URL=https://cart-service-production.up.railway.app/graphql
ORDER_SERVICE_URL=https://order-service-production.up.railway.app/graphql
PAYMENT_SERVICE_URL=https://payment-service-production.up.railway.app/graphql
INVENTORY_SERVICE_URL=https://inventory-service-production.up.railway.app/graphql
PRODUCT_SERVICE_URL=https://product-catalog-service-production.up.railway.app/graphql
PROMOTIONS_SERVICE_URL=https://promotions-service-production.up.railway.app/graphql
MERCHANT_SERVICE_URL=https://merchant-account-service-production.up.railway.app/graphql

# Frontend URLs for CORS
CORS_ORIGINS=https://storefront-eta-six.vercel.app,https://admin-panel-tau-eight.vercel.app

# Environment
NODE_ENV=production

# JWT Secret (use Railway's built-in secrets management)
JWT_SECRET=your-production-jwt-secret-here
"@

$envTemplate | Out-File -FilePath "railway-environment-variables.txt" -Encoding UTF8

Write-Host "📄 Environment variables template saved to 'railway-environment-variables.txt'" -ForegroundColor Green
Write-Host ""

Write-Host "🔍 How to find your actual Railway service URLs:" -ForegroundColor Cyan
Write-Host ""
Write-Host "1. In Railway Dashboard, click on each service" -ForegroundColor White
Write-Host "2. Go to 'Settings' tab" -ForegroundColor White  
Write-Host "3. Look for 'Public Domain' or 'Deployment URL'" -ForegroundColor White
Write-Host "4. Copy the URL and add '/graphql' to the end" -ForegroundColor White
Write-Host "5. Update the environment variables accordingly" -ForegroundColor White
Write-Host ""

Write-Host "✅ Next Steps:" -ForegroundColor Green
Write-Host ""
Write-Host "1. Ensure all 8 backend services are deployed and running" -ForegroundColor White
Write-Host "2. Get the actual Railway URLs for each service" -ForegroundColor White
Write-Host "3. Update the environment variables in the gateway service" -ForegroundColor White
Write-Host "4. Redeploy the gateway service" -ForegroundColor White
Write-Host "5. Test the gateway connectivity" -ForegroundColor White

Write-Host ""
Write-Host "🎯 Expected Result: Gateway should successfully connect to all services!" -ForegroundColor Green