package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	GRPCRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "membership_service_grpc_requests_total",
			Help: "Total number of gRPC requests handled by membership-service",
		},
		[]string{"method", "status"},
	)
)

func Register() {
	prometheus.MustRegister(GRPCRequestsTotal)
}
