# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## 🚨 MANDATORY RESTRICTIONS

1. **NO signatures in commits** - Never include "Signed-off-by", "Co-authored-by", or similar attribution
2. **Commit format required** - All commits MUST use `<type>: <description>` format (feat:, fix:, docs:, etc.)
3. **API docs mandatory** - Update `docs/API.md` immediately when changing any API endpoint
4. **Database docs mandatory** - Update `docs/DATABASE.md` immediately when making any database schema changes
5. **Lowercase commit messages** - Start description with lowercase letter, no period at end
6. **Present tense commits** - Use "add feature" not "added feature"

## Essential Commands

### Backend Development
```bash
# Start the main application
go run main.go

# Database migrations
go run cmd/migrate/main.go -action=up     # Apply migrations
go run cmd/migrate/main.go -action=down   # Rollback migrations  
go run cmd/migrate/main.go -action=status # Check migration status

# Build the application
go build -o pagemail main.go

# Run with Docker Compose (includes database)
docker-compose up -d db    # Database only
docker-compose up -d       # Full stack
```

### Frontend Development
```bash
# Install dependencies
cd web && npm install

# Start development server with Turbopack
cd web && npm run dev

# Build for production
cd web && npm run build

# Run linting
cd web && npm run lint
```

### Environment Setup
```bash
# Generate environment configuration with secure defaults
./scripts/generate-env.sh development  # or production
./scripts/check-env.sh                 # Validate configuration
```

## Project Architecture

### Backend (Go)
- **Entry Point**: `main.go` - Loads environment, connects to database, starts Gin server
- **API Layer**: `internal/api/` - HTTP handlers and routing (Gin framework)
  - `router.go` - Route definitions and middleware
  - `auth.go` - Authentication endpoints (register/login)
  - `scrape.go` - Page scraping endpoints
  - `health.go` - Health check endpoint
- **Data Layer**: `internal/database/` - PostgreSQL connection and migrations (GORM)
- **Models**: `internal/models/` - Database models and schemas
- **Core Services**:
  - `internal/scraper/` - Web scraping (HTTP-first with Chrome fallback)
  - `internal/converter/` - Format conversion (HTML, PDF, screenshots)  
  - `internal/auth/` - JWT authentication and rate limiting
  - `internal/mailer/` - Email sending via SMTP

### Frontend (Next.js)
- **Framework**: Next.js 15 with App Router and Turbopack
- **UI**: React 19, TypeScript, Tailwind CSS 4
- **Structure**: 
  - `web/src/app/` - App router pages and layouts
  - `web/src/components/` - Reusable React components

### Database
- **Primary**: PostgreSQL with GORM ORM
- **Migrations**: Located in database migration system, managed via `cmd/migrate/`
- **Connection**: Configured via environment variables (DB_HOST, DB_PORT, etc.)

## Development Workflow

### Setting Up Development Environment
1. Copy `.env.example` to `.env` or use `./scripts/generate-env.sh development`
2. Start database: `docker-compose up -d db`  
3. Run migrations: `go run cmd/migrate/main.go -action=up`
4. Start backend: `go run main.go`
5. Start frontend: `cd web && npm install && npm run dev`

### Key Environment Variables
- **Database**: DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME, DB_SSLMODE
- **SMTP**: SMTP_HOST, SMTP_PORT, SMTP_USERNAME, SMTP_PASSWORD, SMTP_FROM_NAME
- **Security**: JWT_SECRET (must be strong for production)
- **Optional**: PORT, GIN_MODE, FILES_DIR, LOG_LEVEL, CORS_ORIGINS

### Web Scraping Strategy
The system uses intelligent scraping with HTTP-first approach:
1. **HTTP Client** - Fast scraping for static content
2. **Chrome Fallback** - Automatic fallback for JavaScript-heavy pages via Chrome DevTools Protocol
3. **Format Options** - HTML (with absolute URLs), PDF, or screenshots

### Authentication Flow  
- JWT-based authentication with bcrypt password hashing
- Rate limiting for API endpoints
- Guest mode (limited) and registered user mode (higher quotas)

## Testing and Quality

### Running Tests
```bash
# Backend tests (when available)
go test ./...

# Frontend linting
cd web && npm run lint
```

### Health Checks
- Backend health endpoint: `GET /health`
- Database connectivity verification
- SMTP server connectivity (via scripts)

## Git Commit Standards (MANDATORY)

**All commits MUST follow these strict rules:**

