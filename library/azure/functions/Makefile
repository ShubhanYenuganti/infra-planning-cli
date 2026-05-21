.PHONY: build test lint install clean

build:
	go build -o bin/functions-pp-cli ./cmd/functions-pp-cli

test:
	go test ./...

lint:
	golangci-lint run

install:
	go install ./cmd/functions-pp-cli

clean:
	rm -rf bin/

build-mcp:
	go build -o bin/functions-pp-mcp ./cmd/functions-pp-mcp

install-mcp:
	go install ./cmd/functions-pp-mcp

build-all: build build-mcp
