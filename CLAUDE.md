# Pagemail Project Instructions

## Pre-commit Checklist (MANDATORY)

Before committing any code changes, **MUST** run local checks:

```bash
# Backend lint (required)
make lint-backend

# Frontend lint (if frontend changes)
make lint-frontend

# Full lint check
make lint

# Run tests
make test
```

**DO NOT** commit or push until all lint errors are resolved. This prevents CI failures on GitHub.

## Quick Start

```bash
# 1. Install dependencies
make install

# 2. Start development environment (full stack with database)
make dev-deploy

# 3. Check service status
make dev-status
make dev-health

# 4. View logs
make dev-logs-app
```

## Common Workflows

### Development Workflow
```bash
# After code changes, rebuild and restart
make dev-update

# View real-time logs
make dev-logs-app

# Run tests before commit
make lint
make test
```

### Production Deployment
```bash
# Build and push multi-arch image
make docker-buildx

# On production server: pull and update
make prod-pull
make prod-update
make prod-health
```

### Release Workflow
```bash
# 1. Bump version
make bump-version TYPE=patch  # or minor/major

# 2. Commit version file
git add VERSION
git commit -m "chore: bump version to $(cat VERSION)"

# 3. Create release tag
make release

# 4. Push to remote
git push && git push --tags

# 5. Build and push to registry
make docker-buildx
```

### Database Migration
```bash
# Create new migration
make migrate-new name=add_users_table

# Apply migrations
make migrate-up

# Check migration status
make migrate-status
```

## Makefile Commands Reference

Always use `make` commands for project operations. Run `make help` for full list.

### 帮助信息
| Command | Description |
|---------|-------------|
| `make help` | 显示帮助信息 |

### 依赖检查
| Command | Description |
|---------|-------------|
| `make check-deps` | 检查必需的依赖工具 |
| `make init` | 初始化项目依赖 |
| `make install` | 安装依赖（init 别名）|

### 开发命令
| Command | Description |
|---------|-------------|
| `make dev` | 启动开发环境（后端+前端）|
| `make dev-backend` | 启动后端热重载（需要 air）|
| `make dev-frontend` | 启动前端开发服务器 |
| `make run` | 本地运行服务 |

### 代码质量
| Command | Description |
|---------|-------------|
| `make fmt` | 格式化代码 |
| `make vet` | 运行 go vet 检查 |
| `make lint` | 运行全部 linter |
| `make lint-backend` | 运行 golangci-lint |
| `make lint-frontend` | 运行前端 linter |
| `make lint-fix` | 自动修复 lint 问题 |
| `make test` | 运行全部测试 |
| `make test-backend` | 运行后端测试 |
| `make test-frontend` | 运行前端单元测试 |
| `make test-e2e` | 运行 E2E 测试 |
| `make test-coverage` | 生成测试覆盖率报告 |

### 构建命令
| Command | Description |
|---------|-------------|
| `make build` | 构建全部 |
| `make build-backend` | 构建后端二进制 |
| `make build-frontend` | 构建前端 |

### Docker 构建
| Command | Description |
|---------|-------------|
| `make docker-build` | 构建 Docker 镜像（本地，单架构）|
| `make docker-build-no-cache` | 构建镜像（无缓存）|
| `make docker-buildx` | 构建并推送多架构镜像 |
| `make docker-push` | 推送镜像到仓库 |

### 开发环境 (docker-compose 完整栈)
| Command | Description |
|---------|-------------|
| `make dev-deploy` | 构建并启动开发环境（完整栈）|
| `make dev-up` | 启动开发环境（所有服务）|
| `make dev-down` | 停止开发环境 |
| `make dev-restart` | 重启开发环境 |
| `make dev-update` | 快速更新 pagemail（仅重建并重启）|
| `make dev-logs` | 查看所有日志 |
| `make dev-logs-app` | 查看 pagemail 日志 |
| `make dev-logs-db` | 查看 postgres 日志 |
| `make dev-health` | 开发环境健康检查 |
| `make dev-status` | 查看开发环境服务状态 |
| `make dev-clean` | 清理开发环境（包括 volumes 和镜像）|
| `make dev-reset` | 重置开发环境（删除 volumes，数据丢失！）|

### 生产环境 (仅 pagemail，外部数据库)
| Command | Description |
|---------|-------------|
| `make prod-deploy` | 构建并启动生产环境（仅 pagemail）|
| `make prod-up` | 启动生产环境（仅 pagemail）|
| `make prod-down` | 停止生产环境 |
| `make prod-restart` | 重启生产环境 |
| `make prod-update` | 快速更新 pagemail |
| `make prod-pull` | 拉取最新生产镜像 |
| `make prod-logs` | 查看生产环境日志 |
| `make prod-health` | 生产环境健康检查 |
| `make prod-status` | 查看生产环境服务状态 |
| `make prod-clean` | 清理生产环境（不包括 volumes）|

### 数据库迁移
| Command | Description |
|---------|-------------|
| `make migrate-new name=xxx` | 创建新迁移 |
| `make migrate-up` | 应用所有待执行迁移 |
| `make migrate-down steps=1` | 回滚迁移 |
| `make migrate-status` | 显示迁移状态 |
| `make migrate-lint` | 校验迁移文件（CI 用）|
| `make migrate-hash` | 重新计算迁移 hash |

### 部署辅助
| Command | Description |
|---------|-------------|
| `make env-merge` | 合并 .env.example 到 .env |
| `make compose-merge` | 合并 docker-compose.yaml.example |

### 文档生成
| Command | Description |
|---------|-------------|
| `make docs-api` | 生成 API 文档 |

### 版本管理
| Command | Description |
|---------|-------------|
| `make version` | 显示版本信息 |
| `make get-version` | 获取当前版本号 |
| `make bump-version TYPE=patch\|minor\|major` | 版本递增 |
| `make pre-release TYPE=rc\|beta\|alpha` | 创建预发布 Git 标签 |
| `make release` | 创建正式版本 Git 标签 |
| `make release-help` | 显示发布帮助 |

### 清理命令
| Command | Description |
|---------|-------------|
| `make clean` | 清理构建产物 |
| `make clean-all` | 清理所有（构建产物+Docker 镜像）|

### 预提交钩子
| Command | Description |
|---------|-------------|
| `make pre-commit-install` | 安装 pre-commit 钩子 |
| `make pre-commit-run` | 对所有文件运行 pre-commit |

## Code Style

- Go: Follow golangci-lint rules defined in `.golangci.yml`
- Vue/TS: Follow ESLint rules in `web/eslint.config.ts`
- Use `make fmt` to format code before committing

## Common Lint Issues

- `gocritic/rangeValCopy`: Use pointer or index in range loops for large structs
- `gocritic/hugeParam`: Pass large structs by pointer
- `errcheck`: Always handle error returns
- `goconst`: Extract repeated string literals to constants
- `gofmt`: Run `go fmt ./...` or `make fmt`
