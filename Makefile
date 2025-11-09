.PHONY: build run clean test fmt vet help

# Build the bot
build:
	@echo "Building margonem-bot..."
	@go build -o bin/margonem-bot ./cmd/bot
	@echo "Build complete: bin/margonem-bot"

# Build with version
build-release:
	@echo "Building release..."
	@go build -ldflags "-X main.version=v1.0.0" -o bin/margonem-bot ./cmd/bot
	@echo "Release build complete: bin/margonem-bot"

# Run the bot
run: build
	@./bin/margonem-bot

# Run with custom config
run-config: build
	@./bin/margonem-bot --config $(CONFIG)

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf bin/
	@rm -rf screenshots/
	@echo "Clean complete"

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	@go test -cover ./...

# Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...

# Run go vet
vet:
	@echo "Running go vet..."
	@go vet ./...

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	@go mod download
	@echo "Dependencies downloaded"

# Tidy dependencies
tidy:
	@echo "Tidying dependencies..."
	@go mod tidy
	@echo "Dependencies tidied"

# Install the bot
install: build
	@echo "Installing margonem-bot..."
	@cp bin/margonem-bot $(GOPATH)/bin/
	@echo "Installed to $(GOPATH)/bin/margonem-bot"

# Show help
help:
	@echo "Margonem Bot - Makefile commands:"
	@echo ""
	@echo "  make build          - Build the bot"
	@echo "  make build-release  - Build with version tag"
	@echo "  make run            - Build and run the bot"
	@echo "  make run-config     - Run with custom config (CONFIG=path)"
	@echo "  make clean          - Remove build artifacts"
	@echo "  make test           - Run tests"
	@echo "  make test-coverage  - Run tests with coverage"
	@echo "  make fmt            - Format code"
	@echo "  make vet            - Run go vet"
	@echo "  make deps           - Download dependencies"
	@echo "  make tidy           - Tidy dependencies"
	@echo "  make install        - Install to GOPATH/bin"
	@echo "  make help           - Show this help message"
