# Контейнеризация WardenSolverCore - Полный процесс

## 📋 Что было создано

### Docker

1. **`Dockerfile`** - Multi-stage build для минимизации размера
   - Builder stage: собирает Go приложение
   - Runtime stage: Alpine Linux (~50MB)
   - Health check включен

2. **`docker-compose.yml`** - Локальная разработка
   - PostgreSQL Master + Slave с репликацией
   - Go API Service
   - Nginx Load Balancer
   - Автоматическое управление сетью и томами

3. **`docker/`** - Вспомогательные скрипты и конфигурации
   - `build.sh` - Сборка Docker образа
   - `up.sh` - Запуск docker-compose
   - `nginx.conf` - Load balancer конфигурация
   - `master-init.sql` - PostgreSQL Master инициализация с репликацией
   - `slave-init.sh` - PostgreSQL Slave инициализация (base backup + recovery)

4. **`.dockerignore`** - Оптимизация сборки

### Kubernetes

1. **`k8s/00-namespace-config.yaml`**
   - Namespace `warden`
   - ConfigMap с переменными окружения
   - Secret с учётными данными БД

2. **`k8s/01-postgres-statefulset.yaml`**
   - StatefulSet с 2 контейнерами (Master + Slave)
   - Каждый с PersistentVolume (10GB)
   - Health checks и resource limits
   - Автоматическая репликация WAL

3. **`k8s/02-postgres-service.yaml`**
   - Service для Master
   - Service для Slave (читаемый)
   - Headless service для StatefulSet

4. **`k8s/03-api-deployment.yaml`**
   - Deployment с 3+ replicas по умолчанию
   - LoadBalancer Service для external access
   - ClusterIP Service для internal access
   - Health probes (liveness & readiness)
   - Resource requests/limits

5. **`k8s/04-autoscaling-network.yaml`**
   - HPA (Horizontal Pod Autoscaler): 3-10 replicas
   - NetworkPolicy для ограничения трафика
   - PodDisruptionBudget для high availability

6. **`k8s/05-ingress.yaml`**
   - Ingress для HTTP/HTTPS маршрутизации
   - Let's Encrypt автоматизация (опционально)
   - Rate limiting

7. **`k8s/deploy.sh`** - Скрипт развёртывания
   - Проверка prerequisites
   - Создание namespace и конфигов
   - Запуск PostgreSQL
   - Запуск API с ожиданием готовности
   - Вывод статуса

8. **`k8s/cleanup.sh`** - Скрипт очистки
   - Безопасное удаление всех ресурсов
   - Подтверждение перед удалением

### Документация

1. **`CONTAINERIZATION_PLAN.md`** - План контейнеризации
2. **`DEPLOYMENT_GUIDE.md`** - Подробное руководство (200+ строк)
3. **`CONTAINER_README.md`** - Quick start и архитектура
4. **`.env.example`** - Пример переменных окружения

## 🚀 Процесс развёртывания

### Этап 1: Локальная разработка (Docker Compose)

```bash
cd /workspaces/WardenSolverCore

# Шаг 1: Собрать образ
./docker/build.sh
# Результат: docker.io/warden-api:latest (локально)

# Шаг 2: Запустить контейнеры
./docker/up.sh
# Результат: 4 контейнера (postgres-master, postgres-slave, api, nginx)

# Шаг 3: Проверить статус
docker-compose ps

# Шаг 4: Тестировать
curl http://localhost:8080/
curl http://localhost/  # через Nginx

# Шаг 5: Проверить репликацию
# docker-compose exec postgres-master psql -U postgres -c "SELECT * FROM pg_stat_replication;"
```

### Этап 2: Развёртывание в Kubernetes

```bash
# Шаг 1: Подготовить образ
./docker/build.sh
# Собрать и отправить в registry (опционально для локального тестирования)
docker tag docker.io/warden-api:latest your-registry/warden-api:v1.0
docker push your-registry/warden-api:v1.0

# Шаг 2: Обновить K8s конфиг (если не docker.io)
# Отредактировать k8s/03-api-deployment.yaml:
# image: your-registry/warden-api:v1.0

# Шаг 3: Развернуть
./k8s/deploy.sh
# Результат: Полный K8s stack в namespace 'warden'

# Шаг 4: Проверить статус
kubectl get all -n warden
kubectl get pods -n warden -o wide

# Шаг 5: Проверить логи
kubectl logs -n warden -l app=warden-api -f

# Шаг 6: Port forward для тестирования
kubectl port-forward -n warden svc/warden-api-service 8080:80
curl http://localhost:8080/
```

## 🏗️ Архитектура

### Docker Compose (DEV)

```
┌──────────────────────────────────────┐
│   Docker Compose Network             │
├──────────────────────────────────────┤
│                                      │
│  Nginx (Port 80)                    │
│    └─> API (Port 8080)              │
│          └─> PostgreSQL Master      │
│              └─> PostgreSQL Slave   │
│                  (WAL Replication)  │
│                                      │
└──────────────────────────────────────┘
```

### Kubernetes (PROD)

