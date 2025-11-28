# Руководство по контейнеризации WardenSolverCore

## Содержание

1. [Локальная разработка (Docker Compose)](#локальная-разработка)
2. [Развёртывание в Kubernetes](#развёртывание-в-kubernetes)
3. [Архитектура и масштабирование](#архитектура)
4. [Мониторинг и отладка](#мониторинг)
5. [Production considerations](#production)

---

## Локальная разработка

### Требования

- Docker 20.10+
- Docker Compose 2.0+
- 4GB RAM
- 10GB свободного места на диске

### Быстрый старт

```bash
# Сделать скрипты исполняемыми
chmod +x docker/*.sh

# Собрать Docker образ
./docker/build.sh

# Запустить контейнеры
./docker/up.sh
```

### Проверка статуса

```bash
# Просмотр всех контейнеров
docker-compose ps

# Просмотр логов API
docker-compose logs -f api

# Просмотр логов PostgreSQL Master
docker-compose logs -f postgres-master

# Просмотр логов PostgreSQL Slave
docker-compose logs -f postgres-slave
```

### Тестирование API

```bash
# Проверка здоровья
curl http://localhost:8080/

# Через Nginx Load Balancer
curl http://localhost/

# Добавление в чёрный список
curl -X POST http://localhost:8080/api/post \
  -H "Content-Type: application/json" \
  -d '{
    "social_network": "x",
    "user_id": "550e8400-e29b-41d4-a716-446655440000",
    "nickname": "testuser"
  }'

# Получение списка
curl "http://localhost:8080/api/blacklist?social_network=x&user_id=550e8400-e29b-41d4-a716-446655440000"

# Проверка наличия
curl "http://localhost:8080/api/check?social_network=x&user_id=550e8400-e29b-41d4-a716-446655440000&nickname=testuser"
```

### Проверка репликации PostgreSQL

```bash
# Подключение к Master
docker-compose exec postgres-master psql -U postgres -d warden

# Подключение к Slave
docker-compose exec postgres-slave psql -U postgres -d warden

# В Master: проверить статус репликации
SELECT * FROM pg_stat_replication;

# В Slave: проверить статус восстановления
SELECT pg_last_wal_receive_lsn();
```

### Остановка

```bash
# Остановить контейнеры
docker-compose down

# Остановить и удалить данные
docker-compose down -v
```

---

## Развёртывание в Kubernetes

### Требования

- Kubernetes кластер 1.20+
- kubectl настроен и подключен
- 16GB RAM в кластере
- 50GB свободного места на диске

### Шаг 1: Подготовка Docker образа

```bash
# Собрать образ
./docker/build.sh

# Отправить в registry
export DOCKER_REGISTRY=your-registry.com
export IMAGE_NAME=warden-api
export IMAGE_TAG=v1.0.0

docker tag docker.io/warden-api:latest your-registry.com/warden-api:v1.0.0
docker push your-registry.com/warden-api:v1.0.0
```

### Шаг 2: Обновить K8s конфиги

Отредактировать `k8s/03-api-deployment.yaml`:

```yaml
image: your-registry.com/warden-api:v1.0.0
imagePullPolicy: Always
```

### Шаг 3: Развёртывание

```bash
# Сделать скрипты исполняемыми
chmod +x k8s/*.sh

# Развернуть
./k8s/deploy.sh
```

### Шаг 4: Проверка статуса

```bash
# Просмотр всех ресурсов
kubectl get all -n warden

# Просмотр pods
kubectl get pods -n warden -o wide

# Просмотр services
kubectl get svc -n warden

# Просмотр PersistentVolumeClaims
kubectl get pvc -n warden

# Просмотр logs
kubectl logs -n warden -l app=warden-api -f

# Просмотр events
kubectl get events -n warden --sort-by='.lastTimestamp'
```

### Шаг 5: Проверка связи с БД

```bash
# Port forward к API
kubectl port-forward -n warden svc/warden-api-service 8080:80

# В другом терминале
curl http://localhost:8080/

# Port forward к Master БД
kubectl port-forward -n warden svc/postgres-master 5432:5432

# Подключение из CLI
psql -h localhost -U postgres -d warden
```

### Удаление развёртывания

```bash
./k8s/cleanup.sh
```

---

## Архитектура

### API Tier

- **Deployment**: 3 реплики по умолчанию
- **HPA**: 3-10 реплик в зависимости от нагрузки
- **Anti-Affinity**: pods распределяются по разным узлам
- **Resource Limits**: 100-500m CPU, 128-512Mi Memory
- **Health Checks**: liveness & readiness probes каждые 10 и 5 сек

### Database Tier

- **StatefulSet**: 2 pods (Master + Slave)
- **Storage**: 10GB PersistentVolume на каждый pod
- **Replication**: WAL-based репликация
- **Failover**: Ручной (может быть автоматизирован)

### Network

- **Ingress**: HTTP/HTTPS маршрутизация
- **LoadBalancer Service**: Для external access
- **ClusterIP Services**: Для internal communication
- **NetworkPolicy**: Ограничение трафика между pods

---

## Масштабирование

### Горизонтальное масштабирование (API)

Автоматическое через HPA:

```bash
# Просмотр HPA статуса
kubectl get hpa -n warden -w

# Проверить текущую нагрузку
kubectl top pods -n warden
```

### Вертикальное масштабирование

Отредактировать `k8s/03-api-deployment.yaml`:

```yaml
resources:
  requests:
    cpu: 200m      # Увеличить
    memory: 256Mi   # Увеличить
  limits:
    cpu: 1000m
    memory: 2Gi
```

Затем:

```bash
kubectl apply -f k8s/03-api-deployment.yaml
```

---

## Мониторинг

### Просмотр метрик

```bash
# CPU и Memory использование
kubectl top nodes
kubectl top pods -n warden

# Более подробно
kubectl describe pod <pod-name> -n warden
```

### Логирование

```bash
# Логи одного pod'a
kubectl logs -n warden pod/<pod-name>

# Логи всех pods приложения
kubectl logs -n warden -l app=warden-api -f

# Логи предыдущего pod'a (если был перезагружен)
kubectl logs -n warden pod/<pod-name> --previous

# Логи PostgreSQL
kubectl logs -n warden pod/postgres-0
```

### Health checks

```bash
# Проверить health check endpoints
kubectl exec -n warden pod/<api-pod> -- curl -s http://localhost:8080/

# Проверить PostgreSQL доступность
kubectl exec -n warden pod/postgres-0 -- pg_isready -U postgres
```

---

## Production Considerations

### Безопасность

1. **Secrets Management**
   - Использовать HashiCorp Vault или Kubernetes External Secrets Operator
   - Ротация пароля БД каждые 30 дней

2. **Network Security**
   - Включить Network Policies
   - Использовать TLS для всего трафика (Istio Service Mesh)

3. **RBAC**
   - Ограничить доступ через Role-Based Access Control
   - Создать отдельные сервис-аккаунты

### Отказоустойчивость

1. **Redundancy**
   - Минимум 3 реплики API на разных узлах
   - Минимум 2 PostgreSQL узла (Master + Slave)

2. **Backup & Recovery**
   - Автоматические снимки PersistentVolume каждый час
   - Хранить резервные копии в объектном хранилище (S3)
   - Тестировать восстановление ежемесячно

3. **Failover**
   - Автоматический failover Slave → Master при отказе Master'а
   - Для этого использовать pg_auto_failover или Patroni

### Производительность

1. **Кэширование**
   - Добавить Redis для кэширования результатов GET запросов
   - TTL: 5 минут для списков, 60 минут для проверок

2. **Database Optimization**
   - Индексы уже созданы, но мониторить медленные запросы
   - Использовать EXPLAIN для оптимизации запросов

3. **Load Testing**
   - Использовать Apache JMeter или Locust
   - Регулярное тестирование до 1000 RPS

### Обновления

1. **Rolling Updates**
   - Настроены: maxSurge: 1, maxUnavailable: 0
   - PodDisruptionBudget: минимум 2 pod'а должны быть доступны

2. **Database Updates**
   - Обновления PostgreSQL требуют ручного вмешательства
   - Использовать pg_upgrade для минорных версий
   - Major версии требуют dump & restore

---

## Troubleshooting

### API pods не стартуют

```bash
# Проверить events
kubectl describe pod <pod-name> -n warden

# Проверить логи
kubectl logs -n warden pod/<pod-name>

# Проверить ресурсы кластера
kubectl describe nodes

# Проверить connection к БД
# kubectl exec -n warden pod/<api-pod> -- \
#   psql -h postgres-master -U postgres -d warden -c "SELECT 1"
```

### PostgreSQL не синхронизируется

```bash
# Подключиться к Master
kubectl exec -it -n warden pod/postgres-0 -- psql -U postgres

# Проверить репликацию
SELECT * FROM pg_stat_replication;

# Если не видно Slave'а, перезагрузить Slave pod
kubectl delete pod -n warden postgres-1
```

### High latency на API

```bash
# Проверить CPU/Memory
kubectl top pods -n warden

# Если высокий CPU - увеличить HPA maxReplicas
# Если высокий Memory - искать утечку в коде

# Проверить network latency между pods
kubectl exec -n warden pod/<api-pod> -- ping postgres-master
```

---

## Дополнительные ресурсы

- [Kubernetes Documentation](https://kubernetes.io/docs/)
- [Docker Documentation](https://docs.docker.com/)
- [PostgreSQL Replication](https://www.postgresql.org/docs/current/warm-standby.html)
- [Kubernetes Deployment Best Practices](https://kubernetes.io/docs/concepts/configuration/overview/)
