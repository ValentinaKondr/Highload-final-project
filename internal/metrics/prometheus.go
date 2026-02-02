package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// RequestTotal counts total requests
	RequestTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "endpoint", "status"},
	)

	// RequestDuration tracks request latency
	RequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint"},
	)

	// RPSRate tracks requests per second
	RPSRate = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "rps_rate",
			Help: "Current requests per second rate",
		},
	)

	// AnomalyCount counts detected anomalies
	AnomalyCount = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "anomalies_detected_total",
			Help: "Total number of anomalies detected",
		},
	)

	// AnomalyRate tracks anomaly rate per minute
	AnomalyRate = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "anomaly_rate_per_minute",
			Help: "Current anomaly rate per minute",
		},
	)

	// RollingAverageValue tracks rolling average
	RollingAverageValue = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "rolling_average_value",
			Help: "Current rolling average value",
		},
	)

	// CPUMetric tracks CPU usage
	CPUMetric = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "cpu_usage_percent",
			Help: "CPU usage percentage",
		},
	)
)

