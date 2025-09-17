# Retail OS Automated Railway CLI Deployment Script
# This script deploys all services using Railway CLI commands

Write-Host "🚀 Starting Retail OS Automated Railway Deployment" -ForegroundColor Green
Write-Host "======================================================" -ForegroundColor Green

# Function to deploy a service
function Deploy-RetailOSService {
    param(
        [string]$ServiceName,
        [string]$ServicePath,
        [hashtable]$EnvVars
    )
    
    Write-Host "📦 Deploying $ServiceName..." -ForegroundColor Blue
    
    # Navigate to service directory
    Push-Location $ServicePath
    
    try {
        # Create empty service
        Write-Host "   Creating service: $ServiceName" -ForegroundColor Yellow
        railway add --service
        
        # Set environment variables
        Write-Host "   Setting environment variables..." -ForegroundColor Yellow
        foreach ($env in $EnvVars.GetEnumerator()) {
            railway variables set "$($env.Key)=$($env.Value)"
        }
        
        # Deploy the service
        Write-Host "   Deploying service..." -ForegroundColor Yellow
        railway up --detach
        
        Write-Host "✅ $ServiceName deployed successfully!" -ForegroundColor Green
    }
    catch {
        Write-Host "❌ Failed to deploy $ServiceName : $_" -ForegroundColor Red
    }
    finally {
        Pop-Location
    }
}

# Backend Services Configuration
$backendServices = @(
    @{
        Name = "Identity Service"
        Path = "services/identity"
        EnvVars = @{
            "SERVICE_NAME" = "identity"
            "ENVIRONMENT" = "production"
            "JWT_SECRET" = "prod-secure-jwt-secret-identity-2024"
            "JWT_EXPIRATION" = "86400"
            "LOG_LEVEL" = "info"
        }
    },
    @{
        Name = "Merchant Service"
        Path = "services/merchant"
        EnvVars = @{
            "SERVICE_NAME" = "merchant"
            "ENVIRONMENT" = "production"
            "JWT_SECRET" = "prod-secure-jwt-secret-merchant-2024"
            "LOG_LEVEL" = "info"
        }
    },
    @{
        Name = "Product Service"
        Path = "services/product"
        EnvVars = @{
            "SERVICE_NAME" = "product"
            "ENVIRONMENT" = "production"
            "JWT_SECRET" = "prod-secure-jwt-secret-product-2024"
            "LOG_LEVEL" = "info"
        }
    },
    @{
        Name = "Inventory Service"
        Path = "services/inventory"
        EnvVars = @{
            "SERVICE_NAME" = "inventory"
            "ENVIRONMENT" = "production"
            "JWT_SECRET" = "prod-secure-jwt-secret-inventory-2024"
            "LOG_LEVEL" = "info"
        }
    },
    @{
        Name = "Order Service"
        Path = "services/order"
        EnvVars = @{
            "SERVICE_NAME" = "order"
            "ENVIRONMENT" = "production"
            "JWT_SECRET" = "prod-secure-jwt-secret-order-2024"
            "LOG_LEVEL" = "info"
        }
    },
    @{
        Name = "Payment Service"
        Path = "services/payment"
        EnvVars = @{
            "SERVICE_NAME" = "payment"
            "SERVICE_PORT" = "8005"
            "ENVIRONMENT" = "production"
            "JWT_SECRET" = "prod-secure-jwt-secret-payment-2024"
            "LOG_LEVEL" = "info"
        }
    },
    @{
        Name = "Cart Service"
        Path = "services/cart"
        EnvVars = @{
            "SERVICE_NAME" = "cart"
            "ENVIRONMENT" = "production"
            "JWT_SECRET" = "prod-secure-jwt-secret-cart-2024"
            "LOG_LEVEL" = "info"
        }
    },
    @{
        Name = "Promotions Service"
        Path = "services/promotions"
        EnvVars = @{
            "SERVICE_NAME" = "promotions"
            "ENVIRONMENT" = "production"
            "JWT_SECRET" = "prod-secure-jwt-secret-promotions-2024"
            "LOG_LEVEL" = "info"
        }
    },
    @{
        Name = "Analytics Service"
        Path = "services/analytics"
        EnvVars = @{
            "SERVICE_NAME" = "analytics"
            "ENVIRONMENT" = "production"
            "JWT_SECRET" = "prod-secure-jwt-secret-analytics-2024"
            "LOG_LEVEL" = "info"
        }
    }
)

# Deploy all backend services
Write-Host "🔧 Deploying Backend Services..." -ForegroundColor Cyan
foreach ($service in $backendServices) {
    Deploy-RetailOSService -ServiceName $service.Name -ServicePath $service.Path -EnvVars $service.EnvVars
    Start-Sleep -Seconds 5  # Wait between deployments
}

# Deploy GraphQL Gateway
Write-Host "🌐 Deploying GraphQL Federation Gateway..." -ForegroundColor Cyan
Deploy-RetailOSService -ServiceName "GraphQL Gateway" -ServicePath "gateway" -EnvVars @{
    "NODE_ENV" = "production"
}

# Deploy Frontend Applications
Write-Host "🖥️  Deploying Frontend Applications..." -ForegroundColor Cyan

# Storefront
Deploy-RetailOSService -ServiceName "Storefront" -ServicePath "apps/storefront" -EnvVars @{
    "NODE_ENV" = "production"
    "NEXT_PUBLIC_API_URL" = "https://retail-os-gateway.railway.app"
}

# Admin Panel
Deploy-RetailOSService -ServiceName "Admin Panel" -ServicePath "apps/admin" -EnvVars @{
    "NODE_ENV" = "production"
    "REACT_APP_API_URL" = "https://retail-os-gateway.railway.app"
}

Write-Host "🎉 Retail OS Deployment Complete!" -ForegroundColor Green
Write-Host "======================================" -ForegroundColor Green
Write-Host "📊 Check your services at: https://railway.app/dashboard" -ForegroundColor Yellow
Write-Host "🌍 Your applications will be available at:" -ForegroundColor Yellow
Write-Host "   - Storefront: https://retail-os-storefront.railway.app" -ForegroundColor Cyan
Write-Host "   - Admin Panel: https://retail-os-admin.railway.app" -ForegroundColor Cyan
Write-Host "   - GraphQL Gateway: https://retail-os-gateway.railway.app" -ForegroundColor Cyan