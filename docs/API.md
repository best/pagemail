# PageMail API 文档

## 🔗 基础信息

- **Base URL**: `http://localhost:8080` (开发环境)
- **API Version**: v1
- **Content-Type**: `application/json`
- **认证方式**: JWT Bearer Token

## 📋 接口列表

### 系统接口

#### 健康检查
```http
GET /health
```

**响应示例（成功，HTTP 200）**:
```json
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

**响应示例（失败，HTTP 503）**:
```json
{
  "status": "unhealthy",
  "service": "pagemail", 
  "checks": {
    "database": "query_failed: connection refused",
    "smtp": "not_configured"
  },
  "version": "1.0.0"
}
```

---

### 认证接口

#### 用户注册
```http
POST /api/v1/auth/register
```

**请求参数**:
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

**参数说明**:
- `email`: 用户邮箱地址（必填，需要有效邮箱格式）
- `password`: 用户密码（必填，最少6位字符）

**响应示例（成功，HTTP 201）**:
```json
{
  "message": "Registration successful. Please check your email to verify your account.",
  "user": {
    "id": 1,
    "email": "user@example.com",
    "is_active": false,
    "email_verified": false,
    "daily_limit": 10,
    "monthly_limit": 300,
    "created_at": "2025-09-01T16:00:00Z",
    "updated_at": "2025-09-01T16:00:00Z"
  }
}
```

**错误响应（邮箱已存在，HTTP 409）**:
```json
{
  "error": "User already exists"
}
```

**错误响应（参数验证失败，HTTP 400）**:
```json
{
  "error": "Key: 'RegisterRequest.Email' Error:Field validation for 'Email' failed on the 'required' tag"
}
```

**注意事项**:
- 新注册用户默认为未激活状态 (`is_active: false`)
- 需要通过邮箱验证激活账户才能登录
- 系统会自动发送验证邮件到注册邮箱

#### 用户登录
```http
POST /api/v1/auth/login
```

**请求参数**:
```json
{
  "email": "user@example.com", 
  "password": "password123"
}
```

**参数说明**:
- `email`: 用户邮箱地址（必填）
- `password`: 用户密码（必填）

**响应示例（成功，HTTP 200）**:
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "email": "user@example.com",
    "is_active": true,
    "email_verified": true,
    "daily_limit": 10,
    "monthly_limit": 300,
    "created_at": "2025-09-01T16:00:00Z",
    "updated_at": "2025-09-01T16:00:00Z"
  }
}
```

**错误响应（邮箱未验证，HTTP 403）**:
```json
{
  "error": "Email not verified",
  "code": "EMAIL_NOT_VERIFIED",
  "message": "Please verify your email address before logging in"
}
```

**错误响应（账户被禁用，HTTP 401）**:
```json
{
  "error": "Account is deactivated"
}
```

**错误响应（凭证无效，HTTP 401）**:
```json
{
  "error": "Invalid credentials"
}
```

#### 邮箱验证
```http
GET /api/v1/auth/verify/:token
```

**URL参数**:
- `token`: 邮箱验证令牌（必填，通过邮件发送）

**响应示例（成功，HTTP 200）**:
```json
{
  "message": "Email verified successfully. You can now login.",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "email": "user@example.com",
    "is_active": true,
    "email_verified": true,
    "daily_limit": 10,
    "monthly_limit": 300,
    "created_at": "2025-09-01T16:00:00Z",
    "updated_at": "2025-09-01T16:00:00Z"
  }
}
```

**错误响应（无效令牌，HTTP 400）**:
```json
{
  "error": "Invalid or expired verification token"
}
```

#### 重发验证邮件
```http
POST /api/v1/auth/resend-verification
```

**请求参数**:
```json
{
  "email": "user@example.com"
}
```

**参数说明**:
- `email`: 用户邮箱地址（必填，需要有效邮箱格式）

**响应示例（成功，HTTP 200）**:
```json
{
  "message": "Verification email sent. Please check your inbox."
}
```

**错误响应（用户不存在或已验证，HTTP 404）**:
```json
{
  "error": "User not found or already verified"
}
```

