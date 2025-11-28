-- Migration: Create blacklist table
-- Created at: 2025-11-24

CREATE TABLE IF NOT EXISTS blacklist (
    id SERIAL PRIMARY KEY,
    social_network VARCHAR(50) NOT NULL,
    user_id UUID NOT NULL,
    nickname VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(social_network, user_id, nickname)
);

-- Индекс для быстрого поиска по социальной сети и user_id
CREATE INDEX IF NOT EXISTS idx_blacklist_social_user ON blacklist(social_network, user_id);

-- Индекс для быстрого поиска по социальной сети, user_id и nickname
CREATE INDEX IF NOT EXISTS idx_blacklist_check ON blacklist(social_network, user_id, nickname);
