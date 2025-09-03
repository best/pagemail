# 📊 Database Management

PageMail 使用专业的数据库迁移系统（golang-migrate/migrate）来管理数据库模式变更，确保数据库结构的版本控制和安全部署。

## 📋 数据库架构概览

PageMail 使用 PostgreSQL 数据库，包含以下核心表：

- **users** - 用户账户信息
- **requests** - 页面抓取请求记录
- **email_configs** - SMTP 邮件配置
- **email_verifications** - 邮箱验证记录
- **schema_migrations** - 迁移版本管理（由migrate工具自动创建）

## 🚀 快速开始

### 运行迁移（生产环境）
```bash
# 启动应用（自动执行迁移）
./pagemail

# 或手动执行迁移
go run cmd/migrate/main.go -action=up
```

### 检查迁移状态
```bash
go run cmd/migrate/main.go -action=status
```

### 回滚迁移
```bash
# 回滚最后 1 个迁移
go run cmd/migrate/main.go -action=down -steps=1

# 回滚最后 2 个迁移
go run cmd/migrate/main.go -action=down -steps=2
```

## 📁 目录结构

```
internal/database/
├── database.go         # 数据库连接和配置
├── migrate.go          # 迁移管理器
└── migrations/         # 迁移文件
    ├── 20250902100001_initial_schema.up.sql
    ├── 20250902100001_initial_schema.down.sql
    ├── 20250902100002_add_email_verification.up.sql
    └── 20250902100002_add_email_verification.down.sql

cmd/migrate/
└── main.go            # CLI 迁移工具
```

## 🔧 创建新迁移

### 1. 生成迁移文件
```bash
# 创建新迁移（手动创建，按时间戳命名）
touch internal/database/migrations/$(date +%Y%m%d%H%M%S)_add_new_feature.up.sql
touch internal/database/migrations/$(date +%Y%m%d%H%M%S)_add_new_feature.down.sql
```

### 2. 编写迁移 SQL

**up.sql（正向迁移）**：
```sql
-- Add new table
CREATE TABLE new_table (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Add index
CREATE INDEX idx_new_table_name ON new_table(name);
```

**down.sql（回滚迁移）**：
```sql
-- Remove table and index
DROP INDEX IF EXISTS idx_new_table_name;
DROP TABLE IF EXISTS new_table;
```

## 🛡️ 最佳实践

### ✅ 安全的操作
- ✅ 新增表
- ✅ 新增字段（允许 NULL）
- ✅ 新增索引
- ✅ 新增约束（非严格）

### ⚠️ 需要谨慎的操作
- ⚠️ 修改字段类型（可能导致数据丢失）
- ⚠️ 删除字段（先标记废弃，后续版本删除）
- ⚠️ 重命名字段（分两步：新增+数据迁移+删除）

### ❌ 避免的操作
- ❌ 在生产环境删除表
- ❌ 删除有数据的字段
- ❌ 修改主键或外键结构

## 🔄 部署流程

### 开发环境
```bash
# 方式1：使用迁移系统（推荐）
go run main.go

# 方式2：使用 GORM AutoMigrate（快速原型）
# 需要手动调用 database.AutoMigrate()
```

### 生产环境
```bash
# 1. 备份数据库
pg_dump pagemail > backup_$(date +%Y%m%d_%H%M%S).sql

# 2. 运行迁移（应用启动时自动执行）
./pagemail

# 3. 验证迁移结果
go run cmd/migrate/main.go -action=status
```

## 🚨 故障恢复

### 迁移失败
```bash
# 检查状态
go run cmd/migrate/main.go -action=status

# 如果状态为 dirty，需要手动修复
# 1. 连接数据库检查问题
psql -d pagemail

# 2. 手动修复数据
# 3. 更新迁移表状态
UPDATE schema_migrations SET dirty = false;
```

### 紧急回滚
```bash
# 回滚到上一个版本
go run cmd/migrate/main.go -action=down -steps=1

# 验证回滚结果
go run cmd/migrate/main.go -action=status
```

## 📊 数据表结构详情

### users 表
用户账户和配额管理

