# 📦 PageMail 部署指南

PageMail 是前后端合一的单体应用，使用 Docker 进行统一构建和部署。一个镜像包含完整的 Go 后端 API 和 Next.js 前端界面。

## 🚀 架构说明

- **单一应用**：前后端合一，无需分离部署
- **统一端口**：应用在 8080 端口同时提供 API 和 Web 界面
- **Docker 优化**：多阶段构建，生产镜像仅包含必需组件

## 🐳 Docker 构建

### 构建应用镜像

```bash
# 构建生产镜像
docker build -t pagemail:latest .

# 查看镜像大小
docker images pagemail:latest
```

**构建过程**：
1. **Stage 1**: 编译 Go 后端 (`golang:1.25-alpine`)
2. **Stage 2**: 构建 Next.js 前端 (`node:20-alpine`)  
3. **Stage 3**: 创建最小运行时镜像 (`alpine:latest`)

## 🏃‍♂️ Docker Run 部署

### 快速测试部署

适用于开发测试或快速验证，使用SQLite数据库：

```bash
# 单容器运行（仅用于测试）
docker run -d \
  --name pagemail-test \
  -p 8080:8080 \
  -e DB_HOST=sqlite \
  -e DB_NAME=pagemail.db \
  -e SMTP_HOST=smtp.gmail.com \
  -e SMTP_USERNAME=your@gmail.com \
  -e SMTP_PASSWORD=your_app_password \
  -e JWT_SECRET=your-super-secret-key \
  -v $(pwd)/files:/app/files \
  pagemail:latest

# 查看日志
docker logs -f pagemail-test

# 访问应用
echo "Frontend: http://localhost:8080"
echo "API: http://localhost:8080/api/v1"
echo "Health: http://localhost:8080/health"
```

### 使用外部PostgreSQL

```bash
# 连接到已有PostgreSQL数据库
docker run -d \
  --name pagemail-app \
  -p 8080:8080 \
  -e DB_HOST=your-postgres-host \
  -e DB_PORT=5432 \
  -e DB_USER=postgres \
  -e DB_PASSWORD=your_db_password \
  -e DB_NAME=pagemail \
  -e DB_SSLMODE=require \
  -e SMTP_HOST=smtp.gmail.com \
  -e SMTP_USERNAME=your@gmail.com \
  -e SMTP_PASSWORD=your_app_password \
  -e JWT_SECRET=your-super-secret-key \
  -v $(pwd)/files:/app/files \
  pagemail:latest
```

## 🔧 Docker Compose 部署（推荐）

### 完整生产部署

使用 Docker Compose 一键部署完整的应用栈：

```bash
# 1. 配置环境变量
cp .env.example .env
# 编辑 .env 文件，设置必需的配置

# 2. 启动所有服务
docker-compose up -d

# 3. 查看服务状态
docker-compose ps

# 4. 查看应用日志
docker-compose logs -f pagemail
```

**服务组成**：
- `pagemail-db`: PostgreSQL 16 数据库
- `pagemail`: 统一的应用服务（API + Web）

### 自定义 Docker Compose 配置

创建生产环境配置文件 `docker-compose.prod.yml`：

```yaml
version: '3.8'

services:
  pagemail-db:
    image: postgres:16-alpine
    container_name: pagemail-db-prod
    restart: unless-stopped
    environment:
      POSTGRES_DB: pagemail
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: ${DB_PASSWORD}
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./backups:/backups
    networks:
      - pagemail_network

  pagemail:
    build: .
    container_name: pagemail-app-prod
    restart: unless-stopped
    environment:
      - DB_HOST=pagemail-db
      - DB_PASSWORD=${DB_PASSWORD}
      - SMTP_HOST=${SMTP_HOST}
      - SMTP_USERNAME=${SMTP_USERNAME}
      - SMTP_PASSWORD=${SMTP_PASSWORD}
      - JWT_SECRET=${JWT_SECRET}
      - GIN_MODE=release
    ports:
      - "8080:8080"
    depends_on:
      - pagemail-db
    networks:
      - pagemail_network
    volumes:
      - ./files:/app/files

volumes:
  postgres_data:

networks:
  pagemail_network:
    driver: bridge
```

使用生产配置启动：

```bash
docker-compose -f docker-compose.prod.yml up -d
```

## ⚙️ 环境变量配置

### 必需配置

```bash
# 数据库配置
DB_HOST=pagemail-db           # Docker Compose中使用服务名
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_secure_password  # 生产环境必须修改
DB_NAME=pagemail
DB_SSLMODE=disable           # 生产环境建议用'require'

# SMTP 邮件配置（必填）
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=your@gmail.com
SMTP_PASSWORD=your_app_password    # Gmail需要应用专用密码
SMTP_FROM_NAME=PageMail

# JWT 密钥（必填，生产环境必须使用强密钥）
JWT_SECRET=your-super-secret-jwt-key-change-in-production
```

### 可选配置

```bash
# 运行模式
GIN_MODE=release             # 生产环境
PORT=8080                    # 服务端口

# 应用配置
FILES_DIR=files              # 文件存储目录
LOG_LEVEL=info               # 日志级别

# Chrome浏览器配置（Docker自动配置）
CHROME_BIN=/usr/bin/chromium-browser
CHROME_PATH=/usr/bin/chromium-browser
```

