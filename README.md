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

| 变量名 | 说明 | 默认值 |
|--------|------|--------|
| `DB_HOST` | 数据库主机 | localhost |
| `DB_PORT` | 数据库端口 | 5432 |
| `DB_USER` | 数据库用户 | postgres |
| `DB_PASSWORD` | 数据库密码 | postgres |
| `DB_NAME` | 数据库名称 | pagemail |
| `SMTP_HOST` | SMTP服务器 | smtp.gmail.com |
| `SMTP_PORT` | SMTP端口 | 587 |
| `SMTP_USERNAME` | 邮箱账号 | - |
| `SMTP_PASSWORD` | 邮箱密码 | - |
| `JWT_SECRET` | JWT密钥 | - |

## 许可证

MIT License - 详见 [LICENSE](LICENSE) 文件