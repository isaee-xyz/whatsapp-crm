.PHONY: all build run test clean docker-build docker-up docker-down migrate frontend-dev frontend-build

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_NAME=whatomate
BINARY_PATH=./cmd/server

# Docker parameters
DOCKER_COMPOSE=docker compose -f docker/docker-compose.yml

all: build

# Build the backend
build:
	$(GOBUILD) -o $(BINARY_NAME) $(BINARY_PATH)

# Run the backend locally
run:
	$(GOCMD) run $(BINARY_PATH)/main.go -config config.toml

# Run with migrations
run-migrate:
	$(GOCMD) run $(BINARY_PATH)/main.go -config config.toml -migrate

# Run tests
test:
	$(GOTEST) -v ./...

# Run tests with coverage
test-coverage:
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

# Clean build artifacts
clean:
	rm -f $(BINARY_NAME)
	rm -f coverage.out coverage.html

# Download dependencies
deps:
	$(GOMOD) download
	$(GOMOD) tidy

# Update dependencies
deps-update:
	$(GOMOD) tidy
	$(GOGET) -u ./...

# Docker commands
docker-build:
	$(DOCKER_COMPOSE) build

docker-up:
	$(DOCKER_COMPOSE) up -d

docker-down:
	$(DOCKER_COMPOSE) down

docker-logs:
	$(DOCKER_COMPOSE) logs -f

docker-restart:
	$(DOCKER_COMPOSE) restart

# Database migrations
migrate:
	$(GOCMD) run $(BINARY_PATH)/main.go -config config.toml -migrate

# Frontend commands
frontend-install:
	cd frontend && npm install

frontend-dev:
	cd frontend && npm run dev

frontend-build:
	cd frontend && npm run build

frontend-preview:
	cd frontend && npm run preview

# Development - run both backend and frontend
dev:
	@echo "Starting backend and frontend in development mode..."
	@make run-migrate &
	@make frontend-dev

# Lint
lint:
	golangci-lint run ./...

# Format code
fmt:
	$(GOCMD) fmt ./...

# Generate swagger docs (if using)
swagger:
	swag init -g cmd/server/main.go -o api/docs

# Help
help:
	@echo "Available targets:"
	@echo "  build          - Build the backend binary"
	@echo "  run            - Run the backend locally"
	@echo "  run-migrate    - Run the backend with database migrations"
	@echo "  test           - Run tests"
	@echo "  test-coverage  - Run tests with coverage report"
	@echo "  clean          - Remove build artifacts"
	@echo "  deps           - Download dependencies"
	@echo "  deps-update    - Update dependencies"
	@echo "  docker-build   - Build Docker images"
	@echo "  docker-up      - Start Docker containers"
	@echo "  docker-down    - Stop Docker containers"
	@echo "  docker-logs    - View Docker logs"
	@echo "  frontend-install - Install frontend dependencies"
	@echo "  frontend-dev   - Run frontend in development mode"
	@echo "  frontend-build - Build frontend for production"
	@echo "  dev            - Run both backend and frontend in development mode"
	@echo "  lint           - Run linter"
	@echo "  fmt            - Format code"
