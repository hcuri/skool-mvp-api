package db

import "time"

// User represents a user within the system.
type User struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

// Community represents a community that users can post to.
type Community struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// Post represents a message authored by a user within a community.
type Post struct {
	ID          string    `json:"id"`
	CommunityID string    `json:"communityId"`
	AuthorID    string    `json:"authorId"`
	Title       string    `json:"title"`
	Content     string    `json:"content"`
	CreatedAt   time.Time `json:"createdAt"`
}

// CommunityInput captures the fields needed to create a community.
type CommunityInput struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// PostInput captures the fields needed to create a post.
type PostInput struct {
	AuthorID string `json:"authorId"`
	Title    string `json:"title"`
	Content  string `json:"content"`
}
