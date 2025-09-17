# 🚀 Retail OS Railway CLI Deployment Guide

## Current Status ✅
- ✅ Railway CLI installed and authenticated
- ✅ Project created: `retail-os-platform`
- ✅ Databases provisioned: PostgreSQL, Redis, MongoDB

## Manual CLI Deployment Process

Since Railway CLI requires interactive input, follow these steps for each service:

### 🔄 For Each Service (Repeat 12 times)

#### Step 1: Navigate to Service Directory
```powershell
cd services/identity  # Replace with each service path
```

#### Step 2: Deploy and Create Service
```powershell
railway up
```

**When prompted:**
1. **"Select a service"** → Choose **"Create new service"** (should be an option)
2. **"Service name"** → Enter: `retail-os-identity` (replace with appropriate name)
3. **Wait for deployment** → CLI will build and deploy automatically

#### Step 3: Set Environment Variables
```powershell
# For Identity Service
railway variables set SERVICE_NAME=identity
railway variables set ENVIRONMENT=production
railway variables set JWT_SECRET=prod-secure-jwt-secret-identity-2024
railway variables set JWT_EXPIRATION=86400
railway variables set LOG_LEVEL=info

# Database variables will be auto-linked by Railway
```

#### Step 4: Verify Deployment
```powershell
railway status
railway logs
```

---

## 📋 Complete Service List (Deploy in Order)

### Backend Services (9 services):

1. **Identity Service** (`services/identity/`)
   ```powershell
   cd services/identity
   railway up  # Select "Create new service" → name: retail-os-identity
   railway variables set SERVICE_NAME=identity
   railway variables set ENVIRONMENT=production
   railway variables set JWT_SECRET=prod-secure-jwt-secret-identity-2024
   railway variables set JWT_EXPIRATION=86400
   railway variables set LOG_LEVEL=info
   cd ../..
   ```

2. **Merchant Service** (`services/merchant/`)
   ```powershell
   cd services/merchant
   railway up  # Select "Create new service" → name: retail-os-merchant
   railway variables set SERVICE_NAME=merchant
   railway variables set ENVIRONMENT=production
   railway variables set JWT_SECRET=prod-secure-jwt-secret-merchant-2024
   railway variables set LOG_LEVEL=info
   cd ../..
   ```

3. **Product Service** (`services/product/`)
   ```powershell
   cd services/product
   railway up  # Select "Create new service" → name: retail-os-product
   railway variables set SERVICE_NAME=product
   railway variables set ENVIRONMENT=production
   railway variables set JWT_SECRET=prod-secure-jwt-secret-product-2024
   railway variables set LOG_LEVEL=info
   cd ../..
   ```

4. **Inventory Service** (`services/inventory/`)
   ```powershell
   cd services/inventory
   railway up  # Select "Create new service" → name: retail-os-inventory
   railway variables set SERVICE_NAME=inventory
   railway variables set ENVIRONMENT=production
   railway variables set JWT_SECRET=prod-secure-jwt-secret-inventory-2024
   railway variables set LOG_LEVEL=info
   cd ../..
   ```

5. **Order Service** (`services/order/`)
   ```powershell
   cd services/order
   railway up  # Select "Create new service" → name: retail-os-order
   railway variables set SERVICE_NAME=order
   railway variables set ENVIRONMENT=production
   railway variables set JWT_SECRET=prod-secure-jwt-secret-order-2024
   railway variables set LOG_LEVEL=info
   cd ../..
   ```

6. **Payment Service** (`services/payment/`) - **Important: Port 8005**
   ```powershell
   cd services/payment
   railway up  # Select "Create new service" → name: retail-os-payment
   railway variables set SERVICE_NAME=payment
   railway variables set SERVICE_PORT=8005
   railway variables set ENVIRONMENT=production
   railway variables set JWT_SECRET=prod-secure-jwt-secret-payment-2024
   railway variables set LOG_LEVEL=info
   cd ../..
   ```

7. **Cart Service** (`services/cart/`)
   ```powershell
   cd services/cart
   railway up  # Select "Create new service" → name: retail-os-cart
   railway variables set SERVICE_NAME=cart
   railway variables set ENVIRONMENT=production
   railway variables set JWT_SECRET=prod-secure-jwt-secret-cart-2024
   railway variables set LOG_LEVEL=info
   cd ../..
   ```

