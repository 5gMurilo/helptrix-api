<!-- Parent: ../AGENTS.md -->
<!-- Generated: 2026-05-25 | Updated: 2026-05-25 -->

# adapter/email

## Purpose
Transactional email adapter using the Resend service. Implements the `IEmailSender` interface from `core/interfaces/email/`. Used primarily by the OTP module to send verification codes.

## Key Files

| File | Description |
|------|-------------|
| `resend.go` | Resend SDK wrapper implementing IEmailSender |
| `resend_test.go` | Unit tests for the email adapter |

## For AI Agents

### Working In This Directory
- The Resend API key is read from `RESEND_API_KEY` environment variable
- Sender address is configured via `EMAIL_FROM` environment variable
- Do not add business logic here — only email transport concerns belong in this file

### Common Patterns
- `NewResendEmailSender(apiKey string) IEmailSender`
- `SendEmail(to, subject, htmlBody string) error`

## Dependencies

### Internal
- `core/interfaces/email/IEmail.sender.go` — interface this adapter implements

### External
- `github.com/resend/resend-go/v2`

<!-- MANUAL: -->
