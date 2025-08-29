#!/usr/bin/env pwsh

Write-Host "Testing microservice builds..." -ForegroundColor Yellow

$services = @(
    "cart",
    "identity", 
    "inventory",
    "merchant-account",
    "order",
    "payment",
    "product-catalog",
    "promotions"
)

$results = @{}

foreach ($service in $services) {
    $path = "services/$service/cmd/server"
    if (Test-Path $path) {
        Write-Host "Building $service..." -ForegroundColor Cyan
        $output = go build "./$path" 2>&1
        if ($LASTEXITCODE -eq 0) {
            $results[$service] = "SUCCESS"
            Write-Host "${service}: SUCCESS" -ForegroundColor Green
        } else {
            $results[$service] = "FAILED"
            Write-Host "${service}: FAILED" -ForegroundColor Red
            Write-Host "Error: $output" -ForegroundColor DarkRed
        }
    } else {
        $results[$service] = "NOT_FOUND"
        Write-Host "${service}: NOT_FOUND" -ForegroundColor Yellow
    }
}

Write-Host "`nBuild Summary:" -ForegroundColor Yellow
foreach ($service in $services) {
    $status = $results[$service]
    $color = switch($status) {
        "SUCCESS" { "Green" }
        "FAILED" { "Red" }
        "NOT_FOUND" { "Yellow" }
    }
    Write-Host "${service} : $status" -ForegroundColor $color
}
