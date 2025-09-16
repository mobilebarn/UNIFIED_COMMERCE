# Retail OS - Production Deployment Scripts

## Quick Start Deployment

### Option 1: Vercel + Railway (Recommended for MVP)
**Timeline: 30-60 minutes | Cost: ~$20/month**

#### Prerequisites Setup
```bash
# Install required tools
npm install -g vercel
npm install -g @railway/cli

# Login to services
vercel login
railway login
```

#### 1. Deploy Frontend Applications
```bash
# Deploy Storefront
cd storefront
npm ci
npm run build
vercel --prod --name "retail-os-storefront"

# Deploy Admin Panel
cd ../admin-panel  
npm ci
npm run build
vercel --prod --name "retail-os-admin"

# Deploy Mobile POS (Web Version)
cd ../mobile-pos
npm ci
npx expo export -p web
vercel --prod --name "retail-os-pos"
```

#### 2. Setup Backend Infrastructure
```bash
# Create Railway project
railway new retail-os-backend

# Add databases
railway add postgresql
railway add redis

# MongoDB Atlas (Free Tier)
# Sign up at https://cloud.mongodb.com/
# Create cluster and get connection string
```

#### 3. Deploy Backend Services
```bash
# Deploy each microservice
cd backend/identity-service && railway up
cd ../merchant-account-service && railway up
cd ../product-catalog-service && railway up
cd ../inventory-service && railway up
cd ../order-service && railway up
cd ../cart-checkout-service && railway up
cd ../payments-service && railway up
cd ../promotions-service && railway up
cd ../analytics-service && railway up
cd ../graphql-federation-gateway && railway up
```

### Option 2: Google Cloud Platform (GKE)
**Timeline: 2-4 hours | Cost: ~$100-300/month**

#### Prerequisites
```bash
# Install gcloud CLI
# https://cloud.google.com/sdk/docs/install

# Install kubectl
gcloud components install kubectl

# Install Helm
curl https://get.helm.sh/helm-v3.12.0-linux-amd64.tar.gz | tar xz
sudo mv linux-amd64/helm /usr/local/bin/
```

#### 1. Setup GKE Cluster
```bash
# Create GKE cluster
gcloud container clusters create retail-os-cluster \
  --zone=us-central1-a \
  --num-nodes=3 \
  --enable-autoscaling \
  --max-nodes=10 \
  --min-nodes=1 \
  --machine-type=e2-standard-4

# Get credentials
gcloud container clusters get-credentials retail-os-cluster --zone=us-central1-a
```

#### 2. Deploy Infrastructure
```bash
# Deploy databases and infrastructure
kubectl apply -f kubernetes/infrastructure/

# Deploy services
kubectl apply -f kubernetes/services/

# Deploy ingress
kubectl apply -f kubernetes/ingress/
```

### Option 3: AWS (EKS)
**Timeline: 3-5 hours | Cost: ~$150-400/month**

#### Prerequisites
```bash
# Install AWS CLI
curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip"
unzip awscliv2.zip
sudo ./aws/install

# Install eksctl
curl --silent --location "https://github.com/weaveworks/eksctl/releases/latest/download/eksctl_$(uname -s)_amd64.tar.gz" | tar xz -C /tmp
sudo mv /tmp/eksctl /usr/local/bin
```

#### 1. Create EKS Cluster
```bash
# Create cluster
eksctl create cluster \
  --name retail-os-cluster \
  --region us-west-2 \
  --nodegroup-name linux-nodes \
  --node-type m5.large \
  --nodes 3 \
  --nodes-min 1 \
  --nodes-max 4 \
  --managed
```

#### 2. Deploy Applications
```bash
# Deploy using Helm charts
helm install retail-os ./helm/retail-os
```

## Domain & SSL Configuration

### Custom Domain Setup
```bash
# For Vercel deployments
vercel domains add yourdomain.com
vercel domains add admin.yourdomain.com  
vercel domains add pos.yourdomain.com

# DNS Configuration
# A Record: @ -> Vercel IP
# CNAME: admin -> alias.vercel.app
# CNAME: pos -> alias.vercel.app
```

### SSL Certificates
- **Vercel**: Automatic SSL (Let's Encrypt)
- **GKE/AWS**: cert-manager with Let's Encrypt
- **CloudFlare**: Optional CDN + SSL

## Environment Variables

### Production Environment Configuration
```bash
# Database URLs (from cloud providers)
DATABASE_URL=postgresql://user:pass@host:5432/dbname
MONGODB_URL=mongodb+srv://user:pass@cluster.mongodb.net/dbname
REDIS_URL=redis://user:pass@host:6379

# API Endpoints
GRAPHQL_ENDPOINT=https://api.yourdomain.com/graphql
STOREFRONT_URL=https://yourdomain.com
ADMIN_URL=https://admin.yourdomain.com
POS_URL=https://pos.yourdomain.com

# Payment Processing
STRIPE_PUBLISHABLE_KEY=pk_live_...
STRIPE_SECRET_KEY=sk_live_...

# Authentication
JWT_SECRET=your-production-jwt-secret
AUTH0_DOMAIN=your-auth0-domain.auth0.com
```

## Monitoring & Logging

### Setup Production Monitoring
```bash
# Vercel Analytics (Included)
# Railway Metrics (Included)

# For GKE/AWS - Setup Prometheus + Grafana
helm install prometheus prometheus-community/kube-prometheus-stack
```

## Security Checklist

- [ ] SSL certificates configured
- [ ] Environment variables secured
- [ ] Database connections encrypted
- [ ] API rate limiting enabled
- [ ] CORS policies configured
- [ ] Authentication tokens rotated
- [ ] Firewall rules configured
- [ ] Backup strategy implemented

## Post-Deployment Testing

### Verify All Services
```bash
# Test API endpoints
curl https://api.yourdomain.com/health
curl https://api.yourdomain.com/graphql

# Test frontend applications
curl https://yourdomain.com
curl https://admin.yourdomain.com
curl https://pos.yourdomain.com
```

### Load Testing
```bash
# Install artillery
npm install -g artillery

# Run load tests
artillery quick --count 10 --num 100 https://yourdomain.com
```

## Recommended Deployment Path

**For immediate live deployment (today):**
1. ✅ **Vercel + Railway** - Fastest, most cost-effective
2. 🔄 Later migrate to GKE/AWS for scale

**For enterprise deployment:**
1. 🏢 **GKE or AWS EKS** - Full control, enterprise features
2. 📊 Complete monitoring and logging stack
3. 🔒 Advanced security and compliance features