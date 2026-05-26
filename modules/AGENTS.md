<!-- Parent: ../AGENTS.md -->
<!-- Generated: 2026-05-25 | Updated: 2026-05-25 -->

# modules

## Purpose
Application layer containing all business feature modules. Each subdirectory is a self-contained feature with a controller (HTTP layer) and a service (business logic layer). Modules only depend on `core/interfaces/` abstractions — never on concrete adapters or other modules.

## Subdirectories

| Directory | Purpose |
|-----------|---------|
| `auth/` | Registration and login flows (see `auth/AGENTS.md`) |
| `category/` | Category listing and management (see `category/AGENTS.md`) |
| `helper/` | Helper profile search and listing (see `helper/AGENTS.md`) |
| `otp/` | One-time password generation and verification (see `otp/AGENTS.md`) |
| `proposal/` | Business-to-helper service request flows (see `proposal/AGENTS.md`) |
| `review/` | Rating and review system (see `review/AGENTS.md`) |
| `service/` | Helper service listings CRUD (see `service/AGENTS.md`) |
| `uploader/` | Image upload with strategy pattern per image type (see `uploader/AGENTS.md`) |
| `user/` | User profile CRUD operations (see `user/AGENTS.md`) |

## For AI Agents

### Working In This Directory
- Controllers handle HTTP only: parse request, call service, map error to status code, return JSON
- Services handle business logic only: validate input, orchestrate calls, return plain error strings (no HTTP codes)
- No module imports another module — cross-cutting concerns go through core interfaces
- Each new module requires: controller, service, interfaces in `core/interfaces/<name>/`, and domain types in `core/domain/`

### Testing Requirements
- Every controller and service file must have a `_test.go` counterpart
- Use interface mocks for dependencies in unit tests
- Minimum 60% coverage across all module files

### Common Patterns
```
modules/<name>/
  <name>.controller.go       — HTTP handler, Gin context, status codes
  <name>.controller_test.go  — Controller unit tests with mocked service
  <name>.service.go          — Business logic, calls repository interface
  <name>.service_test.go     — Service unit tests with mocked repository
```

## Dependencies

### Internal
- `core/interfaces/<name>/` — repository and service contracts
- `core/domain/` — request/response DTOs and entities

<!-- MANUAL: -->
