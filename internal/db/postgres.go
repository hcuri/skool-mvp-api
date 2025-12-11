package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
)

// PostgresStore implements Store backed by PostgreSQL.
type PostgresStore struct {
	db     *sql.DB
	logger *zap.Logger
}

// NewPostgresStore initializes a Postgres-backed store and ensures schema exists.
func NewPostgresStore(ctx context.Context, dsn string, logger *zap.Logger) (Store, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(30 * time.Minute)

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("ping db: %w", err)
	}

	store := &PostgresStore{db: db, logger: logger}
	if err := store.initSchema(ctx); err != nil {
		return nil, fmt.Errorf("init schema: %w", err)
	}

	return store, nil
}

func (s *PostgresStore) initSchema(ctx context.Context) error {
	stmts := []string{
		`CREATE TABLE IF NOT EXISTS communities (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			description TEXT DEFAULT ''
		);`,
		`CREATE TABLE IF NOT EXISTS posts (
			id TEXT PRIMARY KEY,
			community_id TEXT NOT NULL REFERENCES communities(id) ON DELETE CASCADE,
			author_id TEXT,
			title TEXT NOT NULL,
			content TEXT NOT NULL,
			created_at TIMESTAMPTZ NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS users (
			id TEXT PRIMARY KEY,
			email TEXT NOT NULL,
			name TEXT NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS community_memberships (
			community_id TEXT NOT NULL REFERENCES communities(id) ON DELETE CASCADE,
			user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			PRIMARY KEY (community_id, user_id)
		);`,
	}

	for _, stmt := range stmts {
		if _, err := s.db.ExecContext(ctx, stmt); err != nil {
			return err
		}
	}
	return nil
}

func (s *PostgresStore) ListCommunities(ctx context.Context) ([]Community, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, name, description FROM communities ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var communities []Community
	for rows.Next() {
		var c Community
		if err := rows.Scan(&c.ID, &c.Name, &c.Description); err != nil {
			return nil, err
		}
		communities = append(communities, c)
	}
	return communities, rows.Err()
}

func (s *PostgresStore) CreateCommunity(ctx context.Context, input CommunityInput) (Community, error) {
	if input.Name == "" {
		return Community{}, fmt.Errorf("name is required")
	}
	community := Community{
		ID:          newID(),
		Name:        input.Name,
		Description: input.Description,
	}

	_, err := s.db.ExecContext(ctx,
		`INSERT INTO communities (id, name, description) VALUES ($1, $2, $3)`,
		community.ID, community.Name, community.Description)
	if err != nil {
		return Community{}, err
	}
	return community, nil
}

func (s *PostgresStore) ListPostsByCommunity(ctx context.Context, communityID string) ([]Post, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, community_id, author_id, title, content, created_at FROM posts WHERE community_id = $1 ORDER BY created_at DESC`, communityID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var p Post
		if err := rows.Scan(&p.ID, &p.CommunityID, &p.AuthorID, &p.Title, &p.Content, &p.CreatedAt); err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(posts) == 0 {
		// Check community existence to mirror in-memory behavior.
		var exists bool
		if err := s.db.QueryRowContext(ctx, `SELECT EXISTS (SELECT 1 FROM communities WHERE id = $1)`, communityID).Scan(&exists); err != nil {
			return nil, err
		}
		if !exists {
			return nil, ErrCommunityNotFound
		}
	}

	return posts, nil
}

func (s *PostgresStore) CreatePost(ctx context.Context, communityID string, input PostInput) (Post, error) {
	if input.Title == "" {
		return Post{}, fmt.Errorf("title is required")
	}
	if input.Content == "" {
		return Post{}, fmt.Errorf("content is required")
	}

	post := Post{
		ID:          newID(),
		CommunityID: communityID,
		AuthorID:    input.AuthorID,
		Title:       input.Title,
		Content:     input.Content,
		CreatedAt:   time.Now().UTC(),
	}

	_, err := s.db.ExecContext(ctx,
		`INSERT INTO posts (id, community_id, author_id, title, content, created_at) VALUES ($1, $2, $3, $4, $5, $6)`,
		post.ID, post.CommunityID, post.AuthorID, post.Title, post.Content, post.CreatedAt)
	if err != nil {
		if isForeignKeyViolation(err) {
			return Post{}, ErrCommunityNotFound
		}
		return Post{}, err
	}

	return post, nil
}

func isForeignKeyViolation(err error) bool {
	var pqErr interface{ SQLState() string }
	if errors.As(err, &pqErr) {
		// Postgres foreign key violation code.
		return pqErr.SQLState() == "23503"
	}
	return false
}
