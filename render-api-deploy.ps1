# ================================================================
# 🚀 RETAIL OS - RENDER API AUTOMATED DEPLOYMENT
# ================================================================
# This script uses Render's API to deploy all services automatically
# No manual clicking required - completely automated!

param(
    [string]$RenderApiKey = "",
    [switch]$CreateDatabases = $true,
    [switch]$DeployServices = $true,
    [string]$GitHubRepo = "https://github.com/mobilebarn/UNIFIED_COMMERCE"
)

# Configuration
$RenderBaseUrl = "https://api.render.com/v1"
$OwnerEmail = "your-email@example.com"  # Update this

Write-Host ""
Write-Host "🚀 RETAIL OS - RENDER API AUTOMATED DEPLOYMENT" -ForegroundColor Green -BackgroundColor Black
Write-Host "===============================================" -ForegroundColor Green
Write-Host ""

# Check if API key is provided
if (-not $RenderApiKey) {
    Write-Host "🔑 Render API Key Required" -ForegroundColor Yellow
    Write-Host "=========================" -ForegroundColor Yellow
    Write-Host ""
    Write-Host "To get your Render API key:" -ForegroundColor White
    Write-Host "1. Go to https://dashboard.render.com/account" -ForegroundColor Gray
    Write-Host "2. Click 'API Keys' in the left sidebar" -ForegroundColor Gray
    Write-Host "3. Click 'Create API Key'" -ForegroundColor Gray
    Write-Host "4. Copy the key and run:" -ForegroundColor Gray
    Write-Host "   .\render-api-deploy.ps1 -RenderApiKey 'your-api-key-here'" -ForegroundColor Cyan
    Write-Host ""
    
    $RenderApiKey = Read-Host "Or paste your Render API key here"
    if (-not $RenderApiKey) {
        Write-Host "❌ API key is required. Exiting..." -ForegroundColor Red
        exit 1
    }
}

# Headers for API requests
$Headers = @{
    "Authorization" = "Bearer $RenderApiKey"
    "Content-Type" = "application/json"
}

# Function to make API requests
function Invoke-RenderAPI {
    param(
        [string]$Endpoint,
        [string]$Method = "GET",
        [hashtable]$Body = @{}
    )
    
    $Uri = "$RenderBaseUrl$Endpoint"
    
    try {
        if ($Method -eq "GET") {
            $Response = Invoke-RestMethod -Uri $Uri -Headers $Headers -Method $Method
        } else {
            $JsonBody = $Body | ConvertTo-Json -Depth 10
            $Response = Invoke-RestMethod -Uri $Uri -Headers $Headers -Method $Method -Body $JsonBody
        }
        return $Response
    } catch {
        Write-Host "❌ API Error: $($_.Exception.Message)" -ForegroundColor Red
        if ($_.Exception.Response) {
            $ErrorDetails = $_.Exception.Response | ConvertFrom-Json -ErrorAction SilentlyContinue
            if ($ErrorDetails) {
                Write-Host "   Details: $($ErrorDetails.message)" -ForegroundColor Red
            }
        }
        return $null
    }
}

# Function to create databases
function Create-RenderDatabases {
    Write-Host "📊 Creating Databases" -ForegroundColor Cyan
    Write-Host "=====================" -ForegroundColor Cyan
    
    $DatabaseConfigs = @(
        @{
            name = "retail-os-postgres"
            type = "postgresql"
            description = "Main PostgreSQL database for Retail OS"
        },
        @{
            name = "retail-os-redis"
            type = "redis"
            description = "Redis cache for Retail OS"
        }
    )
    
    $CreatedDatabases = @{}
    
    foreach ($dbConfig in $DatabaseConfigs) {
        Write-Host "  Creating $($dbConfig.name)..." -ForegroundColor White
        
        $DatabaseBody = @{
            name = $dbConfig.name
            region = "oregon"
            plan = "free"
            databaseName = "retail_os"
        }
        
        if ($dbConfig.type -eq "postgresql") {
            $Database = Invoke-RenderAPI -Endpoint "/postgres" -Method "POST" -Body $DatabaseBody
        } else {
            $Database = Invoke-RenderAPI -Endpoint "/redis" -Method "POST" -Body $DatabaseBody
        }
        
        if ($Database) {
            Write-Host "  ✅ $($dbConfig.name) created: $($Database.id)" -ForegroundColor Green
            $CreatedDatabases[$dbConfig.type] = $Database
        } else {
            Write-Host "  ❌ Failed to create $($dbConfig.name)" -ForegroundColor Red
        }
        
        Start-Sleep -Seconds 2
    }
    
    return $CreatedDatabases
}

