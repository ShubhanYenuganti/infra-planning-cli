.PHONY: build test lint install clean

build:
	go build -o bin/apprunner-pp-cli ./cmd/apprunner-pp-cli

test:
	go test ./...

lint:
	golangci-lint run

install:
	go install ./cmd/apprunner-pp-cli

clean:
	rm -rf bin/

build-mcp:
	go build -o bin/apprunner-pp-mcp ./cmd/apprunner-pp-mcp

install-mcp:
	go install ./cmd/apprunner-pp-mcp

build-all: build build-mcp
