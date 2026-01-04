# ==============================================================================
# Pagemail Makefile
# ==============================================================================

.DEFAULT_GOAL := help

# ==============================================================================
# Variables
# ==============================================================================

VERSION_FILE := VERSION
VERSION ?= $(shell if [ -f $(VERSION_FILE) ] && [ -s $(VERSION_FILE) ]; then cat $(VERSION_FILE); else git describe --tags --always --dirty 2>/dev/null || echo "dev"; fi)
BUILD_TIME := $(shell date -u '+%Y-%m-%dT%H:%M:%SZ')
LDFLAGS := -ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME) -s -w"

DOCKER_IMAGE := astralor/pagemail
DOCKER_TAG ?= $(VERSION)
PLATFORMS := linux/amd64,linux/arm64

# ==============================================================================
# Help
# ==============================================================================

##@ 帮助信息
.PHONY: help
help: ## 显示帮助信息
	@awk 'BEGIN {FS = ":.*##"; printf "\n\033[1m使用方法:\033[0m\n  make \033[36m<target>\033[0m\n"} \
		/^##@/ {printf "\n\033[1m%s\033[0m\n", substr($$0, 5)} \
		/^[a-zA-Z0-9_-]+:.*##/ {printf "  \033[36m%-22s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# ==============================================================================
# Dependencies Check
# ==============================================================================

##@ 依赖检查
.PHONY: check-deps
check-deps: ## 检查必需的依赖工具
	@echo "检查依赖..."
	@missing=""; \
	for cmd in go node npm docker; do \
		command -v $$cmd >/dev/null 2>&1 || missing="$$missing $$cmd"; \
	done; \
	if [ -n "$$missing" ]; then \
		echo "缺少依赖:$$missing"; exit 1; \
	fi; \
	echo "所有依赖已安装 ✓"

.PHONY: init install
init: ## 初始化项目依赖
	go mod download
	cd web && npm install

install: init ## 安装依赖（init 别名）

# ==============================================================================
# Development
# ==============================================================================

##@ 开发命令
.PHONY: dev
dev: ## 启动开发环境（后端+前端）
	@$(MAKE) -j2 dev-backend dev-frontend

.PHONY: dev-backend
dev-backend: ## 启动后端热重载（需要 air）
	@command -v air >/dev/null 2>&1 || { echo "Installing air..."; go install github.com/air-verse/air@latest; }
	cd cmd/pagemail && air

.PHONY: dev-frontend
dev-frontend: ## 启动前端开发服务器
	cd web && npm run dev

.PHONY: run
run: build ## 本地运行服务
	./bin/pagemail

# ==============================================================================
# Code Quality
# ==============================================================================

##@ 代码质量
.PHONY: fmt
fmt: ## 格式化代码
	go fmt ./...
	cd web && npm run format

.PHONY: vet
vet: ## 运行 go vet 检查
	go vet ./...

.PHONY: lint
lint: lint-backend lint-frontend ## 运行全部 linter

.PHONY: lint-backend
lint-backend: ## 运行 golangci-lint
	@command -v golangci-lint >/dev/null 2>&1 || { echo "Installing golangci-lint..."; go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; }
	golangci-lint run ./...

.PHONY: lint-frontend
lint-frontend: ## 运行前端 linter
	cd web && npm run lint

.PHONY: lint-fix
lint-fix: ## 自动修复 lint 问题
	golangci-lint run --fix ./...
	cd web && npm run lint -- --fix

.PHONY: test
test: test-backend test-frontend ## 运行全部测试

.PHONY: test-backend
test-backend: ## 运行后端测试
	go test -race -covermode=atomic ./...

.PHONY: test-frontend
test-frontend: ## 运行前端单元测试
	cd web && npm run test:unit

.PHONY: test-e2e
test-e2e: ## 运行 E2E 测试
	cd web && npm run test:e2e

.PHONY: test-coverage
test-coverage: ## 生成测试覆盖率报告
	go test -race -coverprofile=coverage.out -covermode=atomic ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "覆盖率报告已生成: coverage.html"

# ==============================================================================
# Build
# ==============================================================================

##@ 构建命令
.PHONY: build
build: build-frontend build-backend ## 构建全部

.PHONY: build-backend
build-backend: ## 构建后端二进制
	CGO_ENABLED=0 go build $(LDFLAGS) -o bin/pagemail ./cmd/pagemail

