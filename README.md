# Skool MVP API

Built by **Hector Curi** as a personal project to demonstrate end-to-end DevOps/SRE skills (Go, Docker, Kubernetes, Helm, AWS).

Go HTTP API for a Skool-like community platform, containerized with Docker and deployed to EKS via Helm. This service talks to a Postgres database running in AWS RDS (or local Postgres for development).

## Architecture / Project structure
- `cmd/api/` – main entrypoint.
- `internal/` – config, HTTP handlers, router, DB stores (stateless app; persistence in Postgres).
- `docs/` – OpenAPI spec and architecture notes.
- `charts/skool-mvp-api/` – Helm chart for Deployment/Service.
- `k8s/` – raw Kubernetes manifests (example namespace/Deployment/Service, secret placeholder).
- `docker-compose.yml` – local Postgres for dev; `Dockerfile` / `Makefile` for build/run helpers.

## How this fits into the Skool MVP
- App code + Helm chart live here.
- AWS infra (VPC, EKS, RDS) is provisioned from `skool-mvp-infra` (https://github.com/hcuri/skool-mvp-infra).
- GitOps/ArgoCD config and environment values live in `skool-mvp-gitops` (https://github.com/hcuri/skool-mvp-gitops).
- Images are tagged with Git SHAs; ArgoCD rolls deployments by updating values in the GitOps repo.

## API endpoints
- `GET /healthz` – health check
- `GET /communities` – list communities
- `POST /communities` – create a community
- `GET /communities/{id}/posts` – list posts within a community
- `POST /communities/{id}/posts` – create a post within a community

## Local development / running locally
1) Start Postgres (or use in-memory by omitting `DATABASE_URL`):
   ```bash
   docker-compose up -d db
   ```
2) Run the API:
   ```bash
   DATABASE_URL="postgres://skool:skool_pass@localhost:5432/skool_mvp?sslmode=disable" go run ./cmd/api
   ```
   Or with Docker:
   ```bash
   DATABASE_URL="postgres://skool:skool_pass@host.docker.internal:5432/skool_mvp?sslmode=disable" docker run --rm -p 8080:8080 -e DATABASE_URL hcuri/skool-mvp-api:latest
   ```
3) Env vars:
   - `PORT` (default 8080)
   - `LOG_LEVEL` (default info)
   - `DATABASE_URL` (required for Postgres; falls back to in-memory store if unset)
4) Swagger UI: http://localhost:8080/swagger
5) Logging: structured JSON via `zap` (method, path, status, bytes, duration); adjust verbosity with `LOG_LEVEL`.

## Deployment
- Packaged as a Helm chart (`charts/skool-mvp-api`).
- In EKS, ArgoCD pulls the chart from this repo and applies environment-specific values from the GitOps repo: https://github.com/hcuri/skool-mvp-gitops
- Kubernetes Secrets provide `DATABASE_URL` (e.g., `skool-mvp-db`), not committed to git.

## Security / secrets
- DB passwords and other secrets are **not** stored in this repo. They are supplied via Kubernetes Secrets or env vars.
- Example manifests include placeholders only; create the real `skool-mvp-db` secret in the target namespace before deploying.
