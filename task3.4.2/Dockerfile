# Используем официальный образ Go как базовый
FROM golang:1.19-alpine as builder

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем исходники приложения в рабочую директорию
COPY . .

# Скачиваем все зависимости
RUN go mod tidy

# Собираем приложение
RUN go build -o main

# Начинаем новую стадию сборки на основе минимального образа
FROM alpine:latest

# Копируем исполняемый файл и необходимые ресурсы из первой стадии
COPY --from=builder /app/main /main
COPY --from=builder /app/internal /internal
COPY --from=builder /app/migrations /migrations
COPY --from=builder /app/db.env .

# Открываем порт 8080
EXPOSE 8080

# Запускаем приложение
CMD ["/main"]
