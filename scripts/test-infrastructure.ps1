#!/usr/bin/env pwsh

Write-Host "Testing Infrastructure Connectivity..." -ForegroundColor Yellow

$tests = @()

# Test PostgreSQL
Write-Host "`nTesting PostgreSQL..." -ForegroundColor Cyan
try {
    $pgResult = docker exec unified-commerce-postgres pg_isready -h localhost -p 5432
    if ($LASTEXITCODE -eq 0) {
        Write-Host "✅ PostgreSQL: Connected" -ForegroundColor Green
        $tests += @{Service="PostgreSQL"; Status="✅ Connected"}
    } else {
        Write-Host "❌ PostgreSQL: Failed" -ForegroundColor Red
        $tests += @{Service="PostgreSQL"; Status="❌ Failed"}
    }
} catch {
    Write-Host "❌ PostgreSQL: Error - $($_.Exception.Message)" -ForegroundColor Red
    $tests += @{Service="PostgreSQL"; Status="❌ Error"}
}

# Test MongoDB
Write-Host "`nTesting MongoDB..." -ForegroundColor Cyan
try {
    $mongoResult = docker exec unified-commerce-mongodb mongosh --eval "db.adminCommand('ping')" --quiet
    if ($LASTEXITCODE -eq 0) {
        Write-Host "✅ MongoDB: Connected" -ForegroundColor Green
        $tests += @{Service="MongoDB"; Status="✅ Connected"}
    } else {
        Write-Host "❌ MongoDB: Failed" -ForegroundColor Red
        $tests += @{Service="MongoDB"; Status="❌ Failed"}
    }
} catch {
    Write-Host "❌ MongoDB: Error - $($_.Exception.Message)" -ForegroundColor Red
    $tests += @{Service="MongoDB"; Status="❌ Error"}
}

# Test Redis
Write-Host "`nTesting Redis..." -ForegroundColor Cyan
try {
    $redisResult = docker exec unified-commerce-redis redis-cli -a redis ping
    if ($redisResult -eq "PONG") {
        Write-Host "✅ Redis: Connected" -ForegroundColor Green
        $tests += @{Service="Redis"; Status="✅ Connected"}
    } else {
        Write-Host "❌ Redis: Failed" -ForegroundColor Red
        $tests += @{Service="Redis"; Status="❌ Failed"}
    }
} catch {
    Write-Host "❌ Redis: Error - $($_.Exception.Message)" -ForegroundColor Red
    $tests += @{Service="Redis"; Status="❌ Error"}
}

# Test Elasticsearch
Write-Host "`nTesting Elasticsearch..." -ForegroundColor Cyan
try {
    $esResult = Invoke-RestMethod -Uri "http://localhost:9200/_cluster/health" -Method Get -TimeoutSec 10
    if ($esResult.status -eq "green" -or $esResult.status -eq "yellow") {
        Write-Host "✅ Elasticsearch: Connected (Status: $($esResult.status))" -ForegroundColor Green
        $tests += @{Service="Elasticsearch"; Status="✅ Connected ($($esResult.status))"}
    } else {
        Write-Host "❌ Elasticsearch: Unhealthy (Status: $($esResult.status))" -ForegroundColor Red
        $tests += @{Service="Elasticsearch"; Status="❌ Unhealthy"}
    }
} catch {
    Write-Host "❌ Elasticsearch: Error - $($_.Exception.Message)" -ForegroundColor Red
    $tests += @{Service="Elasticsearch"; Status="❌ Error"}
}

# Test Kafka
Write-Host "`nTesting Kafka..." -ForegroundColor Cyan
try {
    $kafkaResult = docker exec unified-commerce-kafka kafka-broker-api-versions --bootstrap-server localhost:9092 2>$null
    if ($LASTEXITCODE -eq 0) {
        Write-Host "✅ Kafka: Connected" -ForegroundColor Green
        $tests += @{Service="Kafka"; Status="✅ Connected"}
    } else {
        Write-Host "❌ Kafka: Failed" -ForegroundColor Red
        $tests += @{Service="Kafka"; Status="❌ Failed"}
    }
} catch {
    Write-Host "❌ Kafka: Error - $($_.Exception.Message)" -ForegroundColor Red
    $tests += @{Service="Kafka"; Status="❌ Error"}
}

# Summary
Write-Host "`n" + "="*60 -ForegroundColor Yellow
Write-Host "Infrastructure Test Summary:" -ForegroundColor Yellow
Write-Host "="*60 -ForegroundColor Yellow

foreach ($test in $tests) {
    Write-Host "$($test.Service): $($test.Status)"
}

$successCount = ($tests | Where-Object { $_.Status -like "*Connected*" }).Count
$totalCount = $tests.Count

Write-Host "`nOverall: $successCount/$totalCount services connected" -ForegroundColor $(if ($successCount -eq $totalCount) { "Green" } else { "Yellow" })

if ($successCount -eq $totalCount) {
    Write-Host "`n🎉 All infrastructure services are ready!" -ForegroundColor Green
} else {
    Write-Host "`n⚠️  Some services need attention before proceeding." -ForegroundColor Yellow
}
