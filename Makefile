# lirt Makefile

VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS := -ldflags "-X main.version=$(VERSION)"

.PHONY: all build install test lint clean

all: build

build:
	@echo "Building lirt $(VERSION)..."
	@go build $(LDFLAGS) -o bin/lirt .

install:
	@echo "Installing lirt $(VERSION)..."
	@go install $(LDFLAGS) .

test:
	@echo "Running tests..."
	@go test ./...

test-verbose:
	@echo "Running tests (verbose)..."
	@go test -v ./...

test-cover:
	@echo "Running tests with coverage..."
	@go test -cover ./...

lint:
	@echo "Running linters..."
	@golangci-lint run

clean:
	@echo "Cleaning..."
	@rm -rf bin/
	@go clean

deps:
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy

.PHONY: help
help:
	@echo "lirt Makefile targets:"
	@echo "  make build        - Build lirt binary to bin/lirt"
	@echo "  make install      - Install lirt to GOPATH/bin"
	@echo "  make test         - Run all tests"
	@echo "  make test-verbose - Run tests with verbose output"
	@echo "  make test-cover   - Run tests with coverage report"
	@echo "  make lint         - Run golangci-lint"
	@echo "  make clean        - Remove build artifacts"
	@echo "  make deps         - Download and tidy dependencies"
