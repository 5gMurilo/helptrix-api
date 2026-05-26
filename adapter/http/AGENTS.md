<!-- Parent: ../AGENTS.md -->
<!-- Generated: 2026-05-25 | Updated: 2026-05-25 -->

# adapter/http

## Purpose
HTTP layer configuration. `router.go` defines all API routes and groups, registers controller handlers, and applies middleware. The `middleware/` subdirectory contains the PASETO authentication middleware that protects private routes.

## Key Files

| File | Description |
|------|-------------|
| `router.go` | All route definitions, Gin engine setup, middleware application, and Swagger UI registration |

## Subdirectories

| Directory | Purpose |
|-----------|---------|
| `middleware/` | Auth middleware that validates PASETO tokens on protected routes (see `middleware/AGENTS.md`) |

## For AI Agents

### Working In This Directory
- Register new routes in `router.go` following the existing group pattern
- Public routes go under the public group; authenticated routes go under the group that applies `AuthMiddleware`
- Route naming in Swagger annotations: `<Method><Entity>` in camelCase (e.g., `PostAuthRegister`, `GetUserProfile`)
- The Swagger UI route `/swagger/*any` must remain public

### Common Patterns
```go
// Authenticated group pattern
authorized := r.Group("/")
authorized.Use(middleware.AuthMiddleware(tokenMaker))
{
    authorized.GET("/user/profile/:id", userCtrl.GetProfile)
}
```

## Dependencies

### Internal
- `core/interfaces/` — controller interfaces used as handler parameter types
- `adapter/auth/` — token maker passed to auth middleware

### External
- `github.com/gin-gonic/gin`
- `github.com/swaggo/gin-swagger`

<!-- MANUAL: -->
