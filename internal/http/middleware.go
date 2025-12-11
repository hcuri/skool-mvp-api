package apihttp

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

type responseRecorder struct {
	http.ResponseWriter
	status int
	bytes  int64
}

func (rw *responseRecorder) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseRecorder) Write(b []byte) (int, error) {
	if rw.status == 0 {
		rw.status = http.StatusOK
	}
	n, err := rw.ResponseWriter.Write(b)
	rw.bytes += int64(n)
	return n, err
}

func requestLogger(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			rr := &responseRecorder{ResponseWriter: w}

			next.ServeHTTP(rr, r)

			status := rr.status
			if status == 0 {
				status = http.StatusOK
			}

			logger.Info("http_request",
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.Int("status", status),
				zap.Int64("bytes", rr.bytes),
				zap.Duration("duration", time.Since(start)),
			)
		})
	}
}
