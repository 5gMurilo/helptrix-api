<!-- Parent: ../AGENTS.md -->
<!-- Generated: 2026-05-25 | Updated: 2026-05-25 -->

# core/domain

## Purpose
All domain entities and data transfer objects. Entities map directly to database tables (embedding `gorm.Model`). Request DTOs carry input validation tags. Response DTOs shape the JSON output returned to API consumers. No business logic lives here.

## Key Files

| File | Description |
|------|-------------|
| `address.entity.go` | Address entity — embedded or associated with user profiles |
| `auth.request.dto.go` | Login and registration request payloads |
| `auth.response.dto.go` | Authentication response with token and user info |
| `category.entity.go` | Category entity for helper skill classification |
| `category.response.dto.go` | Category response shape for API output |
| `helper.params.go` | Query parameters for helper search and filtering |
| `helper.response.dto.go` | Helper profile response with aggregated data |
| `otp.entity.go` | OTP entity — stores generated codes with expiry |
| `otp.request.dto.go` | OTP generation and verification request payloads |
| `otp.response.dto.go` | OTP operation response shapes |
| `proposal.entity.go` | Proposal entity — service request from client to helper |
| `proposal.request.dto.go` | Proposal creation request payload |
| `proposal.response.dto.go` | Proposal response with status and associated data |
| `review.entity.go` | Review entity — rating and comment for a proposal |
| `review.request.dto.go` | Review creation request payload |
| `review.response.dto.go` | Review response shape |
| `service.entity.go` | Service entity — helper's offered service |
| `service.request.dto.go` | Service creation/update request payload |
| `service.response.dto.go` | Service response shape |
| `user.entity.go` | User entity — core identity and profile data |
| `user.request.dto.go` | User profile update request payload |
| `user.response.dto.go` | User profile response shape |
| `user_category.entity.go` | Join entity linking users to their skill categories |

## For AI Agents

### Working In This Directory
- Entities use `gorm.Model` embedding (provides `ID`, `CreatedAt`, `UpdatedAt`, `DeletedAt`)
- Request DTOs use `binding:"required"` and other Gin validator tags
- Response DTOs must never expose sensitive fields (passwords, internal IDs not needed by clients)
- All new entities/DTOs must be registered with swaggo via `// @name` comments or explicit struct exposure

### Common Patterns
```go
// Entity pattern
type User struct {
    gorm.Model
    Name  string `gorm:"not null"`
    Email string `gorm:"unique;not null"`
}

// Request DTO pattern
type CreateUserRequest struct {
    Name  string `json:"name"  binding:"required"`
    Email string `json:"email" binding:"required,email"`
}

// Response DTO pattern
type UserResponse struct {
    ID    uint   `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}
```

## Dependencies

### External
- `gorm.io/gorm` — for `gorm.Model` embedding in entities only

<!-- MANUAL: -->
