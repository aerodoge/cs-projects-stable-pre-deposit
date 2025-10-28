# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=gofmt
GOLINT=golint
GOVET=$(GOCMD) vet

# Binary name
BINARY_NAME=cs-projects-stable-pre-deposit
BINARY_UNIX=$(BINARY_NAME)_unix

# Build directory
BUILD_DIR=build

.PHONY: all build clean test deps fmt lint vet run help

# Default target
all: deps fmt vet build

# Build the binary
build:
	@echo "Building..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) -v .

# Build for Linux
build-linux:
	@echo "Building for Linux..."
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BUILD_DIR)/$(BINARY_UNIX) -v .

# Run the application
run: build
	@echo "Running application..."
	./$(BUILD_DIR)/$(BINARY_NAME)

# Run without building (useful for development)
run-dev:
	@echo "Running application in development mode..."
	$(GOCMD) run main.go

# Clean build files
clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy

# Format code
fmt:
	@echo "Formatting code..."
	$(GOFMT) -s -w .

# Lint code
lint:
	@echo "Linting code..."
	@if command -v golint > /dev/null; then \
		golint ./...; \
	else \
		echo "golint not installed. Install with: go install golang.org/x/lint/golint@latest"; \
	fi

# Vet code
vet:
	@echo "Vetting code..."
	$(GOVET) ./...

# Run tests
test:
	@echo "Running tests..."
	$(GOTEST) -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

# Install dependencies for development
install-tools:
	@echo "Installing development tools..."
	$(GOCMD) install golang.org/x/lint/golint@latest
	$(GOCMD) install github.com/securecodewarrior/sast-scan@latest

# Check for security vulnerabilities
security:
	@echo "Checking for security vulnerabilities..."
	@if command -v gosec > /dev/null; then \
		gosec ./...; \
	else \
		echo "gosec not installed. Install with: go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest"; \
	fi

# Update dependencies
update:
	@echo "Updating dependencies..."
	$(GOMOD) get -u ./...
	$(GOMOD) tidy

# Build and run
start: build run

# Show help
help:
	@echo "Available targets:"
	@echo "  build         - Build the binary"
	@echo "  build-linux   - Build for Linux"
	@echo "  run           - Build and run the application"
	@echo "  run-dev       - Run without building (development mode)"
	@echo "  clean         - Clean build files"
	@echo "  deps          - Download and tidy dependencies"
	@echo "  fmt           - Format code"
	@echo "  lint          - Lint code"
	@echo "  vet           - Vet code"
	@echo "  test          - Run tests"
	@echo "  test-coverage - Run tests with coverage report"
	@echo "  install-tools - Install development tools"
	@echo "  security      - Check for security vulnerabilities"
	@echo "  update        - Update dependencies"
	@echo "  start         - Build and run"
	@echo "  help          - Show this help message"