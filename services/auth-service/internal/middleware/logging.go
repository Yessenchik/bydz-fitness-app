package middleware

import (
	"context"
	"time"

	"go.uber.org/zap"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func LoggingInterceptor(
	logger *zap.Logger,
) grpc.UnaryServerInterceptor {

	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {

		start := time.Now()

		resp, err := handler(ctx, req)

		duration := time.Since(start)

		code := status.Code(err).String()

		logger.Info(
			"grpc request",
			zap.String("method", info.FullMethod),
			zap.String("status", code),
			zap.Duration("duration", duration),
		)

		return resp, err
	}
}
