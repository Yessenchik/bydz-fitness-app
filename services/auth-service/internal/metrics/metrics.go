package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	GRPCRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "auth_service_grpc_requests_total",
			Help: "Total number of gRPC requests handled by auth-service",
		},
		[]string{"method", "status"},
	)
)

func Register() {
	prometheus.MustRegister(GRPCRequestsTotal)
}
