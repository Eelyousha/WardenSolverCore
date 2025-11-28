#!/bin/bash

# Скрипт для сборки Docker образа

set -e

REGISTRY="${DOCKER_REGISTRY:-docker.io}"
IMAGE_NAME="${IMAGE_NAME:-warden-api}"
IMAGE_TAG="${IMAGE_TAG:-latest}"

echo "=========================================="
echo "Docker Build Script"
echo "=========================================="
echo ""
echo "Registry: $REGISTRY"
echo "Image: $IMAGE_NAME:$IMAGE_TAG"
echo ""

# Проверка Dockerfile
if [ ! -f "Dockerfile" ]; then
    echo "❌ Dockerfile не найден в текущей директории"
    exit 1
fi

# Сборка образа
echo "Собираю Docker образ..."
docker build \
    --tag $REGISTRY/$IMAGE_NAME:$IMAGE_TAG \
    --tag $REGISTRY/$IMAGE_NAME:latest \
    --file Dockerfile \
    .

echo ""
echo "✓ Docker образ успешно собран!"
echo ""
echo "Образы:"
docker images | grep $IMAGE_NAME
echo ""
echo "Для отправки в реестр:"
echo "  docker push $REGISTRY/$IMAGE_NAME:$IMAGE_TAG"
echo ""
