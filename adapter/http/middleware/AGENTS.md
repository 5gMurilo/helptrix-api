<!-- Parent: ../AGENTS.md -->
<!-- Generated: 2026-05-25 | Updated: 2026-05-25 -->

# adapter/http/middleware

## Purpose
Gin middleware functions applied to protected route groups. Currently contains the authentication middleware that validates PASETO tokens and injects the user payload into the Gin context.

## Key Files

| File | Description |
|------|-------------|
| `auth.go` | PASETO token validation middleware — extracts Bearer token, verifies, sets payload in context |

## For AI Agents

### Working In This Directory
- Token is expected in the `Authorization: Bearer <token>` header
- On valid token, the user payload is stored in the Gin context under a well-known key (e.g., `"authPayload"`) for downstream handlers to read
- Returns 401 on missing, expired, or invalid tokens — never 403
- New middleware functions follow the same signature: `func MiddlewareName(deps...) gin.HandlerFunc`

### Common Patterns
```go
// Reading the payload in a controller
payload := c.MustGet("authPayload").(*auth.Payload)
userID := payload.ID
```

## Dependencies

### Internal
- `adapter/auth/` — ITokenMaker.VerifyToken

<!-- MANUAL: -->
