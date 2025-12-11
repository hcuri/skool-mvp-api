APP_NAME := skool-mvp-app

.PHONY: build run docker-build docker-run

build:
	go build ./...

run:
	go run ./cmd/api

docker-build:
	docker build -t $(APP_NAME):local .

docker-run:
	docker run --rm -p 8080:8080 $(APP_NAME):local
