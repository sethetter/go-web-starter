.PHONY: build
build:
	go build -v ./...

.PHONY: test
test:
	go test ./...

.PHONY: start
start:
	docker compose up

.PHONY: psql
psql:
	docker compose exec db psql "postgresql://go-web-starter:supsupsup@localhost:5432/go-web-starter"

.PHONY: migrate
migrate:
	docker compose exec app go run ./cmd/migrate