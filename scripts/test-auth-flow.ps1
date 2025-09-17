#!/usr/bin/env pwsh
# Complete authentication flow test script

Write-Host "=== UNIFIED COMMERCE AUTHENTICATION FLOW TEST ===" -ForegroundColor Green
Write-Host ""

$ErrorActionPreference = "Continue"

# Test 1: Check infrastructure services
Write-Host "1. Testing Infrastructure Services..." -ForegroundColor Yellow
try {
    $docker_status = docker ps --format "table {{.Names}}\t{{.Status}}" | Select-String "postgres|redis|mongodb|kafka"
    if ($docker_status) {
        Write-Host "✓ Docker infrastructure services are running" -ForegroundColor Green
        $docker_status | ForEach-Object { Write-Host "  $_" -ForegroundColor Gray }
    } else {
        Write-Host "✗ Infrastructure services not running" -ForegroundColor Red
        exit 1
    }
} catch {
    Write-Host "✗ Error checking Docker services: $($_.Exception.Message)" -ForegroundColor Red
    exit 1
}

# Test 2: Identity service health check  
Write-Host "`n2. Testing Identity Service Health..." -ForegroundColor Yellow
try {
    $health_response = Invoke-RestMethod -Uri "http://localhost:8001/health" -Method GET -TimeoutSec 5
    if ($health_response.status -eq "ok") {
        Write-Host "✓ Identity service is healthy" -ForegroundColor Green
        Write-Host "  Database: $($health_response.database)" -ForegroundColor Gray
        Write-Host "  Cache: $($health_response.cache)" -ForegroundColor Gray
    } else {
        Write-Host "✗ Identity service health check failed" -ForegroundColor Red
        exit 1
    }
} catch {
    Write-Host "✗ Identity service not accessible: $($_.Exception.Message)" -ForegroundColor Red
    exit 1
}

# Test 3: GraphQL login mutation
Write-Host "`n3. Testing GraphQL Login Mutation..." -ForegroundColor Yellow
try {
    $login_mutation = @{
        query = "mutation(`$email:String!,`$password:String!){ login(input:{email:`$email,password:`$password}){ accessToken expiresIn user { email username firstName lastName } } }"
        variables = @{
            email = "admin@example.com"
            password = "Admin123!"
        }
    }
    
    $json_body = $login_mutation | ConvertTo-Json -Compress
    $login_response = Invoke-RestMethod -Uri "http://localhost:8001/graphql" -Method POST -Body $json_body -ContentType "application/json" -TimeoutSec 10
    
    if ($login_response.data.login.accessToken) {
        Write-Host "✓ GraphQL login mutation successful" -ForegroundColor Green
        Write-Host "  User: $($login_response.data.login.user.email)" -ForegroundColor Gray
        Write-Host "  Token length: $($login_response.data.login.accessToken.Length) chars" -ForegroundColor Gray
        Write-Host "  Expires in: $($login_response.data.login.expiresIn) seconds" -ForegroundColor Gray
        $script:access_token = $login_response.data.login.accessToken
    } else {
        Write-Host "✗ GraphQL login failed - no access token returned" -ForegroundColor Red
        if ($login_response.errors) {
            $login_response.errors | ForEach-Object { Write-Host "  Error: $($_.message)" -ForegroundColor Red }
        }
        exit 1
    }
} catch {
    Write-Host "✗ GraphQL login mutation failed: $($_.Exception.Message)" -ForegroundColor Red
    exit 1
}

# Test 4: Frontend accessibility
Write-Host "`n4. Testing Frontend Accessibility..." -ForegroundColor Yellow
try {
    $frontend_response = Invoke-WebRequest -Uri "http://localhost:3003" -Method GET -TimeoutSec 5
    if ($frontend_response.StatusCode -eq 200) {
        Write-Host "✓ Frontend is accessible on port 3003" -ForegroundColor Green
        $content_length = $frontend_response.Content.Length
        Write-Host "  Response size: $content_length bytes" -ForegroundColor Gray
        
        if ($frontend_response.Content -match "admin-panel|login|unified commerce") {
            Write-Host "✓ Frontend contains expected content" -ForegroundColor Green
        } else {
            Write-Host "? Frontend content may not be fully loaded" -ForegroundColor Yellow
        }
    } else {
        Write-Host "✗ Frontend not accessible - Status: $($frontend_response.StatusCode)" -ForegroundColor Red
        exit 1
    }
} catch {
    Write-Host "✗ Frontend not accessible: $($_.Exception.Message)" -ForegroundColor Red
    exit 1
}

# Test 5: Integration test execution
Write-Host "`n5. Running Integration Tests..." -ForegroundColor Yellow
try {
    Push-Location "c:\Users\dane\OneDrive\Desktop\UNIFIED_COMMERCE\services\identity"
    $test_output = go test ./integration -v 2>&1
    $test_exit_code = $LASTEXITCODE
    Pop-Location
    
    if ($test_exit_code -eq 0) {
        Write-Host "✓ Integration tests passed" -ForegroundColor Green
        # Count passed tests
        $passed_tests = ($test_output | Select-String "--- PASS:").Count
        Write-Host "  Passed tests: $passed_tests" -ForegroundColor Gray
    } else {
        Write-Host "✗ Integration tests failed" -ForegroundColor Red
        Write-Host $test_output -ForegroundColor Red
        exit 1
    }
} catch {
    Write-Host "✗ Error running integration tests: $($_.Exception.Message)" -ForegroundColor Red
    exit 1
}

# Summary
Write-Host "`n=== AUTHENTICATION FLOW TEST SUMMARY ===" -ForegroundColor Green
Write-Host "✓ Infrastructure services running" -ForegroundColor Green
Write-Host "✓ Identity service healthy" -ForegroundColor Green  
Write-Host "✓ GraphQL authentication working" -ForegroundColor Green
Write-Host "✓ Frontend accessible" -ForegroundColor Green
Write-Host "✓ Integration tests passing" -ForegroundColor Green

Write-Host "`nNext Steps:" -ForegroundColor Cyan
Write-Host "1. Open browser to http://localhost:3003" -ForegroundColor White
Write-Host "2. Use demo credentials: admin@example.com / Admin123!" -ForegroundColor White
Write-Host "3. Test the complete login flow" -ForegroundColor White

Write-Host "`nAuthentication flow is ready for testing!" -ForegroundColor Green
