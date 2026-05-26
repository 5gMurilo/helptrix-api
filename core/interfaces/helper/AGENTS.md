<!-- Parent: ../AGENTS.md -->
<!-- Generated: 2026-05-25 | Updated: 2026-05-25 -->

# core/interfaces/helper

## Purpose
Interface contracts for the helper discovery domain: controller, repository, and service.

## Key Files

| File | Description |
|------|-------------|
| `IHelperController.go` | HTTP handler contract for helper search and profile endpoints |
| `IHelperRepository.go` | Data access contract — helper search with filters and joins |
| `IHelperService.go` | Business logic contract — helper listing and profile retrieval |

## Dependencies

### Internal
- `core/domain/` — helper params and helper response DTO used in method signatures

<!-- MANUAL: -->
