# Retail OS - Vercel Deployment Guide

## 🚀 Quick Deployment to Vercel (30 minutes to live)

### Prerequisites
- Vercel account (free tier available)
- GitHub repository
- Railway account for backend services

### Step 1: Frontend Deployment (Vercel)

#### Deploy Storefront
```bash
# In storefront directory
cd storefront
npm run build
vercel --prod
```

#### Deploy Admin Panel  
```bash
# In admin-panel directory
cd admin-panel
npm run build
vercel --prod
```

#### Deploy Mobile POS
```bash
# In mobile-pos directory
cd mobile-pos
npx expo export -p web
vercel --prod
```

### Step 2: Backend Deployment (Railway)

#### Deploy All Microservices
```bash
# Each service can be deployed individually
railway up
```

Services to deploy:
- identity-service
- merchant-account-service  
- product-catalog-service
- inventory-service
- order-service
- cart-checkout-service
- payments-service
- promotions-service
- analytics-service
- graphql-federation-gateway

### Step 3: Database Setup

#### PostgreSQL (Railway)
```bash
railway add postgresql
```

#### MongoDB (MongoDB Atlas)
```bash
# Free tier available
# Connect string: mongodb+srv://username:password@cluster.mongodb.net/
```

#### Redis (Railway)
```bash
railway add redis
```

### Step 4: Environment Configuration

Create `.env.production` files for each service with live database URLs.

### Step 5: DNS & SSL

Vercel automatically provides:
- SSL certificates
- CDN optimization
- Custom domains

### Expected URLs after deployment:
- **Storefront**: `https://retail-os-storefront.vercel.app`
- **Admin Panel**: `https://retail-os-admin.vercel.app`
- **Mobile POS**: `https://retail-os-pos.vercel.app`
- **GraphQL API**: `https://retail-os-api.railway.app/graphql`

### Cost Estimate:
- **Vercel**: Free tier (3 apps)
- **Railway**: ~$20/month for backend services
- **MongoDB Atlas**: Free tier (512MB)
- **Total**: ~$20/month for production-ready deployment