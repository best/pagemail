# Pagemail

Web page capture and delivery service. Capture web pages as PDF, HTML, or screenshots and deliver them via email or webhook.

## Features

- **Web Page Capture**: Headless browser capture with JavaScript rendering
- **Multiple Output Formats**: PDF, HTML (single file), Screenshot (PNG)
- **Flexible Delivery**: Email (SMTP) and Webhook with file attachments
- **User Management**: Role-based access control (admin/user)
- **SSRF Protection**: Built-in security against server-side request forgery
- **Storage Options**: Local filesystem or S3/MinIO compatible object storage

## Quick Start

### Using Docker

```bash
# Pull the image
docker pull astralor/pagemail:latest

# Run with default settings
docker run -d -p 8080:8080 --env-file .env astralor/pagemail:latest
```

### Using Docker Compose

```bash
# Copy example files
cp .env.example .env
cp docker-compose.yaml.example docker-compose.yaml

# Edit configuration
vim .env

# Start services
docker-compose up -d
```

### Building from Source

```bash
# Install dependencies
make install

# Run in development mode
make dev

# Build for production
make build

# Run tests
make test
```

## Configuration

Copy `.env.example` to `.env` and configure:

| Variable | Description | Default |
|----------|-------------|---------|
| `SERVER_ADDR` | Listen address | `:8080` |
| `DB_DRIVER` | Database driver (postgres/sqlite) | `postgres` |
| `DB_URL` | PostgreSQL connection string | - |
| `JWT_SECRET` | JWT signing key | - |
| `ENCRYPTION_KEY` | AES-256 key for sensitive data | - |
| `STORAGE_BACKEND` | Storage backend (local/s3) | `local` |

See `.env.example` for full configuration options.

## API Endpoints

### Authentication

```
POST /api/v1/auth/register    # Register new user
POST /api/v1/auth/login       # Login and get JWT token
POST /api/v1/auth/refresh     # Refresh access token
POST /api/v1/auth/logout      # Logout
GET  /api/v1/auth/me          # Get current user info
```

### Captures

```
POST   /api/v1/captures           # Create capture task
GET    /api/v1/captures           # List user's captures
GET    /api/v1/captures/:id       # Get capture details
POST   /api/v1/captures/:id/retry # Retry failed capture
DELETE /api/v1/captures/:id       # Delete capture
GET    /api/v1/captures/:id/outputs           # List outputs
GET    /api/v1/captures/:id/outputs/:format   # Download output
GET    /api/v1/captures/:id/deliveries        # List deliveries
```

### Settings

```
# SMTP Profiles
GET    /api/v1/settings/smtp           # List SMTP profiles
POST   /api/v1/settings/smtp           # Create SMTP profile
PUT    /api/v1/settings/smtp/:id       # Update SMTP profile
DELETE /api/v1/settings/smtp/:id       # Delete SMTP profile
POST   /api/v1/settings/smtp/:id/test  # Test SMTP profile

# Webhooks
GET    /api/v1/settings/webhooks           # List webhooks
POST   /api/v1/settings/webhooks           # Create webhook
PUT    /api/v1/settings/webhooks/:id       # Update webhook
DELETE /api/v1/settings/webhooks/:id       # Delete webhook
POST   /api/v1/settings/webhooks/:id/test  # Test webhook

# User Settings
PUT  /api/v1/settings/password     # Change password
PUT  /api/v1/settings/profile      # Update profile
```

### Admin

```
GET    /api/v1/admin/users         # List all users
GET    /api/v1/admin/users/:id     # Get user details
PUT    /api/v1/admin/users/:id     # Update user (role, status)
DELETE /api/v1/admin/users/:id     # Delete user
GET    /api/v1/admin/audit-logs    # List audit logs
GET    /api/v1/admin/stats         # System statistics
GET    /api/v1/admin/smtp          # Get global SMTP config
GET    /api/v1/admin/storage       # Get storage config
```

## Usage Examples

### Create a Capture Task

```bash
curl -X POST http://localhost:8080/api/v1/captures \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "url": "https://example.com",
    "formats": ["pdf", "screenshot"],
    "cookies": "session=abc123"
  }'
```

### Response

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "url": "https://example.com",
  "formats": ["pdf", "screenshot"],
  "status": "pending",
  "created_at": "2025-01-01T00:00:00Z"
}
```

## Architecture

```
pagemail/
├── cmd/pagemail/          # Application entry point
├── internal/
│   ├── audit/             # Audit logging
│   ├── capture/           # Browser capture engine
│   ├── config/            # Configuration management
│   ├── db/                # Database migrations
│   ├── handlers/          # HTTP handlers
│   ├── middleware/        # Auth, rate limiting, etc.
│   ├── models/            # Data models
│   ├── notify/            # Email and webhook delivery
│   ├── pkg/               # Shared utilities
│   ├── queue/             # Job queue (goroutine-based)
│   ├── routes/            # Router setup
│   └── storage/           # Local/S3 storage backends
├── web/                   # Vue 3 frontend
└── deploy/                # Docker and deployment files
```

## Development

```bash
# Install tools
make pre-commit-install

# Run linters
make lint

# Run tests with coverage
make test

# Generate API docs
make docs-api
```

## Security

- JWT authentication with access/refresh tokens
- AES-256-GCM encryption for sensitive data (SMTP passwords, webhook secrets)
- SSRF protection (blocks private IPs, localhost, metadata endpoints)
- Rate limiting per user
- Audit logging for security events

## License

MIT
