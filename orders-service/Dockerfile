# Stage 1: сборка
FROM golang:1.24-alpine AS builder
WORKDIR /app

# 1) скачиваем зависимости
COPY go.mod go.sum ./
RUN go mod download

# 2) копируем код
COPY . .
COPY migrations orders-service/migrations

# 3) собираем бинарь
RUN CGO_ENABLED=0 GOOS=linux \
    go build -ldflags="-s -w" -o /usr/local/bin/orders-service ./cmd/orders

# Stage 2: минимальный рантайм
FROM alpine:latest
RUN apk add --no-cache ca-certificates

COPY --from=builder /usr/local/bin/orders-service /usr/local/bin/orders-service
COPY --from=builder /app/migrations /app/migrations

EXPOSE 8081

USER nobody

ENTRYPOINT ["/usr/local/bin/orders-service"]
