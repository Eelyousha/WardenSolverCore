# WardenSolverCore - Контейнеризация (Docker & Kubernetes)

## 📋 Структура файлов

```
WardenSolverCore/
├── docker/                          # Docker конфигурация
│   ├── build.sh                    # Скрипт сборки Docker образа
│   ├── up.sh                       # Скрипт запуска docker-compose
│   ├── nginx.conf                  # Конфигурация Nginx (load balancer)
│   ├── master-init.sql             # PostgreSQL Master инициализация
│   ├── slave-init.sh               # PostgreSQL Slave инициализация
│   └── docker-compose.yml          # (в корне) Docker Compose конфигурация
│
├── k8s/                             # Kubernetes конфигурация
│   ├── 00-namespace-config.yaml    # Namespace, ConfigMap, Secret
│   ├── 01-postgres-statefulset.yaml # PostgreSQL StatefulSet (Master + Slave)
│   ├── 02-postgres-service.yaml    # PostgreSQL Services
│   ├── 03-api-deployment.yaml      # API Deployment и Services
│   ├── 04-autoscaling-network.yaml # HPA и NetworkPolicy
│   ├── 05-ingress.yaml             # Ingress для маршрутизации
│   ├── deploy.sh                   # Скрипт развёртывания
│   └── cleanup.sh                  # Скрипт удаления развёртывания
│
├── Dockerfile                       # Multi-stage Docker build
├── docker-compose.yml               # Docker Compose (dev environment)
├── .dockerignore                    # Docker ignore файл
├── .env.example                     # Пример переменных окружения
├── CONTAINERIZATION_PLAN.md         # План контейнеризации
└── DEPLOYMENT_GUIDE.md              # Полное руководство по развёртыванию
```

## 🚀 Быстрый старт

### Локальная разработка (Docker Compose)

```bash
# 1. Перейти в директорию проекта
cd /workspaces/WardenSolverCore

# 2. Собрать Docker образ
./docker/build.sh

# 3. Запустить контейнеры
./docker/up.sh

# 4. Проверить статус
docker-compose ps

# 5. Тестировать API
curl http://localhost:8080/
```

### Развёртывание в Kubernetes

```bash
# 1. Убедиться, что kubectl подключен к кластеру
kubectl cluster-info

# 2. Развернуть приложение
./k8s/deploy.sh

# 3. Проверить статус
kubectl get all -n warden

# 4. Проверить логи
kubectl logs -n warden -l app=warden-api -f
```

## 🏗️ Архитектура

### Docker Compose (локально)

```
┌─────────────────────────────┐
│    Docker Compose Network   │
├─────────────────────────────┤
│                             │
│  ┌─────────────────────┐   │
│  │   Nginx LB          │   │
│  │   :80 → :8080      │   │
│  └─────────────────────┘   │
│           ↓                 │
│  ┌─────────────────────┐   │
│  │   Go API Service    │   │
│  │   :8080             │   │
│  └─────────────────────┘   │
│           ↓                 │
│  ┌──────────┬──────────┐   │
│  │ PostgreSQL Master   │   │
│  │ :5432              │   │
│  └──────────┬──────────┘   │
│             │               │
│  ┌──────────↓──────────┐   │
│  │ PostgreSQL Slave    │   │
│  │ :5433              │   │
│  └─────────────────────┘   │
│                             │
└─────────────────────────────┘
```

### Kubernetes (production)

```
┌────────────────────────────────────────────┐
│   Kubernetes Cluster (warden namespace)    │
├────────────────────────────────────────────┤
│                                            │
│  ┌──────────────────────────────────────┐ │
│  │    Ingress (HTTP/HTTPS routing)      │ │
│  └──────────┬───────────────────────────┘ │
│             │                              │
│  ┌──────────↓──────────────────────────┐ │
│  │  LoadBalancer Service               │ │
│  └──────────┬──────────────────────────┘ │
│             │                              │
│  ┌──────────↓────────────────────────────────────────┐ │
│  │  Deployment (warden-api) - 3+ replicas          │ │
│  │  ┌────────┐ ┌────────┐ ┌────────┐ ┌────────┐   │ │
│  │  │Pod 1   │ │Pod 2   │ │Pod 3   │ │Pod N   │   │ │
│  │  │:8080   │ │:8080   │ │:8080   │ │:8080   │   │ │
│  │  └────────┘ └────────┘ └────────┘ └────────┘   │ │
│  └────────────────────────────────────────────────┘ │
│             ↓                                        │
│  ┌────────────────────────────────────────────────┐ │
│  │  StatefulSet (postgres) - 2 replicas          │ │
│  │  ┌──────────────────┐ ┌──────────────────────┐│ │
│  │  │PostgreSQL Master │ │PostgreSQL Slave (RO)││ │
│  │  │:5432 (R/W)       │→│:5432 (Read only)    ││ │
│  │  └──────────────────┘ └──────────────────────┘│ │
│  │  PersistentVolumes (10GB each)               │ │
│  └────────────────────────────────────────────────┘ │
│                                                     │
│  HPA (Horizontal Pod Autoscaler)                   │
│  - Min replicas: 3                                  │
│  - Max replicas: 10                                 │
│  - CPU threshold: 70%                               │
│  - Memory threshold: 80%                            │
│                                                     │
└────────────────────────────────────────────────────┘
```

