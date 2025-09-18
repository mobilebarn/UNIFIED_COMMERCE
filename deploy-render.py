#!/usr/bin/env python3
"""
🚀 RETAIL OS - RENDER DEPLOYMENT SCRIPT
======================================
Simple and reliable Python script for deploying to Render
"""

import requests
import json
import time
import sys

def main():
    print("🚀 RETAIL OS - RENDER DEPLOYMENT")
    print("=================================")
    print()
    print("📋 This script will help you deploy your Retail OS services to Render")
    print("   But first, let's do it the right way using Render's web interface")
    print()
    
    # Check Render API availability
    print("🔍 Checking Render API status...")
    try:
        response = requests.get("https://api.render.com/v1", timeout=10)
        if response.status_code == 401:
            print("✅ Render API is accessible")
        else:
            print("✅ Render API is responding")
    except:
        print("❌ Cannot reach Render API. Check your internet connection.")
        return
    
    print()
    print("🎯 RECOMMENDED DEPLOYMENT APPROACH")
    print("==================================")
    print()
    print("Since you've already fixed the Go version issue, here's the best way to deploy:")
    print()
    print("1. 🌐 MANUAL RENDER DEPLOYMENT (Most Reliable)")
    print("   - Go to: https://dashboard.render.com")
    print("   - Click 'New +' → 'Web Service'")
    print("   - Connect your GitHub repository: mobilebarn/UNIFIED_COMMERCE")
    print("   - For each service:")
    print("     * Name: retail-os-[service-name]")
    print("     * Root Directory: services/[service-name] (or 'gateway' for GraphQL)")
    print("     * Build Command: go build -o app ./cmd/server")
    print("     * Start Command: ./app")
    print("     * Add environment variables")
    print()
    print("2. 📊 CREATE DATABASES FIRST")
    print("   - PostgreSQL: Click 'New +' → 'PostgreSQL'")
    print("   - Redis: Click 'New +' → 'Redis'")
    print("   - Copy the connection URLs for use in services")
    print()
    
    # Service deployment order
    services = [
        ("Identity Service", "services/identity", "8001"),
        ("Product Catalog", "services/product-catalog", "8006"),
        ("Inventory Service", "services/inventory", "8005"),
        ("Cart Service", "services/cart", "8002"),
        ("Order Service", "services/order", "8003"),
        ("Payment Service", "services/payment", "8004"),
        ("Promotions Service", "services/promotions", "8007"),
        ("Merchant Account", "services/merchant-account", "8008"),
        ("Analytics Service", "services/analytics", "8009"),
        ("GraphQL Gateway", "gateway", "4000")
    ]
    
    print("3. 🚀 DEPLOY SERVICES IN THIS ORDER")
    print("   (Deploy databases first, then services)")
    print()
    
    for i, (name, path, port) in enumerate(services, 1):
        print(f"   {i:2d}. {name}")
        print(f"       Root Directory: {path}")
        print(f"       Port: {port}")
        if path == "gateway":
            print(f"       Build: npm install")
            print(f"       Start: npm start")
        else:
            print(f"       Build: go build -o app ./cmd/server")
            print(f"       Start: ./app")
        print()
    
    print("4. 🔧 ENVIRONMENT VARIABLES FOR EACH SERVICE")
    print("   Add these environment variables to each service:")
    print("   - PORT=[service-port]")
    print("   - ENVIRONMENT=production")
    print("   - LOG_LEVEL=info")
    print("   - DATABASE_URL=[postgres-connection-string]")
    print("   - REDIS_URL=[redis-connection-string]")
    print()
    
    print("📱 ALTERNATIVE: USE RENDER'S BLUEPRINT FEATURE")
    print("==============================================")
    print()
    print("Want an even easier way? I can create a render.yaml blueprint file")
    print("that will deploy everything automatically!")
    print()
    
    create_blueprint = input("Create Render Blueprint file? (y/N): ").strip().lower()
    
    if create_blueprint == 'y':
        create_render_blueprint()
    
    print()
    print("🌐 QUICK LINKS")
    print("==============")
    print("• Render Dashboard: https://dashboard.render.com")
    print("• Your GitHub Repo: https://github.com/mobilebarn/UNIFIED_COMMERCE")
    print("• Frontend (Live): https://unified-commerce.vercel.app")
    print()
    print("✅ Your services will be available at:")
    print("   https://retail-os-[service-name].onrender.com")
    print()
    print("🎉 Once deployed, your complete Retail OS platform will be live!")

