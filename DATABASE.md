# 📊 Database Management

PageMail 使用专业的数据库迁移系统（golang-migrate/migrate）来管理数据库模式变更。

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

## 📋 迁移历史

| 版本 | 日期 | 描述 |
|------|------|------|
| 20250902100001 | 2025-09-02 | 初始数据库模式 |
| 20250902100002 | 2025-09-02 | 邮箱验证功能 |

## 🔗 相关工具

- [golang-migrate/migrate](https://github.com/golang-migrate/migrate) - 迁移工具
- [GORM](https://gorm.io/) - Go ORM（仅用于业务逻辑）
- PostgreSQL - 数据库