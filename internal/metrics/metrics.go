package metrics

import (
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	once sync.Once

	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests by route, method, and status",
		},
		[]string{"route", "method", "status"},
	)

	httpRequestDurationSeconds = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"route", "method", "status"},
	)
)

// ObserveHTTP records a single HTTP request metric set.
func ObserveHTTP(route, method, status string, duration time.Duration) {
	once.Do(func() {
		// ensure collectors are created; promauto handles registration
	})
	httpRequestsTotal.WithLabelValues(route, method, status).Inc()
	httpRequestDurationSeconds.WithLabelValues(route, method, status).Observe(duration.Seconds())
}
