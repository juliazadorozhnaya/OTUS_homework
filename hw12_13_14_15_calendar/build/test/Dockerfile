# Используем базовый образ golang
FROM golang:1.22.1-alpine as build

# Устанавливаем переменные окружения
ENV BIN_FILE /bin/app
ENV CODE_DIR /go/src/app

# Создаем рабочую директорию
WORKDIR ${CODE_DIR}

# Устанавливаем зависимости
RUN apk update && apk add --no-cache git

# Копируем файлы проекта
COPY ../../ ./

# Устанавливаем зависимости проекта
RUN go mod download

# Копируем папку с тестами
COPY ../../integration_tests ./integration_tests

# Запускаем тесты
CMD go test ./integration_tests
