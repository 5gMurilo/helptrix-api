<!-- Parent: ../AGENTS.md -->
<!-- Generated: 2026-05-25 | Updated: 2026-05-25 -->

# adapter/db

## Purpose
Database infrastructure layer. Manages the PostgreSQL connection, runs schema migrations on startup, houses all GORM repository implementations, and provides a category seeder for initial data.

## Key Files

| File | Description |
|------|-------------|
| `db.go` | Opens the GORM PostgreSQL connection and runs migrations |

## Subdirectories

| Directory | Purpose |
|-----------|---------|
| `migration/` | SQL migration files (up/down) managed by golang-migrate (see `migration/AGENTS.md`) |
| `repository/` | GORM implementations of all repository interfaces (see `repository/AGENTS.md`) |
| `seeder/` | Category seed data loader (see `seeder/AGENTS.md`) |

## For AI Agents

### Working In This Directory
- `db.go` connects using `DATABASE_URL` from environment — never hardcode connection strings
- Migrations run automatically when `db.go` is called; add new migration files before deploying schema changes
- All repository constructors receive `*gorm.DB` injected from `app/main.go`

### Common Patterns
- Migration files follow `NNNN_description.up.sql` / `NNNN_description.down.sql` naming
- Every repository wraps mutations in a GORM transaction: `db.Begin()` → operations → `tx.Rollback()` on error → `tx.Commit()` on success

## Dependencies

### Internal
- `core/interfaces/` — repository contracts implemented here

### External
- `gorm.io/gorm` + `gorm.io/driver/postgres`
- `github.com/golang-migrate/migrate/v4`

<!-- MANUAL: -->
