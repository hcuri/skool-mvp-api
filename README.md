# Skool MVP API

Small Go HTTP API for a Skool-inspired MVP. It will evolve to run on AWS EKS with Postgres, managed via Terraform and GitOps. Today it uses an in-memory store to keep the surface minimal while wiring the HTTP layer.

## Running locally

```bash
go run ./cmd/api
```

Environment variables:
- `PORT` (default `8080`)
- `LOG_LEVEL` (default `info`)

Swagger UI: http://localhost:8080/swagger

## Endpoints

- `GET /healthz` – health check
- `GET /communities` – list communities
- `POST /communities` – create a community
- `GET /communities/{id}/posts` – list posts within a community
- `POST /communities/{id}/posts` – create a post within a community

## Running with Postgres

1) Start Postgres locally:

```bash
docker-compose up -d db
```

2) Run the API pointing to Postgres:

```bash
DATABASE_URL="postgres://skool:skool_pass@localhost:5432/skool_mvp?sslmode=disable" go run ./cmd/api
```

Or with Dockerized API:

```bash
DATABASE_URL="postgres://skool:skool_pass@host.docker.internal:5432/skool_mvp?sslmode=disable" docker run --rm -p 8080:8080 -e DATABASE_URL skool-mvp-app:local
```

If `DATABASE_URL` is unset, the service uses the in-memory store.

## Logging

Structured request logs are emitted with `zap` (method, path, status, bytes, duration). Adjust verbosity via `LOG_LEVEL`.
