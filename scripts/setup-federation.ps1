Write-Host "Proper Federation Setup - Regenerating for each service" -ForegroundColor Green

$projectRoot = "c:\Users\dane\OneDrive\Desktop\UNIFIED_COMMERCE"
$services = @("payment", "inventory", "product-catalog", "promotions", "merchant-account")

foreach ($service in $services) {
    Write-Host "Properly setting up $service federation..." -ForegroundColor Yellow
    
    $servicePath = "$projectRoot\services\$service\graphql"
    
    if (Test-Path $servicePath) {
        Write-Host "  Cleaning old federation files..." -ForegroundColor Cyan
        
        # Remove incorrectly copied files
        @("entity.resolvers.go", "generated.go", "models_gen.go", "resolver.go", "schema.resolvers.go") | ForEach-Object {
            if (Test-Path "$servicePath\$_") {
                Remove-Item "$servicePath\$_" -Force
                Write-Host "    Removed $_" -ForegroundColor Yellow
            }
        }
        
        Write-Host "  Regenerating GraphQL federation files..." -ForegroundColor Cyan
        Push-Location $servicePath
        
        # Generate federation files from scratch based on the service's actual schema
        go run github.com/99designs/gqlgen generate --config gqlgen.yml
        
        if ($LASTEXITCODE -eq 0) {
            Write-Host "  ✅ $service federation generated successfully" -ForegroundColor Green
        } else {
            Write-Host "  ❌ $service federation generation failed" -ForegroundColor Red
        }
        
        Pop-Location
        Write-Host ""
    } else {
        Write-Host "  ❌ $servicePath not found" -ForegroundColor Red
    }
}

Write-Host "Proper federation setup complete!" -ForegroundColor Green