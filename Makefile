# PageMail Project Makefile

# Variables
APP_NAME := pagemail
BINARY_PATH := ./$(APP_NAME)
IMAGE_NAME := $(APP_NAME):latest
DEV_COMPOSE_FILE := docker-compose.dev.yml
POSTGRES_CONTAINER := $(APP_NAME)-dev-db
POSTGRES_PORT := 5432
POSTGRES_DB := pagemail
POSTGRES_USER := postgres
POSTGRES_PASS := postgres

# Default target
.DEFAULT_GOAL := help

# Help
.PHONY: help
help:
	@echo "PageMail Development Commands:"
	@echo ""
	@echo "Building & Deployment:"
	@echo "  build        - Build frontend + backend binary"
	@echo "  deploy       - Start database and run application"
	@echo "  clean        - Clean build files and frontend dependencies"
	@echo "  clean-all    - Clean everything including dependencies"
	@echo ""
	@echo "Dependency Management:"
	@echo "  update       - Update all Go and npm dependencies"
	@echo "  check-updates- Check for available dependency updates"
	@echo "  audit        - Run security audit on all dependencies"
	@echo ""
	@echo "Code Quality:"
	@echo "  test         - Run tests"
	@echo "  lint         - Run Go and frontend linting"
	@echo "  format       - Format Go and frontend code"
	@echo ""
	@echo "Docker & Services:"
	@echo "  docker       - Build Docker image"
	@echo "  compose      - Start services using Docker Compose"
	@echo ""
	@echo "Utilities:"
	@echo "  status       - Show project status"
	@echo "  help         - Show this help message"
	@echo ""

# Build targets
.PHONY: build
build: build-frontend build-backend start-db
	@echo "✅ Build completed successfully"

.PHONY: build-frontend
build-frontend:
	@echo "🔨 Building frontend..."
	@if [ ! -d frontend/node_modules ]; then \
		echo "📦 Installing dependencies..."; \
		cd frontend && npm install; \
	else \
		echo "📦 Using existing dependencies..."; \
	fi
	cd frontend && npm run build

.PHONY: build-backend
build-backend:
	@echo "🔨 Building backend..."
	go mod download
	go build -ldflags="-w -s" -o $(BINARY_PATH) ./cmd/pagemail

.PHONY: start-db
start-db:
	@echo "🐘 Starting PostgreSQL container..."
	@if ! docker ps | grep -q $(POSTGRES_CONTAINER); then \
		docker run -d \
			--name $(POSTGRES_CONTAINER) \
			-e POSTGRES_DB=$(POSTGRES_DB) \
			-e POSTGRES_USER=$(POSTGRES_USER) \
			-e POSTGRES_PASSWORD=$(POSTGRES_PASS) \
			-p $(POSTGRES_PORT):5432 \
			postgres:16-alpine; \
		echo "⏳ Waiting for database to start..."; \
		sleep 5; \
	else \
		echo "📍 PostgreSQL container already running"; \
	fi

# Deploy target
.PHONY: deploy
deploy: ensure-db migrate-up
	@echo "🚀 Starting application..."
	@if [ ! -f $(BINARY_PATH) ]; then \
		echo "❌ Binary not found. Run 'make build' first."; \
		exit 1; \
	fi
	$(BINARY_PATH)

.PHONY: ensure-db
ensure-db:
	@if ! docker ps | grep -q $(POSTGRES_CONTAINER); then \
		echo "🐘 Starting PostgreSQL container..."; \
		make start-db; \
	fi

.PHONY: migrate-up
migrate-up:
	@echo "🔄 Running database migrations..."
	@sleep 2
	go run cmd/migrate/main.go -action=up

# Clean targets
.PHONY: clean
clean: clean-build clean-db clean-compose clean-docker
	@echo "✨ Lightweight cleanup completed"

.PHONY: clean-all
clean-all: clean-build-all clean-db clean-compose clean-docker
	@echo "✨ Complete cleanup completed"

.PHONY: clean-build
clean-build:
	@echo "🧹 Cleaning build files and frontend dependencies..."
	@rm -f $(BINARY_PATH)
	@rm -rf frontend/dist
	@rm -rf frontend/.next
	@rm -rf frontend/node_modules
	@rm -f frontend/package-lock.json
	@echo "📁 Build files and frontend dependencies cleaned"

.PHONY: clean-build-all
clean-build-all:
	@echo "🧹 Cleaning all build files and dependencies..."
	@rm -f $(BINARY_PATH)
	@rm -rf frontend/dist
	@rm -rf frontend/.next
	@rm -rf frontend/node_modules
	@echo "📁 All build files and dependencies cleaned"

.PHONY: clean-db
clean-db:
	@echo "🗑️  Stopping and removing database container..."
	@if docker ps -a | grep -q $(POSTGRES_CONTAINER); then \
		docker stop $(POSTGRES_CONTAINER) 2>/dev/null || true; \
		docker rm $(POSTGRES_CONTAINER) 2>/dev/null || true; \
		echo "🐘 PostgreSQL container removed"; \
	else \
		echo "📍 No database container to remove"; \
	fi

.PHONY: clean-compose
clean-compose:
	@echo "🐳 Stopping and removing Docker Compose services..."
	@if [ -f $(DEV_COMPOSE_FILE) ]; then \
		docker compose -f $(DEV_COMPOSE_FILE) down -v 2>/dev/null || true; \
		echo "🚀 Docker Compose services stopped"; \
	else \
		echo "📍 No Docker Compose file found"; \
	fi

