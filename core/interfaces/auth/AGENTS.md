<!-- Parent: ../AGENTS.md -->
<!-- Generated: 2026-05-25 | Updated: 2026-05-25 -->

# core/interfaces/auth

## Purpose
Interface contracts for the authentication domain: controller, repository, service, and token maker.

## Key Files

| File | Description |
|------|-------------|
| `IAuthController.go` | HTTP handler contract for register and login endpoints |
| `IAuthRepository.go` | Data access contract — user lookup by email, user creation |
| `IAuthService.go` | Business logic contract — register and login operations |
| `ITokenMaker.go` | Token creation and verification contract — implemented by PASETO adapter |

## Dependencies

### Internal
- `core/domain/` — auth request/response DTOs used in method signatures

<!-- MANUAL: -->
