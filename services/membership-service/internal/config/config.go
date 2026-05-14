package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv string

	MembershipGRPCPort    string
	MembershipDB          string
	MembershipMetricsPort string

	PostgresHost     string
	PostgresPort     string
	PostgresUser     string
	PostgresPassword string
	PostgresSSLMode  string

	NATSHost string
	NATSPort string
}

func Load() *Config {
	_ = godotenv.Load()

	return &Config{
		AppEnv: os.Getenv("APP_ENV"),

		MembershipGRPCPort:    os.Getenv("MEMBERSHIP_GRPC_PORT"),
		MembershipDB:          os.Getenv("MEMBERSHIP_DB"),
		MembershipMetricsPort: os.Getenv("MEMBERSHIP_METRICS_PORT"),

		PostgresHost:     os.Getenv("POSTGRES_HOST"),
		PostgresPort:     os.Getenv("POSTGRES_PORT"),
		PostgresUser:     os.Getenv("POSTGRES_USER"),
		PostgresPassword: os.Getenv("POSTGRES_PASSWORD"),
		PostgresSSLMode:  os.Getenv("POSTGRES_SSL_MODE"),

		NATSHost: os.Getenv("NATS_HOST"),
		NATSPort: os.Getenv("NATS_PORT"),
	}
}
