package apihttp

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hcuri/skool-mvp-app/internal/db"
	"go.uber.org/zap/zaptest"
)

func newTestServer(t *testing.T) http.Handler {
	store := db.NewInMemoryStore()
	logger := zaptest.NewLogger(t)
	return NewRouter(store, logger)
}

func decodeResponse[T any](tb testing.TB, body []byte) T {
	tb.Helper()
	var out T
	if err := json.Unmarshal(body, &out); err != nil {
		tb.Fatalf("decode response: %v", err)
	}
	return out
}

func TestHealthz(t *testing.T) {
	ts := newTestServer(t)
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rr := httptest.NewRecorder()

	ts.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	body := decodeResponse[map[string]string](t, rr.Body.Bytes())
	if body["status"] != "ok" {
		t.Fatalf("unexpected body: %v", body)
	}
}

func TestHealthzHead(t *testing.T) {
	ts := newTestServer(t)
	req := httptest.NewRequest(http.MethodHead, "/healthz", nil)
	rr := httptest.NewRecorder()

	ts.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
}

func TestCommunitiesAndPosts(t *testing.T) {
	ts := newTestServer(t)

	// Create a community.
	createBody := bytes.NewBufferString(`{"name":"Go","description":"golang"}`)
	req := httptest.NewRequest(http.MethodPost, "/communities", createBody)
	rr := httptest.NewRecorder()
	ts.ServeHTTP(rr, req)
	if rr.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", rr.Code, rr.Body.String())
	}
	community := decodeResponse[db.Community](t, rr.Body.Bytes())

	// List communities.
	req = httptest.NewRequest(http.MethodGet, "/communities", nil)
	rr = httptest.NewRecorder()
	ts.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	list := decodeResponse[[]db.Community](t, rr.Body.Bytes())
	if len(list) != 1 || list[0].ID != community.ID {
		t.Fatalf("unexpected communities: %+v", list)
	}

	// Create a post.
	postBody := bytes.NewBufferString(`{"authorId":"user-1","title":"Hello","content":"World"}`)
	req = httptest.NewRequest(http.MethodPost, "/communities/"+community.ID+"/posts", postBody)
	rr = httptest.NewRecorder()
	ts.ServeHTTP(rr, req)
	if rr.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", rr.Code, rr.Body.String())
	}
	post := decodeResponse[db.Post](t, rr.Body.Bytes())
	if post.CommunityID != community.ID {
		t.Fatalf("expected community ID %s, got %s", community.ID, post.CommunityID)
	}

	// List posts.
	req = httptest.NewRequest(http.MethodGet, "/communities/"+community.ID+"/posts", nil)
	rr = httptest.NewRecorder()
	ts.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	posts := decodeResponse[[]db.Post](t, rr.Body.Bytes())
	if len(posts) != 1 || posts[0].ID != post.ID {
		t.Fatalf("unexpected posts: %+v", posts)
	}

	// Delete post.
	rr = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodDelete, "/communities/"+community.ID+"/posts/"+post.ID, nil)
	ts.ServeHTTP(rr, req)
	if rr.Code != http.StatusNoContent {
		t.Fatalf("expected 204 on delete post, got %d: %s", rr.Code, rr.Body.String())
	}

	// Delete community.
	rr = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodDelete, "/communities/"+community.ID, nil)
	ts.ServeHTTP(rr, req)
	if rr.Code != http.StatusNoContent {
		t.Fatalf("expected 204 on delete community, got %d: %s", rr.Code, rr.Body.String())
	}

	// Ensure community gone.
	rr = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/communities/"+community.ID+"/posts", nil)
	ts.ServeHTTP(rr, req)
	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected 404 after community delete, got %d", rr.Code)
	}
}

func TestCommunityValidationErrors(t *testing.T) {
	ts := newTestServer(t)

	// Invalid JSON payload.
	req := httptest.NewRequest(http.MethodPost, "/communities", bytes.NewBufferString(`{"name":`))
	rr := httptest.NewRecorder()
	ts.ServeHTTP(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for invalid json, got %d", rr.Code)
	}

	// Missing required name.
	req = httptest.NewRequest(http.MethodPost, "/communities", bytes.NewBufferString(`{"description":"desc"}`))
	rr = httptest.NewRecorder()
	ts.ServeHTTP(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for missing name, got %d", rr.Code)
	}
}

func TestPostEdgeCases(t *testing.T) {
	ts := newTestServer(t)

	// Unknown community should 404.
	req := httptest.NewRequest(http.MethodPost, "/communities/missing/posts", bytes.NewBufferString(`{"title":"t","content":"c"}`))
	rr := httptest.NewRecorder()
	ts.ServeHTTP(rr, req)
	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected 404 for missing community, got %d", rr.Code)
	}

	// Invalid JSON payload.
	req = httptest.NewRequest(http.MethodPost, "/communities/missing/posts", bytes.NewBufferString(`{"title":`))
	rr = httptest.NewRecorder()
	ts.ServeHTTP(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for invalid json, got %d", rr.Code)
	}

	// List posts for missing community should 404.
	req = httptest.NewRequest(http.MethodGet, "/communities/missing/posts", nil)
	rr = httptest.NewRecorder()
	ts.ServeHTTP(rr, req)
	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected 404 for missing community posts, got %d", rr.Code)
	}
}
