# PageMail API 文档

## 🔗 基础信息

- **API Version**: v1
- **Content-Type**: `application/json`
- **认证方式**: JWT Bearer Token

## 📋 标准化错误响应

所有错误响应都遵循统一格式：

```json
{
  "error_code": 1002,
  "error_type": "AUTHENTICATION_ERROR",
  "message": "Invalid email or password",
  "details": "可选的详细信息"
}
```

### 错误码分类
- **1000-1999**: 认证相关错误
- **2000-2999**: 用户相关错误  
- **3000-3999**: 请求验证错误
- **4000-4999**: 配额限制错误
- **5000-5999**: 系统服务错误

## 📋 接口列表

### 系统接口

#### 健康检查
```http
GET /api/v1/health
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
  "error_code": 2002,
  "error_type": "USER_ERROR",
  "message": "User already exists with this email"
}
```

**错误响应（参数验证失败，HTTP 400）**:
```json
{
  "error_code": 3002,
  "error_type": "VALIDATION_ERROR",
  "message": "Request validation failed",
  "details": "Key: 'RegisterRequest.Password' Error:Field validation for 'Password' failed on the 'min' tag"
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
  "error_code": 1004,
  "error_type": "AUTHENTICATION_ERROR",
  "message": "Email address not verified"
}
```

**错误响应（账户被禁用，HTTP 401）**:
```json
{
  "error_code": 1005,
  "error_type": "AUTHENTICATION_ERROR",
  "message": "Account has been deactivated"
}
```

**错误响应（凭证无效，HTTP 401）**:
```json
{
  "error_code": 1002,
  "error_type": "AUTHENTICATION_ERROR",
  "message": "Invalid email or password"
}
```

#### 邮箱验证
```http
POST /api/v1/auth/verification
```

**请求参数**:
```json
{
  "token": "verification-token-from-email"
}
```

**参数说明**:
- `token`: 邮箱验证令牌（必填，通过邮件获取）

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
  "error_code": 1006,
  "error_type": "AUTHENTICATION_ERROR",
  "message": "Invalid or expired verification token"
}
```

#### 重发验证邮件
```http
POST /api/v1/auth/verification/resend
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
  "error_code": 2001,
  "error_type": "USER_ERROR",
  "message": "User not found",
  "details": "User not found or already verified"
}
```

**错误响应（发送频率限制，HTTP 429）**:
```json
{
  "error_code": 4004,
  "error_type": "QUOTA_ERROR",
  "message": "Too many requests"
}
```

---

### 用户接口

#### 获取用户资料
```http
GET /api/v1/users/{user_id}
```

**URL参数**:
- `user_id`: 用户ID（必填，只能访问自己的资料）

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
  "error_code": 1001,
  "error_type": "AUTHENTICATION_ERROR",
  "message": "Authentication required"
}
```

**错误响应（访问他人资源，HTTP 403）**:
```json
{
  "error_code": 1001,
  "error_type": "AUTHENTICATION_ERROR",
  "message": "Authentication required",
  "details": "Access denied: cannot access other user's resources"
}
```

#### 获取用户抓取历史
```http
GET /api/v1/users/{user_id}/scrapes
```

**URL参数**:
- `user_id`: 用户ID（必填，只能访问自己的抓取历史）

**请求头**:
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**响应示例（成功，HTTP 200）**:
```json
{
  "scrapes": [
    {
      "id": 12,
      "user_id": 1,
      "url": "https://example.com",
      "email": "recipient@example.com",
      "format": "html",
      "status": "completed",
      "created_at": "2025-09-01T16:00:00Z",
      "completed_at": "2025-09-01T16:01:30Z"
    },
    {
      "id": 11,
      "user_id": 1,
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
  "error_code": 1001,
  "error_type": "AUTHENTICATION_ERROR",
  "message": "Authentication required"
}
```

#### 获取用户使用情况
```http
GET /api/v1/users/{user_id}/usage
```

**URL参数**:
- `user_id`: 用户ID（必填，只能查看自己的使用情况）

**请求头**:
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**响应示例（成功，HTTP 200）**:
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

---

### 抓取接口

#### 创建抓取任务
```http
POST /api/v1/scrapes
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
  "message": "Scrape request accepted and is being processed",
  "status": "pending"
}
```

**错误响应（参数验证失败，HTTP 400）**:
```json
{
  "error_code": 3002,
  "error_type": "VALIDATION_ERROR",
  "message": "Request validation failed",
  "details": "Key: 'ScrapeRequest.URL' Error:Field validation for 'URL' failed on the 'url' tag"
}
```

