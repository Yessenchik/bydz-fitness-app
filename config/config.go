package config

import (
	"os"
	"time"

	"github.com/nats-io/nats.go"
)

type Config struct {
	GRPC     GRPCConfig
	Postgres PostgresConfig
	Redis    RedisConfig
	NATS     NATSConfig
	SMTP     SMTPConfig
	JWT      JWTConfig
	App      AppConfig
}

type GRPCConfig struct{ Port string }
type PostgresConfig struct{ DSN string }
type RedisConfig struct {
	Addr, Password string
	DB             int
}
type NATSConfig struct{ URL string }
type JWTConfig struct{ AccessSecret, RefreshSecret string }
type AppConfig struct {
	BaseURL      string
	IsDev        bool
	OTLPEndpoint string
}
type SMTPConfig struct {
	Host, Port, Username, Password string
}

func Load() *Config {
	return &Config{
		GRPC:     GRPCConfig{Port: getEnv("GRPC_PORT", ":50051")},
		Postgres: PostgresConfig{DSN: mustEnv("POSTGRES_DSN")},
		Redis: RedisConfig{
			Addr:     getEnv("REDIS_ADDR", "localhost:6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
		},
		NATS: NATSConfig{URL: getEnv("NATS_URL", nats.DefaultURL)},
		JWT: JWTConfig{
			AccessSecret:  mustEnv("JWT_ACCESS_SECRET"),
			RefreshSecret: mustEnv("JWT_REFRESH_SECRET"),
		},
		SMTP: SMTPConfig{
			Host: getEnv("SMTP_HOST", "smtp.gmail.com"),
			Port: getEnv("SMTP_PORT", "587"),
			// ЖЕСТКО ПРОПИСЫВАЕМ ДАННЫЕ БЕЗ getEnv:
			Username: "bekbaulykalymzan@gmail.com",
			Password: "adewiwbkzdvcefqj",
		},
		App: AppConfig{
			BaseURL:      getEnv("APP_BASE_URL", "http://localhost:8080"),
			IsDev:        getEnv("APP_ENV", "production") == "development",
			OTLPEndpoint: getEnv("OTLP_ENDPOINT", "localhost:4317"),
		},
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func mustEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		panic("required env variable not set: " + key)
	}
	return v
}

// Suppress unused import (time used in future health-check timeout)
var _ = time.Second
