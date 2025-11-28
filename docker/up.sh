#!/bin/bash

# Скрипт для запуска docker-compose

set -e

echo "=========================================="
echo "Docker Compose Startup Script"
echo "=========================================="
echo ""

# Проверка docker-compose
if ! command -v docker-compose &> /dev/null; then
    echo "❌ docker-compose не найден"
    exit 1
fi

# Проверка docker
if ! command -v docker &> /dev/null; then
    echo "❌ docker не найден"
    exit 1
fi

# Проверка файла
if [ ! -f "docker-compose.yml" ]; then
    echo "❌ docker-compose.yml не найден в текущей директории"
    exit 1
fi

echo "Запускаю docker-compose..."
echo ""

# Запуск сервисов
docker-compose up -d

echo ""
echo "✓ Сервисы успешно запущены!"
echo ""

# Вывод статуса
echo "Статус контейнеров:"
docker-compose ps

echo ""
echo "Доступные сервисы:"
echo "  API: http://localhost:8080"
echo "  Nginx Load Balancer: http://localhost"
echo "  PostgreSQL Master: localhost:5432"
echo "  PostgreSQL Slave: localhost:5433"
echo ""
echo "Для просмотра логов:"
echo "  docker-compose logs -f api"
echo ""
echo "Для остановки:"
echo "  docker-compose down"
echo ""
