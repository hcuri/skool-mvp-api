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
- AWS infra (VPC, EKS, RDS) is provisioned from `skool-mvp-infra`.
- GitOps/ArgoCD config and environment values live in `skool-mvp-gitops`.
- Images are tagged with Git SHAs; ArgoCD rolls deployments by updating values in the GitOps repo.

## Running locally (high level)
- Start Postgres locally (e.g., `docker-compose up -d db`), or point to any Postgres instance.
- Run: `go run ./cmd/api`
- Env vars: `PORT` (default 8080), `LOG_LEVEL` (default info), `DATABASE_URL` (required for Postgres; otherwise falls back to in-memory store).
- Swagger UI: http://localhost:8080/swagger

## Deployment (high level)
- Packaged as a Helm chart (`charts/skool-mvp-api`).
- In EKS, ArgoCD pulls the chart from this repo and applies environment-specific values from `skool-mvp-gitops`.
- Kubernetes Secrets provide `DATABASE_URL` (e.g., `skool-mvp-db`), not committed to git.

## Security / secrets
- DB passwords and other secrets are **not** stored in this repo. They are supplied via Kubernetes Secrets or env vars.
- Example manifests include placeholders only; create the real `skool-mvp-db` secret in the target namespace before deploying.

## License & attribution
- Licensed under the MIT License (see `LICENSE`).
- You are welcome to copy/adapt the code. If you reuse it, please preserve the original copyright and license notices so attribution to the author, Hector Curi, is retained.
