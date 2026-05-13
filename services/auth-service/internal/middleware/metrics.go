package middleware

import (
	"context"

	"github.com/Yessenchik/bydz-fitness-app/services/auth-service/internal/metrics"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func MetricsInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		resp, err := handler(ctx, req)

		code := status.Code(err).String()

		metrics.GRPCRequestsTotal.WithLabelValues(
			info.FullMethod,
			code,
		).Inc()

		return resp, err
	}
}
