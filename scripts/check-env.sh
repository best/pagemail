#!/bin/bash

# PageMail 环境变量检查脚本
# Usage: ./scripts/check-env.sh

echo "🔍 PageMail Environment Configuration Checker"
echo "=============================================="

# 检查 .env 文件是否存在
if [ ! -f ".env" ]; then
    echo "❌ .env 文件不存在"
    echo "📝 请复制 .env.example 到 .env 并配置环境变量:"
    echo "   cp .env.example .env"
    exit 1
fi

echo "✅ .env 文件存在"

# 加载环境变量
source .env

# 检查必需的环境变量
REQUIRED_VARS=(
    "DB_HOST"
    "DB_PORT" 
    "DB_USER"
    "DB_PASSWORD"
    "DB_NAME"
    "SMTP_HOST"
    "SMTP_PORT"
    "SMTP_USERNAME"
    "SMTP_PASSWORD"
    "JWT_SECRET"
)

echo ""
echo "🔐 检查必需环境变量..."
ALL_SET=true

for var in "${REQUIRED_VARS[@]}"; do
    if [ -z "${!var}" ] || [ "${!var}" = "your-email@gmail.com" ] || [ "${!var}" = "your-app-password" ] || [ "${!var}" = "your-super-secret-jwt-key-change-this-in-production" ]; then
        echo "❌ $var: 未设置或使用默认值"
        ALL_SET=false
    else
        echo "✅ $var: 已设置"
    fi
done

# 检查JWT密钥强度
if [ ${#JWT_SECRET} -lt 32 ]; then
    echo "⚠️  JWT_SECRET 长度少于32字符，建议使用更强的密钥"
    ALL_SET=false
fi

# 检查数据库连接（如果 psql 可用）
echo ""
echo "🗄️  检查数据库连接..."
if command -v psql &> /dev/null; then
    if PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -c '\q' 2>/dev/null; then
        echo "✅ 数据库连接成功"
    else
        echo "❌ 数据库连接失败"
        echo "   请检查数据库是否运行并且连接信息正确"
        ALL_SET=false
    fi
else
    echo "ℹ️  未安装 psql，跳过数据库连接测试"
fi

# 检查SMTP连接（如果 nc 可用）
echo ""
echo "📧 检查SMTP连接..."
if command -v nc &> /dev/null; then
    if nc -z $SMTP_HOST $SMTP_PORT 2>/dev/null; then
        echo "✅ SMTP服务器连接成功"
    else
        echo "❌ SMTP服务器连接失败"
        echo "   请检查 SMTP_HOST 和 SMTP_PORT 配置"
        ALL_SET=false
    fi
else
    echo "ℹ️  未安装 nc，跳过SMTP连接测试"
fi

# 检查可选配置
echo ""
echo "⚙️  检查可选配置..."
OPTIONAL_VARS=(
    "PORT"
    "GIN_MODE"
    "FILES_DIR"
    "LOG_LEVEL"
)

for var in "${OPTIONAL_VARS[@]}"; do
    if [ -z "${!var}" ]; then
        echo "ℹ️  $var: 使用默认值"
    else
        echo "✅ $var: ${!var}"
    fi
done

echo ""
echo "=============================================="

if [ "$ALL_SET" = true ]; then
    echo "🎉 所有必需环境变量配置正确！"
    echo "🚀 可以启动 PageMail 服务"
    exit 0
else
    echo "❌ 配置检查失败，请修复上述问题"
    echo "📖 参考文档: README.md 和 DEPLOYMENT.md"
    exit 1
fi