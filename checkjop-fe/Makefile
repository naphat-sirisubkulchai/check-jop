# Variables
IMAGE_NAME = senior-plan-checking
CONTAINER_NAME = senior-plan-checking-container
PORT = 3000

# Default target
.DEFAULT_GOAL := help

# Help target
.PHONY: help
help: ## Show this help message
	@echo "Available targets:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Initialize project (first time setup)
.PHONY: init
init: build run ## Initialize project: build and run container (first time setup)
	@echo "Project initialized! Container is running at http://localhost:$(PORT)"
	@echo "Use 'make logs' to view logs or 'make help' for more commands"

# Build the Docker image
.PHONY: build
build: ## Build the Docker image
	docker build -t $(IMAGE_NAME) .

# Run the container
.PHONY: run
run: ## Run the container (detached)
	docker run -d --name $(CONTAINER_NAME) -p $(PORT):$(PORT) $(IMAGE_NAME)

# Run the container in foreground
.PHONY: run-fg
run-fg: ## Run the container in foreground
	docker run --rm --name $(CONTAINER_NAME) -p $(PORT):$(PORT) $(IMAGE_NAME)

# Stop the container
.PHONY: stop
stop: ## Stop the running container
	docker stop $(CONTAINER_NAME)

# Remove the container
.PHONY: remove
remove: ## Remove the container
	docker rm $(CONTAINER_NAME)

# Stop and remove the container
.PHONY: clean
clean: ## Stop and remove the container
	-docker stop $(CONTAINER_NAME)
	-docker rm $(CONTAINER_NAME)

# Remove the Docker image
.PHONY: clean-image
clean-image: ## Remove the Docker image
	docker rmi $(IMAGE_NAME)

# Full cleanup (remove container and image)
.PHONY: clean-all
clean-all: clean clean-image ## Remove container and image

# Rebuild (clean build)
.PHONY: rebuild
rebuild: clean-all build ## Clean rebuild (remove image and rebuild)

# Reset everything (stop, remove container, remove image, fix linting issues, and rebuild)
.PHONY: reset
reset: clean-all lint-fix build run ## Complete reset: clean everything and rebuild
	@echo "Reset complete! Container is running at http://localhost:$(PORT)"

# Build and run
.PHONY: up
up: lint-fix build run ## Build and run the container
	@echo "Container is running at http://localhost:$(PORT)"

# View container logs
.PHONY: logs
logs: ## Show container logs
	docker logs $(CONTAINER_NAME)

# Follow container logs
.PHONY: logs-f
logs-f: ## Follow container logs
	docker logs -f $(CONTAINER_NAME)

# Execute bash in running container
.PHONY: shell
shell: ## Execute bash in the running container
	docker exec -it $(CONTAINER_NAME) sh

# Show running containers
.PHONY: ps
ps: ## Show running containers
	docker ps

# Show all containers (including stopped)
.PHONY: ps-all
ps-all: ## Show all containers
	docker ps -a

# Development commands
.PHONY: dev
dev: ## Run development server locally (yarn dev)
	yarn dev

# Install dependencies locally
.PHONY: install
install: ## Install dependencies locally
	yarn install

# Build locally
.PHONY: build-local
build-local: ## Build the project locally
	yarn build

# Lint the code
.PHONY: lint
lint: ## Run ESLint to check code quality
	yarn lint

# Lint and fix auto-fixable issues
.PHONY: lint-fix
lint-fix: ## Run ESLint and fix auto-fixable issues
	yarn lint --fix

# Run built project locally
.PHONY: start-local
start-local: ## Start the built project locally
	yarn start
