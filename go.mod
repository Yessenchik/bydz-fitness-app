module github.com/Yessenchik/bydz-fitness-app

go 1.26

require (
	github.com/golang-jwt/jwt/v5 v5.3.1
	github.com/golang/mock v1.6.0
	github.com/google/uuid v1.6.0
	github.com/jmoiron/sqlx v1.4.0
	github.com/lib/pq v1.12.3
	github.com/nats-io/nats.go v1.52.0
	github.com/prometheus/client_golang v1.23.2
	github.com/redis/go-redis/v9 v9.19.0
	github.com/stretchr/testify v1.11.1

	go.opentelemetry.io/otel v1.43.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.43.0
	go.opentelemetry.io/otel/sdk v1.43.0

	go.uber.org/zap v1.28.0

	golang.org/x/crypto v0.50.0

	google.golang.org/grpc v1.81.0
	google.golang.org/protobuf v1.36.11
)