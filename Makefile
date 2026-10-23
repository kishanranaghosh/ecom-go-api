.PHONY: help build run lint migrate-up migrate-down

help:
	@echo "Available commands:"
	@echo "  build          Build the application"
	@echo "  run            Build and run the application"
	@echo "  lint           Run golangci-lint on the codebase"
	@echo "  migrate-up     Apply database migrations"
	@echo "  migrate-down   Rollback database migrations"

build:
	go build -o bin/app ./cmd/api/main.go
run: build
	./bin/app

lint:
	golangci-lint run ./...

migrate-up:
	migrate -path db/migrations -database "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable" up
migrate-down:
	migrate -path db/migrations -database "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable" down

