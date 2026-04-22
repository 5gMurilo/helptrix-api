# AGENTS.md

This file provides guidance to Codex (Codex.ai/code) when working with code in this repository.

## Commands

```bash
# Start development with hot reload
air

# Run with Docker Compose
docker compose up

# Run tests
go test ./...

# Run a single test
go test ./modules/auth/... -run TestFunctionName

# Generate Swagger docs
swag init -g app/main.go

# Build
go build ./app/...
```

## Architecture

This is a Go REST API following Clean Architecture with Screaming Architecture organization. The codebase is structured around business domains, not technical layers.

### Directory Layout

```
adapter/
  db/
    repository/     # All GORM repositories
  auth/             # Paseto token config and creation/verification
  http/
    middleware/     # Auth middleware
    router.go       # All route definitions
core/
  domain/           # Entities and DTOs per entity
  interfaces/       # Contracts (interfaces) grouped by entity name
  utils/            # Utilities, constants, regex
modules/            # Feature modules — services and controllers per entity
app/
  main.go           # Entry point: dependency injection and app bootstrap
```

### Naming Conventions

- Domain files: `entity-name.entity.go`, `entity-name.request.dto.go`, `entity-name.response.dto.go`
- Interface files: `I<EntityName>.<repository|service|controller>.go` inside `core/interfaces/<entity-name>/`
- Test files mirror the file they test with `_test.go` suffix

### Layer Responsibilities

- **Controller**: Handles HTTP, maps errors to status codes, returns JSON. Never contains business logic. Must not know about DB internals.
- **Service**: Validates input, orchestrates business logic, calls repository. Returns errors as plain messages — no HTTP status codes.
- **Repository**: All DB operations wrapped in a single transaction. Performs rollback on any failure, commit only on full success.

### Key Patterns

**Dependency flow**: `adapter/http` → `modules` → `core/interfaces` ← `adapter/db/repository`

All dependencies are injected in `app/main.go`. Modules depend only on abstractions (`core/interfaces`), never on concrete implementations.

**Transactions**: Every multi-step DB operation uses a single GORM transaction with explicit rollback on error and commit on success.

**Error propagation**: Repositories return errors with context (e.g., `"error to assign categories for this user"`). Services return these messages. Controllers map them to appropriate HTTP status codes.

**Soft deletes**: All delete operations use GORM soft delete (`deleted_at` field).

### Authentication

Paseto symmetric tokens carry `id`, `name`, and `email`. Token expiration is 8 hours. Auth middleware validates tokens on protected routes.

### Swagger

All endpoints must be documented with swaggo annotations: known errors, request/response body examples, and appropriate status codes. Route naming convention: `<Method-Entity>` in camelCase (e.g., `PostAuthRegister`).

All DTOs and entities need to be exposed to swagger documentation

### Modules in Scope

- `auth` — register (`/auth/register`) and login (`/auth/login`)
- `user` — profile CRUD (`/user/profile/:id`)
- `service` — helper service listings
- `proposal` — business-to-helper service requests

### Database Schema

Schema definitions are in `docs/app-specs/db-modeling/`. Key tables: `users`, `categories`, `user_categories`, `proposals`, `services`, `reviews`, `payments`. All tables use `created_at`, `updated_at`, and soft-delete `deleted_at` where applicable.

### Testing

Minimum 60% coverage is required. Each module file must have a corresponding unit test file.
