APP := aegis-mcp

.PHONY: build test fmt run

build:
	go build ./...

test:
	go test ./...

fmt:
	gofmt -w ./cmd ./internal

run:
	go run ./cmd/$(APP)
