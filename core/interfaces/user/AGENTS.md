<!-- Parent: ../AGENTS.md -->
<!-- Generated: 2026-05-25 | Updated: 2026-05-25 -->

# core/interfaces/user

## Purpose
Interface contracts for the user profile domain: controller, repository, and service.

## Key Files

| File | Description |
|------|-------------|
| `IUserController.go` | HTTP handler contract for user profile GET and PUT endpoints |
| `IUserRepository.go` | Data access contract — user profile reads, updates, and category assignment |
| `IUserService.go` | Business logic contract — profile retrieval and update operations |

## Dependencies

### Internal
- `core/domain/` — user entity and request/response DTOs used in method signatures

<!-- MANUAL: -->
