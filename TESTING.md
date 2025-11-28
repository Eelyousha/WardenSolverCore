# Тестирование WardenSolverCore

Этот документ описывает тестирование функционала приложения WardenSolverCore.

## Структура тестов

Тесты размещены согласно стандартам Go в следующих файлах:

```
handlers/v1/
  ├── blacklist_test.go    # Тесты для всех обработчиков (POST, GET, CHECK)
  └── types_test.go        # Тесты для типов и валидации
```

## Запуск тестов

### Все тесты
```bash
cd handlers/v1
go test -v
```

### Конкретный тест
```bash
go test -v -run TestHandlePostBlacklistSuccess
```

### С покрытием кода
```bash
go test -v -cover
go test -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### Бенчмарки (тесты производительности)
```bash
go test -bench=. -benchmem
```

## Описание тестов

### blacklist_test.go - Тесты обработчиков

#### POST /api/post - Добавление в чёрный список

- **TestHandlePostBlacklistSuccess** - Успешное добавление пользователя
- **TestHandlePostBlacklistInvalidMethod** - Проверка некорректного HTTP метода (GET вместо POST)
- **TestHandlePostBlacklistInvalidJSON** - Некорректный JSON в теле запроса
- **TestHandlePostBlacklistInvalidSocialNetwork** - Несуществующая социальная сеть
- **TestHandlePostBlacklistMissingUserID** - Отсутствие поля user_id
- **TestHandlePostBlacklistMissingNickname** - Отсутствие поля nickname
- **BenchmarkHandlePostBlacklist** - Тест производительности

#### GET /api/blacklist - Получение чёрного списка

- **TestHandleGetBlacklistSuccess** - Успешное получение списка
- **TestHandleGetBlacklistInvalidMethod** - Проверка некорректного метода (POST вместо GET)
- **TestHandleGetBlacklistMissingSocialNetwork** - Отсутствие параметра social_network
- **TestHandleGetBlacklistMissingUserID** - Отсутствие параметра user_id
- **TestHandleGetBlacklistInvalidSocialNetwork** - Некорректная социальная сеть
- **BenchmarkHandleGetBlacklist** - Тест производительности

#### GET /api/check - Проверка пользователя в списке

- **TestHandleCheckBlacklistSuccess** - Успешная проверка (пользователь в списке)
- **TestHandleCheckBlacklistUserNotInList** - Пользователь не в списке
- **TestHandleCheckBlacklistInvalidMethod** - Некорректный метод (POST вместо GET)
- **TestHandleCheckBlacklistMissingSocialNetwork** - Отсутствие параметра social_network
- **TestHandleCheckBlacklistMissingUserID** - Отсутствие параметра user_id
- **TestHandleCheckBlacklistMissingNickname** - Отсутствие параметра nickname
- **TestHandleCheckBlacklistInvalidSocialNetwork** - Некорректная социальная сеть
- **BenchmarkHandleCheckBlacklist** - Тест производительности

#### GET / - Корневой endpoint

- **TestHandleRoot** - Проверка корректного ответа

### types_test.go - Тесты типов и валидации

- **TestSocialNetworksValidation** - Проверка валидации допустимых социальных сетей
- **TestResponsePayloadStructure** - Проверка структуры ответа с данными
- **TestResponsePayloadWithoutData** - Проверка структуры ответа без данных

## Mock Database

Для тестирования используется `MockDatabase`, который реализует интерфейс `DatabaseInterface`. Это позволяет тестировать обработчики без подключения к реальной базе данных.

```go
mockDB := &MockDatabase{
    addToBlacklistFunc: func(socialNetwork, userID, nickname string) (int, error) {
        return 1, nil
    },
}
```

## Результаты тестов

Все 31 тест успешно проходит:

```
PASS
ok      handlerv1       0.004s
```

### Производительность (бенчмарки):

| Операция | Вызовов в сек | Наносекунд на операцию | Выделено памяти |
|----------|---------------|----------------------|-----------------|
| POST /api/post | ~179,959 | 6,942 ns/op | 8,858 B |
| GET /api/blacklist | ~208,686 | 5,780 ns/op | 8,449 B |
| GET /api/check | ~208,551 | 5,708 ns/op | 8,465 B |

## Запуск через главный модуль

```bash
# Из корневой директории проекта
go test -v ./handlers/v1/...
```

## Интеграционное тестирование

Для интеграционного тестирования запустите сервер:

```bash
go run main.go
```

И отправьте тестовые запросы:

```bash
# POST запрос
curl -X POST http://localhost:8080/api/post \
  -H "Content-Type: application/json" \
  -d '{"social_network": "telegram", "user_id": "123", "nickname": "baduser"}'

# GET запрос списка
curl "http://localhost:8080/api/blacklist?social_network=telegram&user_id=123"

# GET запрос проверки
curl "http://localhost:8080/api/check?social_network=telegram&user_id=123&nickname=baduser"
```