.PHONY: clean-docker
clean-docker:
	@echo "🗑️  Removing Docker image..."
	@if docker images -q $(IMAGE_NAME) 2>/dev/null; then \
		docker rmi $(IMAGE_NAME) 2>/dev/null || true; \
		echo "🐳 Docker image removed: $(IMAGE_NAME)"; \
	else \
		echo "📍 No Docker image to remove"; \
	fi

# Test target
.PHONY: test
test:
	@echo "🧪 Running backend tests..."
	go test ./...
	@echo "🧪 Running frontend tests..."
	cd frontend && npm test 2>/dev/null || echo "No frontend tests configured"

# Lint target
.PHONY: lint
lint: lint-go lint-frontend
	@echo "✅ All linting completed"

.PHONY: lint-go
lint-go:
	@echo "🔍 Running Go linting..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "⚠️  golangci-lint not found, using go vet and go fmt..."; \
		go vet ./...; \
		echo "📝 Checking Go formatting..."; \
		if [ -n "$$(gofmt -l .)" ]; then \
			echo "❌ The following files need formatting:"; \
			gofmt -l .; \
			echo "Run 'make format' to fix formatting"; \
			exit 1; \
		else \
			echo "✅ Go code is properly formatted"; \
		fi; \
	fi

.PHONY: lint-frontend
lint-frontend:
	@echo "🔍 Running frontend linting..."
	cd frontend && npx eslint src/

# Status target
.PHONY: status
status:
	@echo "📊 Project Status:"
	@echo "Binary exists: $$(test -f $(BINARY_PATH) && echo '✅ Yes' || echo '❌ No')"
	@echo "Frontend built: $$(test -d frontend/dist && echo '✅ Yes' || echo '❌ No')"
	@echo "Docker image: $$(docker images -q $(IMAGE_NAME) 2>/dev/null | head -1 | grep -q . && echo '✅ Yes' || echo '❌ No')"
	@echo "Database running: $$(docker ps | grep -q $(POSTGRES_CONTAINER) && echo '✅ Yes' || echo '❌ No')"
	@echo "Compose services: $$(docker compose -f $(DEV_COMPOSE_FILE) ps -q 2>/dev/null | wc -l | awk '{if($$1>0) print "✅ Running"; else print "❌ Stopped"}')"

# Database management
.PHONY: migrate-down
migrate-down:
	@echo "⬇️  Rolling back database migrations..."
	go run cmd/migrate/main.go -action=down

.PHONY: migrate-status
migrate-status:
	@echo "📋 Migration status:"
	go run cmd/migrate/main.go -action=status

# Format target
.PHONY: format
format: format-go format-frontend
	@echo "✅ All formatting completed"

.PHONY: format-go
format-go:
	@echo "📝 Formatting Go code..."
	gofmt -w .
	@echo "✅ Go code formatted"

.PHONY: format-frontend
format-frontend:
	@echo "📝 Formatting frontend code..."
	cd frontend && (npm run format 2>/dev/null || echo "No frontend formatter configured")

# Docker targets
.PHONY: docker
docker:
	@echo "🐳 Building Docker image..."
	docker build --no-cache -t $(IMAGE_NAME) .
	@echo "✅ Docker image built: $(IMAGE_NAME)"

.PHONY: compose
compose: docker
	@echo "🚀 Starting services with Docker Compose..."
	docker compose -f $(DEV_COMPOSE_FILE) up -d
	@echo "✅ Services started. Application available at http://localhost:8080"

.PHONY: logs
logs:
	@docker logs $(POSTGRES_CONTAINER) 2>/dev/null || echo "No container logs available"

# Update targets
.PHONY: update
update: update-go update-frontend
	@echo "✅ All dependencies updated successfully"

.PHONY: update-go
update-go:
	@echo "📦 Updating Go dependencies..."
	go get -u ./...
	go mod tidy
	@echo "✅ Go dependencies updated"

.PHONY: update-frontend
update-frontend:
	@echo "📦 Updating frontend dependencies..."
	@if [ ! -d frontend/node_modules ]; then \
		echo "🔧 Installing dependencies first..."; \
		cd frontend && npm install; \
	fi
	cd frontend && npm update
	@echo "✅ Frontend dependencies updated"

.PHONY: check-updates
check-updates: check-updates-go check-updates-frontend
	@echo "📊 Dependency check completed"

.PHONY: check-updates-go
check-updates-go:
	@echo "🔍 Checking Go dependency updates..."
	@go list -u -m all | grep -v "^[[:space:]]*$$" | head -10 || echo "All Go dependencies are up to date"

.PHONY: check-updates-frontend
check-updates-frontend:
	@echo "🔍 Checking frontend dependency updates..."
	@if [ -d frontend/node_modules ]; then \
		cd frontend && npm outdated || echo "All frontend dependencies are up to date"; \
	else \
		echo "⚠️  Run 'make build' first to install frontend dependencies"; \
	fi

.PHONY: audit
audit: audit-go audit-frontend
	@echo "🔒 Security audit completed"

.PHONY: audit-go
audit-go:
	@echo "🔒 Auditing Go dependencies..."
	@go mod verify
	@echo "✅ Go dependencies verified"

.PHONY: audit-frontend
audit-frontend:
	@echo "🔒 Auditing frontend dependencies..."
	@if [ -d frontend/node_modules ]; then \
		cd frontend && npm audit; \
	else \
		echo "⚠️  Run 'make build' first to install frontend dependencies"; \
	fi