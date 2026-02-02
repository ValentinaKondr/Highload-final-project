#!/bin/bash
set -e

echo "== Applying Kubernetes manifests =="

echo "[1/7] ConfigMaps and Secrets"
kubectl apply -f k8s/configmaps/app-config.yaml
kubectl apply -f k8s/configmaps/redis-secret.yaml

echo "[2/7] Redis"
kubectl apply -f k8s/deployments/redis-deployment.yaml
kubectl apply -f k8s/services/redis-service.yaml

echo "[3/7] Go service"
kubectl apply -f k8s/deployments/metrics-analyzer-deployment.yaml
kubectl apply -f k8s/services/metrics-analyzer-service.yaml

echo "[4/7] HPA"
kubectl apply -f k8s/hpa/metrics-analyzer-hpa.yaml

echo "[5/7] Ingress"
kubectl apply -f k8s/ingress/metrics-analyzer-ingress.yaml || echo "Ingress skipped"

echo "[6/7] Monitoring config"
kubectl apply -f k8s/monitoring/prometheus-config.yaml || echo "Prometheus config skipped"

echo "[7/7] ServiceMonitor"
kubectl apply -f k8s/monitoring/service-monitor.yaml || echo "ServiceMonitor skipped"

echo "== Done =="
echo
kubectl get pods,svc,hpa
