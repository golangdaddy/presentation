.PHONY: help build run test clean docker-build docker-up docker-down docker-logs

# Default target
help:
	@echo "Available commands:"
	@echo "  build        - Build the Go application"
	@echo "  run          - Run the application locally"
	@echo "  test         - Run tests"
	@echo "  clean        - Clean build artifacts"
	@echo "  docker-build - Build Docker containers"
	@echo "  docker-up    - Start Docker services"
	@echo "  docker-down  - Stop Docker services"
	@echo "  docker-logs  - View Docker logs"
	@echo "  docker-restart - Restart Docker services"

# Build the Go application
build:
	go build -o bin/api cmd/api/main.go

# Run the application locally
run:
	go run cmd/api/main.go

# Run tests
test:
	go test ./...

# Clean build artifacts
clean:
	rm -rf bin/
	go clean

# Docker commands
docker-build:
	docker compose build

docker-up:
	docker compose up -d

docker-down:
	docker compose down

docker-logs:
	docker compose logs -f

docker-restart:
	docker compose down
	docker compose up -d --build

# Database commands
db-connect:
	docker compose exec postgres psql -U fleet_user -d fleet_management

db-reset:
	docker compose down -v
	docker compose up -d

# Development helpers
dev-setup:
	go mod tidy
	docker compose up -d

dev-logs:
	docker compose logs -f api

# API testing
test-api:
	@echo "Testing API endpoints..."
	@curl -s http://localhost:8080/health | jq .
	@echo "Health check completed"

# Show service status
status:
	docker compose ps 