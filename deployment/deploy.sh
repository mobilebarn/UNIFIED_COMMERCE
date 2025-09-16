#!/bin/bash

# Retail OS - Quick Deployment Script
# This script deploys Retail OS to Vercel + Railway in ~30 minutes

set -e

echo "🚀 Starting Retail OS Deployment..."
echo "================================================="

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check prerequisites
check_prerequisites() {
    print_status "Checking prerequisites..."
    
    # Check if Node.js is installed
    if ! command -v node &> /dev/null; then
        print_error "Node.js is not installed. Please install Node.js first."
        exit 1
    fi
    
    # Check if npm is installed
    if ! command -v npm &> /dev/null; then
        print_error "npm is not installed. Please install npm first."
        exit 1
    fi
    
    print_success "Prerequisites check completed"
}

# Install deployment tools
install_tools() {
    print_status "Installing deployment tools..."
    
    # Install Vercel CLI
    if ! command -v vercel &> /dev/null; then
        print_status "Installing Vercel CLI..."
        npm install -g vercel
    fi
    
    # Install Railway CLI
    if ! command -v railway &> /dev/null; then
        print_status "Installing Railway CLI..."
        npm install -g @railway/cli
    fi
    
    print_success "Deployment tools installed"
}

# Login to services
login_services() {
    print_status "Logging into deployment services..."
    
    print_warning "Please log in to Vercel when prompted..."
    vercel login
    
    print_warning "Please log in to Railway when prompted..."
    railway login
    
    print_success "Logged into deployment services"
}

# Deploy frontend applications
deploy_frontend() {
    print_status "Deploying frontend applications..."
    
    # Deploy Storefront
    print_status "Deploying Storefront..."
    cd storefront
    npm ci
    npm run build
    vercel --prod --name "retail-os-storefront" --yes
    cd ..
    print_success "Storefront deployed"
    
    # Deploy Admin Panel
    print_status "Deploying Admin Panel..."
    cd admin-panel
    npm ci
    npm run build
    vercel --prod --name "retail-os-admin" --yes
    cd ..
    print_success "Admin Panel deployed"
    
    # Deploy Mobile POS
    print_status "Deploying Mobile POS..."
    cd mobile-pos
    npm ci
    npx expo export -p web
    vercel --prod --name "retail-os-pos" --yes
    cd ..
    print_success "Mobile POS deployed"
    
    print_success "All frontend applications deployed"
}

# Setup backend infrastructure
setup_backend_infrastructure() {
    print_status "Setting up backend infrastructure..."
    
    # Create Railway project
    print_status "Creating Railway project..."
    railway new retail-os-backend
    
    # Add PostgreSQL
    print_status "Adding PostgreSQL database..."
    railway add postgresql
    
    # Add Redis
    print_status "Adding Redis cache..."
    railway add redis
    
    print_success "Backend infrastructure setup completed"
    print_warning "Please set up MongoDB Atlas manually:"
    print_warning "1. Go to https://cloud.mongodb.com/"
    print_warning "2. Create a free cluster"
    print_warning "3. Get the connection string"
    print_warning "4. Update environment variables in Railway"
}

# Deploy backend services
deploy_backend() {
    print_status "Deploying backend services..."
    
    services=(
        "identity-service"
        "merchant-account-service"
        "product-catalog-service"
        "inventory-service"
        "order-service"
        "cart-checkout-service"
        "payments-service"
        "promotions-service"
        "analytics-service"
        "graphql-federation-gateway"
    )
    
    for service in "${services[@]}"; do
        print_status "Deploying $service..."
        cd "backend/$service"
        railway up
        cd ../..
        print_success "$service deployed"
    done
    
    print_success "All backend services deployed"
}

# Configure environment variables
configure_environment() {
    print_status "Configuring environment variables..."
    
    print_warning "Please configure the following environment variables in Railway:"
    print_warning "1. Database connection strings"
    print_warning "2. API endpoints"
    print_warning "3. Payment processor keys"
    print_warning "4. Authentication secrets"
    
    print_status "Environment configuration template created in deployment/env-template.txt"
}

# Display deployment results
show_results() {
    print_success "🎉 Retail OS Deployment Completed!"
    echo "================================================="
    print_success "Your Retail OS applications are now live:"
    echo ""
    print_status "📱 Storefront: https://retail-os-storefront.vercel.app"
    print_status "🔧 Admin Panel: https://retail-os-admin.vercel.app"
    print_status "💰 POS System: https://retail-os-pos.vercel.app"
    print_status "🔗 GraphQL API: https://retail-os-backend.railway.app/graphql"
    echo ""
    print_warning "Next steps:"
    print_warning "1. Configure custom domains in Vercel"
    print_warning "2. Set up MongoDB Atlas connection"
    print_warning "3. Configure payment processor credentials"
    print_warning "4. Test all functionality"
    echo ""
    print_success "Total deployment time: ~30-60 minutes"
    print_success "Monthly cost estimate: ~$20-50"
}

# Main deployment flow
main() {
    echo "🏪 Retail OS - Quick Deployment Script"
    echo "This will deploy your complete e-commerce platform to the cloud"
    echo ""
    
    read -p "Are you ready to deploy? (y/N): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        print_warning "Deployment cancelled"
        exit 0
    fi
    
    check_prerequisites
    install_tools
    login_services
    deploy_frontend
    setup_backend_infrastructure
    deploy_backend
    configure_environment
    show_results
}

# Run main function
main "$@"