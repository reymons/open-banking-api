GC=go
BUILD_DIR=./bin
CMD_DIR=./cmd

# Build
.PHONY: run-internal run-open-banking run-internal-ws run-internal-broker

run-internal:
	$(GC) run $(CMD_DIR)/internal/main.go

run-internal-ws:
	$(GC) run $(CMD_DIR)/internal-ws/main.go

run-open-banking:
	$(GC) run $(CMD_DIR)/open-banking/main.go

run-internal-broker:
	$(GC) run $(CMD_DIR)/internal-broker/main.go

# Build
.PHONY: build-internal build-open-banking build-internal-ws

build-internal:
	$(GC) build -o $(BUILD_DIR)/internal $(CMD_DIR)/internal/main.go

build-internal-vs:
	$(GC) build -o $(BUILD_DIR)/internal-ws $(CMD_DIR)/internal-ws/main.go

build-open-banking:
	$(GC) build -o $(BUILD_DIR)/open-banking $(CMD_DIR)/open-banking/main.go

# Test
.PHONY: test

test:
	go test -v ./...

# Migrations
.PHONY: migrate-internal migrate-open-banking

migrate-internal:
	go run ./cmd/migrate/main.go

migrate-open-banking:
	go run ./cmd/migrate/main.go

# Docker
.PHONY: docker-upp-all docker-down-all

docker-up-all:
	docker compose -f build/docker-compose.yml up --build -d

docker-down-all:
	docker compose -f build/docker-compose.yml down
