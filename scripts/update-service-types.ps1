#!/usr/bin/env powershell

# Update all service handler and resolver files to use correct service types

$serviceUpdates = @{
    "order" = @{
        "ServiceType" = "OrderService"
        "VarName" = "orderService"
    }
    "payment" = @{
        "ServiceType" = "PaymentService" 
        "VarName" = "paymentService"
    }
    "inventory" = @{
        "ServiceType" = "InventoryService"
        "VarName" = "inventoryService"  
    }
    "product-catalog" = @{
        "ServiceType" = "ProductService"
        "VarName" = "productService"
    }
    "promotions" = @{
        "ServiceType" = "PromotionsService"
        "VarName" = "promotionsService"
    }
    "merchant-account" = @{
        "ServiceType" = "MerchantService" 
        "VarName" = "merchantService"
    }
}

foreach ($service in $serviceUpdates.Keys) {
    $serviceInfo = $serviceUpdates[$service]
    $serviceType = $serviceInfo.ServiceType
    $varName = $serviceInfo.VarName
    
    Write-Host "🔧 Updating $service to use $serviceType..." -ForegroundColor Cyan
    
    # Update handler.go
    $handlerPath = "services\$service\graphql\handler.go"
    if (Test-Path $handlerPath) {
        $content = Get-Content $handlerPath -Raw
        $content = $content -replace 'unified-commerce/services/identity/service', "unified-commerce/services/$service/service"
        $content = $content -replace 'identityService \*service\.IdentityService', "$varName *service.$serviceType"
        $content = $content -replace 'identityService, logger', "$varName, logger"
        Set-Content -Path $handlerPath -Value $content -Encoding UTF8
        Write-Host "  ✅ Updated handler.go for $service" -ForegroundColor Green
    }
    
    # Update resolver.go  
    $resolverPath = "services\$service\graphql\resolver.go"
    if (Test-Path $resolverPath) {
        $content = Get-Content $resolverPath -Raw
        $content = $content -replace 'unified-commerce/services/identity/service', "unified-commerce/services/$service/service"
        $content = $content -replace 'IdentityService \*service\.IdentityService', "$serviceType *service.$serviceType"
        $content = $content -replace 'identityService \*service\.IdentityService', "$varName *service.$serviceType"
        $content = $content -replace 'IdentityService: identityService', "$serviceType`: $varName"
        Set-Content -Path $resolverPath -Value $content -Encoding UTF8
        Write-Host "  ✅ Updated resolver.go for $service" -ForegroundColor Green
    }
}

Write-Host "`n🎯 All services updated to enterprise architecture!" -ForegroundColor Green
