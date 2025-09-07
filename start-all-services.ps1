# Script to start all Unified Commerce services
Write-Host "Starting all Unified Commerce services..." -ForegroundColor Green

# Start Identity Service (8001)
Write-Host "Starting Identity Service on port 8001..." -ForegroundColor Yellow
Start-Process -NoNewWindow -FilePath "go" -ArgumentList "run", "cmd/server/main.go" -WorkingDirectory ".\services\identity"

# Start Cart Service (8002)
Write-Host "Starting Cart Service on port 8002..." -ForegroundColor Yellow
Start-Process -NoNewWindow -FilePath "go" -ArgumentList "run", "cmd/server/main.go" -WorkingDirectory ".\services\cart"

# Start Order Service (8003)
Write-Host "Starting Order Service on port 8003..." -ForegroundColor Yellow
Start-Process -NoNewWindow -FilePath "go" -ArgumentList "run", "cmd/server/main.go" -WorkingDirectory ".\services\order"

# Start Payment Service (8004)
Write-Host "Starting Payment Service on port 8004..." -ForegroundColor Yellow
Start-Process -NoNewWindow -FilePath "go" -ArgumentList "run", "cmd/server/main.go" -WorkingDirectory ".\services\payment"

# Start Inventory Service (8005)
Write-Host "Starting Inventory Service on port 8005..." -ForegroundColor Yellow
Start-Process -NoNewWindow -FilePath "go" -ArgumentList "run", "cmd/server/main.go" -WorkingDirectory ".\services\inventory"

# Start Product Catalog Service (8006)
Write-Host "Starting Product Catalog Service on port 8006..." -ForegroundColor Yellow
Start-Process -NoNewWindow -FilePath "go" -ArgumentList "run", "cmd/server/main.go" -WorkingDirectory ".\services\product-catalog"

# Start Promotions Service (8007)
Write-Host "Starting Promotions Service on port 8007..." -ForegroundColor Yellow
Start-Process -NoNewWindow -FilePath "go" -ArgumentList "run", "cmd/server/main.go" -WorkingDirectory ".\services\promotions"

# Start Merchant Account Service (8008)
Write-Host "Starting Merchant Account Service on port 8008..." -ForegroundColor Yellow
Start-Process -NoNewWindow -FilePath "go" -ArgumentList "run", "cmd/server/main.go" -WorkingDirectory ".\services\merchant-account"

Write-Host "All services started! Waiting 10 seconds for services to initialize..." -ForegroundColor Green
Start-Sleep -Seconds 30

Write-Host "Checking service status..." -ForegroundColor Green
$ports = @(8001, 8002, 8003, 8004, 8005, 8006, 8007, 8008)
foreach ($port in $ports) {
    $status = try {
        Invoke-WebRequest -Uri "http://localhost:$port/health" -UseBasicParsing -TimeoutSec 5
        "RUNNING"
    } catch {
        "NOT RESPONDING"
    }
    Write-Host "Port $port : $status" -ForegroundColor $(if ($status -eq "RUNNING") { "Green" } else { "Red" })
}

Write-Host "Services startup complete!" -ForegroundColor Green