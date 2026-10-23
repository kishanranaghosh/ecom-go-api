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

