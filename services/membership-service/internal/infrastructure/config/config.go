package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	AppEnv string

	GRPCPort string

	PostgresHost     string
	PostgresPort     string
	PostgresUser     string
	PostgresPassword string
	PostgresDB       string
	PostgresSSLMode  string

	RedisHost string
	RedisPort string
	RedisDB   int

	NATSURL string

	JWTSecret string

	SMTPHost     string
	SMTPPort     string
	SMTPUsername string
	SMTPPassword string
	SMTPFrom     string

	RequestTimeout time.Duration
}

func Load() Config {
	return Config{
		AppEnv: getEnv("APP_ENV", "development"),

		GRPCPort: getEnv("GRPC_PORT", "50051"),

		PostgresHost:     getEnv("POSTGRES_HOST", "localhost"),
		PostgresPort:     getEnv("POSTGRES_PORT", "5432"),
		PostgresUser:     getEnv("POSTGRES_USER", "gym"),
		PostgresPassword: getEnv("POSTGRES_PASSWORD", "gym_password"),
		PostgresDB:       getEnv("POSTGRES_DB", "membership_db"),
		PostgresSSLMode:  getEnv("POSTGRES_SSL_MODE", "disable"),

		RedisHost: getEnv("REDIS_HOST", "localhost"),
		RedisPort: getEnv("REDIS_PORT", "6379"),
		RedisDB:   getEnvAsInt("REDIS_DB", 0),

		NATSURL: getEnv("NATS_URL", "nats://localhost:4222"),

		JWTSecret: getEnv("JWT_SECRET", "super-secret-key"),

		SMTPHost:     getEnv("SMTP_HOST", "smtp.gmail.com"),
		SMTPPort:     getEnv("SMTP_PORT", "587"),
		SMTPUsername: getEnv("SMTP_USERNAME", ""),
		SMTPPassword: getEnv("SMTP_PASSWORD", ""),
		SMTPFrom:     getEnv("SMTP_FROM", ""),

		RequestTimeout: time.Duration(
			getEnvAsInt("REQUEST_TIMEOUT_SECONDS", 5),
		) * time.Second,
	}
}

func (c Config) PostgresDSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.PostgresHost,
		c.PostgresPort,
		c.PostgresUser,
		c.PostgresPassword,
		c.PostgresDB,
		c.PostgresSSLMode,
	)
}

func (c Config) RedisAddr() string {
	return fmt.Sprintf("%s:%s", c.RedisHost, c.RedisPort)
}

func (c Config) IsProduction() bool {
	return c.AppEnv == "production"
}

func getEnv(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}

func getEnvAsInt(key string, fallback int) int {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	parsedValue, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}

	return parsedValue
}
