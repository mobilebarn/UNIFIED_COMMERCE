# RETAIL OS DIRECT DEPLOYMENT
Write-Host "DEPLOYING IDENTITY SERVICE NOW..." -ForegroundColor Green

# Navigate and deploy Identity Service
cd services/identity

# Set environment variables first
railway variables set SERVICE_NAME=identity
railway variables set ENVIRONMENT=production  
railway variables set JWT_SECRET=prod-identity-secret-2024

# Deploy directly
railway up --detach

Write-Host "IDENTITY SERVICE DEPLOYED!" -ForegroundColor Green

# Return to root
cd ../..