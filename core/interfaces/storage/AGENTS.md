<!-- Parent: ../AGENTS.md -->
<!-- Generated: 2026-05-25 | Updated: 2026-05-25 -->

# core/interfaces/storage

## Purpose
Interface contract for cloud file storage. Decouples all upload strategies from the concrete Firebase Storage adapter, enabling alternative storage providers without changing business code.

## Key Files

| File | Description |
|------|-------------|
| `IStorageService.go` | File storage contract — upload file and return public URL, delete file |

## Dependencies

### Internal
- No domain types — uses `context.Context`, `[]byte`, and `string` primitives only

<!-- MANUAL: -->
