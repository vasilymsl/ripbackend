# Используем официальный образ Golang (последняя версия)
FROM golang:alpine AS builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем go.mod и go.sum для кеширования зависимостей
COPY go.mod go.sum ./

# Скачиваем зависимости
RUN go mod download

# Копируем весь код приложения
COPY . .

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server ./cmd/lab1/main.go

# Финальный образ
FROM alpine:latest

WORKDIR /app

# Копируем скомпилированное приложение из builder
COPY --from=builder /app/server .

# Копируем статические файлы и шаблоны
COPY templates ./templates
COPY resources ./resources
COPY config/config-docker.toml ./config/config-docker.toml

# Устанавливаем переменные окружения по умолчанию
ENV GIN_MODE=release
ENV DB_HOST=postgres
ENV DB_PORT=5432
ENV DB_USER=myuser
ENV DB_PASSWORD=mypassword
ENV DB_NAME=credits

# Открываем порт
EXPOSE 8080

# Запускаем приложение
CMD ["./server"]

