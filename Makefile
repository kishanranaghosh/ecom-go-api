.PHONY: build run dev lint migrate-up migrate-down help

help:
	@echo "Makefile commands:"
	@echo "  build          - Build the Go application"
	@echo "  run            - Run the built application"
	@echo "  dev            - Run the application in development mode"
	@echo "  lint           - Run golangci-lint on the codebase"
	@echo "  migrate-up     - Apply database migrations"
	@echo "  migrate-down   - Rollback database migrations"

build:
	go build -o bin/api cmd/api/main.go

run: build
	./bin/api

dev:
	go run cmd/api/main.go

lint:
	golangci-lint run ./...

migrate-up:
	migrate -path db/migrations -database "postgres://user:123456@localhost:5432/ecom?sslmode=disable" up

migrate-down:
	migrate -path db/migrations -database "postgres://user:123456@localhost:5432/ecom?sslmode=disable" down

docker-up:
	docker compose -f ./docker/docker-compose.yml up -d

docker-down:
	docker compose -f ./docker/docker-compose.yml down
