# Этап 1: Сборка Go приложения
FROM golang:1.25-alpine AS builder

# Рабочая директория внутри контейнера
WORKDIR /app

# Копируем go.mod
COPY go.mod go.sum ./

# Скачиваем зависимости
RUN go mod download

# Копируем весь проект
COPY . .

# Собираем бинарь из корня проекта
RUN go build -o ./gateway/gateway ./gateway/cmd/gateway

# Этап 2: Финальный образ (runtime)
FROM alpine:latest

# Устанавливаем ca-certificates для HTTPS запросов
RUN apk --no-cache add ca-certificates

# Создаем пользователя для безопасности
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# Копируем собранное приложение из builder этапа
COPY --from=builder /app/gateway/gateway /usr/local/bin/gateway

# Меняем владельца файла на не-root пользователя
RUN chown appuser:appgroup /usr/local/bin/gateway

# Переключаемся на не-root пользователя
USER appuser

# Открываем порт
EXPOSE 3000

# Запускаем приложение
CMD ["/usr/local/bin/gateway"]