<!-- Parent: ../AGENTS.md -->
<!-- Generated: 2026-05-25 | Updated: 2026-05-25 -->

# core/interfaces/service

## Purpose
Interface contracts for the helper service domain: controller, repository, and service. Note: the Go package name avoids collision with the standard library by using a distinct package identifier.

## Key Files

| File | Description |
|------|-------------|
| `IServiceController.go` | HTTP handler contract for service CRUD endpoints |
| `IServiceRepository.go` | Data access contract — service creation, retrieval, update, and soft delete |
| `IServiceService.go` | Business logic contract — service management operations |

## Dependencies

### Internal
- `core/domain/` — service entity and request/response DTOs used in method signatures

<!-- MANUAL: -->
