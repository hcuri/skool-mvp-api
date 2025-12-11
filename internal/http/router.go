package apihttp

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/hcuri/skool-mvp-app/internal/db"
)

// Handler bundles dependencies for HTTP handlers.
type Handler struct {
	store  db.Store
	logger *log.Logger
}

// NewRouter wires routes to handlers and returns an http.Handler.
func NewRouter(store db.Store, logger *log.Logger) http.Handler {
	h := &Handler{
		store:  store,
		logger: logger,
	}

	r := chi.NewRouter()

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
