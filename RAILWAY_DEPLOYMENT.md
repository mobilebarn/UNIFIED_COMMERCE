# Retail OS Railway Deployment Guide

This guide walks you through deploying your Retail OS platform to Railway.

## Prerequisites

1. **Railway Account**: Sign up at [railway.app](https://railway.app)
2. **Railway CLI**: Already installed ✅
3. **Git Repository**: Your code should be in a Git repository

## Authentication

First, authenticate with Railway:

```bash
railway login
```

This will open your browser to complete the authentication.

## Deployment Options

### Option 1: Automated Deployment (Recommended)

Run the automated deployment script:

**Windows PowerShell:**
```powershell
.\deploy-railway.ps1
```

**Linux/Mac:**
```bash
chmod +x deploy-railway.sh
./deploy-railway.sh
```

### Option 2: Manual Deployment

#### Step 1: Create Railway Project

```bash
railway init --name "retail-os-platform"
```

#### Step 2: Add Databases

```bash
railway add postgresql
railway add redis
```

For MongoDB (if needed):
```bash
railway add mongodb
```

#### Step 3: Deploy Backend Services

For each service, navigate to the service directory and deploy:

```bash
# Example for Identity Service
cd services/identity
railway init --name "retail-os-identity"
railway up --detach
```

Repeat for all services:
- retail-os-identity
- retail-os-merchant  
- retail-os-product
- retail-os-inventory
- retail-os-order
- retail-os-payment
- retail-os-cart
- retail-os-promotions
- retail-os-analytics

#### Step 4: Deploy GraphQL Gateway

```bash
cd gateway
railway init --name "retail-os-gateway"
railway up --detach
```

#### Step 5: Deploy Frontend Applications

**Storefront:**
```bash
cd apps/storefront
railway init --name "retail-os-storefront"
railway variables set NEXT_PUBLIC_API_URL="https://retail-os-gateway.railway.app"
railway up --detach
```

**Admin Panel:**
```bash
cd apps/admin
railway init --name "retail-os-admin"
railway variables set REACT_APP_API_URL="https://retail-os-gateway.railway.app"
railway up --detach
```

## Environment Variables

### Required Environment Variables for Each Service

Copy from `.env.railway.template` and set these in Railway dashboard:

1. **Database Variables** (Auto-populated by Railway):
   - `DATABASE_URL`
   - `REDIS_URL`

2. **Service Configuration**:
   - `SERVICE_NAME`
   - `ENVIRONMENT=production`
   - `JWT_SECRET` (Generate a secure secret)

3. **Service Discovery URLs** (Update after deployment):
   - `IDENTITY_SERVICE_URL`
   - `MERCHANT_SERVICE_URL`
   - etc.

### Setting Environment Variables

```bash
# Set environment variables for a service
railway variables set JWT_SECRET="your-secure-secret"
railway variables set ENVIRONMENT="production"
```

## Post-Deployment Configuration

### 1. Update Service URLs

After all services are deployed, update the environment variables with the actual Railway URLs:

```bash
railway variables set IDENTITY_SERVICE_URL="https://retail-os-identity.railway.app"
railway variables set MERCHANT_SERVICE_URL="https://retail-os-merchant.railway.app"
# ... repeat for all services
```

### 2. Database Migrations

Railway will automatically run your database migrations on first deployment.

### 3. Custom Domains (Optional)

Set up custom domains in the Railway dashboard:
1. Go to your service in Railway dashboard
2. Click "Settings" → "Domains"
3. Add your custom domain

## Monitoring and Logs

### View Logs
```bash
railway logs
```

### Check Service Status
```bash
railway status
```

### View Deployments
```bash
railway ps
```

## Troubleshooting

### Common Issues

1. **Build Failures**:
   - Check that `railway.json` is properly configured
   - Verify Go version compatibility
   - Check build commands

2. **Database Connection Issues**:
   - Ensure database environment variables are set
   - Check that services are connecting to Railway databases

3. **Service Communication**:
   - Verify service URLs are correctly set
   - Check that all services are deployed and running

### Getting Help

- Check Railway logs: `railway logs`
- Railway Documentation: [docs.railway.app](https://docs.railway.app)
- Railway Discord: [discord.gg/railway](https://discord.gg/railway)

## Expected URLs

After successful deployment, your services will be available at:

- **Storefront**: `https://retail-os-storefront.railway.app`
- **Admin Panel**: `https://retail-os-admin.railway.app`
- **GraphQL Gateway**: `https://retail-os-gateway.railway.app`
- **Individual Services**: `https://retail-os-[service-name].railway.app`

## Cost Considerations

Railway provides:
- **Starter Plan**: $5/month per service with usage-based pricing
- **Pro Plan**: $20/month per service with more resources

For a full Retail OS deployment (11+ services), budget approximately:
- **Development**: $55-110/month (Starter plan)
- **Production**: $220-440/month (Pro plan)

## Security Recommendations

1. **Environment Variables**: Never commit sensitive data to Git
2. **JWT Secrets**: Use strong, unique secrets for production
3. **Database Passwords**: Use Railway's auto-generated secure passwords
4. **HTTPS**: Railway automatically provides SSL certificates
5. **Access Control**: Configure proper authentication and authorization

## Next Steps

After deployment:
1. Test all services are working
2. Set up monitoring and alerting
3. Configure custom domains
4. Set up CI/CD for automatic deployments
5. Monitor performance and scale as needed