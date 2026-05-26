<!-- Parent: ../AGENTS.md -->
<!-- Generated: 2026-05-25 | Updated: 2026-05-25 -->

# modules/helper

## Purpose
Helper profile search and listing. Aggregates helper data (user info, services, categories, reviews) for the client-facing discovery flow. Supports filtering by category and other query parameters defined in `core/domain/helper.params.go`.

## Key Files

| File | Description |
|------|-------------|
| `helper.controller.go` | Gin handlers for helper search and profile endpoints |
| `helper.controller_test.go` | Controller unit tests with mocked helper service |
| `helper.service.go` | Helper search and aggregation business logic |
| `helper.service_test.go` | Service unit tests with mocked repository |

## For AI Agents

### Working In This Directory
- `helper.response.dto.go` in `core/domain/` includes nested services and reviews — keep response shapes consistent with the spec in `docs/app-specs/profile.md`
- Query parameters for filtering come from `helper.params.go` bound via `c.ShouldBindQuery`
- Search results should include average rating computed from reviews

## Dependencies

### Internal
- `core/interfaces/helper/` — IHelperController, IHelperRepository, IHelperService
- `core/domain/` — helper params, helper response DTO

<!-- MANUAL: -->
