# Ступень 1: Сборка приложения
FROM golang:1.22 AS builder

WORKDIR /app

ENV REDIS_USER=
ENV REDIS_PASSWORD=

# Копируем только файлы go.mod и go.sum, чтобы сначала установить зависимости
COPY go.mod go.sum ./

# Загружаем зависимости
RUN go mod download

# Копируем остальные файлы проекта
COPY . .

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o floodcontrol cmd/main.go

# Ступень 2: Создание минимального образа для запуска приложения
FROM alpine:latest  

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /root/

# Копируем собранное приложение из ступени 1
COPY --from=builder /app/floodcontrol .

# Добавляем файл config.yml в контейнер
COPY ./configs/config.yml ./configs/

# Открываем порт, на котором работает приложение
EXPOSE 8080

# Запускаем приложение
CMD ["./floodcontrol"]