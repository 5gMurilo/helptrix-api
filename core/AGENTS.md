<!-- Parent: ../AGENTS.md -->
<!-- Generated: 2026-05-25 | Updated: 2026-05-25 -->

# core

## Purpose
Domain layer containing the business model and all interface contracts. This is the innermost layer of the Clean Architecture — it has zero dependencies on frameworks, databases, or external services. Everything else depends on core; core depends on nothing.

## Subdirectories

| Directory | Purpose |
|-----------|---------|
| `domain/` | Entities, request DTOs, and response DTOs for every feature (see `domain/AGENTS.md`) |
| `interfaces/` | Go interface contracts grouped by entity — repositories, services, controllers (see `interfaces/AGENTS.md`) |
| `utils/` | Shared utilities and application-wide constants (see `utils/AGENTS.md`) |

## For AI Agents

### Working In This Directory
- Never import from `adapter/` or `modules/` — core is dependency-free
- Adding a new entity requires files in both `domain/` and `interfaces/`
- DTOs are separate from entities: entities map to DB tables, request DTOs carry input validation, response DTOs shape API output

### Common Patterns
- Entity files carry GORM model tags and `gorm.Model` embedding
- Request DTOs carry `binding:"required"` and other Gin validation tags
- Response DTOs contain only exported fields safe to expose to clients

## Dependencies

### External
- `gorm.io/gorm` — only for `gorm.Model` embedding in entities (no DB calls here)

<!-- MANUAL: -->
