package apihttp

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"

	"github.com/hcuri/skool-mvp-app/internal/db"
)

// Handler bundles dependencies for HTTP handlers.
type Handler struct {
	store  db.Store
	logger *zap.Logger
}

// NewRouter wires routes to handlers and returns an http.Handler.
func NewRouter(store db.Store, logger *zap.Logger) http.Handler {
	h := &Handler{
		store:  store,
		logger: logger,
	}

	r := chi.NewRouter()
	r.Use(metricsMiddleware)
	r.Use(requestLogger(h.logger))

	r.Get("/healthz", h.Healthz)
	r.Head("/healthz", h.Healthz)
	r.Get("/metrics", promhttp.Handler().ServeHTTP)

	r.Route("/communities", func(r chi.Router) {
		r.Get("/", h.ListCommunities)
		r.Post("/", h.CreateCommunity)
		r.Delete("/{id}", h.DeleteCommunity)

		r.Route("/{id}/posts", func(r chi.Router) {
			r.Get("/", h.ListPosts)
			r.Post("/", h.CreatePost)
			r.Delete("/{postId}", h.DeletePost)
		})
	})

	r.Get("/swagger", swaggerUIHandler)
	r.Get("/swagger/openapi.yaml", openAPISpecHandler)

	return r
}
