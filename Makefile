ENTRY=main.go
OUT=main

run: main.go
	go run $(ENTRY)

build:
	go build -o $(OUT) build $(ENTRY)
