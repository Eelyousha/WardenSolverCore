#!/bin/bash

# Скрипт для удаления развёртывания из Kubernetes

set -e

NAMESPACE="warden"

echo "=========================================="
echo "Удаление развёртывания из Kubernetes"
echo "=========================================="
echo ""

# Запрос подтверждения
read -p "Вы уверены, что хотите удалить namespace '$NAMESPACE' и все ресурсы? (y/n) " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "Отмена"
    exit 1
fi

echo "Удаляю развёртывание..."

# Удаление Ingress
echo "[1/3] Удаление Ingress..."
kubectl delete -f k8s/05-ingress.yaml --ignore-not-found=true

# Удаление API
echo "[2/3] Удаление API сервиса..."
kubectl delete -f k8s/03-api-deployment.yaml --ignore-not-found=true
kubectl delete -f k8s/04-autoscaling-network.yaml --ignore-not-found=true

# Удаление PostgreSQL и namespace
echo "[3/3] Удаление PostgreSQL и namespace..."
kubectl delete namespace $NAMESPACE --ignore-not-found=true

echo ""
echo "✓ Развёртывание успешно удалено!"
echo ""