```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,                    -- 用户唯一标识
    email VARCHAR(255) UNIQUE NOT NULL,       -- 邮箱地址（唯一）
    password VARCHAR(255) NOT NULL,           -- bcrypt加密密码
    is_active BOOLEAN DEFAULT false,          -- 账户状态（默认未激活）
    email_verified BOOLEAN DEFAULT false,     -- 邮箱验证状态
    email_verify_token VARCHAR(255),          -- 邮箱验证令牌
    email_verify_expires TIMESTAMP,           -- 验证令牌过期时间
    daily_limit INTEGER DEFAULT 10,           -- 日请求配额
    monthly_limit INTEGER DEFAULT 300,        -- 月请求配额
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

**索引**：
- `idx_users_email` - 邮箱快速查询
- `idx_users_email_verify_token` - 验证令牌查询

### requests 表
页面抓取请求记录

```sql
CREATE TABLE requests (
    id SERIAL PRIMARY KEY,                    -- 请求唯一标识
    user_id INTEGER REFERENCES users(id),     -- 用户ID（可为空，支持游客）
    url TEXT NOT NULL,                        -- 目标URL
    email VARCHAR(255) NOT NULL,              -- 接收邮箱
    format VARCHAR(20) DEFAULT 'html',        -- 输出格式：html/pdf/screenshot
    status VARCHAR(20) DEFAULT 'pending',     -- 状态：pending/processing/completed/failed
    file_path TEXT,                           -- 生成文件路径
    error_msg TEXT,                           -- 错误信息
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMP                    -- 完成时间
);
```

**索引**：
- `idx_requests_user_id` - 用户请求查询
- `idx_requests_status` - 状态筛选
- `idx_requests_created_at` - 时间排序

### email_configs 表
SMTP 邮件服务配置

```sql
CREATE TABLE email_configs (
    id SERIAL PRIMARY KEY,
    host VARCHAR(255) NOT NULL,               -- SMTP服务器地址
    port INTEGER NOT NULL,                    -- 端口号
    username VARCHAR(255) NOT NULL,           -- SMTP用户名
    password VARCHAR(255) NOT NULL,           -- SMTP密码
    from_name VARCHAR(255),                   -- 发件人名称
    is_active BOOLEAN DEFAULT true            -- 配置是否启用
);
```

### email_verifications 表
邮箱验证记录（防刷机制）

```sql
CREATE TABLE email_verifications (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL,              -- 验证邮箱
    ip_address VARCHAR(45) NOT NULL,          -- 请求IP（支持IPv6）
    sent_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

**索引**：
- `idx_email_verifications_email` - 邮箱查询
- `idx_email_verifications_ip` - IP查询
- `idx_email_verifications_sent_at` - 时间筛选

## 📋 迁移历史

| 版本 | 日期 | 描述 | 文件 |
|------|------|------|------|
| 20250902100001 | 2025-09-02 | 初始数据库模式 | `initial_schema.{up,down}.sql` |
| 20250902100002 | 2025-09-02 | 邮箱验证功能 | `add_email_verification.{up,down}.sql` |

## 🚀 数据库性能优化

### 索引策略
- **用户查询优化**：email字段建立唯一索引，支持快速登录验证
- **请求历史查询**：user_id + created_at 复合索引，优化用户历史记录查询
- **状态筛选**：status字段索引，支持管理员监控
- **时间范围查询**：created_at字段索引，支持配额统计

### 查询优化建议
```sql
-- ✅ 高效的用户配额查询
SELECT COUNT(*) FROM requests 
WHERE user_id = $1 AND created_at >= $2;

-- ✅ 高效的状态筛选
SELECT * FROM requests 
WHERE status = 'pending' 
ORDER BY created_at ASC LIMIT 10;

-- ❌ 避免全表扫描
SELECT * FROM requests WHERE error_msg LIKE '%timeout%';
```

### 连接池配置
```go
// 生产环境推荐配置
db.SetMaxOpenConns(25)      // 最大连接数
db.SetMaxIdleConns(5)       // 最大空闲连接数
db.SetConnMaxLifetime(5 * time.Minute)  // 连接最大生命周期
```

## 💾 数据备份与恢复

### 定期备份
```bash
# 每日备份脚本
#!/bin/bash
BACKUP_DIR="/backups/pagemail"
DATE=$(date +%Y%m%d_%H%M%S)
DB_NAME="pagemail"

# 创建备份目录
mkdir -p $BACKUP_DIR

# 执行备份
pg_dump -h localhost -U postgres -d $DB_NAME > $BACKUP_DIR/pagemail_$DATE.sql

# 压缩备份文件
gzip $BACKUP_DIR/pagemail_$DATE.sql

# 清理7天前的备份
find $BACKUP_DIR -name "*.sql.gz" -mtime +7 -delete

echo "Backup completed: pagemail_$DATE.sql.gz"
```

### 紧急恢复
```bash
# 从备份恢复数据库
psql -h localhost -U postgres -d pagemail < backup_20250903.sql

# 验证恢复结果
psql -h localhost -U postgres -d pagemail -c "SELECT COUNT(*) FROM users;"
```

## 🔍 监控与维护

### 数据库健康检查
```sql
-- 检查表大小
SELECT 
    schemaname,
    tablename,
    pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) as size
FROM pg_tables 
WHERE schemaname = 'public'
ORDER BY pg_total_relation_size(schemaname||'.'||tablename) DESC;

-- 检查活跃连接
SELECT COUNT(*) as active_connections 
FROM pg_stat_activity 
WHERE state = 'active';

-- 检查慢查询
SELECT query, mean_time, calls 
FROM pg_stat_statements 
ORDER BY mean_time DESC 
LIMIT 10;
```

### 定期维护任务
```bash
# 更新统计信息（每日执行）
psql -d pagemail -c "ANALYZE;"

# 清理过期验证记录（每周执行）
psql -d pagemail -c "DELETE FROM email_verifications WHERE created_at < NOW() - INTERVAL '7 days';"

# 清理失败请求的临时文件
find files/ -name "failed_*" -mtime +7 -delete
```

## 🌍 环境差异处理

### 开发环境
- 使用 Docker Compose 快速启动
- 允许使用 GORM AutoMigrate 快速原型开发
- 数据可以随时重置

### 生产环境
- 严格使用迁移系统
- 必须先备份后迁移
- 禁用 AutoMigrate 功能
- 使用连接池和监控

### 测试环境
```bash
# 快速重置测试数据库
dropdb pagemail_test
createdb pagemail_test
go run cmd/migrate/main.go -action=up
```

## 🔗 相关工具

- **[golang-migrate/migrate](https://github.com/golang-migrate/migrate)** - 数据库迁移工具
- **[GORM](https://gorm.io/)** - Go ORM（仅用于业务逻辑查询）
- **[PostgreSQL](https://www.postgresql.org/)** - 主数据库
- **[pgAdmin](https://www.pgadmin.org/)** - 数据库管理界面
- **[pg_stat_statements](https://www.postgresql.org/docs/current/pgstatstatements.html)** - 查询统计扩展