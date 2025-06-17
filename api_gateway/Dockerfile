# Stage 1: сборка
FROM golang:1.24-alpine AS builder
WORKDIR /app

# 1) скачиваем зависимости
COPY go.mod go.sum ./
RUN go mod download

# 2) копируем код
COPY . .

# 3) собираем бинарь
RUN CGO_ENABLED=0 GOOS=linux \
    go build -ldflags="-s -w" -o /usr/local/bin/api_gateway ./cmd/gateway

# Stage 2: минимальный рантайм
FROM alpine:latest
RUN apk add --no-cache ca-certificates

COPY --from=builder /usr/local/bin/api_gateway /usr/local/bin/api_gateway

EXPOSE 8083

USER nobody

ENTRYPOINT ["/usr/local/bin/api_gateway"]
