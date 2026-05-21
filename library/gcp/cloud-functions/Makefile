.PHONY: build test lint install clean

build:
	go build -o bin/cloud-functions-pp-cli ./cmd/cloud-functions-pp-cli

test:
	go test ./...

lint:
	golangci-lint run

install:
	go install ./cmd/cloud-functions-pp-cli

clean:
	rm -rf bin/

build-mcp:
	go build -o bin/cloud-functions-pp-mcp ./cmd/cloud-functions-pp-mcp

install-mcp:
	go install ./cmd/cloud-functions-pp-mcp

build-all: build build-mcp
