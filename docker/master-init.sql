-- PostgreSQL Master initialization
-- Создаем пользователя репликации
CREATE ROLE replicator WITH REPLICATION PASSWORD 'replicator_password' LOGIN;

-- Даем права на подключение
GRANT CONNECT ON DATABASE warden TO replicator;

-- Убеждаемся, что репликация включена
ALTER SYSTEM SET max_wal_senders = 3;
ALTER SYSTEM SET max_replication_slots = 3;
ALTER SYSTEM SET wal_keep_size = '1GB';

-- Применяем миграции
\c warden

CREATE TABLE IF NOT EXISTS blacklist (
    id SERIAL PRIMARY KEY,
    social_network VARCHAR(50) NOT NULL,
    user_id UUID NOT NULL,
    nickname VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(social_network, user_id, nickname)
);

CREATE INDEX IF NOT EXISTS idx_blacklist_social_user ON blacklist(social_network, user_id);
CREATE INDEX IF NOT EXISTS idx_blacklist_check ON blacklist(social_network, user_id, nickname);