### Commit Message Format
```
<type>: <description>

# Examples:
feat: add user authentication system
fix: resolve database connection timeout
docs: update API documentation
refactor: optimize scraping performance
style: format code according to standards
test: add unit tests for auth module
chore: update dependencies
```

### Commit Types (Required)
- **feat:** New features
- **fix:** Bug fixes
- **docs:** Documentation changes
- **refactor:** Code refactoring without feature changes
- **style:** Code formatting, missing semicolons, etc.
- **test:** Adding or updating tests
- **chore:** Build process, dependency updates, etc.

### Commit Rules (ENFORCED)
1. **NO signatures or attribution** - Never include "Signed-off-by", "Co-authored-by", or similar
2. **Concise descriptions** - Keep messages short and descriptive
3. **Present tense** - Use "add feature" not "added feature" 
4. **Lowercase** - Start description with lowercase letter
5. **No period** - Don't end description with a period

### Invalid Examples (DO NOT USE)
```bash
# ❌ Missing type prefix
git commit -m "update user model"

# ❌ Contains signature
git commit -m "feat: add login system

Co-authored-by: Someone <email@example.com>"

# ❌ Too verbose
git commit -m "feat: implemented a comprehensive user authentication system with JWT tokens and password hashing using bcrypt algorithm"
```

### Valid Examples (CORRECT)
```bash
# ✅ Proper format
git commit -m "feat: add JWT authentication"
git commit -m "fix: handle database connection errors"
git commit -m "docs: update installation guide"
```

## API Documentation (MANDATORY)

**All API changes MUST be documented immediately:**

### Documentation Location
- **API Documentation**: `docs/API.md` - Complete API reference
- **Keep Updated**: API docs MUST be updated with every interface change

### API Documentation Rules (ENFORCED)
1. **Immediate Updates**: Update `docs/API.md` whenever you:
   - Add new endpoints
   - Modify existing endpoints
   - Change request/response formats
   - Update authentication requirements
   - Modify error responses

2. **Required Documentation Elements**:
   - **Endpoint**: HTTP method and URL path
   - **Description**: What the endpoint does
   - **Authentication**: Required auth type (JWT, guest, etc.)
   - **Request Format**: Headers, body schema, parameters
   - **Response Format**: Success and error response schemas
   - **Examples**: Complete request/response examples

3. **Documentation Standards**:
   - Use consistent formatting with existing docs
   - Include all possible error codes and messages
   - Provide realistic example data
   - Update table of contents if needed

### API Change Workflow
```bash
# 1. Make API changes to code
# 2. Update docs/API.md IMMEDIATELY  
# 3. Test the changes
# 4. Commit both code and docs together
git add internal/api/ docs/API.md
git commit -m "feat: add user profile endpoint"
```

### Documentation Validation
Before any API commit, verify:
- [ ] All new/changed endpoints documented in `docs/API.md`
- [ ] Request/response examples are accurate
- [ ] Error cases are documented
- [ ] Authentication requirements are clear

## Database Documentation (MANDATORY)

**All database schema changes MUST be documented immediately:**

### Documentation Location
- **Database Documentation**: `docs/DATABASE.md` - Complete database schema and migration guide
- **Keep Updated**: Database docs MUST be updated with every schema change

### Database Documentation Rules (ENFORCED)
1. **Immediate Updates**: Update `docs/DATABASE.md` whenever you:
   - Create new migrations (up/down SQL files)
   - Add new tables or modify existing tables
   - Add/remove/modify columns, indexes, constraints
   - Change data types or field properties
   - Add foreign key relationships

2. **Required Documentation Elements**:
   - **Schema Changes**: Document all table structure changes
   - **Migration History**: Update migration history table with new entries
   - **Index Strategy**: Document new indexes and their purpose
   - **Data Relationships**: Update table relationship descriptions
   - **Performance Impact**: Note any performance considerations

3. **Database Change Workflow**:
   ```bash
   # 1. Create migration files in internal/database/migrations/
   # 2. Update docs/DATABASE.md immediately with schema changes
   # 3. Test migration up/down commands
   # 4. Commit both migration files and documentation together
   git add internal/database/migrations/ docs/DATABASE.md
   git commit -m "feat: add user preferences table"
   ```

### Database Documentation Standards
- Document table structures with field types and constraints
- Include realistic examples of complex queries
- Update migration history table with every new migration
- Explain the purpose and impact of schema changes
- Include rollback considerations for production deployments