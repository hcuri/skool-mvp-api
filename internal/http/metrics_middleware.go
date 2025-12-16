package apihttp

import (
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/hcuri/skool-mvp-app/internal/metrics"
)

// metricsMiddleware captures request metrics for Prometheus.
func metricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		next.ServeHTTP(ww, r)

		route := "unknown"
		if rctx := chi.RouteContext(r.Context()); rctx != nil {
			if pattern := rctx.RoutePattern(); pattern != "" {
				route = pattern
			}
		}

		status := ww.Status()
		if status == 0 {
			status = http.StatusOK
		}

		metrics.ObserveHTTP(route, r.Method, strconv.Itoa(status), time.Since(start))
	})
}
