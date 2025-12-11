package apihttp

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hcuri/skool-mvp-app/internal/db"
)

func newTestServer() http.Handler {
	store := db.NewInMemoryStore()
	logger := log.New(io.Discard, "", 0)
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
	ts := newTestServer()
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

func TestCommunitiesAndPosts(t *testing.T) {
	ts := newTestServer()

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
}
