<!-- Parent: ../AGENTS.md -->
<!-- Generated: 2026-05-25 | Updated: 2026-05-25 -->

# adapter/db/seeder

## Purpose
Database seeders that populate reference/initial data. Currently seeds the `categories` table with predefined skill categories. Called by `cmd/seed/main.go`.

## Key Files

| File | Description |
|------|-------------|
| `category_seeder.go` | Inserts predefined categories using upsert to stay idempotent |

## For AI Agents

### Working In This Directory
- Seeders must be idempotent — safe to run multiple times without creating duplicates
- Use GORM's `FirstOrCreate` or `OnConflict` clause for upsert behavior
- New seeders should follow the same pattern: accept `*gorm.DB`, return `error`

<!-- MANUAL: -->