### Gmail SMTP 配置说明

1. **启用两步验证**：在Google账户设置中启用两步验证
2. **生成应用密码**：在安全设置中生成专用应用密码
3. **使用应用密码**：在 `SMTP_PASSWORD` 中使用应用密码，不是登录密码

其他邮件服务商配置：

```bash
# Outlook/Hotmail
SMTP_HOST=smtp.live.com
SMTP_PORT=587

# 163邮箱
SMTP_HOST=smtp.163.com
SMTP_PORT=465

# QQ邮箱  
SMTP_HOST=smtp.qq.com
SMTP_PORT=587
```

## 🔍 健康检查

### 服务状态检查

```bash
# 检查应用健康状态
curl http://localhost:8080/health

# 预期响应
{
  "status": "healthy",
  "service": "pagemail", 
  "checks": {
    "database": "connected",
    "smtp": "connected"
  },
  "version": "1.0.0"
}
```

### Docker 服务检查

```bash
# 查看容器状态
docker-compose ps

# 查看容器日志
docker-compose logs pagemail
docker-compose logs pagemail-db

# 查看资源使用
docker stats
```

### 数据库连接测试

```bash
# 连接到PostgreSQL检查
docker-compose exec pagemail-db psql -U postgres -d pagemail -c "SELECT COUNT(*) FROM users;"

# 检查数据库迁移状态
docker-compose exec pagemail go run cmd/migrate/main.go -action=status
```

## 🚨 故障排除

### 常见问题

1. **容器启动失败**
```bash
# 检查容器日志
docker-compose logs pagemail

# 检查环境变量配置
docker-compose config
```

2. **数据库连接失败**
```bash
# 检查数据库容器状态
docker-compose ps pagemail-db

# 测试数据库连接
docker-compose exec pagemail-db pg_isready -U postgres
```

3. **Chrome无法启动（PDF/截图功能异常）**
```bash
# 检查Chrome版本
docker-compose exec pagemail chromium-browser --version

# 查看Chrome相关错误
docker-compose logs pagemail | grep -i chrome
```

4. **邮件发送失败**
```bash
# 检查SMTP配置
docker-compose exec pagemail env | grep SMTP

# 测试SMTP连接
telnet smtp.gmail.com 587
```

5. **文件权限问题**
```bash
# 检查文件目录权限
ls -la files/

# 修复权限（如果需要）
sudo chown -R 1001:1001 files/
```

### 应用重启

```bash
# 重启应用容器
docker-compose restart pagemail

# 重启所有服务
docker-compose restart

# 完全重建和启动
docker-compose down
docker-compose up -d --build
```

## 📊 日志和监控

### 查看日志

```bash
# 实时查看应用日志
docker-compose logs -f pagemail

# 查看最近100行日志
docker-compose logs --tail 100 pagemail

# 查看数据库日志
docker-compose logs pagemail-db
```

### 性能监控

```bash
# 查看容器资源使用
docker stats --no-stream

# 查看磁盘使用
df -h
du -sh files/

# 数据库连接数监控
docker-compose exec pagemail-db psql -U postgres -d pagemail -c "SELECT COUNT(*) FROM pg_stat_activity;"
```

## 🔄 备份和恢复

### 数据库备份

```bash
# 创建备份脚本 backup.sh
#!/bin/bash
DATE=$(date +%Y%m%d_%H%M%S)
docker-compose exec -T pagemail-db pg_dump -U postgres pagemail > backup_pagemail_$DATE.sql
gzip backup_pagemail_$DATE.sql
echo "Backup created: backup_pagemail_$DATE.sql.gz"

# 执行备份
chmod +x backup.sh
./backup.sh
```

### 数据恢复

```bash
# 从备份恢复
gunzip backup_pagemail_20250903_120000.sql.gz
docker-compose exec -T pagemail-db psql -U postgres pagemail < backup_pagemail_20250903_120000.sql
```

### 文件备份

```bash
# 备份用户上传的文件
tar -czf files_backup_$(date +%Y%m%d).tar.gz files/

# 恢复文件
tar -xzf files_backup_20250903.tar.gz
```

## 🎯 生产部署清单

部署前检查清单：

- [ ] 修改默认数据库密码 (`DB_PASSWORD`)
- [ ] 设置强JWT密钥 (`JWT_SECRET`)
- [ ] 配置有效的SMTP邮箱信息
- [ ] 启用SSL模式 (`DB_SSLMODE=require`)
- [ ] 设置运行模式 (`GIN_MODE=release`)
- [ ] 配置防火墙只开放必要端口
- [ ] 设置定期备份计划
- [ ] 配置日志轮转
- [ ] 测试健康检查接口
- [ ] 验证邮件发送功能

## 🔗 相关链接

- [API 文档](API.md) - 完整的API接口说明
- [数据库文档](DATABASE.md) - 数据库架构和迁移管理
- [Docker官方文档](https://docs.docker.com/) - Docker使用指南
- [Docker Compose文档](https://docs.docker.com/compose/) - Compose配置参考