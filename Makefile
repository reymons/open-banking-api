run-internal:
	go run ./internal/main.go

run-open-banking:
	go run ./open-banking/main.go

build-internal:
	go build -o ./internal.bin ./internal/main.go

build-open-banking:
	go build -o ./open-banking.bin ./open-banking/main.go

migrate-internal:
	go run ./scripts/migrate.go

migrate-open-banking:
	go run ./scripts/migrate.go

docker-up-all:
	docker compose -f docker-compose.yml --env-file .env up --build -d

docker-down-all:
	docker compose -f docker-compose.yml --env-file .env down

