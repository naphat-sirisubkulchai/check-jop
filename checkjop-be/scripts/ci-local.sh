#!/bin/bash

# Local CI Script - Mimics GitHub Actions CI pipeline
# Run this script to test your changes before pushing to GitHub

set -e

echo "🚀 Starting Local CI Pipeline..."

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Function to print status
print_status() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if required tools are installed
check_dependencies() {
    print_status "Checking dependencies..."
    
    if ! command -v go &> /dev/null; then
        print_error "Go is not installed"
        exit 1
    fi
    
    if ! command -v docker &> /dev/null; then
        print_error "Docker is not installed"
        exit 1
    fi
    
    if ! command -v docker-compose &> /dev/null; then
        print_error "Docker Compose is not installed"
        exit 1
    fi
    
    print_status "All dependencies are installed ✅"
}

# Format code
format_code() {
    print_status "Formatting code with goimports and gofmt..."
    
    # Run goimports
    if command -v goimports &> /dev/null; then
        goimports -w .
    else
        print_warning "goimports not available, trying to install..."
        go install golang.org/x/tools/cmd/goimports@latest
        export PATH=$PATH:$(go env GOPATH)/bin
        goimports -w .
    fi
    
    # Run gofmt
    go fmt ./...
    
    # Check if code is properly formatted
    if [ "$(gofmt -s -l . | wc -l)" -gt 0 ]; then
        print_error "Code is not formatted properly:"
        gofmt -s -l .
        print_warning "Run formatting again to fix issues"
        return 1
    fi
    
    print_status "Code formatting check passed ✅"
}

# Vet code
vet_code() {
    print_status "Running go vet..."
    go vet ./...
    print_status "Go vet check passed ✅"
}

# Install dependencies
install_deps() {
    print_status "Installing dependencies..."
    go mod download
    go mod verify
    print_status "Dependencies installed ✅"
}

# Install goimports
install_goimports() {
    print_status "Installing goimports..."
    go install golang.org/x/tools/cmd/goimports@latest
    print_status "goimports installed ✅"
}

# Run linter (if available)
run_linter() {
    print_status "Running linter..."
    
    if command -v golangci-lint &> /dev/null; then
        golangci-lint run --timeout=5m
        print_status "Linting passed ✅"
    else
        print_warning "golangci-lint not installed, skipping linting"
        print_warning "Install it with: curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b \$(go env GOPATH)/bin v1.54.2"
    fi
}

# Start test database
start_test_db() {
    print_status "Starting test database..."
    
    # Check if PostgreSQL is already running
    if docker ps | grep -q postgres; then
        print_warning "PostgreSQL container already running"
    else
        docker run -d \
            --name postgres-test \
            -e POSTGRES_USER=postgres \
            -e POSTGRES_PASSWORD=postgres \
            -e POSTGRES_DB=checkjop_test \
            -p 5433:5432 \
            postgres:15-alpine
        
        # Wait for database to be ready
        print_status "Waiting for database to be ready..."
        sleep 10
    fi
    
    # Test database connection
    export DB_HOST=localhost
    export DB_PORT=5433
    export DB_USER=postgres
    export DB_PASSWORD=postgres
    export DB_NAME=checkjop_test
    export DB_SSLMODE=disable
    
    print_status "Test database is ready ✅"
}

# Stop test database
stop_test_db() {
    print_status "Stopping test database..."
    if docker ps | grep -q postgres-test; then
        docker stop postgres-test
        docker rm postgres-test
    fi
    print_status "Test database stopped ✅"
}

# Run tests
run_tests() {
    print_status "Running tests..."
    
    # Set test environment
    export GIN_MODE=test
    
    # Run tests with coverage
    go test -v -race -coverprofile=coverage.out -covermode=atomic ./...
    
    # Generate coverage report
    go tool cover -html=coverage.out -o coverage.html
    
    # Display coverage summary
    coverage_percent=$(go tool cover -func=coverage.out | grep total | awk '{print $3}')
    print_status "Test coverage: $coverage_percent"
    
    print_status "Tests passed ✅"
}

# Build application
build_app() {
    print_status "Building application..."
    
    mkdir -p bin
    go build -v -o bin/main cmd/main.go
    
    print_status "Build completed ✅"
}

# Build Docker image
build_docker() {
    print_status "Building Docker image..."
    
    if ! command -v docker &> /dev/null; then
        print_warning "Docker not installed, skipping Docker build"
        return
    fi
    
    docker build -t checkjop-be:latest . || {
        print_error "Docker build failed"
        return 1
    }
    
    print_status "Docker build completed ✅"
    docker images | grep checkjop-be || print_warning "Image not found in local registry"
}

# Run security scan (if gosec is available)
run_security_scan() {
    print_status "Running security scan..."
    
    if command -v gosec &> /dev/null; then
        gosec -fmt json -out gosec.json ./...
        print_status "Security scan passed ✅"
    else
        print_warning "gosec not installed, skipping security scan"
        print_warning "Install it with: go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest"
    fi
}

# Clean up function
cleanup() {
    print_status "Cleaning up..."
    stop_test_db
    rm -f coverage.out coverage.html gosec.json
}

# Main execution
main() {
    # Set trap to cleanup on exit
    trap cleanup EXIT
    
    check_dependencies
    install_deps
    format_code
    vet_code
    run_linter
    start_test_db
    run_tests
    build_app
    build_docker
    
    print_status "🎉 Local CI Pipeline completed successfully!"
    print_status "Your code is ready to be pushed to GitHub"
}

# Handle script arguments
case "${1:-all}" in
    "deps")
        install_deps
        ;;
    "format")
        format_code
        ;;
    "vet")
        vet_code
        ;;
    "lint")
        run_linter
        ;;
    "test")
        start_test_db
        run_tests
        stop_test_db
        ;;
    "build")
        build_app
        ;;
    "docker")
        build_docker
        ;;
    "clean")
        cleanup
        ;;
    "all")
        main
        ;;
    *)
        echo "Usage: $0 [deps|format|vet|lint|test|build|docker|clean|all]"
        echo "  deps     - Install dependencies"
        echo "  format   - Format code"
        echo "  vet      - Run go vet"
        echo "  lint     - Run linter"
        echo "  test     - Run tests"
        echo "  build    - Build application"
        echo "  docker   - Build Docker image"
        echo "  clean    - Clean up"
        echo "  all      - Run complete CI pipeline (default)"
        exit 1
        ;;
esac