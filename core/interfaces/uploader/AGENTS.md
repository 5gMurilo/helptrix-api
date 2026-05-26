<!-- Parent: ../AGENTS.md -->
<!-- Generated: 2026-05-25 | Updated: 2026-05-25 -->

# core/interfaces/uploader

## Purpose
Interface contracts for the image upload system: the upload strategy pattern, the uploader service, and the uploader controller.

## Key Files

| File | Description |
|------|-------------|
| `IImageUploadStrategy.go` | Strategy contract — one implementation per image type (profile, service) |
| `IUploaderController.go` | HTTP handler contract for the multipart file upload endpoint |
| `IUploaderService.go` | Orchestration contract — selects and delegates to the correct strategy |

## Dependencies

### Internal
- No domain types — uses `multipart.FileHeader` and `context.Context` primitives

<!-- MANUAL: -->
