package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	// ErrorCounter counts error events with service and type labels
	ErrorCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "app_errors_total",
			Help: "Total number of errors by service and type",
		},
		[]string{"service", "type"},
	)

	// RequestCounter counts gRPC method calls with service and method labels
	RequestCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "grpc_requests_total",
			Help: "Total number of gRPC requests by service and method",
		},
		[]string{"service", "method"},
	)
)

// Init registers all metrics with Prometheus
func Init() {
	// Register metrics with Prometheus
	prometheus.MustRegister(ErrorCounter)
	prometheus.MustRegister(RequestCounter)
}

// IncrementErrorCounter increments the error counter for the specified service and error type
func IncrementErrorCounter(service, errorType string) {
	ErrorCounter.WithLabelValues(service, errorType).Inc()
}

// IncrementRequestCounter increments the request counter for the specified service and method
func IncrementRequestCounter(service, method string) {
	RequestCounter.WithLabelValues(service, method).Inc()
}
