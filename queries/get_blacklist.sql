-- Получение всех забаненных пользователей по социальной сети и user_id
SELECT nickname, created_at
FROM blacklist
WHERE social_network = $1 AND user_id = $2;
