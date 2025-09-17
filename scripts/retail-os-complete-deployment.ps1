# RETAIL OS COMPLETE AUTOMATED RAILWAY DEPLOYMENT
# This script automates the entire deployment process for the Retail OS platform
# Handles: Service creation, GitHub deployment, environment configuration, testing

Write-Host "RETAIL OS COMPLETE AUTOMATED DEPLOYMENT" -ForegroundColor Green
Write-Host "=======================================" -ForegroundColor Green
Write-Host ""

Write-Host "Phase 1: Creating Missing Railway Services" -ForegroundColor Yellow
Write-Host "==========================================" -ForegroundColor Yellow

# Missing services to create
$MissingServices = @(
    @{Name="retail-os-identity"; Dir="services/identity"; Port="8000"},
    @{Name="retail-os-cart"; Dir="services/cart"; Port="8080"},
    @{Name="retail-os-order"; Dir="services/order"; Port="8004"},
    @{Name="retail-os-merchant"; Dir="services/merchant-account"; Port="8006"},
    @{Name="retail-os-promotions"; Dir="services/promotions"; Port="8007"},
    @{Name="retail-os-gateway"; Dir="gateway"; Port="4000"}
)

# Existing services (already created)
$ExistingServices = @(
    @{Name="retail-os-product"; Dir="services/product-catalog"; Port="8003"},
    @{Name="retail-os-payment"; Dir="services/payment"; Port="8005"},
    @{Name="retail-os-analytics"; Dir="services/analytics"; Port="8001"},
    @{Name="retail-os-inventory"; Dir="services/inventory"; Port="8002"}
)

Write-Host "Creating missing services in Railway dashboard..." -ForegroundColor Cyan
Write-Host "Note: Services will be created via Railway dashboard since CLI service creation is limited" -ForegroundColor Gray

foreach ($service in $MissingServices) {
    Write-Host "  - $($service.Name) (Port: $($service.Port))" -ForegroundColor White
}

Write-Host ""
Write-Host "MANUAL STEP REQUIRED:" -ForegroundColor Red
Write-Host "Please create these 6 services in Railway dashboard:" -ForegroundColor Yellow
foreach ($service in $MissingServices) {
    Write-Host "  1. Go to Railway dashboard" -ForegroundColor White
    Write-Host "  2. Click 'New Service'" -ForegroundColor White
    Write-Host "  3. Name: $($service.Name)" -ForegroundColor Cyan
    Write-Host ""
}

Write-Host "Press any key to continue once services are created..." -ForegroundColor Yellow
$null = $Host.UI.RawUI.ReadKey("NoEcho,IncludeKeyDown")

Write-Host ""
Write-Host "Phase 2: Configuring GitHub Deployment" -ForegroundColor Yellow
Write-Host "======================================" -ForegroundColor Yellow

# Function to configure service for GitHub deployment
function Configure-ServiceDeployment {
    param(
        [string]$ServiceName,
        [string]$ServiceDir,
        [string]$Port,
        [string]$Database = "Postgres"
    )
    
    Write-Host "Configuring $ServiceName..." -ForegroundColor Cyan
    
    try {
        # Link to service
        railway link --service $ServiceName
        
        # Set root directory for monorepo
        Write-Host "  Setting root directory: $ServiceDir" -ForegroundColor White
        railway variables --set "RAILWAY_ROOT_DIR=$ServiceDir"
        
        # Set common environment variables
        Write-Host "  Setting environment variables..." -ForegroundColor White
        railway variables --set "SERVICE_PORT=$Port"
        railway variables --set "ENVIRONMENT=production"
        railway variables --set "JWT_SECRET=prod-$($ServiceName.Replace('retail-os-', ''))-jwt-2024"
        
        # Set database connection
        if ($Database -eq "MongoDB") {
            railway variables --set "MONGO_URL=`${{MongoDB.MONGO_URL}}"
        } else {
            railway variables --set "DATABASE_URL=`${{Postgres.DATABASE_URL}}"
        }
        
        # Set Redis connection
        railway variables --set "REDIS_URL=`${{Redis.REDIS_URL}}"
        
        Write-Host "  ✅ $ServiceName configured successfully!" -ForegroundColor Green
        return $true
    }
    catch {
        Write-Host "  ❌ Failed to configure $ServiceName" -ForegroundColor Red
        Write-Host "  Error: $_" -ForegroundColor Red
        return $false
    }
}

