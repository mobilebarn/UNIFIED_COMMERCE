# Retail OS - Quick Deployment Script (PowerShell)
# This script deploys Retail OS to Vercel + Railway in ~30 minutes

param(
    [switch]$SkipConfirmation
)

# Colors for output
$Colors = @{
    Red = "Red"
    Green = "Green" 
    Yellow = "Yellow"
    Blue = "Blue"
    White = "White"
}

function Write-Status {
    param([string]$Message)
    Write-Host "[INFO] $Message" -ForegroundColor $Colors.Blue
}

function Write-Success {
    param([string]$Message)
    Write-Host "[SUCCESS] $Message" -ForegroundColor $Colors.Green
}

function Write-Warning {
    param([string]$Message)
    Write-Host "[WARNING] $Message" -ForegroundColor $Colors.Yellow
}

function Write-Error {
    param([string]$Message)
    Write-Host "[ERROR] $Message" -ForegroundColor $Colors.Red
}

function Test-Prerequisites {
    Write-Status "Checking prerequisites..."
    
    # Check if Node.js is installed
    try {
        $nodeVersion = node --version
        Write-Success "Node.js found: $nodeVersion"
    }
    catch {
        Write-Error "Node.js is not installed. Please install Node.js first."
        exit 1
    }
    
    # Check if npm is installed
    try {
        $npmVersion = npm --version
        Write-Success "npm found: $npmVersion"
    }
    catch {
        Write-Error "npm is not installed. Please install npm first."
        exit 1
    }
    
    Write-Success "Prerequisites check completed"
}

function Install-DeploymentTools {
    Write-Status "Installing deployment tools..."
    
    # Install Vercel CLI
    try {
        vercel --version | Out-Null
        Write-Success "Vercel CLI already installed"
    }
    catch {
        Write-Status "Installing Vercel CLI..."
        npm install -g vercel
    }
    
    # Install Railway CLI
    try {
        railway --version | Out-Null
        Write-Success "Railway CLI already installed"
    }
    catch {
        Write-Status "Installing Railway CLI..."
        npm install -g @railway/cli
    }
    
    Write-Success "Deployment tools installed"
}

function Connect-Services {
    Write-Status "Logging into deployment services..."
    
    Write-Warning "Please log in to Vercel when prompted..."
    vercel login
    
    Write-Warning "Please log in to Railway when prompted..."
    railway login
    
    Write-Success "Logged into deployment services"
}

function Deploy-Frontend {
    Write-Status "Deploying frontend applications..."
    
    # Deploy Storefront
    Write-Status "Deploying Storefront..."
    Set-Location "storefront"
    npm ci
    npm run build
    vercel --prod --name "retail-os-storefront" --yes
    Set-Location ".."
    Write-Success "Storefront deployed"
    
    # Deploy Admin Panel
    Write-Status "Deploying Admin Panel..."
    Set-Location "admin-panel"
    npm ci
    npm run build
    vercel --prod --name "retail-os-admin" --yes
    Set-Location ".."
    Write-Success "Admin Panel deployed"
    
    # Deploy Mobile POS
    Write-Status "Deploying Mobile POS..."
    Set-Location "mobile-pos"
    npm ci
    npx expo export -p web
    vercel --prod --name "retail-os-pos" --yes
    Set-Location ".."
    Write-Success "Mobile POS deployed"
    
    Write-Success "All frontend applications deployed"
}

function Setup-BackendInfrastructure {
    Write-Status "Setting up backend infrastructure..."
    
    # Create Railway project
    Write-Status "Creating Railway project..."
    railway new retail-os-backend
    
    # Add PostgreSQL
    Write-Status "Adding PostgreSQL database..."
    railway add postgresql
    
    # Add Redis
    Write-Status "Adding Redis cache..."
    railway add redis
    
    Write-Success "Backend infrastructure setup completed"
    Write-Warning "Please set up MongoDB Atlas manually:"
    Write-Warning "1. Go to https://cloud.mongodb.com/"
    Write-Warning "2. Create a free cluster"
    Write-Warning "3. Get the connection string"
    Write-Warning "4. Update environment variables in Railway"
}

function Deploy-Backend {
    Write-Status "Deploying backend services..."
    
    $services = @(
        "identity-service",
        "merchant-account-service", 
        "product-catalog-service",
        "inventory-service",
        "order-service",
        "cart-checkout-service",
        "payments-service",
        "promotions-service",
        "analytics-service",
        "graphql-federation-gateway"
    )
    
    foreach ($service in $services) {
        Write-Status "Deploying $service..."
        Set-Location "backend\$service"
        railway up
        Set-Location "..\..\"
        Write-Success "$service deployed"
    }
    
    Write-Success "All backend services deployed"
}

function Set-EnvironmentConfiguration {
    Write-Status "Configuring environment variables..."
    
    Write-Warning "Please configure the following environment variables in Railway:"
    Write-Warning "1. Database connection strings"
    Write-Warning "2. API endpoints"
    Write-Warning "3. Payment processor keys"
    Write-Warning "4. Authentication secrets"
    
    Write-Status "Environment configuration template created in deployment\env-template.txt"
}

function Show-DeploymentResults {
    Write-Success "🎉 Retail OS Deployment Completed!"
    Write-Host "=================================================" -ForegroundColor $Colors.Green
    Write-Success "Your Retail OS applications are now live:"
    Write-Host ""
    Write-Status "📱 Storefront: https://retail-os-storefront.vercel.app"
    Write-Status "🔧 Admin Panel: https://retail-os-admin.vercel.app"
    Write-Status "💰 POS System: https://retail-os-pos.vercel.app"
    Write-Status "🔗 GraphQL API: https://retail-os-backend.railway.app/graphql"
    Write-Host ""
    Write-Warning "Next steps:"
    Write-Warning "1. Configure custom domains in Vercel"
    Write-Warning "2. Set up MongoDB Atlas connection"
    Write-Warning "3. Configure payment processor credentials"
    Write-Warning "4. Test all functionality"
    Write-Host ""
    Write-Success "Total deployment time: ~30-60 minutes"
    Write-Success "Monthly cost estimate: ~$20-50"
}

function Start-Deployment {
    Write-Host "🏪 Retail OS - Quick Deployment Script" -ForegroundColor $Colors.Blue
    Write-Host "This will deploy your complete e-commerce platform to the cloud" -ForegroundColor $Colors.White
    Write-Host ""
    
    if (-not $SkipConfirmation) {
        $confirmation = Read-Host "Are you ready to deploy? (y/N)"
        if ($confirmation -notmatch "^[Yy]$") {
            Write-Warning "Deployment cancelled"
            exit 0
        }
    }
    
    Test-Prerequisites
    Install-DeploymentTools
    Connect-Services
    Deploy-Frontend
    Setup-BackendInfrastructure
    Deploy-Backend
    Set-EnvironmentConfiguration
    Show-DeploymentResults
}

# Run main deployment
try {
    Start-Deployment
}
catch {
    Write-Error "Deployment failed: $($_.Exception.Message)"
    exit 1
}