<!-- Parent: ../AGENTS.md -->
<!-- Generated: 2026-05-25 | Updated: 2026-05-25 -->

# core/interfaces

## Purpose
All Go interface contracts grouped by entity. These are the ports in the Hexagonal Architecture — they decouple the application logic (modules) from the infrastructure (adapters). Every repository, service, controller, and external service has a corresponding interface here.

## Subdirectories

| Directory | Purpose |
|-----------|---------|
| `auth/` | IAuthController, IAuthRepository, IAuthService, ITokenMaker (see `auth/AGENTS.md`) |
| `category/` | ICategoryController, ICategoryRepository, ICategoryService |
| `email/` | IEmailSender |
| `helper/` | IHelperController, IHelperRepository, IHelperService |
| `otp/` | IOtpController, IOtpRepository, IOtpService |
| `proposal/` | IProposalController, IProposalRepository, IProposalService |
| `review/` | IReviewController, IReviewRepository, IReviewService |
| `service/` | IServiceController, IServiceRepository, IServiceService |
| `storage/` | IStorageService |
| `uploader/` | IImageUploadStrategy, IUploaderController, IUploaderService |
| `user/` | IUserController, IUserRepository, IUserService |

## For AI Agents

### Working In This Directory
- Each file defines exactly one interface
- Interface method signatures must use only types from `core/domain/` or Go standard library — no framework types
- When adding a method to an interface, update all implementations in `adapter/db/repository/` and `modules/`
- Naming: `I<EntityName>.<role>.go` (e.g., `IUser.repository.go`, `IUserController.go`)

### Common Patterns
```go
// Repository interface pattern
type IUserRepository interface {
    Create(user *domain.User) error
    FindByID(id uint) (*domain.User, error)
    Update(user *domain.User) error
    Delete(id uint) error
}

// Service interface pattern
type IUserService interface {
    GetProfile(id uint) (*domain.UserResponse, error)
    UpdateProfile(id uint, req *domain.UpdateUserRequest) error
}
```

<!-- MANUAL: -->
