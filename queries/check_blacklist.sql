-- Проверка наличия конкретного пользователя в чёрном списке
SELECT EXISTS(
    SELECT 1 FROM blacklist
    WHERE social_network = $1 AND user_id = $2 AND nickname = $3
) AS is_banned;
