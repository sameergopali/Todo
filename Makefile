# Todo API Makefile

.PHONY: help build run test clean deps example

# Default target
help:
	@echo "Available commands:"
	@echo "  make run      - Run the todo API server"
	@echo "  make build    - Build the application"
	@echo "  make test     - Run tests"
	@echo "  make deps     - Download dependencies"
	@echo "  make example  - Run the client example"
	@echo "  make clean    - Clean build artifacts"
	@echo "  make help     - Show this help message"

# Download dependencies
deps:
	go mod tidy
	go mod download

# Build the application
build:
	go build -o bin/todo-api cmd/main.go

# Run the application
run:
	go run cmd/main.go

# Run tests
test:
	go test -v ./...

# Run the client example
example:
	@echo "Make sure the server is running in another terminal with 'make run'"
	@echo "Starting client example in 3 seconds..."
	@sleep 3
	go run examples/client_example.go

# Clean build artifacts
clean:
	rm -rf bin/
	go clean

# Format code
fmt:
	go fmt ./...

# Run linter
lint:
	golangci-lint run

# Create build directory
bin:
	mkdir -p bin
