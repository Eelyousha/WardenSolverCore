# Multi-stage build для минимизации размера образа
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Копируем go.mod и go.sum
COPY go.mod go.sum ./

# Копируем локальные пакеты (они указаны в replace directives)
COPY generators ./generators
COPY handlers ./handlers
COPY migrations ./migrations
COPY queries ./queries

# Скачиваем зависимости
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Финальный образ
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Копируем бинарник из builder образа
COPY --from=builder /app/main .

# Копируем миграции и запросы
COPY --from=builder /app/migrations ./migrations
COPY --from=builder /app/queries ./queries

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --quiet --tries=1 --spider http://localhost:8080/ || exit 1

# Запуск приложения
CMD ["./main"]
