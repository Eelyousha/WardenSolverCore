# WardenSolverCore - Quick Start Guide

## 🚀 За 5 минут в production

### Требования

- Docker & Docker Compose (для локального тестирования)
- Kubernetes кластер (для production)
- kubectl (для управления K8s)

### 1️⃣ Локально (Docker Compose)

```bash
# Собрать образ
./docker/build.sh

# Запустить контейнеры
./docker/up.sh

# Проверить
curl http://localhost:8080/
```

**Результат**: 4 контейнера запущены + данные синхронизированы

### 2️⃣ В Kubernetes

```bash
# Развернуть все
./k8s/deploy.sh

# Проверить статус
kubectl get all -n warden

# Тестировать
kubectl port-forward -n warden svc/warden-api-service 8080:80
curl http://localhost:8080/
```

**Результат**: Production-ready deployment с 3+ replicas

---

## 📋 Что получилось

| Компонент | Локально | K8s | Особенности |
|-----------|----------|-----|------------|
| **API** | 1 pod | 3-10 replicas | Auto-scaling, Load balancing |
| **PostgreSQL Master** | 1 container | StatefulSet | Read/Write, WAL replication |
| **PostgreSQL Slave** | 1 container | StatefulSet | Read-only, auto-sync |
| **Load Balancer** | Nginx | Ingress | Rate limiting, SSL ready |
| **Storage** | Volumes | PersistentVolumes | 10GB each |
| **Health Checks** | ✓ | ✓ | Liveness & Readiness probes |

---

## 🔍 Структура

```
Project Root/
├── docker/              # Docker configs
│   ├── build.sh        # Build script
│   ├── up.sh           # Run script
│   └── ... configs
├── k8s/                # Kubernetes configs
│   ├── deploy.sh       # Deploy script
│   ├── cleanup.sh      # Cleanup script
│   └── *.yaml          # K8s manifests
├── Dockerfile          # App image
├── docker-compose.yml  # Dev environment
└── ... docs
```

---

## 🎯 Основные особенности

✅ **Master-Slave Реplication** - PostgreSQL синхронизируется автоматически  
✅ **Auto-scaling** - HPA масштабирует API при нагрузке (3-10 replicas)  
✅ **Load Balancing** - Nginx + K8s Services  
✅ **High Availability** - Pod Anti-Affinity + PodDisruptionBudget  
✅ **Security** - NetworkPolicy, Secrets, Resource limits  
✅ **Monitoring** - Health checks, liveness/readiness probes  

---

## 📖 Документация

- `CONTAINER_README.md` - Полный обзор
- `DEPLOYMENT_GUIDE.md` - Подробное руководство
- `CONTAINERIZATION_PLAN.md` - Архитектура
- `CONTAINERIZATION_SUMMARY.md` - Итоговый процесс

---

## 🐛 Troubleshooting

### API pods не стартуют?
```bash
kubectl describe pod -n warden <pod-name>
kubectl logs -n warden -l app=warden-api
```

### БД не синхронизируется?
```bash
kubectl exec -it -n warden postgres-0 -- psql -U postgres
SELECT * FROM pg_stat_replication;
```

### Изменить количество replicas?
```bash
kubectl scale deployment warden-api -n warden --replicas=5
```

---

Для более подробной информации см. документацию в проекте.
