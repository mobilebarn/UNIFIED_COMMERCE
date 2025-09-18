#!/usr/bin/env python3
"""
🚀 RETAIL OS - AUTOMATED RENDER DEPLOYMENT
==========================================
Python script for automated deployment to Render platform
Uses Render API to deploy all services without manual intervention
"""

import requests
import json
import time
import sys
from typing import Dict, List, Optional

class RenderDeployer:
    def __init__(self, api_key: str):
        self.api_key = api_key
        self.base_url = "https://api.render.com/v1"
        self.headers = {
            "Authorization": f"Bearer {api_key}",
            "Content-Type": "application/json"
        }
        self.github_repo = "https://github.com/mobilebarn/UNIFIED_COMMERCE"
        
    def make_request(self, endpoint: str, method: str = "GET", data: Dict = None) -> Optional[Dict]:
        """Make API request to Render"""
        url = f"{self.base_url}{endpoint}"
        
        try:
            if method == "GET":
                response = requests.get(url, headers=self.headers)
            elif method == "POST":
                response = requests.post(url, headers=self.headers, json=data)
            elif method == "PUT":
                response = requests.put(url, headers=self.headers, json=data)
            
            if response.status_code in [200, 201]:
                return response.json()
            else:
                print(f"❌ API Error: {response.status_code} - {response.text}")
                return None
                
        except Exception as e:
            print(f"❌ Request Error: {str(e)}")
            return None
    
    def test_connection(self) -> bool:
        """Test API connection"""
        print("🔐 Testing API connection...")
        user = self.make_request("/users/me")
        if user:
            print(f"✅ Connected as: {user.get('email', 'Unknown')}")
            return True
        else:
            print("❌ Failed to authenticate with Render API")
            return False
    
    def create_databases(self) -> Dict[str, Dict]:
        """Create PostgreSQL and Redis databases"""
        print("\n📊 Creating Databases")
        print("=====================")
        
        databases = {}
        
        # Create PostgreSQL database
        print("  Creating PostgreSQL database...")
        postgres_config = {
            "name": "retail-os-postgres",
            "region": "oregon",
            "plan": "free",
            "databaseName": "retail_os"
        }
        
        postgres_db = self.make_request("/postgres", "POST", postgres_config)
        if postgres_db:
            print(f"  ✅ PostgreSQL created: {postgres_db.get('id')}")
            databases['postgresql'] = postgres_db
        else:
            print("  ❌ Failed to create PostgreSQL database")
        
        time.sleep(2)
        
        # Create Redis database
        print("  Creating Redis database...")
        redis_config = {
            "name": "retail-os-redis",
            "region": "oregon",
            "plan": "free"
        }
        
        redis_db = self.make_request("/redis", "POST", redis_config)
        if redis_db:
            print(f"  ✅ Redis created: {redis_db.get('id')}")
            databases['redis'] = redis_db
        else:
            print("  ❌ Failed to create Redis database")
        
        return databases
    
    def deploy_service(self, service_config: Dict, databases: Dict) -> Dict:
        """Deploy a single service"""
        print(f"  Deploying {service_config['name']}...")
        
        # Environment variables
        env_vars = {
            "PORT": str(service_config['port']),
            "ENVIRONMENT": "production",
            "LOG_LEVEL": "info",
            "SERVICE_NAME": service_config['name']
        }
        
        # Add database URLs
        if service_config['database'] == "PostgreSQL" and 'postgresql' in databases:
            env_vars["DATABASE_URL"] = databases['postgresql'].get('connectionString', '')
        elif service_config['database'] == "MongoDB":
            # Use PostgreSQL for MongoDB services temporarily
            env_vars["MONGO_URL"] = databases.get('postgresql', {}).get('connectionString', '')
        
        if service_config['name'] != "GraphQL Gateway" and 'redis' in databases:
            env_vars["REDIS_URL"] = databases['redis'].get('connectionString', '')
        
        # Service configuration
        service_name = f"retail-os-{service_config['name'].lower().replace(' ', '-')}"
        
        service_data = {
            "name": service_name,
            "type": "web_service",
            "repo": self.github_repo,
            "rootDir": service_config['path'],
            "region": "oregon",
            "plan": "free",
            "branch": "master",
            "buildCommand": "npm install" if service_config['path'] == "gateway" else "go build -o app ./cmd/server",
            "startCommand": "npm start" if service_config['path'] == "gateway" else "./app",
            "envVars": [{"key": k, "value": v} for k, v in env_vars.items()]
        }
        
        service = self.make_request("/services", "POST", service_data)
        
        if service:
            url = f"https://{service_name}.render.com"
            print(f"  ✅ {service_config['name']} deployed: {url}")
            return {
                "success": True,
                "name": service_config['name'],
                "url": url,
                "id": service.get('id')
            }
        else:
            print(f"  ❌ Failed to deploy {service_config['name']}")
            return {
                "success": False,
                "name": service_config['name'],
                "url": "",
                "id": ""
            }
    
    def deploy_all_services(self) -> List[Dict]:
        """Deploy all Retail OS services"""
        services = [
            {"name": "Identity Service", "path": "services/identity", "port": 8001, "database": "PostgreSQL"},
            {"name": "Product Catalog", "path": "services/product-catalog", "port": 8006, "database": "MongoDB"},
            {"name": "Inventory Service", "path": "services/inventory", "port": 8005, "database": "PostgreSQL"},
            {"name": "Cart Service", "path": "services/cart", "port": 8002, "database": "PostgreSQL"},
            {"name": "Order Service", "path": "services/order", "port": 8003, "database": "PostgreSQL"},
            {"name": "Payment Service", "path": "services/payment", "port": 8004, "database": "PostgreSQL"},
            {"name": "Promotions Service", "path": "services/promotions", "port": 8007, "database": "PostgreSQL"},
            {"name": "Merchant Account", "path": "services/merchant-account", "port": 8008, "database": "PostgreSQL"},
            {"name": "Analytics Service", "path": "services/analytics", "port": 8009, "database": "PostgreSQL"},
            {"name": "GraphQL Gateway", "path": "gateway", "port": 4000, "database": "None"}
        ]
        
        print("\n🚀 Deploying Services")
        print("=====================")
        
        # Create databases first
        databases = self.create_databases()
        
        # Deploy services
        results = []
        successful = 0
        
        for service in services:
            result = self.deploy_service(service, databases)
            results.append(result)
            
            if result['success']:
                successful += 1
            
            time.sleep(3)  # Wait between deployments
        
        return results, successful, len(services)
    
    def show_deployment_summary(self, results: List[Dict], successful: int, total: int):
        """Show deployment summary"""
        print("\n🎉 DEPLOYMENT COMPLETE!")
        print("=======================")
        print(f"\n📊 Deployment Summary:")
        print(f"  Total Services: {total}")
        print(f"  Successful: {successful}")
        print(f"  Failed: {total - successful}")
        print(f"\n🌐 Service URLs:")
        
        for result in results:
            if result['success']:
                print(f"  ✅ {result['name']}: {result['url']}")
            else:
                print(f"  ❌ {result['name']}: Deployment failed")
        
        print(f"\n🎯 Complete Platform URLs:")
        print(f"  Frontend (Storefront): https://unified-commerce.vercel.app")
        print(f"  Backend (Services): See URLs above")
        print(f"\n🚀 Your Retail OS platform is now LIVE!")


