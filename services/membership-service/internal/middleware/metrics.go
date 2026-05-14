package middleware

import (
	"context"

	"github.com/Yessenchik/bydz-fitness-app/services/membership-service/internal/metrics"

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

		metrics.GRPCRequestsTotal.WithLabelValues(
			info.FullMethod,
			status.Code(err).String(),
		).Inc()

		return resp, err
	}
}
