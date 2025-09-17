# Kubernetes Deployment Guide for Retail OS

## Overview
This document provides instructions for deploying Retail OS on Kubernetes. The platform consists of 8 microservices, a GraphQL Federation Gateway, and two frontend applications (storefront and admin panel).

## Prerequisites
- Kubernetes cluster (v1.21+)
- kubectl CLI
- Helm v3+ (optional)
- Ingress controller (NGINX recommended)
- At least 8GB RAM and 20GB storage available

## Architecture
The platform is composed of the following Kubernetes resources:

### Core Components
- Namespace: `retail-os`
- ConfigMap: `retail-os-config`
- Secret: `retail-os-secrets`
- PersistentVolumeClaims for databases

### Infrastructure Services
- PostgreSQL (Stateful)
- MongoDB (Stateful)
- Redis (Stateful)
- Kafka
- Zookeeper

### Application Services
- Identity Service
- Merchant Account Service
- Product Catalog Service
- Inventory Service
- Order Service
- Cart Service
- Payment Service
- Promotions Service
- GraphQL Federation Gateway
- Storefront Application
- Admin Panel Application

### Networking
- Services for internal communication
- Ingress for external access

## Deployment Steps

### 1. Create Namespace
```bash
kubectl apply -f k8s/manifests/namespace.yaml
```

### 2. Deploy Configuration
```bash
kubectl apply -f k8s/manifests/configmap.yaml
kubectl apply -f k8s/manifests/secrets.yaml
```

### 3. Deploy Persistent Storage
```bash
kubectl apply -f k8s/manifests/pvcs.yaml
```

### 4. Deploy Infrastructure Services
```bash
kubectl apply -f k8s/manifests/infrastructure.yaml
```

### 5. Deploy Application Services
```bash
kubectl apply -f k8s/manifests/microservices.yaml
kubectl apply -f k8s/manifests/gateway-and-frontend.yaml
```

### 6. Deploy Ingress
```bash
kubectl apply -f k8s/manifests/ingress.yaml
```

### 7. Verify Deployment
```bash
# Check all pods
kubectl get pods -n retail-os

# Check services
kubectl get services -n retail-os

# Check ingress
kubectl get ingress -n retail-os
```

## Configuration Details

### Environment Variables
The ConfigMap contains non-sensitive configuration:
- Database names and users
- Service ports
- Environment settings

### Secrets
The Secret contains sensitive data:
- Database passwords
- JWT secret
- Payment gateway credentials

### Persistent Storage
PersistentVolumeClaims are created for:
- PostgreSQL (5Gi)
- MongoDB (5Gi)
- Redis (1Gi)
- Elasticsearch (10Gi)
- Prometheus (5Gi)
- Grafana (1Gi)

## Resource Management

### CPU and Memory Requests/Limits
All services have defined resource requests and limits:
- Microservices: 100m CPU / 128Mi RAM (request), 200m CPU / 256Mi RAM (limit)
- Frontend apps: 200m CPU / 256Mi RAM (request), 500m CPU / 512Mi RAM (limit)
- Gateway: 100m CPU / 128Mi RAM (request), 200m CPU / 256Mi RAM (limit)

## Health Checks

### Liveness Probes
All services include liveness probes to detect unresponsive applications:
- HTTP GET requests to /health endpoints
- Database connection checks
- Service-specific health endpoints

### Readiness Probes
All services include readiness probes to control traffic flow:
- HTTP GET requests to /health endpoints
- Database readiness checks
- Service-specific readiness endpoints

## Scaling

### Horizontal Pod Autoscaler
To enable auto-scaling, create HPA resources:
```yaml
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: identity-service-hpa
  namespace: retail-os
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: identity-service
  minReplicas: 1
  maxReplicas: 10
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 70
```

### Manual Scaling
```bash
# Scale a specific deployment
kubectl scale deployment identity-service --replicas=3 -n retail-os
```

## Networking

### Service Discovery
Services communicate through Kubernetes DNS:
- `identity-service.retail-os.svc.cluster.local`
- `graphql-gateway.retail-os.svc.cluster.local`
- etc.

### Ingress Configuration
Ingress resources provide external access:
- storefront.local → Storefront application
- admin.local → Admin panel application
- api.local → GraphQL Gateway

## Monitoring and Logging

### Built-in Monitoring
- Prometheus metrics endpoints
- Health check endpoints
- Kubernetes resource monitoring

### Logging
- Standard output logging
- Structured JSON logs
- Kubernetes log aggregation compatibility

## Security

### Network Policies
Implement network policies for enhanced security:
```yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: default-deny
  namespace: retail-os
spec:
  podSelector: {}
  policyTypes:
  - Ingress
```

### RBAC
Create role-based access control:
```yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  namespace: retail-os
  name: service-reader
rules:
- apiGroups: [""]
  resources: ["pods", "services"]
  verbs: ["get", "list"]
```

## Backup and Recovery

### Database Backups
Implement regular database backups:
- PostgreSQL logical backups
- MongoDB dumps
- Redis snapshots

### Configuration Backups
Version control all Kubernetes manifests:
- Git repository
- CI/CD integration
- Regular backups

## Troubleshooting

### Common Issues
1. **Pods stuck in Pending state**: Check resource quotas and node capacity
2. **Pods crashing**: Check logs and resource limits
3. **Services not accessible**: Check service selectors and endpoints
4. **Ingress not working**: Check ingress controller and DNS configuration

### Debugging Commands
```bash
# View pod logs
kubectl logs <pod-name> -n retail-os

# Describe pod for detailed status
kubectl describe pod <pod-name> -n retail-os

# Exec into pod
kubectl exec -it <pod-name> -n retail-os -- sh

# Port forward for local testing
kubectl port-forward service/graphql-gateway 4000:4000 -n retail-os
```

## Production Considerations

### High Availability
- Deploy multiple replicas for critical services
- Use Kubernetes node affinity and anti-affinity
- Implement proper load balancing

### Security Hardening
- Use Kubernetes secrets for sensitive data
- Implement network policies
- Enable RBAC
- Use pod security policies

### Performance Optimization
- Configure appropriate resource requests/limits
- Implement horizontal pod autoscaling
- Use readiness/liveness probes
- Optimize database queries

### Disaster Recovery
- Regular backups of databases and configurations
- Cross-cluster replication
- Automated restore procedures
- Regular disaster recovery testing

## Cleanup
To remove all resources:
```bash
kubectl delete namespace retail-os
```

To remove specific components:
```bash
kubectl delete -f k8s/manifests/<specific-file>.yaml
```