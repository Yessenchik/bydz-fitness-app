package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	RegistrationAttempts = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "auth_registration_attempts_total",
		Help: "Total registration attempts",
	})
	RegistrationSuccess = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "auth_registration_success_total",
		Help: "Successful registrations",
	})
	LoginAttempts = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "auth_login_attempts_total",
		Help: "Total login attempts",
	})
	LoginSuccess = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "auth_login_success_total",
		Help: "Successful logins",
	})
	LoginFailure = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "auth_login_failure_total",
		Help: "Failed logins (wrong password)",
	})
	CacheHits = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "auth_cache_hits_total",
		Help: "Redis cache hits for user profiles",
	})
	CacheMisses = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "auth_cache_misses_total",
		Help: "Redis cache misses for user profiles",
	})
	GRPCRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "auth_grpc_request_duration_seconds",
			Help:    "gRPC request latency",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)
)

func init() {
	prometheus.MustRegister(
		RegistrationAttempts,
		RegistrationSuccess,
		LoginAttempts,
		LoginSuccess,
		LoginFailure,
		CacheHits,
		CacheMisses,
		GRPCRequestDuration,
	)
}
