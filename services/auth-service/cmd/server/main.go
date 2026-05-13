package main

import (
	"context"
	"net"
	"net/http"
	"os/signal"
	"syscall"

	"github.com/jmoiron/sqlx"
	"github.com/nats-io/nats.go"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/Yessenchik/bydz-fitness-app/user-auth-service/config"
	pb "github.com/Yessenchik/bydz-fitness-app/user-auth-service/gen/userauth/v1"
	deliveryGRPC "github.com/Yessenchik/bydz-fitness-app/user-auth-service/internal/delivery/grpc"
	infraEmail "github.com/Yessenchik/bydz-fitness-app/user-auth-service/internal/infrastructure/email"
	infraNATS "github.com/Yessenchik/bydz-fitness-app/user-auth-service/internal/infrastructure/nats"
	infraPG "github.com/Yessenchik/bydz-fitness-app/user-auth-service/internal/infrastructure/postgres"
	infraRedis "github.com/Yessenchik/bydz-fitness-app/user-auth-service/internal/infrastructure/redis"
	"github.com/Yessenchik/bydz-fitness-app/user-auth-service/internal/usecase"
	"github.com/Yessenchik/bydz-fitness-app/user-auth-service/pkg/logger"
	"github.com/Yessenchik/bydz-fitness-app/user-auth-service/pkg/tracing"
)

func main() {
	cfg := config.Load()

	// ── Логгер ──────────────────────────────
	log, err := logger.New(cfg.App.IsDev)
	if err != nil {
		panic(err)
	}
	defer log.Sync()

	// ── Tracing ─────────────────────────────
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	tp, err := tracing.InitTracer(ctx, "user-auth-service", cfg.App.OTLPEndpoint)
	if err != nil {
		log.Fatal("init tracer", zap.Error(err))
	}
	defer tp.Shutdown(ctx)

	// ── PostgreSQL ───────────────────────────
	db, err := sqlx.ConnectContext(ctx, "postgres", cfg.Postgres.DSN)
	if err != nil {
		log.Fatal("connect postgres", zap.Error(err))
	}
	defer db.Close()

	// ── Redis ────────────────────────────────
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})
	defer rdb.Close()

	// ── NATS ─────────────────────────────────
	nc, err := nats.Connect(cfg.NATS.URL)
	if err != nil {
		log.Fatal("connect nats", zap.Error(err))
	}
	defer nc.Drain()

	// ── Infrastructure ───────────────────────
	userRepo := infraPG.NewUserRepository(db, log)
	cacheRepo := infraRedis.NewCacheRepository(rdb)
	publisher := infraNATS.NewPublisher(nc, log)
	emailSvc := infraEmail.NewSMTPService(
		cfg.SMTP.Host, cfg.SMTP.Port,
		cfg.SMTP.Username, cfg.SMTP.Password,
		cfg.App.BaseURL, log,
	)
	tokenMgr := infraPG.NewJWTManager(cfg.JWT.AccessSecret, cfg.JWT.RefreshSecret)

	// ── Usecases ─────────────────────────────
	authUC := usecase.NewAuthUsecase(userRepo, cacheRepo, publisher, emailSvc, tokenMgr, log)
	userUC := usecase.NewUserUsecase(userRepo, cacheRepo, publisher, log)

	// ── gRPC Server ──────────────────────────
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			deliveryGRPC.UnaryTracingInterceptor(),
			deliveryGRPC.UnaryLoggingInterceptor(log),
			deliveryGRPC.UnaryAuthInterceptor(tokenMgr, cacheRepo),
		),
	)
	pb.RegisterUserAuthServiceServer(grpcServer, deliveryGRPC.NewHandler(authUC, userUC, log))

	// ── Prometheus HTTP ───────────────────────
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		if err := http.ListenAndServe(":9090", nil); err != nil {
			log.Error("metrics server", zap.Error(err))
		}
	}()

	// ── Listen ───────────────────────────────
	lis, err := net.Listen("tcp", cfg.GRPC.Port)
	if err != nil {
		log.Fatal("listen", zap.Error(err))
	}
	log.Info("gRPC server started", zap.String("addr", cfg.GRPC.Port))

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatal("serve grpc", zap.Error(err))
		}
	}()

	<-ctx.Done()
	log.Info("shutting down...")
	grpcServer.GracefulStop()
}
