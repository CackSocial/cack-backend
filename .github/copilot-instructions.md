# Copilot Instructions — SocialConnect Backend

## Code Style

- Go standard formatting (`gofmt`). No external linter config beyond golangci-lint defaults.
- Use `error` returns, not panics. Wrap errors with context when useful.
- Keep functions focused — one responsibility per function.
- Handlers should be thin: validate input → call use case → return response.
- Business logic belongs in use cases, never in handlers or repositories.
- Repository implementations should only contain data access logic.

## Conventions

- **Package naming:** lowercase, single word (e.g., `user`, `post`, `follow`)
- **Constructor pattern:** `NewXxx(deps...) *Xxx` — accepts interfaces, returns concrete
- **Repository interfaces** live in `internal/repository/`, implementations in `internal/infrastructure/database/repository/`
- **Use case structs** are in `internal/usecase/{feature}/usecase.go`
- **DTOs** are separate from domain models — never expose GORM models directly in API responses
- **Error mapping:** Use sentinel errors from `internal/usecase/errors/errors.go`. Map them to HTTP status in handler layer.

## Auth

- JWT HS256 tokens. Secret from `JWT_SECRET` env var.
- `c.Get("userID")` in handlers returns the authenticated user's UUID string.
- Protected routes use `middleware.AuthMiddleware`. Optional auth uses `middleware.OptionalAuth`.
- Never trust client-supplied user IDs for ownership checks — always use the JWT-derived userID.

## Database

- GORM with PostgreSQL. Auto-migration on startup.
- UUIDs as primary keys for User, Post, Comment, Message. `gen_random_uuid()` default.
- Use `Preload` for eager loading related entities. Use pagination for all list queries.
- Join table `post_tags` is managed by GORM's `many2many` tag.

## Testing

- Use Go's standard `testing` package. No external mock libraries.
- Create manual mock structs that implement repository interfaces.
- Test files go next to the code they test (e.g., `usecase_test.go` alongside `usecase.go`).

## Common Tasks

### Adding a new API endpoint to an existing feature
1. Add method to use case
2. Add DTO if needed
3. Add handler method with Swagger annotation
4. Register route in handler's `RegisterRoutes`
5. Run `make swagger` to regenerate docs

### Adding a new feature
Follow the checklist in AGENTS.md

### Modifying the database schema
1. Update the domain model in `internal/domain/`
2. GORM auto-migrate handles schema changes on restart
3. For destructive changes, manual SQL migration may be needed

## Do NOT

- Do not add fields to domain models without updating the corresponding DTOs and response mappers
- Do not bypass the use case layer — handlers must not access repositories directly
- Do not store secrets in code — use environment variables via `pkg/config`
- Do not use `panic` for recoverable errors
- Do not commit `.env` files (gitignored)
