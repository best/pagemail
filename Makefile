# ==============================================================================
# Pagemail Makefile
# ==============================================================================

.PHONY: help
help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# ==============================================================================
# Development
# ==============================================================================

.PHONY: dev
dev: ## Start development environment (backend + frontend)
	@$(MAKE) -j2 dev-backend dev-frontend

.PHONY: dev-backend
dev-backend: ## Start backend with hot reload (requires air)
	@command -v air >/dev/null 2>&1 || { echo "Installing air..."; go install github.com/air-verse/air@latest; }
	cd cmd/pagemail && air

.PHONY: dev-frontend
dev-frontend: ## Start frontend dev server
	cd web && npm run dev

.PHONY: install
install: ## Install all dependencies
	go mod download
	cd web && npm install

# ==============================================================================
# Code Quality
# ==============================================================================

.PHONY: lint
lint: lint-backend lint-frontend ## Run all linters

.PHONY: lint-backend
lint-backend: ## Run Go linter
	@command -v golangci-lint >/dev/null 2>&1 || { echo "Installing golangci-lint..."; go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; }
	golangci-lint run ./...

.PHONY: lint-frontend
lint-frontend: ## Run frontend linter
	cd web && npm run lint

.PHONY: lint-fix
lint-fix: ## Fix linting issues
	golangci-lint run --fix ./...
	cd web && npm run lint -- --fix

.PHONY: fmt
fmt: ## Format code
	go fmt ./...
	cd web && npm run format

.PHONY: test
test: test-backend test-frontend ## Run all tests

.PHONY: test-backend
test-backend: ## Run Go tests with coverage
	go test -race -coverprofile=coverage.out -covermode=atomic ./...
	go tool cover -html=coverage.out -o coverage.html

.PHONY: test-frontend
test-frontend: ## Run frontend unit tests
	cd web && npm run test:unit

.PHONY: test-e2e
test-e2e: ## Run E2E tests
	cd web && npm run test:e2e

# ==============================================================================
# Database Migrations (Atlas)
# ==============================================================================

ATLAS := atlas
DB_URL ?= $(shell grep DB_URL .env 2>/dev/null | cut -d '=' -f2-)
ifeq ($(DB_URL),)
	DB_URL := postgres://pagemail:pagemail@localhost:5432/pagemail?sslmode=disable
endif

.PHONY: migrate-new
migrate-new: ## Create a new migration: make migrate-new name=xxx
	@if [ -z "$(name)" ]; then echo "Usage: make migrate-new name=migration_name"; exit 1; fi
	$(ATLAS) migrate diff $(name) \
		--dir "file://internal/db/migrations" \
		--to "file://internal/db/schema.sql" \
		--dev-url "docker://postgres/16/dev?search_path=public"

.PHONY: migrate-up
migrate-up: ## Apply all pending migrations
	$(ATLAS) migrate apply \
		--dir "file://internal/db/migrations" \
		--url "$(DB_URL)"

.PHONY: migrate-down
migrate-down: ## Rollback migrations: make migrate-down steps=1
	@if [ -z "$(steps)" ]; then steps=1; fi
	$(ATLAS) migrate down $(steps) \
		--dir "file://internal/db/migrations" \
		--url "$(DB_URL)"

.PHONY: migrate-status
migrate-status: ## Show migration status
	$(ATLAS) migrate status \
		--dir "file://internal/db/migrations" \
		--url "$(DB_URL)"

.PHONY: migrate-lint
migrate-lint: ## Lint migration files (for CI)
	$(ATLAS) migrate lint \
		--dir "file://internal/db/migrations" \
		--dev-url "docker://postgres/16/dev?search_path=public" \
		--latest 1

.PHONY: migrate-hash
migrate-hash: ## Recalculate migration hash
	$(ATLAS) migrate hash --dir "file://internal/db/migrations"

# ==============================================================================
# Build
# ==============================================================================

VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME := $(shell date -u '+%Y-%m-%dT%H:%M:%SZ')
LDFLAGS := -ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME) -s -w"

.PHONY: build
build: build-frontend build-backend ## Build all

.PHONY: build-backend
build-backend: ## Build Go binary
	CGO_ENABLED=0 go build $(LDFLAGS) -o bin/pagemail ./cmd/pagemail

.PHONY: build-frontend
build-frontend: ## Build frontend
	cd web && npm run build

