<!-- Parent: ../AGENTS.md -->
<!-- Generated: 2026-05-25 | Updated: 2026-05-25 -->

# modules/proposal

## Purpose
Business-to-helper service request flows. A client creates a proposal targeting a helper's service. The proposal goes through a lifecycle: pending → accepted/rejected → completed. Proposals are the prerequisite for creating a review.

## Key Files

| File | Description |
|------|-------------|
| `proposal.controller.go` | Gin handlers for proposal CRUD and status transitions |
| `proposal.controller_test.go` | Controller unit tests with mocked proposal service |
| `proposal.service.go` | Proposal creation, status transitions, and retrieval logic |
| `proposal.service_test.go` | Service unit tests with mocked repository |

## For AI Agents

### Working In This Directory
- Only authenticated users (clients) can create proposals
- Status transitions must follow the defined lifecycle — validate current status before transitioning
- A proposal can only receive a review once it reaches `completed` status
- See `docs/app-specs/proposals.md` for the full business rules

## Dependencies

### Internal
- `core/interfaces/proposal/` — IProposalController, IProposalRepository, IProposalService
- `core/domain/` — proposal entity and request/response DTOs

<!-- MANUAL: -->
