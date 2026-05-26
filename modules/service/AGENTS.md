<!-- Parent: ../AGENTS.md -->
<!-- Generated: 2026-05-25 | Updated: 2026-05-25 -->

# modules/service

## Purpose
Helper service listings CRUD. A helper can create, update, and delete the services they offer (e.g., "House Cleaning - 2h - R$150"). Services are returned as part of the helper profile and are the target of proposals.

## Key Files

| File | Description |
|------|-------------|
| `service.controller.go` | Gin handlers for service CRUD endpoints |
| `service.controller_test.go` | Controller unit tests with mocked service |
| `service.service.go` | Service creation, update, retrieval, and deletion business logic |
| `service.service_test.go` | Service unit tests with mocked repository |

## For AI Agents

### Working In This Directory
- Only authenticated helpers can create/update/delete their own services
- Service images are handled by the uploader module — `service.controller.go` delegates to the uploader
- Soft delete must be used — never hard delete service records
- See `docs/app-specs/service.md` for field requirements

## Dependencies

### Internal
- `core/interfaces/service/` — IServiceController, IServiceRepository, IServiceService
- `core/domain/` — service entity and request/response DTOs

<!-- MANUAL: -->
