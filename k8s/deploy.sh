#!/bin/bash

# Скрипт для развёртывания в Kubernetes

set -e

NAMESPACE="warden"
REGISTRY="${DOCKER_REGISTRY:-docker.io}"
IMAGE_NAME="${IMAGE_NAME:-warden-api}"
IMAGE_TAG="${IMAGE_TAG:-latest}"

echo "=========================================="
echo "Kubernetes Deployment Script"
echo "=========================================="
echo ""

# 1. Проверка prerequisites
echo "[1/5] Проверка prerequisites..."
if ! command -v kubectl &> /dev/null; then
    echo "❌ kubectl не найден. Пожалуйста, установите kubectl"
    exit 1
fi

CLUSTER=$(kubectl cluster-info 2>&1 | head -n 1)
echo "✓ Kubernetes кластер: $CLUSTER"

# 2. Создание namespace и конфиги
echo "[2/5] Создание namespace и конфигов..."
kubectl apply -f k8s/00-namespace-config.yaml
echo "✓ Namespace и ConfigMaps созданы"

# 3. Развёртывание PostgreSQL
echo "[3/5] Развёртывание PostgreSQL (Master + Slave)..."
kubectl apply -f k8s/01-postgres-statefulset.yaml
kubectl apply -f k8s/02-postgres-service.yaml
echo "✓ PostgreSQL развёрнут"
echo "  Ожидаю готовности PostgreSQL pods (это может занять несколько минут)..."
kubectl wait --for=condition=ready pod -l app=postgres -n $NAMESPACE --timeout=300s

# 4. Развёртывание API
echo "[4/5] Развёртывание API сервиса..."
kubectl apply -f k8s/03-api-deployment.yaml
kubectl apply -f k8s/04-autoscaling-network.yaml
echo "✓ API сервис развёрнут"
echo "  Ожидаю готовности API pods..."
kubectl wait --for=condition=available --timeout=120s \
  deployment/warden-api -n $NAMESPACE || true

# 5. Развёртывание Ingress
echo "[5/5] Развёртывание Ingress..."
kubectl apply -f k8s/05-ingress.yaml
echo "✓ Ingress развёрнут"

# Вывод статуса
echo ""
echo "=========================================="
echo "✓ Развёртывание завершено успешно!"
echo "=========================================="
echo ""
echo "Статус сервиса:"
echo ""
kubectl get svc -n $NAMESPACE
echo ""
echo "Статус pods:"
echo ""
kubectl get pods -n $NAMESPACE
echo ""
echo "Для проверки статуса используйте:"
echo "  kubectl get all -n $NAMESPACE"
echo ""
echo "Для просмотра логов API:"
echo "  kubectl logs -n $NAMESPACE -l app=warden-api -f"
echo ""
echo "Для доступа к API:"
echo "  kubectl port-forward -n $NAMESPACE svc/warden-api-service 8080:80"
echo "  Затем: curl http://localhost:8080/"
echo ""
