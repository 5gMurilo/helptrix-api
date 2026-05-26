<!-- Parent: ../AGENTS.md -->
<!-- Generated: 2026-05-25 | Updated: 2026-05-25 -->

# adapter/auth

## Purpose
PASETO (Platform-Agnostic Security TOkens) implementation. Provides token creation and verification for the authentication system. Implements the `ITokenMaker` interface from `core/interfaces/auth/`.

## Key Files

| File | Description |
|------|-------------|
| `paseto.go` | PASETO symmetric token maker — creates and verifies tokens carrying user payload |

## For AI Agents

### Working In This Directory
- Tokens carry `id`, `name`, and `email` claims; expiration is 8 hours
- The symmetric key must be exactly 32 bytes — validated at construction time
- Do not change the payload structure without updating all consumers that extract claims from the token

### Common Patterns
- `NewPasetoMaker(symmetricKey string)` — constructor, returns `ITokenMaker`
- `CreateToken(id, name, email string, duration time.Duration) (string, error)`
- `VerifyToken(token string) (*Payload, error)`

## Dependencies

### Internal
- `core/interfaces/auth/ITokenMaker.go` — interface this adapter implements

### External
- PASETO/chacha20poly1305 signing library

<!-- MANUAL: -->
