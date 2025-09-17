#!/bin/bash

# Retail OS Railway Deployment Script
# This script deploys the entire Retail OS platform to Railway

echo "🚀 Starting Retail OS Platform Deployment to Railway"
echo "=================================================="

# Check if Railway CLI is authenticated
railway whoami > /dev/null 2>&1
if [ $? -ne 0 ]; then
    echo "❌ You need to login to Railway first"
    echo "Please run: railway login"
    exit 1
fi

echo "✅ Railway CLI authenticated"

# Create main project
echo "📦 Creating Railway project..."
railway init --name "retail-os-platform" || true

# Deploy databases first
echo "🗄️  Setting up databases..."
railway add postgresql || true
railway add redis || true

# Wait for databases to be ready
echo "⏳ Waiting for databases to initialize..."
sleep 30

# Deploy backend services
echo "🔧 Deploying backend microservices..."

services=("identity" "merchant" "product" "inventory" "order" "payment" "cart" "promotions" "analytics")

for service in "${services[@]}"; do
    echo "📤 Deploying $service service..."
    cd "services/$service"
    
    # Create a new Railway service for each microservice
    railway init --name "retail-os-$service" || true
    
    # Link to existing databases
    railway link
    
    # Set environment variables
    railway variables set SERVICE_NAME="$service"
    railway variables set ENVIRONMENT="production"
    
    # Deploy the service
    railway up --detach
    
    cd "../.."
    echo "✅ $service service deployed"
done

# Deploy GraphQL Gateway
echo "🌐 Deploying GraphQL Federation Gateway..."
cd "gateway"
railway init --name "retail-os-gateway" || true
railway link
railway variables set NODE_ENV="production"
railway up --detach
cd ".."

# Deploy frontend applications
echo "🖥️  Deploying frontend applications..."

# Deploy Storefront
cd "apps/storefront"
railway init --name "retail-os-storefront" || true
railway link
railway variables set NODE_ENV="production"
railway variables set NEXT_PUBLIC_API_URL="https://retail-os-gateway.railway.app"
railway up --detach
cd "../.."

# Deploy Admin Panel
cd "apps/admin"
railway init --name "retail-os-admin" || true
railway link
railway variables set NODE_ENV="production"
railway variables set REACT_APP_API_URL="https://retail-os-gateway.railway.app"
railway up --detach
cd "../.."

echo "🎉 Retail OS Platform deployment to Railway completed!"
echo "=================================================="
echo "Your services are being deployed. You can check their status with:"
echo "railway status"
echo ""
echo "🌍 Your applications will be available at:"
echo "- Storefront: https://retail-os-storefront.railway.app"
echo "- Admin Panel: https://retail-os-admin.railway.app"
echo "- GraphQL Gateway: https://retail-os-gateway.railway.app"
echo ""
echo "📊 Monitor your deployments in the Railway dashboard:"
echo "https://railway.app/dashboard"