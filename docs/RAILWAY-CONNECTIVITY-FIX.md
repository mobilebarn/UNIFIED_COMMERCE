# Railway Connectivity Fix Guide

## 🔍 Problem Identified

The GraphQL Federation Gateway is failing because it's trying to connect to services using `localhost` URLs, but in Railway's container environment, each service has its own internal URL.

**Error:** `connect ECONNREFUSED ::1:8001` - Gateway can't reach Identity service at localhost:8001

## ✅ Solution

Update the gateway to use Railway's internal service URLs via environment variables.

## 🚀 Step-by-Step Fix

### Step 1: Get Your Railway Service URLs

1. Go to [Railway Dashboard](https://railway.app/dashboard)
2. Select your project
3. For each service, click on it and note the **Public Domain URL**
4. Your service URLs will look like:
   - `https://identity-service-production.up.railway.app`
   - `https://cart-service-production.up.railway.app`
   - etc.

### Step 2: Configure Gateway Environment Variables

1. In Railway Dashboard, click on your **GraphQL Gateway** service
2. Go to **Variables** tab
3. Add these environment variables (replace URLs with your actual ones):

```
IDENTITY_SERVICE_URL=https://your-identity-service.up.railway.app/graphql
CART_SERVICE_URL=https://your-cart-service.up.railway.app/graphql
ORDER_SERVICE_URL=https://your-order-service.up.railway.app/graphql
PAYMENT_SERVICE_URL=https://your-payment-service.up.railway.app/graphql
INVENTORY_SERVICE_URL=https://your-inventory-service.up.railway.app/graphql
PRODUCT_SERVICE_URL=https://your-product-service.up.railway.app/graphql
PROMOTIONS_SERVICE_URL=https://your-promotions-service.up.railway.app/graphql
MERCHANT_SERVICE_URL=https://your-merchant-service.up.railway.app/graphql
CORS_ORIGINS=https://storefront-eta-six.vercel.app,https://admin-panel-tau-eight.vercel.app
NODE_ENV=production
JWT_SECRET=your-production-jwt-secret
```

### Step 3: Verify All Services Are Running

Before configuring the gateway, ensure all backend services are deployed and showing "**Deployed**" status:

- ✅ Identity Service
- ✅ Cart & Checkout Service  
- ✅ Order Service
- ✅ Payment Service
- ✅ Inventory Service
- ✅ Product Catalog Service
- ✅ Promotions Service
- ✅ Merchant Account Service

### Step 4: Redeploy Gateway

1. After adding all environment variables
2. Click **Deploy** on the gateway service
3. Monitor the deployment logs

### Step 5: Test Connectivity

1. Once gateway is deployed, check its logs
2. You should see successful connections to all services
3. Test the gateway endpoint: `https://your-gateway.up.railway.app/graphql`

## 📋 Environment Variables Checklist

Copy and modify these for your Railway dashboard:

```bash
# Service URLs (replace with your actual Railway URLs)
IDENTITY_SERVICE_URL=https://identity-service-production.up.railway.app/graphql
CART_SERVICE_URL=https://cart-service-production.up.railway.app/graphql
ORDER_SERVICE_URL=https://order-service-production.up.railway.app/graphql
PAYMENT_SERVICE_URL=https://payment-service-production.up.railway.app/graphql
INVENTORY_SERVICE_URL=https://inventory-service-production.up.railway.app/graphql
PRODUCT_SERVICE_URL=https://product-catalog-service-production.up.railway.app/graphql
PROMOTIONS_SERVICE_URL=https://promotions-service-production.up.railway.app/graphql
MERCHANT_SERVICE_URL=https://merchant-account-service-production.up.railway.app/graphql

# Frontend CORS
CORS_ORIGINS=https://storefront-eta-six.vercel.app,https://admin-panel-tau-eight.vercel.app

# Environment
NODE_ENV=production

# Security
JWT_SECRET=your-production-jwt-secret
```

## 🔧 How the Fix Works

1. **Before:** Gateway hardcoded `localhost:8001`, `localhost:8002`, etc.
2. **After:** Gateway uses environment variables that point to actual Railway URLs
3. **Fallback:** If environment variables aren't set, falls back to localhost (for local development)

## 🎯 Expected Results

After applying this fix:

1. **Gateway logs** should show successful connections
2. **No more ECONNREFUSED errors**
3. **All 8 services** should be accessible through the gateway
4. **Frontend apps** can connect to the unified GraphQL endpoint

## 🚨 Common Issues & Solutions

### Issue: Service URLs don't match
**Solution:** Double-check service names in Railway dashboard

### Issue: Services still not connecting  
**Solution:** Ensure all backend services are fully deployed and healthy

### Issue: CORS errors
**Solution:** Verify CORS_ORIGINS includes your Vercel frontend URLs

### Issue: 500 errors
**Solution:** Check that all services have proper database connections

## 📞 Report Progress

After implementing this fix, please report:

1. **How many services** show "Deployed" (green) status
2. **Gateway deployment** status 
3. **Any error messages** you see in gateway logs
4. **Gateway URL** for testing

This should resolve the connectivity issues and get your Retail OS platform fully operational on Railway! 🎉