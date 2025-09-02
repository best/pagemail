#!/bin/bash

# PageMail 环境变量生成脚本
# Usage: ./scripts/generate-env.sh [production|development]

MODE=${1:-development}

echo "🔧 PageMail Environment Generator"
echo "================================="
echo "模式: $MODE"
echo ""

# 生成强随机JWT密钥
generate_jwt_secret() {
    if command -v openssl &> /dev/null; then
        openssl rand -base64 64 | tr -d "\n="
    elif command -v head &> /dev/null && [ -f /dev/urandom ]; then
        head -c 48 /dev/urandom | base64 | tr -d "\n="
    else
        echo "$(date +%s)_$(whoami)_$(hostname)_random_jwt_key_change_this"
    fi
}

# 生成数据库密码
generate_db_password() {
    if command -v openssl &> /dev/null; then
        openssl rand -base64 32 | tr -d "\n="
    else
        echo "postgres_$(date +%s)"
    fi
}

JWT_SECRET=$(generate_jwt_secret)
DB_PASSWORD=$(generate_db_password)

# 根据模式设置不同的配置
if [ "$MODE" = "production" ]; then
    GIN_MODE="release"
    NODE_ENV="production"
    DB_SSLMODE="require"
    LOG_LEVEL="info"
    API_URL="https://your-domain.com"
    CORS_ORIGINS="https://your-domain.com"
else
    GIN_MODE="debug"
    NODE_ENV="development"
    DB_SSLMODE="disable"
    LOG_LEVEL="debug"
    API_URL="http://localhost:8080"
    CORS_ORIGINS="http://localhost:3000"
fi

# 生成 .env 文件
cat > .env << EOF
# PageMail Environment Configuration
# Generated on $(date)
# Mode: $MODE

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=$DB_PASSWORD
DB_NAME=pagemail
DB_SSLMODE=$DB_SSLMODE

# SMTP Email Configuration
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=your-email@gmail.com
SMTP_PASSWORD=your-app-password
SMTP_FROM_NAME=PageMail

# JWT Secret Key
JWT_SECRET=$JWT_SECRET

# Server Configuration
PORT=8080
GIN_MODE=$GIN_MODE

# Chrome/Chromium Configuration
CHROME_BIN=/usr/bin/chromium-browser
CHROME_PATH=/usr/bin/chromium-browser

# Frontend Configuration
NEXT_PUBLIC_API_URL=$API_URL
NODE_ENV=$NODE_ENV

# File Storage
FILES_DIR=files

# Logging
LOG_LEVEL=$LOG_LEVEL
LOG_FORMAT=json

# Security
CORS_ORIGINS=$CORS_ORIGINS
TRUSTED_PROXIES=127.0.0.1

# Rate Limiting
GUEST_DAILY_LIMIT=1
GUEST_MONTHLY_LIMIT=5
DEFAULT_DAILY_LIMIT=10
DEFAULT_MONTHLY_LIMIT=300

# Performance
DB_MAX_OPEN_CONNS=25
DB_MAX_IDLE_CONNS=5
DB_CONN_MAX_LIFETIME=1h
SCRAPER_TIMEOUT=60s
EMAIL_TIMEOUT=30s
EOF

echo "✅ .env 文件已生成"
echo ""
echo "⚠️  重要提醒:"
echo "1. 请修改 SMTP_USERNAME 和 SMTP_PASSWORD"
echo "2. 如果是Gmail，需要使用应用专用密码"
echo "3. 生产环境请修改 DB_PASSWORD"
if [ "$MODE" = "production" ]; then
echo "4. 请将 API_URL 和 CORS_ORIGINS 改为您的域名"
fi
echo ""
echo "🔍 运行环境检查:"
echo "   ./scripts/check-env.sh"
echo ""
echo "📖 详细配置说明请参考:"
echo "   - README.md"
echo "   - DEPLOYMENT.md"