# Этап 1: Сборка Go приложения
FROM golang:1.25-alpine AS builder

# Рабочая директория внутри контейнера
WORKDIR /app

# Копируем go.mod и go.sum из корня проекта
COPY go.mod go.sum ./

# Скачиваем зависимости
RUN go mod download

# Копируем весь проект
COPY . .

# Переходим в папку с main.go
WORKDIR /app/auth/cmd/auth

# Собираем бинарь
RUN go build -o auth main.go

# Этап 2: Финальный образ (runtime)
FROM alpine:latest

# Устанавливаем ca-certificates для HTTPS запросов
RUN apk --no-cache add ca-certificates

# Создаем пользователя для безопасности
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# Копируем собранное приложение из builder этапа
COPY --from=builder /app/auth/cmd/auth/auth /usr/local/bin/auth

# Меняем владельца файла на не-root пользователя
RUN chown appuser:appgroup /usr/local/bin/auth

# Переключаемся на не-root пользователя
USER appuser

# Открываем порт
EXPOSE 8082

# Запускаем приложение
CMD ["/usr/local/bin/auth"]