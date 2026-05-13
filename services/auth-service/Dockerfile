FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY . .

WORKDIR /app/services/auth-service

RUN go mod download

RUN go build -o auth-service ./cmd/server

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/services/auth-service/auth-service .

CMD ["./auth-service"]