def main():
    print("🚀 RETAIL OS - AUTOMATED RENDER DEPLOYMENT")
    print("===========================================")
    print()
    print("📋 Before we start, please note:")
    print("  - You need a Render account (free tier available)")
    print("  - Get your API key from: https://dashboard.render.com/account")
    print("  - Your GitHub repository must be public or connected to Render")
    print()
    
    # Get API key
    api_key = input("Enter your Render API key (starts with 'rnd_'): ").strip()
    
    if not api_key:
        print("❌ API key is required. Get it from: https://dashboard.render.com/account")
        print("   1. Go to Render Dashboard")
        print("   2. Click 'Account' → 'API Keys'")
        print("   3. Click 'Create API Key'")
        print("   4. Copy the key and run this script again")
        sys.exit(1)
    
    if not api_key.startswith('rnd_'):
        print("⚠️  Warning: Render API keys usually start with 'rnd_'")
        proceed = input("Continue anyway? (y/N): ").strip().lower()
        if proceed != 'y':
            sys.exit(1)
    
    # Initialize deployer
    deployer = RenderDeployer(api_key)
    
    # Test connection
    if not deployer.test_connection():
        print("❌ Cannot proceed without valid API connection")
        sys.exit(1)
    
    # Start deployment
    print("\n🏗️  Starting Automated Deployment")
    print("=================================")
    
    try:
        results, successful, total = deployer.deploy_all_services()
        deployer.show_deployment_summary(results, successful, total)
        
        if successful == total:
            print("\n🎉 All services deployed successfully!")
        else:
            print(f"\n⚠️  {total - successful} services failed to deploy. Check Render dashboard for details.")
            
    except KeyboardInterrupt:
        print("\n❌ Deployment cancelled by user")
        sys.exit(1)
    except Exception as e:
        print(f"\n❌ Deployment failed: {str(e)}")
        sys.exit(1)


if __name__ == "__main__":
    main()