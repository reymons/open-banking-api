GC=go
BUILD_DIR=./bin
CMD_DIR=./cmd

# Build
.PHONY: build-internal build-open-banking

build-internal:
	$(GC) build -o $(BUILD_DIR)/internal $(CMD_DIR)/internal/main.go

build-open-banking:
	$(GC) build -o $(BUILD_DIR)/open-banking $(CMD_DIR)/open-banking/main.go

# Migrations
.PHONY: migrate-internal migrate-open-banking

migrate-internal:
	go run ./scripts/migrate.go

migrate-open-banking:
	go run ./scripts/migrate.go

# Docker
.PHONY: docker-upp-all docker-down-all

docker-up-all:
	docker compose -f build/docker-compose.yml --env-file .env up --build -d

docker-down-all:
	docker compose -f docker-compose.yml --env-file .env down

