package apihttp

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/hcuri/skool-mvp-app/internal/db"
)

func (h *Handler) Healthz(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Handler) ListCommunities(w http.ResponseWriter, r *http.Request) {
	communities, err := h.store.ListCommunities(r.Context())
	if err != nil {
		h.logger.Error("list communities failed", zap.Error(err))
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

func (h *Handler) DeleteCommunity(w http.ResponseWriter, r *http.Request) {
	communityID := chi.URLParam(r, "id")
	if communityID == "" {
		http.Error(w, "community id required", http.StatusBadRequest)
		return
	}
	if err := h.store.DeleteCommunity(r.Context(), communityID); err != nil {
		if errors.Is(err, db.ErrCommunityNotFound) {
			http.Error(w, "community not found", http.StatusNotFound)
			return
		}
		h.logger.Error("delete community failed", zap.Error(err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) ListPosts(w http.ResponseWriter, r *http.Request) {
	communityID := chi.URLParam(r, "id")
	posts, err := h.store.ListPostsByCommunity(r.Context(), communityID)
	if err != nil {
		if errors.Is(err, db.ErrCommunityNotFound) {
			http.Error(w, "community not found", http.StatusNotFound)
			return
		}
		h.logger.Error("list posts failed", zap.Error(err))
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

func (h *Handler) DeletePost(w http.ResponseWriter, r *http.Request) {
	communityID := chi.URLParam(r, "id")
	postID := chi.URLParam(r, "postId")
	if communityID == "" || postID == "" {
		http.Error(w, "community id and post id required", http.StatusBadRequest)
		return
	}

	if err := h.store.DeletePost(r.Context(), communityID, postID); err != nil {
		if errors.Is(err, db.ErrCommunityNotFound) {
			http.Error(w, "community not found", http.StatusNotFound)
			return
		}
		if errors.Is(err, db.ErrPostNotFound) {
			http.Error(w, "post not found", http.StatusNotFound)
			return
		}
		h.logger.Error("delete post failed", zap.Error(err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		zap.L().Error("write json response failed", zap.Error(err))
	}
}
