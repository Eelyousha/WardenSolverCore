#!/bin/bash
# PostgreSQL Slave initialization

# Ждём, пока Master будет доступен
echo "Waiting for PostgreSQL Master..."
until pg_basebackup -h postgres-master -U replicator -D /var/lib/postgresql/data -Fp -Xs -P; do
  echo "Master is not ready, retrying..."
  sleep 1
done

# Создаём файл recovery.conf для включения режима standby
cat > /var/lib/postgresql/data/recovery.conf <<EOF
standby_mode = 'on'
primary_conninfo = 'host=postgres-master port=5432 user=replicator password=replicator_password'
recovery_target_timeline = 'latest'
EOF

# Устанавливаем правильные права доступа
chmod 600 /var/lib/postgresql/data/recovery.conf

echo "Slave initialization completed"
