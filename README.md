# Skool MVP API

Small Go HTTP API for a Skool-inspired MVP. It will evolve to run on AWS EKS with Postgres, managed via Terraform and GitOps. Today it uses an in-memory store to keep the surface minimal while wiring the HTTP layer.

## Running locally

```bash
go run ./cmd/api
```

Environment variables:
- `PORT` (default `8080`)
- `LOG_LEVEL` (default `info`)

## Endpoints

- `GET /healthz` – health check
- `GET /communities` – list communities
- `POST /communities` – create a community
- `GET /communities/{id}/posts` – list posts within a community
- `POST /communities/{id}/posts` – create a post within a community
