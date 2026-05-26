<!-- Parent: ../AGENTS.md -->
<!-- Generated: 2026-05-25 | Updated: 2026-05-25 -->

# core/interfaces/otp

## Purpose
Interface contracts for the OTP verification domain: controller, repository, and service.

## Key Files

| File | Description |
|------|-------------|
| `IOtp.controller.go` | HTTP handler contract for OTP generation and verification endpoints |
| `IOtp.repository.go` | Data access contract — OTP persistence, lookup, and expiry check |
| `IOtp.service.go` | Business logic contract — generate, send, and verify OTP operations |

## Dependencies

### Internal
- `core/domain/` — OTP entity and request/response DTOs used in method signatures

<!-- MANUAL: -->
