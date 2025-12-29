# Pagemail Project Instructions

## Pre-commit Checklist (MANDATORY)

Before committing any code changes, **MUST** run local checks:

```bash
# Backend lint (required)
make lint-backend

# Frontend lint (if frontend changes)
make lint-frontend

# Full lint check
make lint

# Run tests
make test
```

**DO NOT** commit or push until all lint errors are resolved. This prevents CI failures on GitHub.

## Quick Commands

| Command | Description |
|---------|-------------|
| `make lint` | Run all linters (backend + frontend) |
| `make lint-backend` | Run golangci-lint only |
| `make lint-fix` | Auto-fix linting issues |
| `make test` | Run all tests |
| `make dev` | Start development servers |

## Code Style

- Go: Follow golangci-lint rules defined in `.golangci.yml`
- Vue/TS: Follow ESLint rules in `web/eslint.config.ts`
- Use `make fmt` to format code before committing

## Common Lint Issues

- `gocritic/rangeValCopy`: Use pointer or index in range loops for large structs
- `gocritic/hugeParam`: Pass large structs by pointer
- `errcheck`: Always handle error returns
- `goconst`: Extract repeated string literals to constants
- `gofmt`: Run `go fmt ./...` or `make fmt`
