package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"

	_ "github.com/lib/pq"

	authv1 "github.com/Yessenchik/bydz-fitness-app/services/auth-service/gen/proto/auth/v1"

	"github.com/Yessenchik/bydz-fitness-app/services/auth-service/internal/config"
	grpcHandler "github.com/Yessenchik/bydz-fitness-app/services/auth-service/internal/delivery/grpc"
	jwtManager "github.com/Yessenchik/bydz-fitness-app/services/auth-service/internal/infrastructure/jwt"
	natsPublisher "github.com/Yessenchik/bydz-fitness-app/services/auth-service/internal/infrastructure/nats"
	redisClient "github.com/Yessenchik/bydz-fitness-app/services/auth-service/internal/infrastructure/redis"
	"github.com/Yessenchik/bydz-fitness-app/services/auth-service/internal/logger"
	"github.com/Yessenchik/bydz-fitness-app/services/auth-service/internal/metrics"
	authMiddleware "github.com/Yessenchik/bydz-fitness-app/services/auth-service/internal/middleware"
	postgresRepo "github.com/Yessenchik/bydz-fitness-app/services/auth-service/internal/repository/postgres"
	"github.com/Yessenchik/bydz-fitness-app/services/auth-service/internal/usecase"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	cfg := config.Load()

	zapLogger, err := logger.New()
	if err != nil {
		log.Fatal(err)
	}
	defer zapLogger.Sync()

	metrics.Register()

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		log.Println("metrics server running on :9101")

		if err := http.ListenAndServe(":9101", nil); err != nil {
			log.Fatal(err)
		}
	}()

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.PostgresHost,
		cfg.PostgresPort,
		cfg.PostgresUser,
		cfg.PostgresPassword,
		cfg.AuthDB,
		cfg.PostgresSSLMode,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	userRepo := postgresRepo.NewUserRepository(db)

	jwtMgr := jwtManager.NewJWTManager(cfg.JWTSecret)

	redisAddr := fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort)
	redis := redisClient.NewClient(redisAddr)

	natsURL := fmt.Sprintf("nats://%s:%s", cfg.NATSHost, cfg.NATSPort)
	publisher, err := natsPublisher.NewPublisher(natsURL)
	if err != nil {
		log.Fatal(err)
	}
	defer publisher.Close()

	authUC := usecase.NewAuthUsecase(
		userRepo,
		jwtMgr,
		redis,
		publisher,
	)

	handler := grpcHandler.NewHandler(authUC)

	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			authMiddleware.LoggingInterceptor(zapLogger),
			authMiddleware.MetricsInterceptor(),
			authMiddleware.AuthInterceptor(jwtMgr),
		),
	)

	authv1.RegisterAuthServiceServer(server, handler)

	reflection.Register(server)

	lis, err := net.Listen("tcp", ":"+cfg.AuthGRPCPort)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("auth-service running on :%s\n", cfg.AuthGRPCPort)

	if err := server.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
