<!-- Parent: ../AGENTS.md -->
<!-- Generated: 2026-05-25 | Updated: 2026-05-25 -->

# core/interfaces/email

## Purpose
Interface contract for the email sending capability. Decouples the OTP module (and any other email consumer) from the concrete Resend adapter.

## Key Files

| File | Description |
|------|-------------|
| `IEmail.sender.go` | Email sending contract — send transactional emails with subject and HTML body |

## Dependencies

### Internal
- No domain types — method signatures use only primitives (`string`)

<!-- MANUAL: -->