def create_render_blueprint():
    """Create a render.yaml blueprint file for automatic deployment"""
    blueprint = """# Render Blueprint for Retail OS
# This file automatically deploys all services when you connect it to Render

services:
  # Databases
  - type: pserv
    name: retail-os-postgres
    env: docker
    plan: free
    dockerfilePath: ./infrastructure/postgres/Dockerfile
    
  - type: redis
    name: retail-os-redis
    plan: free
    
  # Microservices
  - type: web
    name: retail-os-identity
    env: go
    plan: free
    buildCommand: go build -o app ./cmd/server
    startCommand: ./app
    rootDir: services/identity
    envVars:
      - key: PORT
        value: 8001
      - key: ENVIRONMENT
        value: production
      - key: DATABASE_URL
        fromService:
          type: pserv
          name: retail-os-postgres
          property: connectionString
      - key: REDIS_URL
        fromService:
          type: redis
          name: retail-os-redis
          property: connectionString
          
  - type: web
    name: retail-os-product-catalog
    env: go
    plan: free
    buildCommand: go build -o app ./cmd/server
    startCommand: ./app
    rootDir: services/product-catalog
    envVars:
      - key: PORT
        value: 8006
      - key: ENVIRONMENT
        value: production
      - key: MONGO_URL
        fromService:
          type: pserv
          name: retail-os-postgres
          property: connectionString
      - key: REDIS_URL
        fromService:
          type: redis
          name: retail-os-redis
          property: connectionString
          
  - type: web
    name: retail-os-inventory
    env: go
    plan: free
    buildCommand: go build -o app ./cmd/server
    startCommand: ./app
    rootDir: services/inventory
    envVars:
      - key: PORT
        value: 8005
      - key: ENVIRONMENT
        value: production
      - key: DATABASE_URL
        fromService:
          type: pserv
          name: retail-os-postgres
          property: connectionString
      - key: REDIS_URL
        fromService:
          type: redis
          name: retail-os-redis
          property: connectionString
          
  - type: web
    name: retail-os-cart
    env: go
    plan: free
    buildCommand: go build -o app ./cmd/server
    startCommand: ./app
    rootDir: services/cart
    envVars:
      - key: PORT
        value: 8002
      - key: ENVIRONMENT
        value: production
      - key: DATABASE_URL
        fromService:
          type: pserv
          name: retail-os-postgres
          property: connectionString
      - key: REDIS_URL
        fromService:
          type: redis
          name: retail-os-redis
          property: connectionString
          
  - type: web
    name: retail-os-order
    env: go
    plan: free
    buildCommand: go build -o app ./cmd/server
    startCommand: ./app
    rootDir: services/order
    envVars:
      - key: PORT
        value: 8003
      - key: ENVIRONMENT
        value: production
      - key: DATABASE_URL
        fromService:
          type: pserv
          name: retail-os-postgres
          property: connectionString
      - key: REDIS_URL
        fromService:
          type: redis
          name: retail-os-redis
          property: connectionString
          
  - type: web
    name: retail-os-payment
    env: go
    plan: free
    buildCommand: go build -o app ./cmd/server
    startCommand: ./app
    rootDir: services/payment
    envVars:
      - key: PORT
        value: 8004
      - key: ENVIRONMENT
        value: production
      - key: DATABASE_URL
        fromService:
          type: pserv
          name: retail-os-postgres
          property: connectionString
      - key: REDIS_URL
        fromService:
          type: redis
          name: retail-os-redis
          property: connectionString
          
  - type: web
    name: retail-os-promotions
    env: go
    plan: free
    buildCommand: go build -o app ./cmd/server
    startCommand: ./app
    rootDir: services/promotions
    envVars:
      - key: PORT
        value: 8007
      - key: ENVIRONMENT
        value: production
      - key: DATABASE_URL
        fromService:
          type: pserv
          name: retail-os-postgres
          property: connectionString
      - key: REDIS_URL
        fromService:
          type: redis
          name: retail-os-redis
          property: connectionString
          
  - type: web
    name: retail-os-merchant-account
    env: go
    plan: free
    buildCommand: go build -o app ./cmd/server
    startCommand: ./app
    rootDir: services/merchant-account
    envVars:
      - key: PORT
        value: 8008
      - key: ENVIRONMENT
        value: production
      - key: DATABASE_URL
        fromService:
          type: pserv
          name: retail-os-postgres
          property: connectionString
      - key: REDIS_URL
        fromService:
          type: redis
          name: retail-os-redis
          property: connectionString
          
  - type: web
    name: retail-os-analytics
    env: go
    plan: free
    buildCommand: go build -o app ./cmd/server
    startCommand: ./app
    rootDir: services/analytics
    envVars:
      - key: PORT
        value: 8009
      - key: ENVIRONMENT
        value: production
      - key: DATABASE_URL
        fromService:
          type: pserv
          name: retail-os-postgres
          property: connectionString
      - key: REDIS_URL
        fromService:
          type: redis
          name: retail-os-redis
          property: connectionString
          
  - type: web
    name: retail-os-gateway
    env: node
    plan: free
    buildCommand: npm install
    startCommand: npm start
    rootDir: gateway
    envVars:
      - key: PORT
        value: 4000
      - key: NODE_ENV
        value: production
"""

    try:
        with open("render.yaml", "w") as f:
            f.write(blueprint)
        print("✅ Created render.yaml blueprint file!")
        print()
        print("🎯 TO USE THE BLUEPRINT:")
        print("1. Commit and push render.yaml to your GitHub repository")
        print("2. Go to Render Dashboard → 'New +' → 'Blueprint'")
        print("3. Connect your GitHub repository")
        print("4. Render will automatically deploy everything!")
        print()
        print("📁 Blueprint file saved as: render.yaml")
        return True
    except Exception as e:
        print(f"❌ Failed to create blueprint: {e}")
        return False

if __name__ == "__main__":
    main()