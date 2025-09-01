# PageMail 部署指南

## 🚀 快速部署

### 开发环境

1. **启动数据库**
```bash
docker-compose up -d db
```

2. **配置环境变量**
```bash
cp .env.example .env
# 编辑 .env 文件，设置数据库和SMTP配置
```

3. **启动后端服务**
```bash
go run main.go
```

4. **启动前端开发服务器**
```bash
cd web && npm run dev
```

### 生产环境部署

#### 方式一：Docker Compose（推荐）

```bash
# 配置环境变量
cp .env.example .env
vim .env  # 设置生产环境配置

# 启动所有服务
docker-compose up -d

# 查看日志
docker-compose logs -f backend
```

#### 方式二：分离式部署

1. **数据库部署**
```bash
# PostgreSQL 16
docker run -d \
  --name pagemail-db \
  -e POSTGRES_DB=pagemail \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=your_password \
  -v pagemail_data:/var/lib/postgresql/data \
  -p 5432:5432 \
  postgres:16-alpine
```

2. **后端部署**
```bash
# 构建二进制文件
go build -ldflags="-w -s" -o pagemail ./main.go

# 设置环境变量并运行
export DB_HOST=localhost
export DB_PASSWORD=your_password
export SMTP_USERNAME=your@email.com
export SMTP_PASSWORD=your_app_password
export JWT_SECRET=your-super-secret-key

./pagemail
```

3. **前端部署**
```bash
cd web
npm run build
npm start
```

## ⚙️ 环境变量配置

### 必需配置
```bash
# 数据库配置
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_db_password
DB_NAME=pagemail

# SMTP 邮件配置
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=your-email@gmail.com
SMTP_PASSWORD=your-app-password
SMTP_FROM_NAME=PageMail

# JWT 密钥
JWT_SECRET=your-super-secret-jwt-key-change-in-production
```

### 可选配置
```bash
# 服务器端口
PORT=8080

# Chrome 浏览器路径（Docker 环境自动配置）
CHROME_BIN=/usr/bin/chromium-browser
CHROME_PATH=/usr/bin/chromium-browser

# 前端配置
NEXT_PUBLIC_API_URL=http://localhost:8080
```

## 🔧 性能优化

### 1. 数据库优化
```sql
-- 创建索引优化查询性能
CREATE INDEX IF NOT EXISTS idx_requests_user_created ON requests(user_id, created_at);
CREATE INDEX IF NOT EXISTS idx_requests_status ON requests(status);
CREATE INDEX IF NOT EXISTS idx_requests_email ON requests(email);
```

### 2. Redis 缓存（可选升级）
```bash
# 添加 Redis 用于速率限制和会话管理
docker run -d \
  --name pagemail-redis \
  -p 6379:6379 \
  redis:7-alpine
```

### 3. Nginx 反向代理
```nginx
server {
    listen 80;
    server_name your-domain.com;

    # 后端 API
    location /api/ {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    # 健康检查
    location /health {
        proxy_pass http://localhost:8080;
    }

    # 前端应用
    location / {
        proxy_pass http://localhost:3000;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }
}
```

## 🐳 Docker 优化

### 多阶段构建优化
```dockerfile
# 使用更小的基础镜像
FROM golang:1.25-alpine AS backend-builder
# ... 构建步骤

FROM node:20-alpine AS frontend-builder  
# ... 前端构建

FROM alpine:3.18 AS runtime
# 只安装必需的运行时依赖
RUN apk --no-cache add ca-certificates chromium wget
# ... 其他配置
```

### Docker Compose 生产配置
```yaml
# docker-compose.prod.yml
version: '3.8'
services:
  backend:
    build: .
    restart: unless-stopped
    environment:
      - DB_HOST=db
      - JWT_SECRET=${JWT_SECRET}
    depends_on:
      - db
    networks:
      - pagemail_network

  db:
    image: postgres:16-alpine
    restart: unless-stopped
    volumes:
      - postgres_data:/var/lib/postgresql/data
    environment:
      - POSTGRES_DB=${DB_NAME}
      - POSTGRES_PASSWORD=${DB_PASSWORD}
    networks:
      - pagemail_network

volumes:
  postgres_data:

networks:
  pagemail_network:
    driver: bridge
```

## 🔐 安全配置

### 1. 防火墙设置
```bash
# 只开放必要端口
ufw allow 80
ufw allow 443
ufw allow 22
ufw enable
```

### 2. HTTPS 配置 (Let's Encrypt)
```bash
# 安装 Certbot
apt install certbot python3-certbot-nginx

# 获取证书
certbot --nginx -d your-domain.com

# 自动续期
crontab -e
# 添加: 0 12 * * * /usr/bin/certbot renew --quiet
```

### 3. 环境变量安全
```bash
# 使用 Docker secrets 或环境变量文件
echo "your_jwt_secret" | docker secret create jwt_secret -
```

## 📊 监控和日志

### 1. 健康检查
```bash
# 检查服务状态
curl http://localhost:8080/health

# 预期响应
{
  "status": "healthy",
  "service": "pagemail",
  "checks": {
    "database": "connected",
    "smtp": "configured"
  },
  "version": "1.0.0"
}
```

### 2. 日志管理
```bash
# Docker 日志
docker-compose logs -f backend

# 系统日志 (systemd)
journalctl -u pagemail -f
```

### 3. 性能监控
```bash
# 资源使用情况
docker stats

# 数据库连接监控
SELECT * FROM pg_stat_activity WHERE datname = 'pagemail';
```

## 🚨 故障排除

### 常见问题

1. **数据库连接失败**
```bash
# 检查数据库状态
docker-compose ps db
docker-compose logs db

# 测试连接
psql -h localhost -U postgres -d pagemail
```

2. **Chrome 无法启动**
```bash
# 检查 Chrome 依赖
chromium-browser --version

# Docker 中的权限问题
docker run --rm --security-opt seccomp=unconfined your-image
```

3. **邮件发送失败**
```bash
# 测试 SMTP 连接
telnet smtp.gmail.com 587

# 检查应用密码（Gmail）
# 需要启用两步验证并生成应用密码
```

4. **PDF 生成失败**
```bash
# 检查 wkhtmltopdf 安装
wkhtmltopdf --version

# Docker 中安装额外字体
RUN apk add --no-cache \
    fontconfig \
    ttf-dejavu \
    ttf-liberation
```

## 🔄 备份和恢复

### 数据库备份
```bash
# 备份
docker exec pagemail_db pg_dump -U postgres pagemail > backup.sql

# 恢复
docker exec -i pagemail_db psql -U postgres pagemail < backup.sql
```

### 文件备份
```bash
# 备份用户文件
tar -czf files_backup.tar.gz files/

# 恢复
tar -xzf files_backup.tar.gz
```

## 📈 扩展建议

### 水平扩展
- 使用负载均衡器分发请求
- 分离文件存储到对象存储（S3/MinIO）
- 使用 Redis 集群管理会话

### 功能扩展
- 添加 WebSocket 支持实时状态更新
- 集成队列系统（RabbitMQ/Redis）处理大量请求
- 添加管理后台界面
- 支持更多输出格式（EPUB、Markdown等）