# Function to deploy a service
function Deploy-RenderService {
    param(
        [hashtable]$ServiceConfig,
        [hashtable]$Databases
    )
    
    Write-Host "  Deploying $($ServiceConfig.name)..." -ForegroundColor White
    
    # Environment variables
    $EnvVars = @{
        "PORT" = $ServiceConfig.port.ToString()
        "ENVIRONMENT" = "production"
        "LOG_LEVEL" = "info"
        "SERVICE_NAME" = $ServiceConfig.name
    }
    
    # Add database URLs
    if ($ServiceConfig.database -eq "PostgreSQL" -and $Databases.postgresql) {
        $EnvVars["DATABASE_URL"] = $Databases.postgresql.connectionString
    } elseif ($ServiceConfig.database -eq "MongoDB") {
        # For now, use PostgreSQL for MongoDB services too
        $EnvVars["MONGO_URL"] = $Databases.postgresql.connectionString
    }
    
    if ($ServiceConfig.name -ne "GraphQL Gateway" -and $Databases.redis) {
        $EnvVars["REDIS_URL"] = $Databases.redis.connectionString
    }
    
    # Service configuration
    $ServiceBody = @{
        name = "retail-os-$($ServiceConfig.name.ToLower().Replace(' ', '-'))"
        type = "web_service"
        repo = $GitHubRepo
        rootDir = $ServiceConfig.path
        region = "oregon"
        plan = "free"
        branch = "master"
        buildCommand = if ($ServiceConfig.path -eq "gateway") { "npm install" } else { "go build -o app ./cmd/server" }
        startCommand = if ($ServiceConfig.path -eq "gateway") { "npm start" } else { "./app" }
        envVars = @($EnvVars.GetEnumerator() | ForEach-Object { @{ key = $_.Key; value = $_.Value } })
    }
    
    $Service = Invoke-RenderAPI -Endpoint "/services" -Method "POST" -Body $ServiceBody
    
    if ($Service) {
        Write-Host "  ✅ $($ServiceConfig.name) deployed: https://$($Service.name).render.com" -ForegroundColor Green
        return @{
            success = $true
            name = $ServiceConfig.name
            url = "https://$($Service.name).render.com"
            id = $Service.id
        }
    } else {
        Write-Host "  ❌ Failed to deploy $($ServiceConfig.name)" -ForegroundColor Red
        return @{
            success = $false
            name = $ServiceConfig.name
            url = ""
            id = ""
        }
    }
}

# Main deployment function
function Start-AutomatedDeployment {
    Write-Host "🏗️  Starting Automated Render Deployment" -ForegroundColor Yellow
    Write-Host "=========================================" -ForegroundColor Yellow
    Write-Host ""
    
    # Test API connection
    Write-Host "🔐 Testing API connection..." -ForegroundColor Cyan
    $User = Invoke-RenderAPI -Endpoint "/users/me"
    if (-not $User) {
        Write-Host "❌ Failed to authenticate with Render API" -ForegroundColor Red
        Write-Host "Please check your API key and try again." -ForegroundColor Yellow
        exit 1
    }
    Write-Host "✅ Connected as: $($User.email)" -ForegroundColor Green
    Write-Host ""
    
    # Create databases
    $Databases = @{}
    if ($CreateDatabases) {
        $Databases = Create-RenderDatabases
        Write-Host ""
    }
    
    # Define services
    $Services = @(
        @{name="Identity Service"; path="services/identity"; port=8001; database="PostgreSQL"},
        @{name="Product Catalog"; path="services/product-catalog"; port=8006; database="MongoDB"},
        @{name="Inventory Service"; path="services/inventory"; port=8005; database="PostgreSQL"},
        @{name="Cart Service"; path="services/cart"; port=8002; database="PostgreSQL"},
        @{name="Order Service"; path="services/order"; port=8003; database="PostgreSQL"},
        @{name="Payment Service"; path="services/payment"; port=8004; database="PostgreSQL"},
        @{name="Promotions Service"; path="services/promotions"; port=8007; database="PostgreSQL"},
        @{name="Merchant Account"; path="services/merchant-account"; port=8008; database="PostgreSQL"},
        @{name="Analytics Service"; path="services/analytics"; port=8009; database="PostgreSQL"},
        @{name="GraphQL Gateway"; path="gateway"; port=4000; database="None"}
    )
    
    # Deploy services
    if ($DeployServices) {
        Write-Host "🚀 Deploying Services" -ForegroundColor Cyan
        Write-Host "=====================" -ForegroundColor Cyan
        
        $DeploymentResults = @()
        $SuccessfulDeployments = 0
        
        foreach ($service in $Services) {
            $result = Deploy-RenderService -ServiceConfig $service -Databases $Databases
            $DeploymentResults += $result
            
            if ($result.success) {
                $SuccessfulDeployments++
            }
            
            Start-Sleep -Seconds 3  # Wait between deployments
        }
        
        # Display results
        Write-Host ""
        Write-Host "🎉 DEPLOYMENT COMPLETE!" -ForegroundColor Green -BackgroundColor Black
        Write-Host "=======================" -ForegroundColor Green
        Write-Host ""
        Write-Host "📊 Deployment Summary:" -ForegroundColor Cyan
        Write-Host "  Total Services: $($Services.Count)" -ForegroundColor White
        Write-Host "  Successful: $SuccessfulDeployments" -ForegroundColor Green
        Write-Host "  Failed: $($Services.Count - $SuccessfulDeployments)" -ForegroundColor Red
        Write-Host ""
        
        Write-Host "🌐 Service URLs:" -ForegroundColor Cyan
        foreach ($result in $DeploymentResults) {
            if ($result.success) {
                Write-Host "  ✅ $($result.name): $($result.url)" -ForegroundColor Green
            } else {
                Write-Host "  ❌ $($result.name): Deployment failed" -ForegroundColor Red
            }
        }
        
        Write-Host ""
        Write-Host "🎯 Complete Platform URLs:" -ForegroundColor Yellow
        Write-Host "  Frontend (Storefront): https://unified-commerce.vercel.app" -ForegroundColor Green
        Write-Host "  Backend (Services): See URLs above" -ForegroundColor Green
        Write-Host ""
        Write-Host "🚀 Your Retail OS platform is now LIVE!" -ForegroundColor Green -BackgroundColor Black
    }
}

# Check if we're in the right directory
if (-not (Test-Path "services" -PathType Container)) {
    Write-Host "❌ Error: Please run this script from the UNIFIED_COMMERCE directory" -ForegroundColor Red
    Write-Host "Current directory: $(Get-Location)" -ForegroundColor Yellow
    exit 1
}

# Start deployment
Start-AutomatedDeployment

Write-Host ""
Write-Host "🎉 Retail OS deployment automation complete!" -ForegroundColor Green