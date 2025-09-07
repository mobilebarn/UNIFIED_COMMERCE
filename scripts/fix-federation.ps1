# Federation Fix Script
# This script standardizes all services to use proper gqlgen federation

param(
    [switch]$Force
)

Write-Host "🔧 Fixing GraphQL Federation for all services..." -ForegroundColor Yellow

$services = @("identity", "cart", "order", "payment", "inventory", "product-catalog", "promotions", "merchant-account")
$rootDir = "C:\Users\dane\OneDrive\Desktop\UNIFIED_COMMERCE"

foreach ($service in $services) {
    Write-Host "`n📦 Processing $service service..." -ForegroundColor Cyan
    
    $serviceDir = "$rootDir\services\$service"
    $graphqlDir = "$serviceDir\graphql"
    
    if (Test-Path $graphqlDir) {
        Set-Location $graphqlDir
        
        # Remove existing generated files to force clean regeneration
        $filesToRemove = @("generated.go", "federation.go", "models_gen.go")
        foreach ($file in $filesToRemove) {
            if (Test-Path $file) {
                Write-Host "  🗑️ Removing old $file" -ForegroundColor Gray
                Remove-Item $file -Force
            }
        }
        
        # Regenerate federation files
        Write-Host "  🔄 Regenerating federation files..." -ForegroundColor Yellow
        try {
            $output = & go run github.com/99designs/gqlgen generate --config gqlgen.yml 2>&1
            if ($LASTEXITCODE -eq 0) {
                Write-Host "  ✅ $service federation files generated successfully" -ForegroundColor Green
            } else {
                Write-Host "  ❌ $service federation generation failed: $output" -ForegroundColor Red
            }
        } catch {
            Write-Host "  ❌ $service federation generation error: $_" -ForegroundColor Red
        }
    } else {
        Write-Host "  ⚠️ GraphQL directory not found for $service" -ForegroundColor Yellow
    }
}

Set-Location $rootDir
Write-Host "`n🏁 Federation fix completed!" -ForegroundColor Green
