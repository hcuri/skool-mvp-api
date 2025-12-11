# Skool MVP API – Architecture

## Overview
Minimal Go HTTP API that exposes Skool-style entities with an in-memory store by default and optional Postgres backing. Focus is on a small surface area to support infra/GitOps demos.

## Endpoints
- `GET /healthz`
- `GET /communities`
- `POST /communities`
- `GET /communities/{id}/posts`
- `POST /communities/{id}/posts`

## Data model
- `users`: id, email, name (not yet used in handlers but present in schema)
- `communities`: id, name, description
- `community_memberships`: community_id, user_id
- `posts`: id, community_id, author_id, title, content, created_at

## Storage
- Default: in-memory store (thread-safe maps).
- Optional: Postgres store (`DATABASE_URL`) auto-creates tables on startup and enforces FK between posts and communities.

## Runtime
- Config via env (`PORT`, `LOG_LEVEL`, optional `DATABASE_URL`).
- Router: chi with structured zap request logging middleware.
- Swagger UI at `/swagger` serving embedded `docs/openapi.yaml`.
