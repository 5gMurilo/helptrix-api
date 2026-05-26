<!-- Parent: ../AGENTS.md -->
<!-- Generated: 2026-05-25 | Updated: 2026-05-25 -->

# adapter/storage

## Purpose
Firebase Cloud Storage adapter for file uploads. Implements the `IStorageService` interface from `core/interfaces/storage/`. Used by the uploader module to store profile and service images.

## Key Files

| File | Description |
|------|-------------|
| `firebase.go` | Firebase Storage SDK wrapper implementing IStorageService |
| `firebase_test.go` | Unit tests for the storage adapter |

## For AI Agents

### Working In This Directory
- Firebase credentials are loaded from the path specified in `FIREBASE_CREDENTIALS_PATH` environment variable
- Storage bucket name is read from `FIREBASE_STORAGE_BUCKET` environment variable
- Uploaded files return a public download URL; do not change the URL generation logic without testing the uploader module

### Common Patterns
- `NewFirebaseStorageService(credentialsPath, bucket string) (IStorageService, error)`
- `UploadFile(ctx, fileName, contentType string, data []byte) (string, error)` — returns public URL
- `DeleteFile(ctx, fileName string) error`

## Dependencies

### Internal
- `core/interfaces/storage/IStorageService.go` — interface this adapter implements

### External
- `firebase.google.com/go/v4`
- `cloud.google.com/go/storage`

<!-- MANUAL: -->
