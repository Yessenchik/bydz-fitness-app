package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv string

	AuthGRPCPort string

	PostgresHost     string
	PostgresPort     string
	PostgresUser     string
	PostgresPassword string
	PostgresSSLMode  string

	AuthDB string

	RedisHost string
	RedisPort string

	NATSHost string
	NATSPort string

	JWTSecret string
}

func Load() *Config {

	env := os.Getenv("APP_ENV")

	switch env {

	case "docker":
		_ = godotenv.Load(".env.docker")

	case "production":
		_ = godotenv.Load(".env.production")

	default:
		_ = godotenv.Load(".env.development")
	}

	return &Config{
		AppEnv: os.Getenv("APP_ENV"),

		AuthGRPCPort: os.Getenv("AUTH_GRPC_PORT"),

		PostgresHost:     os.Getenv("POSTGRES_HOST"),
		PostgresPort:     os.Getenv("POSTGRES_PORT"),
		PostgresUser:     os.Getenv("POSTGRES_USER"),
		PostgresPassword: os.Getenv("POSTGRES_PASSWORD"),
		PostgresSSLMode:  os.Getenv("POSTGRES_SSL_MODE"),

		AuthDB: os.Getenv("AUTH_DB"),

		RedisHost: os.Getenv("REDIS_HOST"),
		RedisPort: os.Getenv("REDIS_PORT"),

		NATSHost: os.Getenv("NATS_HOST"),
		NATSPort: os.Getenv("NATS_PORT"),

		JWTSecret: os.Getenv("JWT_SECRET"),
	}
}
