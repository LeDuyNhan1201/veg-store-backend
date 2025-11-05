# Makefile

TEST_FLAGS=-v -cover -count=1
PKG=./...

.PHONY: test test-coverage lint tidy

test:
	@echo "Running unit tests..."
	@go test $(PKG) $(TEST_FLAGS)

test-coverage:
	@echo "Running tests with coverage report..."
	@go test $(PKG) -coverprofile=coverage.out -covermode=atomic
	@go tool cover -func=coverage.out | tail -n 1
	@go tool cover -html=coverage.out -o coverage.html
	@echo "✅ Coverage report saved to coverage.html"

lint:
	@echo "Running go vet and lint..."
	@go vet $(PKG)
	@if command -v golangci-lint >/dev/null 2>&1; then golangci-lint run; fi

tidy:
	@echo "Tidying modules..."
	@go mod tidy
