<!-- Parent: ../AGENTS.md -->
<!-- Generated: 2026-05-25 | Updated: 2026-05-25 -->

# modules/auth

## Purpose
Handles user registration and login. The service validates credentials, hashes passwords, and delegates token creation to the `ITokenMaker` adapter. Public routes — no authentication middleware applied.

## Key Files

| File | Description |
|------|-------------|
| `auth.controller.go` | Gin handlers for POST /auth/register and POST /auth/login |
| `auth.controller_test.go` | Controller unit tests with mocked auth service |
| `auth.service.go` | Registration and login business logic |
| `auth.service_test.go` | Service unit tests with mocked repository and token maker |

## For AI Agents

### Working In This Directory
- Registration creates a user and sends an OTP verification email
- Login verifies credentials and returns a PASETO token
- Passwords must be hashed using bcrypt — never store plaintext
- Service returns plain error strings; controller maps them to 400/401/409/500

### Testing Requirements
- Mock `IAuthRepository` and `ITokenMaker` in service tests
- Mock `IAuthService` in controller tests
- Test happy path, duplicate email, invalid credentials, and server error cases

## Dependencies

### Internal
- `core/interfaces/auth/` — IAuthController, IAuthRepository, IAuthService, ITokenMaker
- `core/domain/` — auth request/response DTOs, user entity

<!-- MANUAL: -->
