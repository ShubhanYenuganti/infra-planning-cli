.PHONY: build test lint install clean

build:
	go build -o bin/lambda-pp-cli ./cmd/lambda-pp-cli

test:
	go test ./...

lint:
	golangci-lint run

install:
	go install ./cmd/lambda-pp-cli

clean:
	rm -rf bin/

build-mcp:
	go build -o bin/aws-lambda-pp-mcp ./cmd/aws-lambda-pp-mcp

install-mcp:
	go install ./cmd/aws-lambda-pp-mcp

build-all: build build-mcp
