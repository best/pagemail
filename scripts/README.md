# PageMail 脚本工具

这个目录包含用于 PageMail 项目管理和部署的实用脚本。

## 🔧 可用脚本

### `generate-env.sh` - 环境变量生成器

自动生成 `.env` 配置文件，包含安全的随机密钥。

**用法:**
```bash
# 生成开发环境配置
./scripts/generate-env.sh development

# 生成生产环境配置
./scripts/generate-env.sh production
```

**功能:**
- 自动生成强随机 JWT 密钥
- 生成安全的数据库密码
- 根据环境设置合适的默认值
- 包含所有必需和可选配置项

### `check-env.sh` - 环境配置检查器

验证 `.env` 文件中的配置是否正确完整。

**用法:**
```bash
./scripts/check-env.sh
```

**检查项目:**
- ✅ 必需环境变量是否设置
- ✅ JWT 密钥强度验证
- ✅ 数据库连接测试（如果 psql 可用）
- ✅ SMTP 服务器连接测试（如果 nc 可用）
- ✅ 可选配置项检查

**输出示例:**
```
🔍 PageMail Environment Configuration Checker
==============================================
✅ .env 文件存在

🔐 检查必需环境变量...
✅ DB_HOST: 已设置
✅ SMTP_USERNAME: 已设置
❌ JWT_SECRET: 未设置或使用默认值

🗄️  检查数据库连接...
✅ 数据库连接成功

📧 检查SMTP连接...
✅ SMTP服务器连接成功

🎉 所有必需环境变量配置正确！
```

## 📋 使用流程

### 新项目设置
```bash
# 1. 生成环境配置
./scripts/generate-env.sh development

# 2. 编辑邮箱配置
vim .env  # 设置 SMTP_USERNAME 和 SMTP_PASSWORD

# 3. 验证配置
./scripts/check-env.sh

# 4. 启动服务
docker-compose up -d db
go run main.go
```

### 生产部署
```bash
# 1. 生成生产配置
./scripts/generate-env.sh production

# 2. 完善配置
vim .env  # 设置所有生产环境参数

# 3. 检查配置
./scripts/check-env.sh

# 4. 部署
docker-compose -f docker-compose.prod.yml up -d
```

## 🔒 安全说明

- `generate-env.sh` 生成的密钥使用系统随机源
- JWT 密钥长度至少 64 字符
- 生产环境的数据库密码自动生成
- 所有生成的密钥都是唯一的

## 🛠️ 依赖工具

### 必需
- `bash` - 脚本执行环境

### 可选（增强功能）
- `openssl` - 用于生成强随机密钥
- `psql` - 用于数据库连接测试
- `nc` (netcat) - 用于网络连接测试

### 安装可选依赖

**Ubuntu/Debian:**
```bash
sudo apt-get update
sudo apt-get install openssl postgresql-client netcat-openbsd
```

**macOS:**
```bash
brew install openssl postgresql netcat
```

**Alpine Linux (Docker):**
```bash
apk add --no-cache openssl postgresql-client netcat-openbsd
```

## 🐛 故障排除

### 权限错误
```bash
chmod +x scripts/*.sh
```

### OpenSSL 不可用
脚本会自动回退到基础的密钥生成方法，但建议安装 OpenSSL 以获得更强的安全性。

### 数据库连接测试失败
确保：
1. PostgreSQL 服务正在运行
2. 网络连接正常
3. 数据库配置正确

### SMTP 连接测试失败
检查：
1. SMTP 服务器地址和端口
2. 防火墙设置
3. 网络连接