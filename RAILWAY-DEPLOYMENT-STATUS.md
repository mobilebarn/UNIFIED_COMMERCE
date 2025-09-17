# RETAIL OS RAILWAY DEPLOYMENT GUIDE
# Manual deployment guide for backend services

## AVAILABLE RAILWAY SERVICES (from CLI output):
✅ retail-os-product
✅ retail-os-payment  
✅ retail-os-analytics
✅ retail-os-inventory
✅ retail-os-identity
✅ retail-os-cart
❌ retail-os-admin (frontend - skip, already on Vercel)
❌ retail-os-storefront (frontend - skip, already on Vercel)

## BACKEND SERVICES TO DEPLOY:

### 1. PRODUCT SERVICE
**Directory:** services/product-catalog
**Railway Service:** retail-os-product
**Port:** 8003
**Commands:**
```powershell
cd services\product-catalog
railway up --detach
# Select: retail-os-product
```

### 2. PAYMENT SERVICE  
**Directory:** services/payment
**Railway Service:** retail-os-payment
**Port:** 8005
**Commands:**
```powershell
cd services\payment
railway up --detach
# Select: retail-os-payment
```

### 3. ANALYTICS SERVICE
**Directory:** services/analytics
**Railway Service:** retail-os-analytics
**Port:** 8001
**Commands:**
```powershell
cd services\analytics
railway up --detach
# Select: retail-os-analytics
```

### 4. INVENTORY SERVICE
**Directory:** services/inventory
**Railway Service:** retail-os-inventory
**Port:** 8002
**Commands:**
```powershell
cd services\inventory
railway up --detach
# Select: retail-os-inventory
```

### 5. IDENTITY SERVICE
**Directory:** services/identity
**Railway Service:** retail-os-identity
**Port:** 8000
**Commands:**
```powershell
cd services\identity
railway up --detach
# Select: retail-os-identity
```

### 6. CART SERVICE
**Directory:** services/cart
**Railway Service:** retail-os-cart
**Port:** 8080
**Commands:**
```powershell
cd services\cart
railway up --detach
# Select: retail-os-cart
```

## MISSING SERVICES TO CREATE:
- retail-os-order (Order Service)
- retail-os-merchant (Merchant Service) 
- retail-os-promotions (Promotions Service)
- retail-os-gateway (GraphQL Gateway)

## CURRENT STATUS:
✅ Railway project linked
✅ Database services ready (PostgreSQL, MongoDB, Redis)
⏳ WAITING: Manual service selection in terminal (366)

## NEXT ACTIONS:
1. Complete current Product service deployment (select retail-os-product)
2. Deploy remaining 5 services manually
3. Create missing services for Order, Merchant, Promotions, Gateway
4. Configure environment variables for all services