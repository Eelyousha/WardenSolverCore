#!/bin/bash
# Примеры тестовых запросов к API WardenSolverCore
# Используйте валидные UUID для user_id

# UUID для тестирования (замените на реальные UUID при необходимости)
UUID_USER1="550e8400-e29b-41d4-a716-446655440000"
UUID_USER2="f47ac10b-58cc-4372-a567-0e02b2c3d479"

echo "=== Пример 1: Добавить пользователя в чёрный список (Telegram) ==="
curl -X POST http://localhost:8080/api/post \
  -H "Content-Type: application/json" \
  -d "{
    \"social_network\": \"telegram\",
    \"user_id\": \"$UUID_USER1\",
    \"nickname\": \"baduser123\"
  }"
echo -e "\n"

echo "=== Пример 2: Добавить другого пользователя (X/Twitter) ==="
curl -X POST http://localhost:8080/api/post \
  -H "Content-Type: application/json" \
  -d "{
    \"social_network\": \"x\",
    \"user_id\": \"$UUID_USER1\",
    \"nickname\": \"troll_account\"
  }"
echo -e "\n"

echo "=== Пример 3: Добавить пользователя (Twitch) ==="
curl -X POST http://localhost:8080/api/post \
  -H "Content-Type: application/json" \
  -d "{
    \"social_network\": \"twitch\",
    \"user_id\": \"$UUID_USER2\",
    \"nickname\": \"spammer123\"
  }"
echo -e "\n"

echo "=== Пример 4: Получить весь чёрный список для пользователя (Telegram) ==="
curl -X GET "http://localhost:8080/api/blacklist?social_network=telegram&user_id=$UUID_USER1"
echo -e "\n"

echo "=== Пример 5: Получить весь чёрный список для пользователя (X/Twitter) ==="
curl -X GET "http://localhost:8080/api/blacklist?social_network=x&user_id=$UUID_USER1"
echo -e "\n"

echo "=== Пример 6: Получить весь чёрный список для другого пользователя ==="
curl -X GET "http://localhost:8080/api/blacklist?social_network=twitch&user_id=$UUID_USER2"
echo -e "\n"

echo "=== Пример 7: Проверить конкретного пользователя (в списке) ==="
curl -X GET "http://localhost:8080/api/check?social_network=telegram&user_id=$UUID_USER1&nickname=baduser123"
echo -e "\n"

echo "=== Пример 8: Проверить конкретного пользователя (не в списке) ==="
curl -X GET "http://localhost:8080/api/check?social_network=telegram&user_id=$UUID_USER1&nickname=gooduser"
echo -e "\n"

echo "=== Пример 9: Проверить другого пользователя ==="
curl -X GET "http://localhost:8080/api/check?social_network=x&user_id=$UUID_USER1&nickname=troll_account"
echo -e "\n"

echo "=== Пример 10: Получить информацию о сервере ==="
curl -X GET "http://localhost:8080/"
echo -e "\n"

echo "=== Примеры ошибок ==="

echo -e "\n--- Ошибка 1: Некорректная социальная сеть ---"
curl -X POST http://localhost:8080/api/post \
  -H "Content-Type: application/json" \
  -d "{
    \"social_network\": \"facebook\",
    \"user_id\": \"$UUID_USER1\",
    \"nickname\": \"baduser\"
  }"
echo -e "\n"

echo -e "\n--- Ошибка 2: Пустой user_id ---"
curl -X POST http://localhost:8080/api/post \
  -H "Content-Type: application/json" \
  -d "{
    \"social_network\": \"telegram\",
    \"user_id\": \"\",
    \"nickname\": \"baduser\"
  }"
echo -e "\n"

echo -e "\n--- Ошибка 3: Пустой nickname ---"
curl -X POST http://localhost:8080/api/post \
  -H "Content-Type: application/json" \
  -d "{
    \"social_network\": \"telegram\",
    \"user_id\": \"$UUID_USER1\",
    \"nickname\": \"\"
  }"
echo -e "\n"

echo -e "\n--- Ошибка 4: Некорректный JSON ---"
curl -X POST http://localhost:8080/api/post \
  -H "Content-Type: application/json" \
  -d "invalid json"
echo -e "\n"

echo -e "\n--- Ошибка 5: Отсутствует обязательный параметр (GET) ---"
curl -X GET "http://localhost:8080/api/blacklist?user_id=$UUID_USER1"
echo -e "\n"
