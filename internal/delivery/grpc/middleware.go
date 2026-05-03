package grpc

import (
	"context"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/Yessenchik/bydz-fitness-app/user-auth-service/internal/domain"
	"github.com/Yessenchik/bydz-fitness-app/user-auth-service/pkg/metrics"
)

// UnaryLoggingInterceptor — логирует каждый gRPC-вызов с latency и статусом
func UnaryLoggingInterceptor(logger *zap.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		start := time.Now()
		resp, err := handler(ctx, req)
		duration := time.Since(start)

		code := codes.OK
		if err != nil {
			code = status.Code(err)
		}

		logger.Info("grpc call",
			zap.String("method", info.FullMethod),
			zap.Duration("duration", duration),
			zap.String("code", code.String()),
			zap.Error(err),
		)

		// Prometheus — histogram латентности
		metrics.GRPCRequestDuration.
			WithLabelValues(info.FullMethod, code.String()).
			Observe(duration.Seconds())

		return resp, err
	}
}

// UnaryTracingInterceptor — создаёт OpenTelemetry span для каждого RPC
func UnaryTracingInterceptor() grpc.UnaryServerInterceptor {
	tracer := otel.Tracer("user-auth-service/grpc")
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		ctx, span := tracer.Start(ctx, info.FullMethod)
		defer span.End()
		span.SetAttributes(attribute.String("rpc.method", info.FullMethod))
		return handler(ctx, req)
	}
}

// UnaryAuthInterceptor — проверяет JWT для защищённых методов.
// Публичные методы пропускаются без проверки токена.
func UnaryAuthInterceptor(tokens domain.TokenManager, cache domain.CacheRepository) grpc.UnaryServerInterceptor {
	// Методы, доступные без токена
	publicMethods := map[string]struct{}{
		"/userauth.v1.UserAuthService/Register":             {},
		"/userauth.v1.UserAuthService/Login":                {},
		"/userauth.v1.UserAuthService/RefreshToken":         {},
		"/userauth.v1.UserAuthService/VerifyEmail":          {},
		"/userauth.v1.UserAuthService/RequestPasswordReset": {},
		"/userauth.v1.UserAuthService/ConfirmPasswordReset": {},
	}

	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		if _, ok := publicMethods[info.FullMethod]; ok {
			return handler(ctx, req)
		}

		// Извлекаем токен из metadata
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "missing metadata")
		}
		authHeader := md.Get("authorization")
		if len(authHeader) == 0 {
			return nil, status.Error(codes.Unauthenticated, "missing authorization header")
		}

		tokenStr := authHeader[0]
		if len(tokenStr) > 7 && tokenStr[:7] == "Bearer " {
			tokenStr = tokenStr[7:]
		}

		claims, err := tokens.ParseAccessToken(tokenStr)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, "invalid token")
		}

		// Проверяем блэклист в Redis
		blacklisted, err := cache.IsTokenBlacklisted(ctx, claims.TokenID)
		if err == nil && blacklisted {
			return nil, status.Error(codes.Unauthenticated, "token revoked")
		}

		// Прокидываем claims в контекст для использования в handler
		ctx = context.WithValue(ctx, ctxKeyUserID{}, claims.UserID)
		ctx = context.WithValue(ctx, ctxKeyUserRole{}, claims.Role)

		return handler(ctx, req)
	}
}

type ctxKeyUserID struct{}
type ctxKeyUserRole struct{}
