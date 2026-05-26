<!-- Parent: ../AGENTS.md -->
<!-- Generated: 2026-05-25 | Updated: 2026-05-25 -->

# adapter/db/migration

## Purpose
SQL migration files managed by `golang-migrate`. Each migration has an `up` file (apply change) and a `down` file (revert change). Migrations run automatically on startup via `adapter/db/db.go`.

## Key Files

| File | Description |
|------|-------------|
| `0042_create_reviews_table.up.sql` | Creates the `reviews` table |
| `0042_create_reviews_table.down.sql` | Drops the `reviews` table |
| `0043_add_proposal_id_to_reviews.up.sql` | Adds `proposal_id` foreign key to `reviews` |
| `0043_add_proposal_id_to_reviews.down.sql` | Removes `proposal_id` from `reviews` |

## For AI Agents

### Working In This Directory
- Migration files are numbered sequentially — never reuse or skip numbers
- Always create both `.up.sql` and `.down.sql` for every migration
- Down migrations must exactly reverse the up migration
- Never modify existing migration files — create a new migration to correct mistakes
- Schema definitions are documented in `docs/app-specs/db-modeling/`

### Common Patterns
- Naming: `NNNN_description_of_change.up.sql` / `NNNN_description_of_change.down.sql`
- New tables must include `created_at`, `updated_at`, and `deleted_at TIMESTAMPTZ` columns for GORM compatibility

<!-- MANUAL: -->
