DATABASE_URL ?= postgres://shortly:shortly@localhost:5432/shortly?sslmode=disable

.PHONY: run test race vet sqlc db-up migrate-up migrate-down docker

run:            ## run the API locally (uses DATABASE_URL if set)
	go run ./cmd/api

sqlc:           ## regenerate the db package from db/query.sql (go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest)
	sqlc generate

test:           ## run all tests
	go test ./...

race:           ## tests with the race detector
	go test -race ./...

vet:            ## static analysis
	go vet ./...

db-up:          ## start just postgres via docker
	docker compose up -d db

migrate-up:     ## apply migrations (requires: go install github.com/pressly/goose/v3/cmd/goose@latest)
	goose -dir migrations postgres "$(DATABASE_URL)" up

migrate-down:   ## roll back last migration
	goose -dir migrations postgres "$(DATABASE_URL)" down

docker:         ## build + run the full stack
	docker compose up --build
