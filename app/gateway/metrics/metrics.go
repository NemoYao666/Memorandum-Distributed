package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	GateWayRequestCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "gateway_http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"service", "method", "endpoint", "status"},
	)
	GateWayRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "gateway_http_request_duration_seconds",
			Help:    "Histogram of HTTP request durations",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"service", "method", "endpoint"},
	)
)

func init() {
	prometheus.MustRegister(GateWayRequestCounter)
	prometheus.MustRegister(GateWayRequestDuration)
}
