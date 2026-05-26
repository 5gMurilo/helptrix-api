<!-- Parent: ../AGENTS.md -->
<!-- Generated: 2026-05-25 | Updated: 2026-05-25 -->

# modules/otp

## Purpose
4-digit one-time password flow for email verification. Generates a short-lived OTP, persists it to the database, sends it via email, and verifies it on submission. Used after registration to confirm user email addresses.

## Key Files

| File | Description |
|------|-------------|
| `otp.controller.go` | Gin handlers for OTP generation and verification endpoints |
| `otp.controller_test.go` | Controller unit tests with mocked OTP service |
| `otp.service.go` | OTP generation, email dispatch, and verification business logic |
| `otp.service_test.go` | Service unit tests with mocked repository and email sender |

## For AI Agents

### Working In This Directory
- OTPs are 4 digits, expire after a configurable duration (see `core/utils/constants.go`)
- Expiry is stored in the database and checked at verification time
- The service calls `IEmailSender` to dispatch the OTP — never call the email adapter directly
- Verification marks the OTP as used to prevent replay attacks

## Dependencies

### Internal
- `core/interfaces/otp/` — IOtpController, IOtpRepository, IOtpService
- `core/interfaces/email/` — IEmailSender injected into the service
- `core/domain/` — OTP entity and request/response DTOs

<!-- MANUAL: -->