8. **Promotions Service** (`services/promotions/`)
   ```powershell
   cd services/promotions
   railway up  # Select "Create new service" → name: retail-os-promotions
   railway variables set SERVICE_NAME=promotions
   railway variables set ENVIRONMENT=production
   railway variables set JWT_SECRET=prod-secure-jwt-secret-promotions-2024
   railway variables set LOG_LEVEL=info
   cd ../..
   ```

9. **Analytics Service** (`services/analytics/`)
   ```powershell
   cd services/analytics
   railway up  # Select "Create new service" → name: retail-os-analytics
   railway variables set SERVICE_NAME=analytics
   railway variables set ENVIRONMENT=production
   railway variables set JWT_SECRET=prod-secure-jwt-secret-analytics-2024
   railway variables set LOG_LEVEL=info
   cd ../..
   ```

### Gateway and Frontend (3 services):

10. **GraphQL Gateway** (`gateway/`)
    ```powershell
    cd gateway
    railway up  # Select "Create new service" → name: retail-os-gateway
    railway variables set NODE_ENV=production
    cd ..
    ```

11. **Storefront** (`apps/storefront/`)
    ```powershell
    cd apps/storefront
    railway up  # Select "Create new service" → name: retail-os-storefront
    railway variables set NODE_ENV=production
    railway variables set NEXT_PUBLIC_API_URL=https://retail-os-gateway.railway.app
    cd ../..
    ```

12. **Admin Panel** (`apps/admin/`)
    ```powershell
    cd apps/admin
    railway up  # Select "Create new service" → name: retail-os-admin
    railway variables set NODE_ENV=production
    railway variables set REACT_APP_API_URL=https://retail-os-gateway.railway.app
    cd ../..
    ```

---

## 🔗 Post-Deployment: Update Service URLs

After all backend services are deployed, update the Gateway with service URLs:

```powershell
cd gateway
railway variables set IDENTITY_SERVICE_URL=https://retail-os-identity.railway.app
railway variables set MERCHANT_SERVICE_URL=https://retail-os-merchant.railway.app
railway variables set PRODUCT_SERVICE_URL=https://retail-os-product.railway.app
railway variables set INVENTORY_SERVICE_URL=https://retail-os-inventory.railway.app
railway variables set ORDER_SERVICE_URL=https://retail-os-order.railway.app
railway variables set PAYMENT_SERVICE_URL=https://retail-os-payment.railway.app
railway variables set CART_SERVICE_URL=https://retail-os-cart.railway.app
railway variables set PROMOTIONS_SERVICE_URL=https://retail-os-promotions.railway.app
railway variables set ANALYTICS_SERVICE_URL=https://retail-os-analytics.railway.app

# Redeploy gateway with new URLs
railway up
cd ..
```

---

## 🎯 Expected Results

After deployment, your Retail OS platform will be available at:

- **🏪 Storefront**: `https://retail-os-storefront.railway.app`
- **⚙️ Admin Panel**: `https://retail-os-admin.railway.app`  
- **🔗 GraphQL Gateway**: `https://retail-os-gateway.railway.app/graphql`
- **📊 Railway Dashboard**: https://railway.com/project/1f76e4aa-b36a-4670-a445-1bfde6bba0a3

## 📝 Deployment Notes

- **Database Connections**: Railway automatically provides database URLs via environment variables
- **Service Discovery**: Update Gateway URLs after all backend services are deployed
- **Port Configuration**: Payment service uses port 8005 (important for avoiding conflicts)
- **Security**: All services use production JWT secrets
- **Monitoring**: Use `railway logs` to monitor deployment progress

## 🚀 Quick Start Command

Start with the first service:
```powershell
cd services/identity
railway up
```

Then follow the interactive prompts to create "retail-os-identity" service.

---

**Total Services**: 12 (9 backend + 1 gateway + 2 frontend)
**Estimated Time**: 30-45 minutes for full deployment
**Cost**: ~$65-85/month (development) | ~$240-300/month (production)