#!/bin/bash

set -e

echo "Deploying to Kubernetes..."

# Apply secrets first
kubectl apply -f k8s/configmaps/redis-secret.yaml

# Apply configmaps
kubectl apply -f k8s/configmaps/app-config.yaml

# Deploy Redis
kubectl apply -f k8s/deployments/redis-deployment.yaml
kubectl apply -f k8s/services/redis-service.yaml

# Wait for Redis to be ready
echo "Waiting for Redis to be ready..."
kubectl wait --for=condition=available --timeout=300s deployment/redis

# Deploy Go service
kubectl apply -f k8s/deployments/go-service-deployment.yaml
kubectl apply -f k8s/services/go-service-service.yaml

# Wait for Go service to be ready
echo "Waiting for Go service to be ready..."
kubectl wait --for=condition=available --timeout=300s deployment/go-service

# Apply HPA
kubectl apply -f k8s/hpa/go-service-hpa.yaml

# Apply Ingress (if NGINX Ingress Controller is installed)
kubectl apply -f k8s/ingress/go-service-ingress.yaml || echo "Ingress not applied (NGINX Ingress Controller may not be installed)"

echo "Deployment completed!"
echo "Check status with: kubectl get pods,svc,hpa"

