# 🚀 Railway Final Deployment Guide - Retail OS

## 🎯 MISSION: Complete Railway Deployment in One Go

This guide will help you finalize your Retail OS deployment on Railway with exact values and steps.

## 📋 PHASE 1: Check Your Current Railway Services

1. Go to [Railway Dashboard](https://railway.app/dashboard)
2. Open your project
3. **Count your services** - you should see these services:

### ✅ Required Services (8 Backend + 1 Gateway = 9 Total)
- [ ] Identity Service
- [ ] Cart Service  
- [ ] Order Service
- [ ] Payment Service
- [ ] Inventory Service
- [ ] Product Catalog Service
- [ ] Promotions Service
- [ ] Merchant Account Service
- [ ] GraphQL Gateway

## 📝 PHASE 2: Get Your Exact Service URLs

**For each service above, click on it and copy the "Public Domain" URL.**

Your URLs will follow this pattern:
- `https://[service-name]-production-[random].up.railway.app`

## 🔧 PHASE 3: Configure Gateway Environment Variables

**Go to your GraphQL Gateway service → Variables tab**

**Copy and paste these EXACT environment variables:**

```bash
# === BACKEND SERVICE URLS ===
# Replace the domain parts with YOUR actual Railway URLs

IDENTITY_SERVICE_URL=https://identity-service-production-xxxx.up.railway.app/graphql
CART_SERVICE_URL=https://cart-service-production-xxxx.up.railway.app/graphql
ORDER_SERVICE_URL=https://order-service-production-xxxx.up.railway.app/graphql
PAYMENT_SERVICE_URL=https://payment-service-production-xxxx.up.railway.app/graphql
INVENTORY_SERVICE_URL=https://inventory-service-production-xxxx.up.railway.app/graphql
PRODUCT_SERVICE_URL=https://product-catalog-service-production-xxxx.up.railway.app/graphql
PROMOTIONS_SERVICE_URL=https://promotions-service-production-xxxx.up.railway.app/graphql
MERCHANT_SERVICE_URL=https://merchant-account-service-production-xxxx.up.railway.app/graphql

# === FRONTEND CORS ORIGINS ===
CORS_ORIGINS=https://storefront-eta-six.vercel.app,https://admin-panel-tau-eight.vercel.app

# === ENVIRONMENT ===
NODE_ENV=production

# === SECURITY ===
JWT_SECRET=retail-os-production-jwt-secret-2024
```

## 🎯 PHASE 4: Smart URL Replacement Strategy

**Instead of manually copying each URL, use this pattern:**

1. **Find your project's base URL pattern** in Railway
2. **Most likely patterns:**
   - `https://[service-name]-production.up.railway.app`
   - `https://[service-name]-production-[hash].up.railway.app`

3. **Quick replacement method:**
   - Copy ONE service URL from Railway
   - Replace `[service-name]` part with each service name
   - Add `/graphql` to the end

### Example:
If your Identity Service URL is:
`https://identity-service-production-abc123.up.railway.app`

Then your other URLs would be:
```
CART_SERVICE_URL=https://cart-service-production-abc123.up.railway.app/graphql
ORDER_SERVICE_URL=https://order-service-production-abc123.up.railway.app/graphql
```

## 📲 PHASE 5: One-Click Environment Variable Setup

**Copy this template and fill in YOUR Railway project's URL pattern:**

```bash
# === REPLACE 'YOUR-PROJECT-HASH' WITH YOUR ACTUAL HASH ===

IDENTITY_SERVICE_URL=https://identity-service-production-YOUR-PROJECT-HASH.up.railway.app/graphql
CART_SERVICE_URL=https://cart-service-production-YOUR-PROJECT-HASH.up.railway.app/graphql
ORDER_SERVICE_URL=https://order-service-production-YOUR-PROJECT-HASH.up.railway.app/graphql
PAYMENT_SERVICE_URL=https://payment-service-production-YOUR-PROJECT-HASH.up.railway.app/graphql
INVENTORY_SERVICE_URL=https://inventory-service-production-YOUR-PROJECT-HASH.up.railway.app/graphql
PRODUCT_SERVICE_URL=https://product-catalog-service-production-YOUR-PROJECT-HASH.up.railway.app/graphql
PROMOTIONS_SERVICE_URL=https://promotions-service-production-YOUR-PROJECT-HASH.up.railway.app/graphql
MERCHANT_SERVICE_URL=https://merchant-account-service-production-YOUR-PROJECT-HASH.up.railway.app/graphql
CORS_ORIGINS=https://storefront-eta-six.vercel.app,https://admin-panel-tau-eight.vercel.app
NODE_ENV=production
JWT_SECRET=retail-os-production-jwt-secret-2024
```

## 🚀 PHASE 6: Deploy and Test

1. **Add all environment variables** to GraphQL Gateway service
2. **Click Deploy** on the gateway service
3. **Wait for deployment** to complete
4. **Check logs** - should see successful connections
5. **Test the gateway** at your gateway URL + `/graphql`

## 📊 PHASE 7: Verification Checklist

After deployment, verify:

- [ ] All 8 backend services show "Deployed" status
- [ ] Gateway shows "Deployed" status  
- [ ] Gateway logs show no ECONNREFUSED errors
- [ ] Gateway GraphQL playground loads successfully
- [ ] Frontend apps can connect to gateway

## 🎯 FINAL RESULT

**Your Retail OS platform will be fully operational with:**

✅ **Backend:** 8 microservices on Railway  
✅ **Frontend:** Storefront + Admin Panel on Vercel  
✅ **Gateway:** Unified GraphQL API  
✅ **Databases:** PostgreSQL, MongoDB, Redis on Railway  

## 🆘 Quick Troubleshooting

**If gateway still fails:**
1. Double-check service URLs are correct
2. Ensure all backend services are "Deployed" 
3. Verify environment variable names are exact
4. Check that `/graphql` is added to service URLs

**Need the exact URLs?** 
Tell me one of your Railway service URLs and I'll generate all the others for you!

---

## 🎉 Success Criteria

Your deployment is complete when:
- Gateway logs show: "✅ All 8 services connected"
- Frontend apps load without errors
- You can query the GraphQL playground
- No more ECONNREFUSED errors

**Let's get your Retail OS fully deployed! 🚀**