## 🔄 Репликация PostgreSQL

### Master → Slave синхронизация

```
Master (Read/Write)
│
├─ WAL (Write-Ahead Logs)
│  ├─ Отправка логов на Slave
│  └─ Проверка получения
│
Slave (Read-only)
│
├─ Восстановление данных из WAL
├─ Постоянная синхронизация
└─ Готов для чтения (SELECT запросы)
```

### Использование в приложении

- **Writes**: Всегда идут на Master
- **Reads**: Могут распределяться между Master и Slave
- **Failover**: При отказе Master, Slave может быть промоувлен

## 📊 Масштабирование

### Горизонтальное (API)

Автоматическое через HPA:
- Минимум: 3 pod'а
- Максимум: 10 pod'ов
- Срабатывает при: CPU > 70% или Memory > 80%

### Вертикальное (ресурсы)

Отредактировать в `k8s/03-api-deployment.yaml`:

```yaml
resources:
  requests:
    cpu: 100m → 200m
    memory: 128Mi → 256Mi
  limits:
    cpu: 500m → 1000m
    memory: 512Mi → 2Gi
```

## 🔒 Безопасность

### Реализовано

- ✅ Secret для учётных данных БД
- ✅ NetworkPolicy для ограничения трафика
- ✅ PodDisruptionBudget для высокой доступности
- ✅ Health checks (liveness & readiness probes)
- ✅ Resource limits (CPU & Memory)

### Рекомендуется добавить

- 🔲 RBAC (Role-Based Access Control)
- 🔲 Istio Service Mesh для mTLS
- 🔲 Pod Security Policy
- 🔲 Secrets encryption at rest
- 🔲 Regular vulnerability scanning

## 📈 Production Checklist

```
Pre-deployment:
  [ ] Образ протестирован локально (docker-compose)
  [ ] Образ загружен в приватный registry
  [ ] K8s кластер имеет >= 16GB RAM
  [ ] PersistentVolume провайдер настроен
  [ ] Ingress controller установлен

Deployment:
  [ ] Namespace создан
  [ ] Secrets и ConfigMaps настроены
  [ ] PostgreSQL готов (Master + Slave синхронизированы)
  [ ] API pods здоровы и готовы
  [ ] Ingress маршруты работают

Post-deployment:
  [ ] API отвечает на запросы
  [ ] Database репликация синхронизирована
  [ ] HPA работает корректно
  [ ] Логирование и мониторинг активны
  [ ] Backup'ы настроены
```

## 🛠️ Полезные команды

### Docker

```bash
# Собрать образ
./docker/build.sh

# Запустить контейнеры
./docker/up.sh

# Просмотр логов
docker-compose logs -f api

# Остановить контейнеры
docker-compose down

# Остановить и удалить данные
docker-compose down -v
```

### Kubernetes

```bash
# Развернуть
./k8s/deploy.sh

# Просмотр статуса
kubectl get all -n warden

# Просмотр логов
kubectl logs -n warden -l app=warden-api -f

# Port forwarding
kubectl port-forward -n warden svc/warden-api-service 8080:80

# Удалить развёртывание
./k8s/cleanup.sh
```

## 📚 Документация

Подробную информацию см. в:
- `CONTAINERIZATION_PLAN.md` - План контейнеризации
- `DEPLOYMENT_GUIDE.md` - Полное руководство (100+ строк)

## ⚠️ Known Issues & Limitations

1. **PostgreSQL Failover**: Требует ручного вмешательства. Для автоматического используйте pg_auto_failover или Patroni.

2. **Data Persistence**: При `docker-compose down -v` все данные теряются. Используйте backup'ы.

3. **Image Registry**: Укажите правильный registry перед развёртыванием в K8s.

4. **Ingress**: Требует настройки DNS для production.

## 🤝 Поддержка

По вопросам развёртывания и конфигурации см. `DEPLOYMENT_GUIDE.md`.

## 📝 Лицензия

Та же, что и основной проект.