.PHONY: build-frontend
build-frontend: ## 构建前端
	cd web && npm run build

# ==============================================================================
# Docker Build (CI/CD)
# ==============================================================================

##@ Docker 构建
.PHONY: docker-build
docker-build: ## 构建 Docker 镜像（本地，单架构）
	docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) -f deploy/Dockerfile .

.PHONY: docker-build-no-cache
docker-build-no-cache: ## 构建镜像（无缓存）
	docker build --no-cache -t $(DOCKER_IMAGE):$(DOCKER_TAG) -f deploy/Dockerfile .

.PHONY: docker-buildx
docker-buildx: ## 构建并推送多架构镜像
	docker buildx build \
		--platform $(PLATFORMS) \
		--tag $(DOCKER_IMAGE):$(DOCKER_TAG) \
		--tag $(DOCKER_IMAGE):latest \
		--push \
		-f deploy/Dockerfile .

.PHONY: docker-push
docker-push: ## 推送镜像到仓库
	docker push $(DOCKER_IMAGE):$(DOCKER_TAG)

# ==============================================================================
# Development Environment (docker-compose full stack)
# ==============================================================================

COMPOSE_FILE ?= docker-compose.yaml
# 开发环境镜像名称（需与 docker-compose.yaml 中一致）
DEV_IMAGE ?= pagemail:latest

##@ 开发环境
.PHONY: dev-deploy
dev-deploy: ## 构建并启动开发环境（完整栈）
	docker build -t $(DEV_IMAGE) -f deploy/Dockerfile .
	docker compose -f $(COMPOSE_FILE) up -d

.PHONY: dev-up
dev-up: ## 启动开发环境（所有服务）
	docker compose -f $(COMPOSE_FILE) up -d

.PHONY: dev-down
dev-down: ## 停止开发环境
	docker compose -f $(COMPOSE_FILE) down

.PHONY: dev-restart
dev-restart: ## 重启开发环境
	docker compose -f $(COMPOSE_FILE) restart

.PHONY: dev-update
dev-update: ## 快速更新 pagemail（仅重建并重启）
	docker build -t $(DEV_IMAGE) -f deploy/Dockerfile .
	docker compose -f $(COMPOSE_FILE) up -d --no-deps --force-recreate pagemail

.PHONY: dev-logs
dev-logs: ## 查看所有日志
	docker compose -f $(COMPOSE_FILE) logs -f

.PHONY: dev-logs-app
dev-logs-app: ## 查看 pagemail 日志
	docker compose -f $(COMPOSE_FILE) logs -f pagemail

.PHONY: dev-logs-db
dev-logs-db: ## 查看 postgres 日志
	docker compose -f $(COMPOSE_FILE) logs -f postgres

.PHONY: dev-health
dev-health: ## 开发环境健康检查
	@curl -fsS http://127.0.0.1:8080/v1/health && echo " ✓ healthy" || echo " ✗ unhealthy"

.PHONY: dev-status
dev-status: ## 查看开发环境服务状态
	docker compose -f $(COMPOSE_FILE) ps

.PHONY: dev-clean
dev-clean: ## 清理开发环境（包括 volumes 和镜像）
	docker compose -f $(COMPOSE_FILE) down -v --rmi all --remove-orphans

.PHONY: dev-reset
dev-reset: ## 重置开发环境（删除 volumes，数据丢失！）
	@echo "⚠️  警告: 将删除所有数据卷，数据将丢失！"
	@read -p "确认继续? [y/N] " confirm; \
	if [ "$$confirm" = "y" ] || [ "$$confirm" = "Y" ]; then \
		docker compose -f $(COMPOSE_FILE) down -v --remove-orphans; \
		echo "开发环境已重置"; \
	else \
		echo "操作已取消"; \
	fi

# ==============================================================================
# Production Environment (pagemail only, external DB)
# ==============================================================================

# 生产环境镜像（默认使用 registry 镜像）
PROD_IMAGE ?= $(DOCKER_IMAGE):$(DOCKER_TAG)

##@ 生产环境
.PHONY: prod-deploy
prod-deploy: ## 构建并启动生产环境（仅 pagemail）
	docker build -t $(DEV_IMAGE) -f deploy/Dockerfile .
	docker compose -f $(COMPOSE_FILE) up -d --no-deps pagemail

