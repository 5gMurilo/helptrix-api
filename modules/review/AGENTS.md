<!-- Parent: ../AGENTS.md -->
<!-- Generated: 2026-05-25 | Updated: 2026-05-25 -->

# modules/review

## Purpose
Rating and review system. Clients submit a rating and optional comment after a completed proposal. Reviews are associated with both the proposal and the helper, enabling aggregate rating calculations on helper profiles.

## Key Files

| File | Description |
|------|-------------|
| `review.controller.go` | Gin handlers for review creation and listing |
| `review.controller_test.go` | Controller unit tests with mocked review service |
| `review.service.go` | Review creation validation and retrieval logic |
| `review.service_test.go` | Service unit tests with mocked repository |

## For AI Agents

### Working In This Directory
- A review requires an existing `proposalId` — validate the proposal exists and belongs to the authenticated user
- One review per proposal — reject duplicates
- `proposal_id` is nullable in early schema (migration 0042); migration 0043 adds the column if needed
- See `docs/app-specs/review.md` for business rules

## Dependencies

### Internal
- `core/interfaces/review/` — IReviewController, IReviewRepository, IReviewService
- `core/domain/` — review entity and request/response DTOs

<!-- MANUAL: -->
