.PHONY: help dev build test clean docker-build docker-up docker-down

# Variables
GO_VERSION := 1.21
PYTHON_VERSION := 3.11
DOCKER_REGISTRY := ghcr.io

help:
	@echo "HLFA - Heuristic Logic & Flow Analyzer"
	@echo ""
	@echo "Commands:"
	@echo "  make dev           - Run development environment"
	@echo "  make build         - Build all binaries"
	@echo "  make test          - Run tests"
	@echo "  make clean         - Clean build artifacts"
	@echo "  make docker-build  - Build Docker images"
	@echo "  make docker-up     - Start Docker containers (dev)"
	@echo "  make docker-down   - Stop Docker containers"

# Development
dev: docker-up
	@echo "HLFA development environment started"

dev-stop: docker-down

build: build-backend build-frontend build-python

build-backend:
	@echo "Building backend..."
	cd cmd/hlfa && go build -o ../../bin/hlfa
	cd cmd/hlfa-server && go build -o ../../bin/hlfa-server

build-frontend:
	@echo "Building frontend..."
	cd frontend && npm install && npm run build

build-python:
	@echo "Building Python engine..."
	cd python && pip install -r requirements.txt

test:
	@echo "Running tests..."
	go test -v ./...

clean:
	@echo "Cleaning..."
	rm -rf bin/
	cd frontend && rm -rf dist node_modules
	find . -name "*.pyc" -delete

docker-build:
	@echo "Building Docker images..."
	docker-compose -f deploy/docker-compose.yml build

docker-up:
	@echo "Starting Docker containers..."
	docker-compose -f deploy/docker-compose.yml up -d
	@echo "Services running:"
	@echo "  Backend: http://localhost:8080"
	@echo "  Frontend: http://localhost:3000"
	@echo "  API Docs: http://localhost:8080/swagger"

docker-down:
	@echo "Stopping Docker containers..."
	docker-compose -f deploy/docker-compose.yml down

install-deps:
	@echo "Installing dependencies..."
	go mod download
	cd frontend && npm install
	cd python && pip install -r requirements.txt
