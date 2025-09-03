# PageMail

一个用于抓取网页内容并发送到邮箱的多用户服务。支持将网页保存为HTML、PDF或截图格式。

## 功能特性

- 📄 **多种输出格式**: HTML文件、PDF文档、网页截图
- 🔐 **用户系统**: 支持用户注册登录，配额管理
- 📧 **邮件发送**: 自动将抓取的内容发送到指定邮箱  
- 🌐 **响应式前端**: 支持桌面、平板、手机多端访问
- 🐳 **容器化部署**: 使用Docker和docker-compose快速部署
- 🔄 **多架构支持**: 支持AMD64和ARM64架构

## 技术栈

- **后端**: Go + Gin + GORM + PostgreSQL
- **前端**: Next.js + React + TypeScript + Tailwind CSS  
- **网页抓取**: Chrome DevTools Protocol + HTTP客户端
- **部署**: Docker + GitHub Actions

## 快速开始

### 开发环境

1. 克隆项目
```bash
git clone https://github.com/best/pagemail.git
cd pagemail
```

2. 复制环境变量配置
```bash
cp .env.example .env
# 编辑 .env 文件，配置数据库和SMTP信息
```

3. 启动开发环境
```bash
# 启动数据库
docker-compose up -d db

# 启动后端服务
go run main.go

# 启动前端服务 (新终端)
cd web && npm install && npm run dev
```

### 生产部署

使用Docker Compose一键部署:

```bash
# 配置环境变量
cp .env.example .env
# 编辑 .env 文件

# 启动所有服务
docker-compose up -d
```

服务将在以下端口启动:
- 后端API: http://localhost:8080
- 前端界面: http://localhost:3000 (开发模式)
- PostgreSQL: localhost:5432

## 项目结构

```
pagemail/
├── main.go                 # 程序入口
├── internal/               # 内部包
│   ├── api/               # API处理器
│   ├── auth/              # 认证模块
│   ├── database/          # 数据库连接
│   ├── models/            # 数据模型
│   ├── scraper/           # 网页抓取
│   ├── converter/         # 文件转换
│   └── mailer/            # 邮件发送
├── migrations/            # 数据库迁移
├── web/                   # Next.js前端
│   ├── src/app/          # 页面路由
│   ├── src/components/   # React组件
│   └── src/lib/          # 工具函数
├── docker-compose.yml     # 开发环境
├── Dockerfile            # 生产镜像
└── .github/workflows/    # CI/CD配置
```

## API 接口

### 认证接口
- `POST /api/v1/auth/register` - 用户注册
- `POST /api/v1/auth/login` - 用户登录

### 页面抓取接口  
- `POST /api/v1/pages/scrape` - 提交抓取任务
- `GET /api/v1/pages/history` - 查看历史记录

### 系统接口
- `GET /health` - 健康检查

## 环境变量配置

### 必需配置
| 变量名 | 说明 | 默认值 | 备注 |
|--------|------|--------|------|
| `DB_HOST` | 数据库主机 | localhost | - |
| `DB_PORT` | 数据库端口 | 5432 | - |
| `DB_USER` | 数据库用户 | postgres | - |
| `DB_PASSWORD` | 数据库密码 | postgres | 生产环境必须修改 |
| `DB_NAME` | 数据库名称 | pagemail | - |
| `DB_SSLMODE` | SSL模式 | disable | 生产环境建议enable |
| `SMTP_HOST` | SMTP服务器 | smtp.gmail.com | 支持各邮件服务商 |
| `SMTP_PORT` | SMTP端口 | 587 | 或465(SSL) |
| `SMTP_USERNAME` | 邮箱账号 | - | **必填** |
| `SMTP_PASSWORD` | 邮箱密码 | - | **必填**，建议用应用密码 |
| `SMTP_FROM_NAME` | 发件人名称 | PageMail | - |
| `JWT_SECRET` | JWT密钥 | - | **必填**，生产环境用强密钥 |

