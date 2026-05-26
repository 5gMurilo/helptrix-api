<!-- Parent: ../AGENTS.md -->
<!-- Generated: 2026-05-25 | Updated: 2026-05-25 -->

# modules/category

## Purpose
Category listing and management. Categories classify helper skills (e.g., cleaning, plumbing, tutoring). Populated via the seeder; typically read-only from the API perspective.

## Key Files

| File | Description |
|------|-------------|
| `category.controller.go` | Gin handlers for category endpoints |
| `category.controller_test.go` | Controller unit tests with mocked category service |
| `category.service.go` | Category retrieval business logic |
| `category.service_test.go` | Service unit tests with mocked repository |

## For AI Agents

### Working In This Directory
- Categories are seeded via `adapter/db/seeder/category_seeder.go` — not created through the API
- List endpoint is public (no auth required)

## Dependencies

### Internal
- `core/interfaces/category/` — ICategoryController, ICategoryRepository, ICategoryService
- `core/domain/` — category entity and response DTO

<!-- MANUAL: -->
