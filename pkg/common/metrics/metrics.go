package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
)

type Metrics struct {
	Registry *prometheus.Registry

	HTTPRequestsTotal      *prometheus.CounterVec
	HTTPRequestDurationSec *prometheus.HistogramVec
}

func New() *Metrics {
	registry := prometheus.NewRegistry()
	registry.MustRegister(
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
	)

	httpRequestsTotal := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of API requests by operation and status.",
		},
		[]string{"operation", "status"},
	)
	httpRequestDurationSec := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "API request duration in seconds by operation.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"operation"},
	)
	registry.MustRegister(httpRequestsTotal, httpRequestDurationSec)

	return &Metrics{
		Registry:               registry,
		HTTPRequestsTotal:      httpRequestsTotal,
		HTTPRequestDurationSec: httpRequestDurationSec,
	}
}
