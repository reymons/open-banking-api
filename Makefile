run-internal:
	go run ./internal/main.go

run-open-banking:
	go run ./internal/open-banking.go

build-internal:
	go build -o ./internal.bin ./internal/main.go

build-open-banking:
	go build -o ./open-banking.bin ./open-banking/main.go
