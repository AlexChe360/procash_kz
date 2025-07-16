# 1. Используем Alpine с Go
FROM golang:1.23-alpine

# 2. Устанавливаем C-компилятор и зависимости
RUN apk add --no-cache gcc musl-dev

# 3. Переменные окружения
ENV CGO_ENABLED=1 \
    GO111MODULE=on \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /app

# 4. Скачиваем зависимости
COPY go.mod go.sum ./
RUN go mod download

# 5. Копируем всё и билдим
COPY . .
COPY .env .env

RUN go build -o server ./cmd/api/main.go

EXPOSE 3000

CMD ["./server"]