**错误响应（发送频率限制，HTTP 429）**:
```json
{
  "error": "Too many verification emails sent. Please wait before requesting another."
}
```

---

### 用户接口

#### 获取用户资料
```http
GET /api/v1/user/profile
```

**请求头**:
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**响应示例（成功，HTTP 200）**:
```json
{
  "id": 1,
  "email": "user@example.com",
  "is_active": true,
  "email_verified": true,
  "daily_limit": 10,
  "monthly_limit": 300,
  "created_at": "2025-09-01T16:00:00Z",
  "updated_at": "2025-09-01T16:00:00Z"
}
```

**错误响应（未认证，HTTP 401）**:
```json
{
  "error": "User not authenticated"
}
```

**错误响应（用户不存在，HTTP 404）**:
```json
{
  "error": "User not found"
}
```

---

### 抓取接口

#### 提交抓取请求
```http
POST /api/v1/pages/scrape
```

**请求参数**:
```json
{
  "url": "https://example.com",
  "email": "recipient@example.com",
  "format": "html"
}
```

**参数说明**:
- `url`: 要抓取的网页URL（必填，需要有效的URL格式）
- `email`: 接收邮件的地址（必填，需要有效邮箱格式）
- `format`: 输出格式（必填，可选值：`html`、`pdf`、`screenshot`）

**请求头（可选）**:
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**响应示例（成功，HTTP 202）**:
```json
{
  "request_id": 12,
  "message": "Request accepted and is being processed",
  "status": "pending"
}
```

**错误响应（参数验证失败，HTTP 400）**:
```json
{
  "error": "Key: 'ScrapeRequest.URL' Error:Field validation for 'URL' failed on the 'url' tag"
}
```

**错误响应（游客配额超限，HTTP 429）**:
```json
{
  "error": "Daily limit exceeded for guests",
  "message": "Please register for higher limits",
  "limit": 1,
  "used": 1
}
```

**错误响应（用户配额超限，HTTP 429）**:
```json
{
  "error": "Daily limit exceeded",
  "limit": 10,
  "used": 10,
  "reset_time": 1725148800
}
```

**限制说明**:
- 游客用户：每日1次，每月5次
- 注册用户：每日10次，每月300次（默认值，可由管理员自定义）

**认证要求**:
- 游客用户可直接使用，无需认证
- 已注册用户建议携带 JWT Token 以享受更高配额

#### 查看请求历史
```http
GET /api/v1/pages/history
```

**请求头**:
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**响应示例（成功，HTTP 200）**:
```json
{
  "requests": [
    {
      "id": 12,
      "url": "https://example.com",
      "email": "recipient@example.com",
      "format": "html",
      "status": "completed",
      "created_at": "2025-09-01T16:00:00Z",
      "completed_at": "2025-09-01T16:01:30Z"
    },
    {
      "id": 11,
      "url": "https://news.ycombinator.com",
      "email": "recipient@example.com", 
      "format": "pdf",
      "status": "failed",
      "error_msg": "Failed to convert content: PDF generation failed",
      "created_at": "2025-09-01T15:30:00Z"
    }
  ],
  "total": 2
}
```

**错误响应（未认证，HTTP 401）**:
```json
{
  "error": "User not authenticated"
}
```

---

### 使用情况接口

#### 查看配额使用情况
```http
GET /api/v1/usage/
```

**请求头（可选）**:
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**响应示例（认证用户，HTTP 200）**:
```json
{
  "usage": {
    "type": "authenticated",
    "daily": {
      "used": 3,
      "limit": 10,
      "remaining": 7
    },
    "monthly": {
      "used": 25,
      "limit": 300,
      "remaining": 275
    }
  }
}
```

**响应示例（游客用户，HTTP 200）**:
```json
{
  "usage": {
    "type": "guest",
    "daily": {
      "used": 1,
      "limit": 1,
      "remaining": 0
    },
    "monthly": {
      "used": 3,
      "limit": 5,
      "remaining": 2
    }
  }
}
```

**错误响应（内部错误，HTTP 500）**:
```json
{
  "error": "Failed to retrieve usage information"
}
```

---

## 🚨 错误响应

### 标准错误格式
```json
{
  "error": "错误描述信息"
}
```

### HTTP 状态码

