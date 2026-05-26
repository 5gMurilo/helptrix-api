<!-- Parent: ../AGENTS.md -->
<!-- Generated: 2026-05-25 | Updated: 2026-05-25 -->

# core/interfaces/proposal

## Purpose
Interface contracts for the proposal domain: controller, repository, and service.

## Key Files

| File | Description |
|------|-------------|
| `IProposal.controller.go` | HTTP handler contract for proposal CRUD and status transition endpoints |
| `IProposal.repository.go` | Data access contract — proposal creation, retrieval, and status updates |
| `IProposal.service.go` | Business logic contract — proposal lifecycle management |

## Dependencies

### Internal
- `core/domain/` — proposal entity and request/response DTOs used in method signatures

<!-- MANUAL: -->
