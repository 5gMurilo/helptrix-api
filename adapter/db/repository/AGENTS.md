<!-- Parent: ../AGENTS.md -->
<!-- Generated: 2026-05-25 | Updated: 2026-05-25 -->

# adapter/db/repository

## Purpose
GORM implementations of all repository interfaces defined in `core/interfaces/`. Each file implements one repository interface for one entity. All mutations are wrapped in explicit GORM transactions.

## Key Files

| File | Description |
|------|-------------|
| `auth.repository.go` | User lookup and creation for auth flows |
| `category.repository.go` | Category listing and management |
| `helper.repository.go` | Helper search queries with joins and filters |
| `otp.repository.go` | OTP creation, retrieval, and expiry management |
| `proposal.repository.go` | Proposal CRUD and status transitions |
| `review.repository.go` | Review creation and listing |
| `review.repository_test.go` | Integration-style tests for review repository |
| `service.repository.go` | Helper service CRUD |
| `user.repository.go` | User profile reads and updates |

## For AI Agents

### Working In This Directory
- Every write operation uses a transaction: `tx := db.Begin()` → operations → `tx.Rollback()` on error → `tx.Commit()` on success
- Query column and table names must match the actual DB schema in `docs/app-specs/db-modeling/` — verify before writing queries
- Use GORM soft delete for all `Delete` operations (GORM handles this automatically with `deleted_at`)
- Never write raw SQL unless GORM cannot express the query — always prefer GORM chainable methods first

### Common Patterns
```go
func (r *userRepository) Update(user *domain.User) error {
    tx := r.db.Begin()
    if err := tx.Save(user).Error; err != nil {
        tx.Rollback()
        return fmt.Errorf("error updating user: %w", err)
    }
    return tx.Commit().Error
}
```

## Dependencies

### Internal
- `core/interfaces/` — repository interfaces implemented here
- `core/domain/` — entity types used in method signatures

### External
- `gorm.io/gorm`

<!-- MANUAL: -->
