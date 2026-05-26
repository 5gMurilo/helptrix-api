<!-- Parent: ../AGENTS.md -->
<!-- Generated: 2026-05-25 | Updated: 2026-05-25 -->

# core/utils

## Purpose
Shared utilities and application-wide constants. Provides a single source of truth for magic values, regex patterns, and helper functions used across all layers.

## Key Files

| File | Description |
|------|-------------|
| `constants.go` | Application-wide constants — no magic numbers or strings elsewhere |

## For AI Agents

### Working In This Directory
- All magic numbers and string literals used in more than one place must be declared here
- Never add business logic to this package — only pure utilities and constants
- Constants use `SCREAMING_SNAKE_CASE` for exported values

### Common Patterns
```go
const (
    OTPExpirationMinutes = 5
    TokenDuration        = 8 * time.Hour
    MaxImageSizeMB       = 5
)
```

<!-- MANUAL: -->
