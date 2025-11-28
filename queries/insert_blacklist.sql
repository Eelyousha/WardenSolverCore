-- Добавление пользователя в чёрный список
INSERT INTO blacklist (social_network, user_id, nickname)
VALUES ($1, $2, $3)
ON CONFLICT (social_network, user_id, nickname) DO NOTHING
RETURNING id, social_network, user_id, nickname, created_at;
