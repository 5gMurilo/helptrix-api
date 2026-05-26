<!-- Parent: ../AGENTS.md -->
<!-- Generated: 2026-05-25 | Updated: 2026-05-25 -->

# cmd/seed

## Purpose
Standalone CLI binary that connects to the database and runs the category seeder to populate reference data. Run once during initial setup or when resetting a development/staging environment.

## Key Files

| File | Description |
|------|-------------|
| `main.go` | Entry point: opens DB connection and calls the category seeder |

## For AI Agents

### Working In This Directory
- Build and run with: `go run ./cmd/seed/`
- Requires `DATABASE_URL` environment variable to be set
- Seeder is idempotent — safe to run multiple times (uses upsert or existence checks)

## Dependencies

### Internal
- `adapter/db/db.go` — database connection
- `adapter/db/seeder/category_seeder.go` — seeder implementation

<!-- MANUAL: -->
