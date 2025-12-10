package db

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

var (
	// ErrCommunityNotFound indicates the requested community does not exist.
	ErrCommunityNotFound = errors.New("community not found")
)

// Store defines the persistence contract for the application.
type Store interface {
	ListCommunities(ctx context.Context) ([]Community, error)
	CreateCommunity(ctx context.Context, input CommunityInput) (Community, error)
	ListPostsByCommunity(ctx context.Context, communityID string) ([]Post, error)
	CreatePost(ctx context.Context, communityID string, input PostInput) (Post, error)
}

// InMemoryStore is a simple, concurrency-safe store backed by in-memory maps.
type InMemoryStore struct {
	mu             sync.RWMutex
	communities    map[string]Community
	communityOrder []string
	posts          map[string][]Post
}

// NewInMemoryStore initializes an empty in-memory store.
func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		communities: make(map[string]Community),
		posts:       make(map[string][]Post),
	}
}

func (s *InMemoryStore) ListCommunities(_ context.Context) ([]Community, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	communities := make([]Community, 0, len(s.communities))
	for _, id := range s.communityOrder {
		communities = append(communities, s.communities[id])
	}
	return communities, nil
}

func (s *InMemoryStore) CreateCommunity(_ context.Context, input CommunityInput) (Community, error) {
	if input.Name == "" {
		return Community{}, fmt.Errorf("name is required")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	id := uuid.NewString()
	community := Community{
		ID:          id,
		Name:        input.Name,
		Description: input.Description,
	}
	s.communities[id] = community
	s.communityOrder = append(s.communityOrder, id)

	return community, nil
}

func (s *InMemoryStore) ListPostsByCommunity(_ context.Context, communityID string) ([]Post, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if _, ok := s.communities[communityID]; !ok {
		return nil, ErrCommunityNotFound
	}

	posts := s.posts[communityID]
	out := make([]Post, len(posts))
	copy(out, posts)
	return out, nil
}

func (s *InMemoryStore) CreatePost(_ context.Context, communityID string, input PostInput) (Post, error) {
	if input.Title == "" {
		return Post{}, fmt.Errorf("title is required")
	}
	if input.Content == "" {
		return Post{}, fmt.Errorf("content is required")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.communities[communityID]; !ok {
		return Post{}, ErrCommunityNotFound
	}

	post := Post{
		ID:          uuid.NewString(),
		CommunityID: communityID,
		AuthorID:    input.AuthorID,
		Title:       input.Title,
		Content:     input.Content,
		CreatedAt:   time.Now().UTC(),
	}
	s.posts[communityID] = append(s.posts[communityID], post)

	return post, nil
}
