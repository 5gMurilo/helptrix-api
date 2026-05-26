<!-- Parent: ../AGENTS.md -->
<!-- Generated: 2026-05-25 | Updated: 2026-05-25 -->

# modules/uploader/strategies

## Purpose
Concrete implementations of the `IImageUploadStrategy` interface, one per image type. Each strategy encapsulates the validation rules, Firebase Storage path, and post-upload side effect (e.g., writing the URL back to the user or service record) specific to that image type.

## Key Files

| File | Description |
|------|-------------|
| `profile_image.strategy.go` | Uploads user profile images; updates `users.avatar_url` after upload |
| `profile_image.strategy_test.go` | Unit tests for profile image strategy |
| `service_image.strategy.go` | Uploads service cover images; updates `services.image_url` after upload |
| `service_image.strategy_test.go` | Unit tests for service image strategy |

## For AI Agents

### Working In This Directory
- Each strategy implements `IImageUploadStrategy.Upload(ctx, file, ownerID) (string, error)`
- File size and MIME type validation must happen inside the strategy before calling IStorageService
- Storage paths should be namespaced by type: `profile-images/<userID>/<filename>`, `service-images/<serviceID>/<filename>`
- Adding a new image type: create a new `<type>.strategy.go` + `<type>.strategy_test.go`, then register in `uploader.service.go`

## Dependencies

### Internal
- `core/interfaces/uploader/IImageUploadStrategy.go` — interface this implements
- `core/interfaces/storage/IStorageService.go` — injected for actual file upload
- `core/interfaces/user/IUserRepository.go` or `core/interfaces/service/IServiceRepository.go` — for post-upload record updates

<!-- MANUAL: -->
