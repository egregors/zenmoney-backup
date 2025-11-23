PROJECT_NAME := "zenmoney-backup"
PKG := "github.com/egregors/$(PROJECT_NAME)"
PKG_LIST := $(shell go list ${PKG}/... | grep -v /vendor/)
GO_FILES := $(shell find . -name '*.go' | grep -v /vendor/ | grep -v _test.go)

# Build variables
BINARY_NAME := zenb
BUILD_DIR := ./build
LDFLAGS := -s -w -extldflags '-static'
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev-$(shell date +%Y%m%d-%H%M%S)")

# Docker variables
IMAGE_NAME := zenb
DOCKER_TAG := latest
REGISTRY := ghcr.io/egregors
FULL_IMAGE_NAME := $(REGISTRY)/zenmoney-backup

.PHONY: all build clean docker docker-run docker-push run lint test help deps update-deps

all: build

## Build commands

build: clean  ## Build the binary
	@echo "Building $(BINARY_NAME) version $(VERSION)..."
	@mkdir -p $(BUILD_DIR)
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
		-ldflags "$(LDFLAGS) -X main.revision=$(VERSION)" \
		-o $(BUILD_DIR)/$(BINARY_NAME) \
		./cmd/main.go

build-local: clean  ## Build the binary for local OS
	@echo "Building $(BINARY_NAME) for local OS..."
	@mkdir -p $(BUILD_DIR)
	@CGO_ENABLED=0 go build \
		-ldflags "$(LDFLAGS) -X main.revision=$(VERSION)" \
		-o $(BUILD_DIR)/$(BINARY_NAME) \
		./cmd/main.go

clean:  ## Clean build artifacts
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)

## Docker commands
## Note: Docker images are automatically built and published to GitHub Container Registry
## (ghcr.io/egregors/zenmoney-backup) on every merge to main and for version tags.
## Local build commands below are for development and testing purposes.

docker: ## Build Docker image locally
	@echo "Building Docker image $(IMAGE_NAME):$(DOCKER_TAG)..."
	@docker build -t $(IMAGE_NAME):$(DOCKER_TAG) .
	@echo "Note: Production images are auto-built and available at $(FULL_IMAGE_NAME)"

docker-run: ## Run the application in Docker
	@echo "Running $(IMAGE_NAME):$(DOCKER_TAG)..."
	@docker run --rm -it $(IMAGE_NAME):$(DOCKER_TAG)

docker-push: docker ## Build and push Docker image to GHCR
	@echo "Tagging Docker image for GHCR..."
	@docker tag $(IMAGE_NAME):$(DOCKER_TAG) $(FULL_IMAGE_NAME):$(DOCKER_TAG)
	@docker tag $(IMAGE_NAME):$(DOCKER_TAG) $(FULL_IMAGE_NAME):$(VERSION)
	@echo "Pushing Docker image to GHCR $(FULL_IMAGE_NAME):$(DOCKER_TAG)..."
	@docker push $(FULL_IMAGE_NAME):$(DOCKER_TAG)
	@docker push $(FULL_IMAGE_NAME):$(VERSION)

docker-clean: ## Clean Docker images
	@echo "Cleaning Docker images..."
	@docker rmi $(IMAGE_NAME):$(DOCKER_TAG) 2>/dev/null || true

## Development commands

run: build-local  ## Build and run the application locally
	@echo "Running $(BINARY_NAME)..."
	@$(BUILD_DIR)/$(BINARY_NAME)

dev:  ## Run in development mode
	@echo "Running in development mode..."
	@go run ./cmd/main.go

lint:  ## Lint the code
	@echo "Linting code..."
	golangci-lint run ./...; 

test:  ## Run tests
	@echo "Running tests..."
	@go test -v -race -cover ./...

test-short:  ## Run short tests
	@echo "Running short tests..."
	@go test -short ./...

bench:  ## Run benchmarks
	@echo "Running benchmarks..."
	@go test -bench=. -benchmem ./...

## Dependencies

deps:  ## Download dependencies
	@echo "Downloading dependencies..."
	@go mod download

update-deps:  ## Update Go dependencies
	@echo "Updating Go dependencies..."
	@go get -u ./...
	@go mod tidy

vendor:  ## Create vendor directory
	@echo "Creating vendor directory..."
	@go mod vendor

## Utilities

fmt:  ## Format code
	@echo "Formatting code..."
	@go fmt ./...

vet:  ## Run go vet
	@echo "Running go vet..."
	@go vet ./...

mod-tidy:  ## Tidy go.mod
	@echo "Tidying go.mod..."
	@go mod tidy

version:  ## Show version
	@echo "Version: $(VERSION)"

## Help

help:  ## Show this help message
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'