```
Requests
  │
  ↓
Ingress (HTTP/HTTPS + Rate Limit)
  │
  ↓
LoadBalancer Service
  │
  ↓
Deployment (3+ pods, auto-scaled 3-10)
  ├─ Pod 1 (100-500m CPU, 128-512Mi RAM)
  ├─ Pod 2
  └─ Pod 3
  │
  ↓
  ├─→ PostgreSQL Master (Write)
  │   └─ PersistentVolume (10GB)
  │
  └─→ PostgreSQL Slave (Read)
      └─ PersistentVolume (10GB)
      └─ WAL Replication
```

## ⚙️ Репликация PostgreSQL

### Как это работает

1. **Master** - Основная БД (Read/Write)
   - Записывает все изменения в WAL (Write-Ahead Logs)
   - Отправляет логи на Slave

2. **Slave** - Реплика БД (Read-only)
   - Восстанавливает данные из WAL
   - Постоянно синхронизируется с Master
   - Готов принять роль Master при отказе

### Мониторинг репликации

```bash
# Docker Compose
docker-compose exec postgres-master psql -U postgres

# В psql:
SELECT * FROM pg_stat_replication;  # На Master
SELECT pg_last_wal_receive_lsn();   # На Slave

# Kubernetes
kubectl exec -it -n warden postgres-0 -- psql -U postgres
```

## 🔄 Масштабирование

### Горизонтальное (API)

**Автоматическое через HPA** ✅

- Текущие replicas: 3
- Min replicas: 3
- Max replicas: 10
- CPU trigger: 70%
- Memory trigger: 80%

```bash
# Мониторить масштабирование
kubectl get hpa -n warden -w

# Проверить нагрузку
kubectl top pods -n warden
```

### Вертикальное (ресурсы)

Отредактировать `k8s/03-api-deployment.yaml`:

```yaml
resources:
  requests:
    cpu: 100m → 200m      # Увеличить
    memory: 128Mi → 256Mi  # Увеличить
  limits:
    cpu: 500m → 1000m
    memory: 512Mi → 2Gi
```

## 🔒 Безопасность

### ✅ Реализовано

- Secret для учётных данных БД
- NetworkPolicy для ограничения трафика
- PodDisruptionBudget для высокой доступности
- Health checks (liveness & readiness probes)
- Resource limits (CPU & Memory)

### 🔲 Рекомендуется добавить

- RBAC (Role-Based Access Control)
- Istio для mTLS между pods
- Pod Security Policy
- Secrets encryption at rest
- Regular vulnerability scanning

## 📊 Ресурсы

### API Pod (x3-10)

```
Requests:
  - CPU: 100m (0.1 core)
  - Memory: 128Mi

Limits:
  - CPU: 500m (0.5 core)
  - Memory: 512Mi
```

### PostgreSQL Master

```
Requests:
  - CPU: 200m
  - Memory: 256Mi

Limits:
  - CPU: 1000m
  - Memory: 2Gi

Storage: 10GB
```

### PostgreSQL Slave

```
Requests:
  - CPU: 100m
  - Memory: 256Mi

Limits:
  - CPU: 500m
  - Memory: 1Gi

Storage: 10GB
```

## 🧪 Тестирование

### Docker Compose

```bash
# Добавление в чёрный список
curl -X POST http://localhost:8080/api/post \
  -H "Content-Type: application/json" \
  -d '{"social_network":"x","user_id":"550e8400-e29b-41d4-a716-446655440000","nickname":"user1"}'

# Получение списка
curl "http://localhost:8080/api/blacklist?social_network=x&user_id=550e8400-e29b-41d4-a716-446655440000"

# Проверка наличия
curl "http://localhost:8080/api/check?social_network=x&user_id=550e8400-e29b-41d4-a716-446655440000&nickname=user1"
```

### Kubernetes

```bash
# Port forward
kubectl port-forward -n warden svc/warden-api-service 8080:80

# Затем те же curl команды
curl http://localhost:8080/
```

## 🛑 Остановка и очистка

### Docker Compose

```bash
# Остановить контейнеры
docker-compose down

# Остановить и удалить данные
docker-compose down -v
```

### Kubernetes

```bash
# Удалить развёртывание
./k8s/cleanup.sh

# Или вручную
kubectl delete namespace warden
```

## 📚 Дополнительные ресурсы

- `CONTAINERIZATION_PLAN.md` - Детальный план
- `DEPLOYMENT_GUIDE.md` - Полное руководство (troubleshooting, best practices)
- `CONTAINER_README.md` - Quick reference

## ✅ Чек-лист развёртывания

```
Pre-deployment:
  [ ] Dockerfile собирается без ошибок
  [ ] docker-compose работает локально
  [ ] K8s кластер доступен (kubectl cluster-info)

Deployment:
  [ ] Namespace создан
  [ ] PostgreSQL синхронизирована
  [ ] API pods healthy
  [ ] Services созданы

Post-deployment:
  [ ] API отвечает на запросы
  [ ] Репликация БД работает
  [ ] HPA работает корректно
  [ ] Логи пишутся в stdout
```

## 🎯 Итог

Проект полностью контейнеризирован с поддержкой:

✅ Docker для локальной разработки  
✅ Kubernetes для production  
✅ Master-Slave PostgreSQL репликация  
✅ Горизонтальное масштабирование (HPA)  
✅ Load balancing (Nginx + K8s Services)  
✅ High availability (Pod Anti-Affinity, PDB)  
✅ Network security (NetworkPolicy)  
✅ Health checks и probes  
✅ Resource management (requests/limits)  
✅ Автоматическое развёртывание (deploy scripts)  
