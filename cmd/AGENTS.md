<!-- Parent: ../AGENTS.md -->
<!-- Generated: 2026-05-25 | Updated: 2026-05-25 -->

# cmd

## Purpose
CLI utility commands for the project. Currently contains the database seeder command used to populate the database with initial/reference data during development and staging setup.

## Subdirectories

| Directory | Purpose |
|-----------|---------|
| `seed/` | Database seed command entry point (see `seed/AGENTS.md`) |

## For AI Agents

### Working In This Directory
- Each subdirectory is an independent Go binary with its own `main.go`
- CLI commands are not part of the main HTTP server — they run independently
- New CLI commands follow the same pattern: separate `main.go` in a named subdirectory

<!-- MANUAL: -->