**错误响应（游客配额超限，HTTP 429）**:
```json
{
  "error_code": 4001,
  "error_type": "QUOTA_ERROR",
  "message": "Daily request limit exceeded for guests - please register for higher limits",
  "used": 1,
  "limit": 1
}
```

**错误响应（用户配额超限，HTTP 429）**:
```json
{
  "error_code": 4001,
  "error_type": "QUOTA_ERROR",
  "message": "Daily request limit exceeded",
  "used": 10,
  "limit": 10,
  "reset_time": 1725148800
}
```

**限制说明**:
- 游客用户：每日1次，每月5次
- 注册用户：每日10次，每月300次（默认值，可由管理员自定义）

**认证要求**:
- 游客用户可直接使用，无需认证
- 已注册用户建议携带 JWT Token 以享受更高配额

#### 获取抓取任务详情
```http
GET /api/v1/scrapes/{scrape_id}
```

**URL参数**:
- `scrape_id`: 抓取任务ID（必填）

**请求头**:
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**响应示例（成功，HTTP 200）**:
```json
{
  "id": 12,
  "user_id": 1,
  "url": "https://example.com",
  "email": "recipient@example.com",
  "format": "html",
  "status": "completed",
  "created_at": "2025-09-01T16:00:00Z",
  "completed_at": "2025-09-01T16:01:30Z"
}
```

**错误响应（未认证，HTTP 401）**:
```json
{
  "error_code": 1001,
  "error_type": "AUTHENTICATION_ERROR",
  "message": "Authentication required"
}
```

**错误响应（任务不存在，HTTP 404）**:
```json
{
  "error_code": 2001,
  "error_type": "USER_ERROR",
  "message": "User not found",
  "details": "Scrape not found"
}
```

---

## 🚨 错误响应详情

### HTTP 状态码对照表

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

### 常见错误码说明

| 错误码 | 错误类型 | 消息 |
|--------|----------|------|
| 1001 | AUTHENTICATION_ERROR | Authentication required |
| 1002 | AUTHENTICATION_ERROR | Invalid email or password |
| 1004 | AUTHENTICATION_ERROR | Email address not verified |
| 1005 | AUTHENTICATION_ERROR | Account has been deactivated |
| 1006 | AUTHENTICATION_ERROR | Invalid or expired verification token |
| 2001 | USER_ERROR | User not found |
| 2002 | USER_ERROR | User already exists with this email |
| 3002 | VALIDATION_ERROR | Request validation failed |
| 4001 | QUOTA_ERROR | Daily request limit exceeded |
| 4002 | QUOTA_ERROR | Monthly request limit exceeded |
| 5001 | SYSTEM_ERROR | Internal server error |

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
curl -X POST http://localhost:8080/api/v1/auth/verification \
  -H "Content-Type: application/json" \
  -d '{
    "token": "verification-token-from-email"
  }'
```

#### 重发验证邮件
```bash
curl -X POST http://localhost:8080/api/v1/auth/verification/resend \
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

#### 获取用户资料
```bash
curl -X GET http://localhost:8080/api/v1/users/1 \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

#### 创建抓取任务（需认证）
```bash
curl -X POST http://localhost:8080/api/v1/scrapes \
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
curl -X POST http://localhost:8080/api/v1/scrapes \
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

// 获取用户资料
const profileResponse = await fetch('http://localhost:8080/api/v1/users/1', {
  headers: {
    'Authorization': `Bearer ${token}`
  }
});

// 创建抓取任务
const scrapeResponse = await fetch('http://localhost:8080/api/v1/scrapes', {
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
2. **邮箱验证** → 使用邮件中的token调用验证接口
3. **用户登录** → 获取访问令牌（JWT Token）
4. **创建抓取任务** → 获得任务ID和pending状态
5. **系统异步处理**：
   - 抓取网页内容
   - 转换为指定格式
   - 发送邮件附件
6. **查看抓取历史** → 确认处理状态
7. **检查邮箱** → 接收处理结果

**游客用户**:
1. **直接创建抓取任务** → 无需注册，受限配额（每日1次）
2. **系统处理** → 同上
3. **检查邮箱** → 接收处理结果

### 状态流转

```
pending → processing → completed
                   ↘ → failed
```

- `pending`: 任务已提交，等待处理
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
- 用户只能访问自己的资源
- 所有敏感数据通过环境变量配置

## 🔄 RESTful API 设计

本API遵循RESTful设计原则：
- 使用标准HTTP方法（GET, POST, PUT, DELETE）
- 资源导向的URL设计
- 统一的错误处理机制
- 标准化的状态码使用
- 一致的请求/响应格式