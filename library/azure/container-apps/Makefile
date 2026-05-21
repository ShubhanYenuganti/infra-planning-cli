.PHONY: build test lint install clean

build:
	go build -o bin/container-apps-pp-cli ./cmd/container-apps-pp-cli

test:
	go test ./...

lint:
	golangci-lint run

install:
	go install ./cmd/container-apps-pp-cli

clean:
	rm -rf bin/

build-mcp:
	go build -o bin/container-apps-pp-mcp ./cmd/container-apps-pp-mcp

install-mcp:
	go install ./cmd/container-apps-pp-mcp

build-all: build build-mcp