# ==============================================================================
# Docker
# ==============================================================================

DOCKER_IMAGE := astralor/pagemail
DOCKER_TAG ?= $(VERSION)
PLATFORMS := linux/amd64,linux/arm64

.PHONY: docker-build
docker-build: ## Build Docker image
	docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) -f deploy/Dockerfile .

.PHONY: docker-push
docker-push: ## Push Docker image to registry
	docker push $(DOCKER_IMAGE):$(DOCKER_TAG)

.PHONY: docker-buildx
docker-buildx: ## Build and push multi-arch Docker image
	docker buildx build \
		--platform $(PLATFORMS) \
		--tag $(DOCKER_IMAGE):$(DOCKER_TAG) \
		--tag $(DOCKER_IMAGE):latest \
		--push \
		-f deploy/Dockerfile .

.PHONY: docker-run
docker-run: ## Run Docker container locally
	docker run --rm -p 8080:8080 --env-file .env $(DOCKER_IMAGE):$(DOCKER_TAG)

# ==============================================================================
# Deployment Helpers
# ==============================================================================

.PHONY: env-merge
env-merge: ## Merge .env.example into .env, showing differences
	@echo "Checking .env configuration..."
	@if [ ! -f .env ]; then \
		echo "Creating .env from .env.example..."; \
		cp .env.example .env; \
		echo "Done. Please edit .env with your settings."; \
	else \
		echo "Comparing .env.example with .env..."; \
		echo ""; \
		echo "=== Keys in .env.example but missing in .env ==="; \
		comm -23 <(grep -E '^[A-Z_]+=' .env.example | cut -d= -f1 | sort) \
		         <(grep -E '^[A-Z_]+=' .env | cut -d= -f1 | sort) || true; \
		echo ""; \
		echo "=== Keys in .env but not in .env.example (custom) ==="; \
		comm -13 <(grep -E '^[A-Z_]+=' .env.example | cut -d= -f1 | sort) \
		         <(grep -E '^[A-Z_]+=' .env | cut -d= -f1 | sort) || true; \
		echo ""; \
		read -p "Add missing keys from .env.example to .env? [y/N] " confirm; \
		if [ "$$confirm" = "y" ] || [ "$$confirm" = "Y" ]; then \
			for key in $$(comm -23 <(grep -E '^[A-Z_]+=' .env.example | cut -d= -f1 | sort) \
			                       <(grep -E '^[A-Z_]+=' .env | cut -d= -f1 | sort)); do \
				grep "^$$key=" .env.example >> .env; \
				echo "Added: $$key"; \
			done; \
			echo "Done."; \
		fi; \
	fi

.PHONY: compose-merge
compose-merge: ## Merge docker-compose.yaml.example into docker-compose.yaml
	@echo "Checking docker-compose.yaml configuration..."
	@if [ ! -f docker-compose.yaml ]; then \
		echo "Creating docker-compose.yaml from docker-compose.yaml.example..."; \
		cp docker-compose.yaml.example docker-compose.yaml; \
		echo "Done. Please review docker-compose.yaml."; \
	else \
		echo "docker-compose.yaml already exists."; \
		echo ""; \
		echo "=== Differences ==="; \
		diff -u docker-compose.yaml.example docker-compose.yaml || true; \
		echo ""; \
		echo "Please manually merge any needed changes."; \
	fi

# ==============================================================================
# Documentation
# ==============================================================================

.PHONY: docs-api
docs-api: ## Generate API documentation
	@command -v swag >/dev/null 2>&1 || { echo "Installing swag..."; go install github.com/swaggo/swag/cmd/swag@latest; }
	swag init -g cmd/pagemail/main.go -o docs/api

# ==============================================================================
# Clean
# ==============================================================================

.PHONY: clean
clean: ## Clean build artifacts
	rm -rf bin/
	rm -rf web/dist/
	rm -rf coverage.out coverage.html
	rm -rf tmp/

# ==============================================================================
# Pre-commit
# ==============================================================================

.PHONY: pre-commit-install
pre-commit-install: ## Install pre-commit hooks
	@command -v pre-commit >/dev/null 2>&1 || { echo "Installing pre-commit..."; pip install pre-commit; }
	pre-commit install

.PHONY: pre-commit-run
pre-commit-run: ## Run pre-commit on all files
	pre-commit run --all-files
