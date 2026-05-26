<!-- Parent: ../AGENTS.md -->
<!-- Generated: 2026-05-25 | Updated: 2026-05-25 -->

# modules/user

## Purpose
User profile CRUD operations. Authenticated users can read and update their profile data including name, bio, address, and category assignments. Profile images are handled via the uploader module.

## Key Files

| File | Description |
|------|-------------|
| `user.controller.go` | Gin handlers for GET/PUT /user/profile/:id |
| `user.controller_test.go` | Controller unit tests with mocked user service |
| `user.service.go` | Profile retrieval and update business logic |
| `user.service_test.go` | Service unit tests with mocked repository |

## For AI Agents

### Working In This Directory
- Users can only update their own profile — validate that the token subject matches the `:id` param
- Category updates use the `user_categories` join table — replace all entries in a single transaction
- See `docs/app-specs/profile.md` for the full field specification

## Dependencies

### Internal
- `core/interfaces/user/` — IUserController, IUserRepository, IUserService
- `core/domain/` — user entity and request/response DTOs, user_category entity

<!-- MANUAL: -->
