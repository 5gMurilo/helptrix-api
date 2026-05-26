<!-- Parent: ../AGENTS.md -->
<!-- Generated: 2026-05-25 | Updated: 2026-05-25 -->

# adapter

## Purpose
Infrastructure layer implementing all external integrations. Contains concrete implementations of the interfaces defined in `core/interfaces/`. This layer is the only place that knows about GORM, Firebase, Resend, PASETO, and Gin middleware — all other layers remain framework-agnostic.

## Subdirectories

| Directory | Purpose |
|-----------|---------|
| `auth/` | PASETO token creation and verification (see `auth/AGENTS.md`) |
| `db/` | Database connection, migrations, repositories, and seeders (see `db/AGENTS.md`) |
| `email/` | Resend transactional email adapter (see `email/AGENTS.md`) |
| `http/` | Gin router definitions and authentication middleware (see `http/AGENTS.md`) |
| `storage/` | Firebase Storage adapter for file uploads (see `storage/AGENTS.md`) |

## For AI Agents

### Working In This Directory
- Adapters implement interfaces from `core/interfaces/` — always check the contract before adding methods
- Never import from `modules/` — the dependency flow is one-directional
- Adapters are instantiated and wired in `app/main.go` only

### Common Patterns
- Each adapter file implements exactly one interface from `core/interfaces/`
- Error messages include context (e.g., `"error creating user: " + err.Error()`)

## Dependencies

### Internal
- `core/interfaces/` — contracts this layer implements
- `core/domain/` — domain entities used in method signatures

### External
- `gorm.io/gorm` — ORM (db adapters)
- `github.com/gin-gonic/gin` — HTTP framework (http adapter)
- `firebase.google.com/go/v4` — Firebase SDK (storage adapter)
- `github.com/resend/resend-go/v2` — Email SDK (email adapter)

<!-- MANUAL: -->
