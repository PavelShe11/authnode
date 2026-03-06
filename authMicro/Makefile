.PHONY: test test-unit test-integration test-e2e test-all test-verbose test-coverage clean build

# Default: run all tests (unit + integration + e2e)
test:
	go test -tags=integration,e2e ./...

# Unit tests only (no tags needed)
test-unit:
	go test ./...

# Integration tests only (requires testcontainers)
test-integration:
	go test -tags=integration ./...

# E2E tests only (requires testcontainers)
test-e2e:
	go test -tags=e2e ./...

# All tests with verbose output (hides empty directories)
test-verbose:
	go test -tags=integration,e2e -v ./... 2>&1 | grep -v '\[no test files\]'

# Unit tests verbose
test-unit-verbose:
	go test -v ./... 2>&1 | grep -v '\[no test files\]'

# Test with coverage report
test-coverage:
	go test -tags=integration,e2e -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

# Run specific test by name (usage: make test-run TEST=TestName)
test-run:
	go test -tags=integration,e2e -v -run $(TEST) ./...

# Build the application
build:
	go build -o bin/auth-app ./cmd/app

# Build for debug (with symbols)
build-debug:
	go build -gcflags="all=-N -l" -o bin/auth-app-debug ./cmd/app

# Clean build artifacts
clean:
	rm -rf bin/ coverage.out coverage.html

# Format code
fmt:
	go fmt ./...

# Run linter (requires golangci-lint)
lint:
	golangci-lint run ./...

# Download dependencies
deps:
	go mod download
	go mod tidy

# Generate swagger docs (requires swag)
swagger:
	swag init -g cmd/app/main.go

# Help
help:
	@echo "Available targets:"
	@echo "  test              - Run all tests (unit + integration + e2e)"
	@echo "  test-unit         - Run unit tests only"
	@echo "  test-integration  - Run integration tests only"
	@echo "  test-e2e          - Run E2E tests only"
	@echo "  test-verbose      - Run all tests with verbose output"
	@echo "  test-unit-verbose - Run unit tests with verbose output"
	@echo "  test-coverage     - Run tests with coverage report"
	@echo "  test-run TEST=X   - Run specific test by name"
	@echo "  build             - Build the application"
	@echo "  build-debug       - Build with debug symbols"
	@echo "  clean             - Remove build artifacts"
	@echo "  fmt               - Format code"
	@echo "  lint              - Run linter"
	@echo "  deps              - Download and tidy dependencies"
	@echo "  swagger           - Generate swagger docs"