.PHONY: prod-up
prod-up: ## 启动生产环境（仅 pagemail）
	docker compose -f $(COMPOSE_FILE) up -d --no-deps pagemail

.PHONY: prod-down
prod-down: ## 停止生产环境
	docker compose -f $(COMPOSE_FILE) stop pagemail

.PHONY: prod-restart
prod-restart: ## 重启生产环境
	docker compose -f $(COMPOSE_FILE) restart pagemail

.PHONY: prod-update
prod-update: ## 快速更新 pagemail
	docker build -t $(DEV_IMAGE) -f deploy/Dockerfile .
	docker compose -f $(COMPOSE_FILE) up -d --no-deps --force-recreate pagemail

.PHONY: prod-pull
prod-pull: ## 拉取最新生产镜像
	docker pull $(PROD_IMAGE)
	docker tag $(PROD_IMAGE) $(DEV_IMAGE)

.PHONY: prod-logs
prod-logs: ## 查看生产环境日志
	docker compose -f $(COMPOSE_FILE) logs -f pagemail

.PHONY: prod-health
prod-health: ## 生产环境健康检查
	@curl -fsS http://127.0.0.1:8080/v1/health && echo " ✓ healthy" || echo " ✗ unhealthy"

.PHONY: prod-status
prod-status: ## 查看生产环境服务状态
	docker compose -f $(COMPOSE_FILE) ps pagemail

.PHONY: prod-clean
prod-clean: ## 清理生产环境（不包括 volumes）
	docker compose -f $(COMPOSE_FILE) rm -sf pagemail
	-docker rmi $(DEV_IMAGE) 2>/dev/null || true

# ==============================================================================
# Deployment Helpers
# ==============================================================================

##@ 部署辅助
.PHONY: env-merge
env-merge: ## 合并 .env.example 到 .env
	@echo "检查 .env 配置..."
	@if [ ! -f .env ]; then \
		echo "从 .env.example 创建 .env..."; \
		cp .env.example .env; \
		echo "完成，请编辑 .env 配置"; \
	else \
		echo "对比 .env.example 与 .env..."; \
		echo ""; \
		tmp_example=$$(mktemp); tmp_env=$$(mktemp); \
		grep -E '^[A-Z_]+=' .env.example | cut -d= -f1 | sort > "$$tmp_example"; \
		grep -E '^[A-Z_]+=' .env | cut -d= -f1 | sort > "$$tmp_env"; \
		echo "=== .env.example 中存在但 .env 缺失的配置 ==="; \
		comm -23 "$$tmp_example" "$$tmp_env" || true; \
		echo ""; \
		echo "=== .env 中的自定义配置（不在 .env.example 中）==="; \
		comm -13 "$$tmp_example" "$$tmp_env" || true; \
		rm -f "$$tmp_example" "$$tmp_env"; \
	fi

.PHONY: compose-merge
compose-merge: ## 合并 docker-compose.yaml.example
	@echo "检查 docker-compose.yaml 配置..."
	@if [ ! -f docker-compose.yaml ]; then \
		echo "从 docker-compose.yaml.example 创建..."; \
		cp docker-compose.yaml.example docker-compose.yaml; \
		echo "完成，请检查 docker-compose.yaml"; \
	else \
		echo "docker-compose.yaml 已存在"; \
		echo ""; \
		echo "=== 差异 ==="; \
		diff -u docker-compose.yaml.example docker-compose.yaml || true; \
	fi

# ==============================================================================
# Documentation
# ==============================================================================

##@ 文档生成
.PHONY: docs-api
docs-api: ## 生成 API 文档
	@command -v swag >/dev/null 2>&1 || { echo "Installing swag..."; go install github.com/swaggo/swag/cmd/swag@latest; }
	swag init -g cmd/pagemail/main.go -o docs/api

# ==============================================================================
# Version Management
# ==============================================================================

##@ 版本管理
.PHONY: version
version: ## 显示版本信息
	@echo "Version: $(VERSION)"
	@echo "Docker Image: $(DOCKER_IMAGE):$(DOCKER_TAG)"

.PHONY: get-version
get-version: ## 获取当前版本号
	@echo $(VERSION)

