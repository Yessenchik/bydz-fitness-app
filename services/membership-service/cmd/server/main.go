package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"

	membershipv1 "github.com/Yessenchik/bydz-fitness-app/services/membership-service/gen/proto/membership/v1"
	"github.com/Yessenchik/bydz-fitness-app/services/membership-service/internal/config"
	grpcHandler "github.com/Yessenchik/bydz-fitness-app/services/membership-service/internal/delivery/grpc"
	natsPublisher "github.com/Yessenchik/bydz-fitness-app/services/membership-service/internal/infrastructure/nats"
	"github.com/Yessenchik/bydz-fitness-app/services/membership-service/internal/logger"
	"github.com/Yessenchik/bydz-fitness-app/services/membership-service/internal/metrics"
	middleware "github.com/Yessenchik/bydz-fitness-app/services/membership-service/internal/middleware"
	postgresRepo "github.com/Yessenchik/bydz-fitness-app/services/membership-service/internal/repository/postgres"
	"github.com/Yessenchik/bydz-fitness-app/services/membership-service/internal/usecase"

	_ "github.com/lib/pq"

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
		log.Printf("membership-service metrics running on :%s\n", cfg.MembershipMetricsPort)

		if err := http.ListenAndServe(":"+cfg.MembershipMetricsPort, nil); err != nil {
			log.Fatal(err)
		}
	}()

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.PostgresHost,
		cfg.PostgresPort,
		cfg.PostgresUser,
		cfg.PostgresPassword,
		cfg.MembershipDB,
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

	planRepo := postgresRepo.NewPlanRepository(db)
	subRepo := postgresRepo.NewSubscriptionRepository(db)

	natsURL := fmt.Sprintf("nats://%s:%s", cfg.NATSHost, cfg.NATSPort)

	publisher, err := natsPublisher.NewPublisher(natsURL)
	if err != nil {
		log.Fatal(err)
	}
	defer publisher.Close()

	uc := usecase.NewMembershipUsecase(
		planRepo,
		subRepo,
		publisher,
	)

	handler := grpcHandler.NewHandler(uc)

	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			middleware.LoggingInterceptor(zapLogger),
			middleware.MetricsInterceptor(),
		),
	)

	membershipv1.RegisterMembershipServiceServer(server, handler)
	reflection.Register(server)

	lis, err := net.Listen("tcp", ":"+cfg.MembershipGRPCPort)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("membership-service running on :%s\n", cfg.MembershipGRPCPort)

	if err := server.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
