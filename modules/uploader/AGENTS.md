<!-- Parent: ../AGENTS.md -->
<!-- Generated: 2026-05-25 | Updated: 2026-05-25 -->

# modules/uploader

## Purpose
Image upload orchestration using the Strategy pattern. A single uploader service delegates to the correct upload strategy based on image type (profile image vs. service image). Each strategy handles its own validation rules, storage path, and post-upload side effects (e.g., updating the user or service record).

## Key Files

| File | Description |
|------|-------------|
| `uploader.controller.go` | Gin handler for multipart file upload endpoint |
| `uploader.service.go` | Orchestrator — selects and executes the correct strategy |
| `uploader.service_test.go` | Service unit tests with mocked strategies and storage |

## Subdirectories

| Directory | Purpose |
|-----------|---------|
| `strategies/` | Concrete upload strategies per image type (see `strategies/AGENTS.md`) |

## For AI Agents

### Working In This Directory
- Adding a new image type requires: a new strategy in `strategies/`, implementing `IImageUploadStrategy`, and registering it in the service
- File size and MIME type validation happens inside each strategy
- The controller receives the `imageType` query parameter to determine which strategy to invoke
- Storage is delegated to `IStorageService` — the uploader never calls Firebase directly

### Common Patterns
```go
// Strategy selection pattern
strategy := service.getStrategy(imageType)
result, err := strategy.Upload(ctx, file, userID)
```

## Dependencies

### Internal
- `core/interfaces/uploader/` — IUploaderController, IUploaderService, IImageUploadStrategy
- `core/interfaces/storage/` — IStorageService injected into strategies

<!-- MANUAL: -->
