<!-- Parent: ../AGENTS.md -->
<!-- Generated: 2026-05-25 | Updated: 2026-05-25 -->

# core/interfaces/review

## Purpose
Interface contracts for the review domain: controller, repository, and service.

## Key Files

| File | Description |
|------|-------------|
| `IReviewController.go` | HTTP handler contract for review creation and listing endpoints |
| `IReviewRepository.go` | Data access contract — review persistence and retrieval by helper or proposal |
| `IReviewService.go` | Business logic contract — review creation validation and listing |

## Dependencies

### Internal
- `core/domain/` — review entity and request/response DTOs used in method signatures

<!-- MANUAL: -->