.PHONY: bump-version
bump-version: ## 版本递增 (TYPE=patch|minor|major)
	@if [ -z "$(TYPE)" ]; then echo "Usage: make bump-version TYPE=patch|minor|major"; exit 1; fi
	@current=$$(cat $(VERSION_FILE) 2>/dev/null || echo "0.0.0"); \
	current=$${current#v}; \
	major=$$(echo "$$current" | cut -d. -f1); \
	minor=$$(echo "$$current" | cut -d. -f2); \
	patch=$$(echo "$$current" | cut -d. -f3); \
	major=$${major:-0}; minor=$${minor:-0}; patch=$${patch:-0}; \
	case "$(TYPE)" in \
		major) major=$$((major+1)); minor=0; patch=0 ;; \
		minor) minor=$$((minor+1)); patch=0 ;; \
		patch) patch=$$((patch+1)) ;; \
		*) echo "TYPE 必须为 patch|minor|major"; exit 1 ;; \
	esac; \
	new="$$major.$$minor.$$patch"; \
	echo "$$new" > $(VERSION_FILE); \
	echo "版本更新: $$current -> $$new"

.PHONY: pre-release
pre-release: ## 创建预发布 Git 标签 (TYPE=rc|beta|alpha)
	@if [ -z "$(TYPE)" ]; then echo "Usage: make pre-release TYPE=rc|beta|alpha"; exit 1; fi
	@version=$$(cat $(VERSION_FILE) 2>/dev/null || echo "0.0.0"); \
	version=$${version#v}; \
	base=$${version%%-*}; \
	case "$(TYPE)" in \
		rc|beta|alpha) ;; \
		*) echo "TYPE 必须为 rc|beta|alpha"; exit 1 ;; \
	esac; \
	tag="v$$base-$(TYPE).1"; \
	existing=$$(git tag -l "v$$base-$(TYPE).*" --sort=version:refname | tail -1); \
	if [ -n "$$existing" ]; then \
		n=$${existing##*.}; n=$$((n+1)); \
		tag="v$$base-$(TYPE).$$n"; \
	fi; \
	git tag -a "$$tag" -m "Pre-release $$tag"; \
	echo "创建预发布标签: $$tag"

.PHONY: release
release: ## 创建正式版本 Git 标签
	@version=$$(cat $(VERSION_FILE) 2>/dev/null); \
	if [ -z "$$version" ]; then echo "VERSION 文件不存在或为空"; exit 1; fi; \
	version=$${version#v}; \
	tag="v$$version"; \
	if git rev-parse "$$tag" >/dev/null 2>&1; then \
		echo "标签 $$tag 已存在"; exit 1; \
	fi; \
	git tag -a "$$tag" -m "Release $$tag"; \
	echo "创建发布标签: $$tag"

.PHONY: release-help
release-help: ## 显示发布帮助
	@echo "版本发布流程:"
	@echo "  1. make bump-version TYPE=patch|minor|major  # 递增版本号"
	@echo "  2. git add VERSION && git commit -m 'Bump version'"
	@echo "  3. make release                              # 创建 Git 标签"
	@echo "  4. git push && git push --tags               # 推送到远程"
	@echo ""
	@echo "预发布流程:"
	@echo "  1. make pre-release TYPE=rc|beta|alpha       # 创建预发布标签"
	@echo "  2. git push --tags                           # 推送标签"

# ==============================================================================
# Clean
# ==============================================================================

##@ 清理命令
.PHONY: clean
clean: ## 清理构建产物
	rm -rf bin/
	rm -rf web/dist/
	rm -rf coverage.out coverage.html
	rm -rf tmp/

.PHONY: clean-all
clean-all: clean ## 清理所有（构建产物+Docker 镜像）
	-docker rmi $(DOCKER_IMAGE):$(DOCKER_TAG) 2>/dev/null || true
	-docker rmi $(DOCKER_IMAGE):latest 2>/dev/null || true

# ==============================================================================
# Pre-commit
# ==============================================================================

##@ 预提交钩子
.PHONY: pre-commit-install
pre-commit-install: ## 安装 pre-commit 钩子
	@command -v pre-commit >/dev/null 2>&1 || { echo "Installing pre-commit..."; pip install pre-commit; }
	pre-commit install

.PHONY: pre-commit-run
pre-commit-run: ## 对所有文件运行 pre-commit
	pre-commit run --all-files
