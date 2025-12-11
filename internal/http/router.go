package apihttp

import (
	"net/http"

	"github.com/go-chi/chi/v5"
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
	r.Use(requestLogger(h.logger))

	r.Get("/healthz", h.Healthz)

	r.Route("/communities", func(r chi.Router) {
		r.Get("/", h.ListCommunities)
		r.Post("/", h.CreateCommunity)

		r.Route("/{id}/posts", func(r chi.Router) {
			r.Get("/", h.ListPosts)
			r.Post("/", h.CreatePost)
		})
	})

	r.Get("/swagger", swaggerUIHandler)
	r.Get("/swagger/openapi.yaml", openAPISpecHandler)

	return r
}
