# Warden Solver Core

## Описание
API для управления чёрным списком пользователей в социальных сетях (X, Telegram, Twitch).

## Структура проекта

```
.
├── main.go                          # Точка входа и маршруты
├── database.go                      # Функции для работы с PostgreSQL
├── migrations.go                    # Функция для запуска миграций
├── go.mod                          # Go зависимости
├── go.sum                          # Go checksums
│
├── migrations/                     # Миграции БД
│   └── 001_create_blacklist_table.sql
│
├── queries/                        # SQL-запросы к БД
│   ├── insert_blacklist.sql       # Добавление в чёрный список
│   ├── get_blacklist.sql          # Получение чёрного списка
│   └── check_blacklist.sql        # Проверка наличия в чёрном списке
│
├── generators/                     # Логика генерации данных
│   ├── go.mod
│   └── blacklist.go               # Генератор записей чёрного списка
│
└── handlers/v1/                    # Обработчики HTTP запросов
    ├── go.mod
    ├── types.go                   # Типы и интерфейсы
    └── blacklist.go               # Обработчики для работы с чёрным списком
```

## Настройка

### 1. Переменные окружения

```bash
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=postgres
export DB_NAME=warden
```

### 2. Создание базы данных

```sql
CREATE DATABASE warden;
```

### 3. Установка зависимостей

```bash
go mod download
go mod tidy
```

### 4. Запуск приложения

Миграции выполняются автоматически при запуске:

```bash
go run *.go
```

Или используя собранный бинарник:

```bash
go build -o main
./main
```

## API Endpoints

### 1. POST /api/post - Добавить пользователя в чёрный список

**Запрос:**
```bash
curl -X POST http://localhost:8080/api/post \
  -H "Content-Type: application/json" \
  -d '{
    "social_network": "x",
    "user_id": "550e8400-e29b-41d4-a716-446655440000",
    "nickname": "testuser"
  }'
```

**Ответ:**
```json
{
  "status": "success",
  "message": "User testuser added to blacklist in x",
  "data": {
    "id": 1,
    "social_network": "x",
    "user_id": "550e8400-e29b-41d4-a716-446655440000",
    "nickname": "testuser"
  }
}
```

### 2. GET /api/blacklist - Получить чёрный список пользователя

**Запрос:**
```bash
curl "http://localhost:8080/api/blacklist?social_network=x&user_id=550e8400-e29b-41d4-a716-446655440000"
```

**Ответ:**
```json
{
  "status": "success",
  "message": "Blacklist retrieved successfully",
  "data": {
    "social_network": "x",
    "user_id": "550e8400-e29b-41d4-a716-446655440000",
    "banned_count": 2,
    "banned_users": [
      {
        "nickname": "testuser2",
        "created_at": "2025-11-24T10:30:00Z"
      },
      {
        "nickname": "testuser",
        "created_at": "2025-11-24T10:20:00Z"
      }
    ]
  }
}
```

### 3. GET /api/check - Проверить наличие пользователя в чёрном списке

**Запрос:**
```bash
curl "http://localhost:8080/api/check?social_network=x&user_id=550e8400-e29b-41d4-a716-446655440000&nickname=testuser"
```

**Ответ:**
```json
{
  "status": "success",
  "message": "Check completed",
  "data": {
    "social_network": "x",
    "user_id": "550e8400-e29b-41d4-a716-446655440000",
    "nickname": "testuser",
    "is_banned": true
  }
}
```

## Допустимые социальные сети

- `x` - X (Twitter)
- `telegram` - Telegram
- `twitch` - Twitch

## Структура проекта

```
.
├── main.go              # Основной файл с обработчиками HTTP
├── database.go          # Функции для работы с PostgreSQL
├── schema.sql           # Схема базы данных
├── go.mod              # Go зависимости
└── queries/            # SQL-запросы
    ├── insert_blacklist.sql  # Добавление пользователя в чёрный список
    ├── get_blacklist.sql     # Получение чёрного списка
    └── check_blacklist.sql   # Проверка наличия в чёрном списке
```

## Настройка

### 1. Переменные окружения

```bash
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=postgres
export DB_NAME=warden
```

### 2. Создание базы данных

```sql
CREATE DATABASE warden;
```

### 3. Применение схемы

```bash
psql -U postgres -d warden -f schema.sql
```

### 4. Установка зависимостей

```bash
go mod download
```

### 5. Запуск приложения

```bash
go run main.go database.go
```

## API Endpoints

### 1. POST /api/post - Добавить пользователя в чёрный список

**Запрос:**
```bash
curl -X POST http://localhost:8080/api/post \
  -H "Content-Type: application/json" \
  -d '{
    "social_network": "x",
    "user_id": "550e8400-e29b-41d4-a716-446655440000",
    "nickname": "testuser"
  }'
```

**Ответ:**
```json
{
  "status": "success",
  "message": "User testuser added to blacklist in x",
  "data": {
    "id": 1,
    "social_network": "x",
    "user_id": "550e8400-e29b-41d4-a716-446655440000",
    "nickname": "testuser"
  }
}
```

### 2. GET /api/blacklist - Получить чёрный список пользователя

**Запрос:**
```bash
curl "http://localhost:8080/api/blacklist?social_network=x&user_id=550e8400-e29b-41d4-a716-446655440000"
```

**Ответ:**
```json
{
  "status": "success",
  "message": "Blacklist retrieved successfully",
  "data": {
    "social_network": "x",
    "user_id": "550e8400-e29b-41d4-a716-446655440000",
    "banned_count": 2,
    "banned_users": [
      {
        "nickname": "testuser2",
        "created_at": "2025-11-24T10:30:00Z"
      },
      {
        "nickname": "testuser",
        "created_at": "2025-11-24T10:20:00Z"
      }
    ]
  }
}
```

### 3. GET /api/check - Проверить наличие пользователя в чёрном списке

**Запрос:**
```bash
curl "http://localhost:8080/api/check?social_network=x&user_id=550e8400-e29b-41d4-a716-446655440000&nickname=testuser"
```

**Ответ:**
```json
{
  "status": "success",
  "message": "Check completed",
  "data": {
    "social_network": "x",
    "user_id": "550e8400-e29b-41d4-a716-446655440000",
    "nickname": "testuser",
    "is_banned": true
  }
}
```

## Допустимые социальные сети

- `x` - X (Twitter)
- `telegram` - Telegram
- `twitch` - Twitch
