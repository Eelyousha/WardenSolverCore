# Контейнеризация WardenSolverCore

## Архитектура

```
┌─────────────────────────────────────────────────────────────────┐
│                         Kubernetes Cluster                       │
├─────────────────────────────────────────────────────────────────┤
│                                                                   │
│  ┌──────────────────────────────────────────────────────────┐   │
│  │                  API Service (Deployment)               │   │
│  │  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐   │   │
│  │  │  Replica 1   │  │  Replica 2   │  │  Replica N   │   │   │
│  │  │ Go API Pod   │  │ Go API Pod   │  │ Go API Pod   │   │   │
│  │  └──────────────┘  └──────────────┘  └──────────────┘   │   │
│  └──────────────────────────────────────────────────────────┘   │
│                              ↓                                    │
│  ┌──────────────────────────────────────────────────────────┐   │
│  │              PostgreSQL StatefulSet                       │   │
│  │  ┌──────────────────┐      ┌──────────────────┐          │   │
│  │  │   PostgreSQL     │      │   PostgreSQL     │          │   │
│  │  │    MASTER        │ ───→ │    SLAVE (RO)    │          │   │
│  │  │  (Read/Write)    │      │  (Read Only)     │          │   │
│  │  └──────────────────┘      └──────────────────┘          │   │
│  │        Port: 5432                 Port: 5432            │   │
│  └──────────────────────────────────────────────────────────┘   │
│                                                                   │
└─────────────────────────────────────────────────────────────────┘
```

## Этапы контейнеризации

### 1. Docker образы
- `Dockerfile` для Go приложения (multi-stage build)
- Dockerfile для PostgreSQL Master
- Dockerfile для PostgreSQL Slave

### 2. Docker Compose (для локальной разработки)
- docker-compose.yml с:
  - PostgreSQL Master
  - PostgreSQL Slave
  - Go API Service
  - Nginx (load balancer)

### 3. Kubernetes конфигурация
- ConfigMap для переменных окружения
- Secret для учётных данных БД
- Deployment для API сервиса (масштабируемый)
- StatefulSet для PostgreSQL (Master + Slave)
- Service для API (LoadBalancer)
- Service для БД (внутренний)
- PersistentVolumeClaim для хранения БД

### 4. Репликация PostgreSQL
- Master: принимает всё (чтение + запись)
- Slave: только чтение, синхронизируется с Master
- WAL (Write-Ahead Logging) для репликации

### 5. Масштабирование
- Несколько реплик API (Deployment)
- Балансировка нагрузки через Service
- Автомасштабирование (HPA)

## Требуемые ресурсы

```
API Pod:
  CPU: 100m - 500m
  Memory: 128Mi - 512Mi

PostgreSQL Master:
  CPU: 200m - 1000m
  Memory: 256Mi - 2Gi
  
PostgreSQL Slave:
  CPU: 100m - 500m
  Memory: 256Mi - 1Gi
```

## Процесс развёртывания

1. Создать Docker образы
2. Запустить локально через docker-compose
3. Развернуть в Kubernetes кластер
4. Настроить репликацию
5. Протестировать масштабирование и отказоустойчивость