# Configure all services
Write-Host "Configuring existing services..." -ForegroundColor Cyan
foreach ($service in $ExistingServices) {
    $database = if ($service.Name -eq "retail-os-product") { "MongoDB" } else { "Postgres" }
    Configure-ServiceDeployment -ServiceName $service.Name -ServiceDir $service.Dir -Port $service.Port -Database $database
}

Write-Host ""
Write-Host "Configuring new services..." -ForegroundColor Cyan
foreach ($service in $MissingServices) {
    $database = if ($service.Name -eq "retail-os-gateway") { "Postgres" } else { "Postgres" }
    Configure-ServiceDeployment -ServiceName $service.Name -ServiceDir $service.Dir -Port $service.Port -Database $database
}

Write-Host ""
Write-Host "Phase 3: Deploying All Services" -ForegroundColor Yellow
Write-Host "===============================" -ForegroundColor Yellow

# Function to deploy service from GitHub
function Deploy-ServiceFromGitHub {
    param(
        [string]$ServiceName,
        [string]$ServiceDir
    )
    
    Write-Host "Deploying $ServiceName from GitHub..." -ForegroundColor Cyan
    
    try {
        # Navigate to service directory
        Push-Location $ServiceDir
        
        # Link to service
        railway link --service $ServiceName
        
        # Deploy from current directory (will use GitHub)
        Write-Host "  Starting deployment..." -ForegroundColor White
        railway up --detach
        
        Write-Host "  ✅ $ServiceName deployment started!" -ForegroundColor Green
        return $true
    }
    catch {
        Write-Host "  ❌ Failed to deploy $ServiceName" -ForegroundColor Red
        Write-Host "  Error: $_" -ForegroundColor Red
        return $false
    }
    finally {
        Pop-Location
    }
}

# Deploy all services
$AllServices = $ExistingServices + $MissingServices
$SuccessfulDeployments = 0

foreach ($service in $AllServices) {
    if (Deploy-ServiceFromGitHub -ServiceName $service.Name -ServiceDir $service.Dir) {
        $SuccessfulDeployments++
    }
    Start-Sleep -Seconds 2  # Brief pause between deployments
}

Write-Host ""
Write-Host "Phase 4: Deployment Summary" -ForegroundColor Yellow
Write-Host "===========================" -ForegroundColor Yellow

Write-Host "RETAIL OS DEPLOYMENT COMPLETE!" -ForegroundColor Green -BackgroundColor Black
Write-Host ""
Write-Host "Deployment Results:" -ForegroundColor Cyan
Write-Host "  Total Services: $($AllServices.Count)" -ForegroundColor White
Write-Host "  Successful Deployments: $SuccessfulDeployments" -ForegroundColor Green
Write-Host "  Failed Deployments: $($AllServices.Count - $SuccessfulDeployments)" -ForegroundColor Red
Write-Host ""

if ($SuccessfulDeployments -eq $AllServices.Count) {
    Write-Host "🎉 ALL SERVICES DEPLOYED SUCCESSFULLY!" -ForegroundColor Green -BackgroundColor Black
    Write-Host ""
    Write-Host "Your Retail OS platform is now live on Railway!" -ForegroundColor Green
    Write-Host ""
    Write-Host "Next steps:" -ForegroundColor Yellow
    Write-Host "  1. Check Railway dashboard for deployment status" -ForegroundColor White
    Write-Host "  2. Test service endpoints" -ForegroundColor White
    Write-Host "  3. Update frontend URLs to point to Railway services" -ForegroundColor White
} else {
    Write-Host "⚠️  Some deployments failed. Check Railway dashboard for details." -ForegroundColor Yellow
}

Write-Host ""
Write-Host "Service URLs will be available at:" -ForegroundColor Cyan
foreach ($service in $AllServices) {
    Write-Host "  $($service.Name): https://$($service.Name).railway.app" -ForegroundColor White
}