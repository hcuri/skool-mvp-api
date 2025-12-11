package db

import (
	"context"
	"testing"
)

func TestInMemoryStoreCommunities(t *testing.T) {
	store := NewInMemoryStore()

	// Validate required fields.
	if _, err := store.CreateCommunity(context.Background(), CommunityInput{}); err == nil {
		t.Fatalf("expected error for missing name")
	}

	created, err := store.CreateCommunity(context.Background(), CommunityInput{
		Name:        "Go Fans",
		Description: "Community for Go enthusiasts",
	})
	if err != nil {
		t.Fatalf("create community: %v", err)
	}
	if created.ID == "" {
		t.Fatalf("expected generated ID")
	}

	list, err := store.ListCommunities(context.Background())
	if err != nil {
		t.Fatalf("list communities: %v", err)
	}
	if len(list) != 1 || list[0].ID != created.ID {
		t.Fatalf("unexpected list result: %+v", list)
	}
}

func TestInMemoryStorePosts(t *testing.T) {
	store := NewInMemoryStore()
	ctx := context.Background()

	community, err := store.CreateCommunity(ctx, CommunityInput{Name: "Tech"})
	if err != nil {
		t.Fatalf("create community: %v", err)
	}

	// Missing required fields.
	if _, err := store.CreatePost(ctx, community.ID, PostInput{}); err == nil {
		t.Fatalf("expected error for missing fields")
	}

	post, err := store.CreatePost(ctx, community.ID, PostInput{
		AuthorID: "user-1",
		Title:    "First",
		Content:  "Hello",
	})
	if err != nil {
		t.Fatalf("create post: %v", err)
	}
	if post.ID == "" {
		t.Fatalf("expected generated post ID")
	}
	if post.CommunityID != community.ID {
		t.Fatalf("expected community ID %s, got %s", community.ID, post.CommunityID)
	}

	// Nonexistent community.
	if _, err := store.CreatePost(ctx, "missing", PostInput{
		Title:   "Oops",
		Content: "No community",
	}); err == nil {
		t.Fatalf("expected error for missing community")
	}

	posts, err := store.ListPostsByCommunity(ctx, community.ID)
	if err != nil {
		t.Fatalf("list posts: %v", err)
	}
	if len(posts) != 1 || posts[0].ID != post.ID {
		t.Fatalf("unexpected posts: %+v", posts)
	}

	// Listing for missing community should error.
	if _, err := store.ListPostsByCommunity(ctx, "missing"); err == nil {
		t.Fatalf("expected error for missing community")
	}
}
