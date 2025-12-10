package apihttp

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/hcuri/skool-mvp-app/internal/db"
)

func (h *Handler) Healthz(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Handler) ListCommunities(w http.ResponseWriter, r *http.Request) {
	communities, err := h.store.ListCommunities(r.Context())
	if err != nil {
		h.logger.Printf("list communities: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, communities)
}

func (h *Handler) CreateCommunity(w http.ResponseWriter, r *http.Request) {
	var input db.CommunityInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	community, err := h.store.CreateCommunity(r.Context(), input)
	if err != nil {
		http.Error(w, fmt.Sprintf("unable to create community: %v", err), http.StatusBadRequest)
		return
	}

	writeJSON(w, http.StatusCreated, community)
}

func (h *Handler) ListPosts(w http.ResponseWriter, r *http.Request) {
	communityID := chi.URLParam(r, "id")
	posts, err := h.store.ListPostsByCommunity(r.Context(), communityID)
	if err != nil {
		if errors.Is(err, db.ErrCommunityNotFound) {
			http.Error(w, "community not found", http.StatusNotFound)
			return
		}
		h.logger.Printf("list posts: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, posts)
}

func (h *Handler) CreatePost(w http.ResponseWriter, r *http.Request) {
	communityID := chi.URLParam(r, "id")
	var input db.PostInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	post, err := h.store.CreatePost(r.Context(), communityID, input)
	if err != nil {
		if errors.Is(err, db.ErrCommunityNotFound) {
			http.Error(w, "community not found", http.StatusNotFound)
			return
		}
		http.Error(w, fmt.Sprintf("unable to create post: %v", err), http.StatusBadRequest)
		return
	}

	writeJSON(w, http.StatusCreated, post)
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("write json response: %v", err)
	}
}