| 状态码 | 含义 | 示例场景 |
|--------|------|----------|
| `200` | 成功 | 请求处理成功 |
| `201` | 创建成功 | 用户注册成功 |
| `202` | 已接受 | 抓取请求已提交 |
| `400` | 请求错误 | 参数格式错误 |
| `401` | 未认证 | Token无效或过期 |
| `403` | 权限不足 | 访问受保护资源 |
| `404` | 未找到 | 用户或资源不存在 |
| `409` | 冲突 | 用户已存在 |
| `429` | 频率限制 | 超出配额限制 |
| `500` | 服务器错误 | 内部系统错误 |
| `503` | 服务不可用 | 依赖服务故障 |

### 频率限制错误示例

**注册用户限制（HTTP 429）**:
```json
{
  "error": "Daily limit exceeded",
  "limit": 10,
  "used": 10,
  "reset_time": 1725148800
}
```

**游客用户限制（HTTP 429）**:
```json
{
  "error": "Daily limit exceeded for guests",
  "message": "Please register for higher limits", 
  "limit": 1,
  "used": 1
}
```

---

## 🔧 请求示例

### cURL 示例

#### 用户注册
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }'
```

#### 邮箱验证
```bash
curl -X GET http://localhost:8080/api/v1/auth/verify/TOKEN_FROM_EMAIL
```

#### 重发验证邮件
```bash
curl -X POST http://localhost:8080/api/v1/auth/resend-verification \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com"
  }'
```

#### 用户登录
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }'
```

#### 提交抓取请求（需认证）
```bash
curl -X POST http://localhost:8080/api/v1/pages/scrape \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -d '{
    "url": "https://example.com",
    "email": "recipient@example.com",
    "format": "pdf"
  }'
```

#### 游客抓取请求（无需认证）
```bash
curl -X POST http://localhost:8080/api/v1/pages/scrape \
  -H "Content-Type: application/json" \
  -d '{
    "url": "https://news.ycombinator.com",
    "email": "myemail@example.com", 
    "format": "html"
  }'
```

### JavaScript/Fetch 示例

```javascript
// 用户登录
const loginResponse = await fetch('http://localhost:8080/api/v1/auth/login', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
  },
  body: JSON.stringify({
    email: 'user@example.com',
    password: 'password123'
  })
});

const { token } = await loginResponse.json();

// 提交抓取请求
const scrapeResponse = await fetch('http://localhost:8080/api/v1/pages/scrape', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
    'Authorization': `Bearer ${token}`
  },
  body: JSON.stringify({
    url: 'https://example.com',
    email: 'recipient@example.com',
    format: 'screenshot'
  })
});

const result = await scrapeResponse.json();
console.log('Request ID:', result.request_id);
```

---

## 🔀 工作流程

### 典型使用流程

**首次使用用户**:
1. **用户注册** → 填写邮箱和密码
2. **邮箱验证** → 点击验证邮件中的链接激活账户
3. **用户登录** → 获取访问令牌（JWT Token）
4. **提交抓取请求** → 获得请求ID和pending状态
5. **系统异步处理**：
   - 抓取网页内容
   - 转换为指定格式
   - 发送邮件附件
6. **查看请求历史** → 确认处理状态
7. **检查邮箱** → 接收处理结果

**游客用户**:
1. **直接提交抓取请求** → 无需注册，受限配额（每日1次）
2. **系统处理** → 同上
3. **检查邮箱** → 接收处理结果

### 状态流转

```
pending → processing → completed
                   ↘ → failed
```

- `pending`: 请求已提交，等待处理
- `processing`: 正在抓取和转换
- `completed`: 处理完成，邮件已发送
- `failed`: 处理失败，查看error_msg字段

---

## 📧 邮件通知

处理完成后，系统会自动发送包含以下内容的邮件：
- 原始URL信息
- 处理格式和时间
- 文件附件（HTML/PDF/PNG）
- 精美的HTML邮件模板

## 🛡️ 安全说明

- JWT Token 有效期：24小时
- 密码使用bcrypt加密存储
- 支持CORS跨域访问控制
- 实施速率限制防止滥用
- 所有敏感数据通过环境变量配置