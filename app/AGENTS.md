<!-- Parent: ../AGENTS.md -->
<!-- Generated: 2026-05-25 | Updated: 2026-05-25 -->

# app

## Purpose
Application entry point. `main.go` is the single composition root: it loads environment variables, connects to the database, runs migrations, instantiates all adapters and modules, wires dependencies via constructor injection, and starts the Gin HTTP server.

## Key Files

| File | Description |
|------|-------------|
| `main.go` | Dependency injection container, server bootstrap, and startup sequence |

## For AI Agents

### Working In This Directory
- This is the only place where concrete adapter implementations are instantiated
- When adding a new module, follow the existing wiring pattern: repository → service → controller → register routes
- All environment variables must be read here and passed down — no `os.Getenv` calls in other layers
- Database migrations run automatically on startup via `golang-migrate`

### Common Patterns
```go
// Typical wiring pattern
repo := repository.NewXxxRepository(db)
svc  := modules_xxx.NewXxxService(repo, ...)
ctrl := modules_xxx.NewXxxController(svc)
router.RegisterXxxRoutes(ctrl)
```

## Dependencies

### Internal
- `adapter/` — all concrete infrastructure implementations
- `modules/` — all feature controllers and services
- `core/interfaces/` — types used in wiring signatures

<!-- MANUAL: -->
