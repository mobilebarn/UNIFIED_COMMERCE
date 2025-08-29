#!/usr/bin/env pwsh

param(
    [string]$ServiceName = "",
    [switch]$All = $false
)

$services = @(
    "identity",
    "cart", 
    "order",
    "payment",
    "inventory",
    "product-catalog",
    "promotions",
    "merchant-account"
)

function Start-Service {
    param([string]$service)
    
    $servicePath = "services\$service\cmd\server"
    $envFile = "services\$service\.env"
    
    if (-not (Test-Path $servicePath)) {
        Write-Host "❌ Service path not found: $servicePath" -ForegroundColor Red
        return $false
    }
    
    if (-not (Test-Path $envFile)) {
        Write-Host "❌ Environment file not found: $envFile" -ForegroundColor Red
        return $false
    }
    
    Write-Host "🚀 Starting $service service..." -ForegroundColor Cyan
    
    # Load environment variables
    Get-Content $envFile | ForEach-Object {
        if ($_ -and $_ -notmatch '^#') {
            $name, $value = $_ -split '=', 2
            if ($name -and $value) {
                [Environment]::SetEnvironmentVariable($name.Trim(), $value.Trim())
            }
        }
    }
    
    # Start the service
    $process = Start-Process -FilePath "go" -ArgumentList "run", ".\$servicePath\main.go" -NoNewWindow -PassThru
    
    if ($process) {
        Write-Host "✅ $service service started (PID: $($process.Id))" -ForegroundColor Green
        Start-Sleep -Seconds 2
        return $true
    } else {
        Write-Host "❌ Failed to start $service service" -ForegroundColor Red
        return $false
    }
}

function Test-ServiceHealth {
    param([string]$service, [int]$port)
    
    $healthUrl = "http://localhost:$port/health"
    
    try {
        $response = Invoke-RestMethod -Uri $healthUrl -Method Get -TimeoutSec 5
        if ($response.status -eq "healthy") {
            Write-Host "✅ $service health check passed" -ForegroundColor Green
            return $true
        } else {
            Write-Host "⚠️  $service health check warning: $($response.status)" -ForegroundColor Yellow
            return $false
        }
    } catch {
        Write-Host "❌ $service health check failed: $($_.Exception.Message)" -ForegroundColor Red
        return $false
    }
}

if ($All) {
    Write-Host "Starting all services..." -ForegroundColor Yellow
    
    foreach ($service in $services) {
        $result = Start-Service -service $service
        if (-not $result) {
            Write-Host "❌ Failed to start $service, stopping here" -ForegroundColor Red
            break
        }
    }
} elseif ($ServiceName) {
    if ($services -contains $ServiceName) {
        Start-Service -service $ServiceName
    } else {
        Write-Host "❌ Unknown service: $ServiceName" -ForegroundColor Red
        Write-Host "Available services: $($services -join ', ')" -ForegroundColor Yellow
    }
} else {
    Write-Host "Service Startup Script" -ForegroundColor Yellow
    Write-Host "Usage:"
    Write-Host "  .\start-services.ps1 -ServiceName <service>"
    Write-Host "  .\start-services.ps1 -All"
    Write-Host ""
    Write-Host "Available services: $($services -join ', ')"
}