### 可选配置
| 变量名 | 说明 | 默认值 | 备注 |
|--------|------|--------|------|
| `PORT` | 服务端口 | 8080 | - |
| `GIN_MODE` | 运行模式 | debug | 生产环境用release |
| `CHROME_BIN` | Chrome路径 | - | Docker自动配置 |
| `CHROME_PATH` | Chrome路径 | - | Docker自动配置 |
| `NEXT_PUBLIC_API_URL` | 前端API地址 | http://localhost:8080 | - |
| `NODE_ENV` | Node环境 | development | 生产用production |
| `FILES_DIR` | 文件存储目录 | files | - |
| `LOG_LEVEL` | 日志级别 | info | debug/info/warn/error |
| `CORS_ORIGINS` | CORS允许源 | http://localhost:3000 | 多个用逗号分隔 |

## 📚 文档

- [API 接口文档](docs/API.md) - 详细的API使用说明
- [部署指南](DEPLOYMENT.md) - 生产环境部署配置

## 🚀 快速体验

### 1. 启动开发环境
```bash
# 启动数据库
docker-compose up -d db

# 方式1: 自动生成环境配置（推荐）
./scripts/generate-env.sh development
# 然后编辑 .env 设置 SMTP 配置

# 方式2: 手动配置环境变量
cp .env.example .env
# 编辑 .env 设置所有必需配置

# 启动后端
go run main.go

# 启动前端（新终端）
cd web && npm install && npm run dev
```

### 2. 测试 API
```bash
# 健康检查
curl http://localhost:8080/health

# 游客抓取请求
curl -X POST http://localhost:8080/api/v1/pages/scrape \
  -H "Content-Type: application/json" \
  -d '{
    "url": "https://example.com",
    "email": "your@email.com",
    "format": "html"
  }'
```

### 3. 检查环境配置（推荐）
```bash
# 运行环境检查脚本
./scripts/check-env.sh
```

### 4. 访问前端
打开浏览器访问 http://localhost:3000 使用 Web 界面

## 🎯 核心特性

### ✨ 智能抓取策略
- **HTTP优先**: 快速抓取静态内容
- **Chrome后备**: 自动处理JavaScript渲染页面
- **域名识别**: 自动选择最佳抓取方式

### 📄 多格式支持
- **HTML**: 完整页面归档，相对链接转换
- **PDF**: 高质量文档生成，支持分页和页脚
- **Screenshot**: 全页面高清截图

### 👥 用户管理
- **游客模式**: 无需注册，每日1次限制
- **注册用户**: 更高配额，历史记录管理
- **JWT认证**: 安全的无状态认证

### 📧 邮件投递
- **HTML模板**: 精美的邮件设计
- **自动附件**: 处理结果自动发送
- **SMTP支持**: 兼容各种邮件服务商

## 🛠️ 技术亮点

- **异步处理**: 非阻塞请求处理，支持高并发
- **智能重试**: HTTP失败自动降级到Chrome抓取
- **资源优化**: 相对URL转绝对URL，离线可访问
- **安全防护**: bcrypt密码、JWT认证、CORS保护
- **监控友好**: 结构化日志、健康检查接口

## 🔄 更新日志

### v1.0.0 (2025-09-01)
- ✅ 完整的用户认证系统
- ✅ 智能网页抓取引擎
- ✅ 多格式文件转换
- ✅ 自动邮件发送
- ✅ 速率限制和配额管理
- ✅ 响应式Web界面
- ✅ Docker容器化部署
- ✅ CI/CD自动构建

## 🤝 贡献指南

1. Fork 项目
2. 创建功能分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启 Pull Request

## 📞 支持

- 🐛 [报告 Bug](https://github.com/best/pagemail/issues)
- 💡 [功能建议](https://github.com/best/pagemail/discussions)
- 📖 [查看文档](https://github.com/best/pagemail/wiki)

## 许可证

MIT License - 详见 [LICENSE](LICENSE) 文件