.PHONY: build run test clean clean-image clean-all docker-build docker-run docker-down migrate migrate-up migrate-down migrate-reset init reset

# Build the application
build:
	go build -o bin/main cmd/main.go

# Run the application
run:
	go run cmd/main.go

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	go clean
	rm -f bin/main

# Clean Docker images
clean-image:
	docker-compose down --rmi all --volumes --remove-orphans

# Remove everything: containers, images, volumes, and build artifacts
clean-all: clean clean-image ## Remove container and image

# Install dependencies
deps:
	go mod tidy
	go mod download

# Format code
fmt:
	go fmt ./...

# Vet code
vet:
	go vet ./...

# Run linter
lint:
	golangci-lint run

# Docker commands
docker-build:
	docker-compose build

docker-run:
	docker-compose up -d

docker-down:
	docker-compose down

docker-logs:
	docker-compose logs -f

# Database migration commands
migrate:
	@echo "Running SQL migrations..."
	@for f in migrations/*.sql; do \
		echo "Applying $$f..."; \
		docker-compose exec -T postgres psql -U postgres -d checkjop -f /dev/stdin < $$f; \
	done
	@echo "Migrations complete."

migrate-up:
	migrate -path migrations -database "postgresql://postgres:postgres@localhost:5432/checkjop?sslmode=disable" up

migrate-down:
	migrate -path migrations -database "postgresql://postgres:postgres@localhost:5432/checkjop?sslmode=disable" down

migrate-create:
	migrate create -ext sql -dir migrations -seq $(name)

# Reset and update database schema using migrations via Docker
migrate-reset:
	@echo "Resetting database schema via Docker Compose..."
	docker-compose exec postgres psql -U postgres -d checkjop -c "DROP SCHEMA public CASCADE; CREATE SCHEMA public;"
	docker-compose restart app

# Development commands
dev: docker-run
	sleep 5
	make run

dev-down: docker-down

# Initialize project: build and run container (first time setup)
init: build docker-run ## Initialize project: build and run container (first time setup)
	@echo "Waiting for database to be ready..."
	sleep 10
	@echo "Project initialized successfully!"
	@echo "API will be available at http://localhost:8080"
	@echo "Health check: http://localhost:8080/health"

reset:
	@echo "🔁 Resetting project..."
	# Stop & remove docker containers/images/volumes (if running)
	@docker-compose down --rmi all --volumes --remove-orphans || true
	# Clean Go build artifacts (safe even if already clean)
	@go clean || true
	@rm -f bin/main || true
	@docker builder prune -f 
	# Build app
	@echo "🔨 Building Go application..."
	@go build -o bin/main cmd/main.go

	# Start containers
	@echo "🐳 Starting Docker containers..."
	@docker-compose up -d

	# Wait for DB to be ready
	@sleep 5
	@echo "🗑 Resetting database schema..."
	@docker-compose exec postgres psql -U postgres -d checkjop -c "DROP SCHEMA IF EXISTS public CASCADE; CREATE SCHEMA public;" || true
	@docker-compose restart app

	# Run SQL migrations
	@echo "📦 Running SQL migrations..."
	@make migrate

	# Final message
	@echo "✅ Reset complete! Project is ready to use."
	@echo "🌐 API: http://localhost:8080"
	@echo "🔎 Health: http://localhost:8080/health"

# Help
help:
	@echo "Available commands:"
	@echo "  build        - Build the application"
	@echo "  run          - Run the application"
	@echo "  test         - Run tests"
	@echo "  clean        - Clean build artifacts"
	@echo "  deps         - Install dependencies"
	@echo "  fmt          - Format code"
	@echo "  vet          - Vet code"
	@echo "  lint         - Run linter"
	@echo "  docker-build - Build Docker containers"
	@echo "  docker-run   - Run Docker containers"
	@echo "  docker-down  - Stop Docker containers"
	@echo "  docker-logs  - Show Docker logs"
	@echo "  migrate      - Run all SQL migration files against the Docker DB"
	@echo "  migrate-up   - Run database migrations up"
	@echo "  migrate-down - Run database migrations down"
	@echo "  migrate-create name=<name> - Create new migration"
	@echo "  migrate-reset - Reset database schema via Docker"
	@echo "  init         - Initialize project (first time setup)"
	@echo "  reset        - Complete reset: clean everything and rebuild"
	@echo "  clean-all    - Remove containers, images, and build artifacts"
	@echo "  dev          - Start development environment"
	@echo "  dev-down     - Stop development environment"
	@echo "  ci           - Run local CI pipeline"
	@echo "  ci-test      - Run tests in CI environment"
	@echo "  coverage     - Generate test coverage report"

# CI/CD Commands
ci:
	@./scripts/ci-local.sh all

ci-test:
	@./scripts/ci-local.sh test

coverage:
	go test -v -race -coverprofile=coverage.out -covermode=atomic ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"