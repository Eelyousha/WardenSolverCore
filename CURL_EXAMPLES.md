# Примеры тестовых запросов с UUID

## Об использовании UUID

`user_id` должен быть передан в формате UUID (UUID v4):
- Формат: `xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx`
- Пример: `550e8400-e29b-41d4-a716-446655440000`

## Быстрый старт

### 1. Запустите сервер
```bash
go run main.go
```

### 2. В отдельном терминале, отправьте тестовый запрос

**POST - Добавить пользователя в чёрный список**
```bash
curl -X POST http://localhost:8080/api/post \
  -H "Content-Type: application/json" \
  -d '{
    "social_network": "telegram",
    "user_id": "550e8400-e29b-41d4-a716-446655440000",
    "nickname": "baduser123"
  }'
```

**GET - Получить чёрный список**
```bash
curl "http://localhost:8080/api/blacklist?social_network=telegram&user_id=550e8400-e29b-41d4-a716-446655440000"
```

**GET - Проверить пользователя в списке**
```bash
curl "http://localhost:8080/api/check?social_network=telegram&user_id=550e8400-e29b-41d4-a716-446655440000&nickname=baduser123"
```

## Запуск всех примеров

Для удобства создан скрипт со всеми примерами запросов:

```bash
# Убедитесь, что сервер работает
go run main.go &

# В другом терминале запустите тестовые запросы
./test_requests.sh
```

Скрипт содержит:
- 10 успешных примеров (добавление, получение, проверка)
- 5 примеров ошибок (некорректные данные)

## Доступные социальные сети

- `telegram` - Telegram
- `x` - X (Twitter)
- `twitch` - Twitch

## UUID для тестирования

Используйте любой валидный UUID v4:

```bash
# Пример 1
550e8400-e29b-41d4-a716-446655440000

# Пример 2
f47ac10b-58cc-4372-a567-0e02b2c3d479

# Пример 3 (можно сгенерировать)
uuidgen  # на macOS/Linux
```

## Примеры с curl

### Успешные запросы

```bash
# 1. Добавить в Telegram
curl -X POST http://localhost:8080/api/post \
  -H "Content-Type: application/json" \
  -d '{"social_network":"telegram","user_id":"550e8400-e29b-41d4-a716-446655440000","nickname":"baduser"}'

# 2. Получить список
curl "http://localhost:8080/api/blacklist?social_network=telegram&user_id=550e8400-e29b-41d4-a716-446655440000"

# 3. Проверить пользователя
curl "http://localhost:8080/api/check?social_network=telegram&user_id=550e8400-e29b-41d4-a716-446655440000&nickname=baduser"
```

### Ошибки

```bash
# Некорректная социальная сеть
curl -X POST http://localhost:8080/api/post \
  -H "Content-Type: application/json" \
  -d '{"social_network":"facebook","user_id":"550e8400-e29b-41d4-a716-446655440000","nickname":"baduser"}'

# Пустой user_id
curl -X POST http://localhost:8080/api/post \
  -H "Content-Type: application/json" \
  -d '{"social_network":"telegram","user_id":"","nickname":"baduser"}'

# Некорректный JSON
curl -X POST http://localhost:8080/api/post \
  -H "Content-Type: application/json" \
  -d 'invalid'
```

## Тестирование с Postman

1. **POST /api/post**
   - URL: `http://localhost:8080/api/post`
   - Body (JSON):
   ```json
   {
     "social_network": "telegram",
     "user_id": "550e8400-e29b-41d4-a716-446655440000",
     "nickname": "baduser123"
   }
   ```

2. **GET /api/blacklist**
   - URL: `http://localhost:8080/api/blacklist?social_network=telegram&user_id=550e8400-e29b-41d4-a716-446655440000`

3. **GET /api/check**
   - URL: `http://localhost:8080/api/check?social_network=telegram&user_id=550e8400-e29b-41d4-a716-446655440000&nickname=baduser123`
