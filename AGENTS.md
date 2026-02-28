# AGENTS.md — SocialConnect Backend

## Project Overview

**SocialConnect** is a social networking platform backend built with Go. It provides REST APIs and WebSocket support for a social media application with posts, follows, messaging, likes, comments, and trending tags.

- **Module:** `github.com/CackSocial/cack-backend`
- **Framework:** Gin (HTTP), Gorilla WebSocket (real-time DMs)
- **Database:** PostgreSQL via GORM (auto-migration, no manual SQL migrations)
- **Auth:** JWT (HS256) with username/password (no email)
- **Storage:** Local filesystem with abstraction layer (interface at `internal/infrastructure/storage/storage.go`)

## Architecture — Clean Architecture (4 layers)

```
Domain → Repository (interfaces) → Use Cases (business logic) → Handlers (HTTP/WS)
```

| Layer | Location | Purpose |
|-------|----------|---------|
| Domain entities | `internal/domain/` | GORM models: User, Post, Tag, Follow, Like, Comment, Message |
| Repository interfaces | `internal/repository/` | Abstract data access contracts |
| Repository implementations | `internal/infrastructure/database/repository/` | GORM-based concrete implementations |
| Use cases | `internal/usecase/{feature}/usecase.go` | Business logic, depends on repo interfaces |
| Handlers | `internal/handler/` | Gin HTTP handlers, depends on use cases |
| WebSocket | `internal/handler/ws/` | Hub + handler for real-time DMs |
| Middleware | `internal/middleware/` | Auth (required + optional), CORS, Prometheus metrics |
| Shared packages | `pkg/` | Config, JWT utils, bcrypt, response helpers |

## Key Patterns

- **Dependency injection:** Constructors accept interfaces, wired in `cmd/server/main.go`
- **Error handling:** Domain errors in `internal/usecase/errors/errors.go`, mapped to HTTP status in `internal/handler/helpers.go`
- **DTOs:** Request/response structs in `internal/dto/`, separate from domain models
- **Pagination:** All list endpoints accept `?page=N&limit=N` query params
- **Auth context:** `c.Get("userID")` returns authenticated user ID (set by middleware)
- **Optional auth:** Some endpoints use `OptionalAuth` middleware — sets userID if token present, empty string otherwise

## Running the Application

```bash
# Full stack with Docker
docker compose up -d --build

# Local development (requires PostgreSQL running)
cp .env.example .env    # edit DB credentials
make run

# Tests
make test

# Swagger docs regeneration (after changing handler annotations)
make swagger
```

## API Routes (26 endpoints)

All REST routes are prefixed with `/api/v1`. See `docs/swagger.json` or run the app and visit `/swagger/index.html`.

**Public:** Register, Login, Get profile, Get posts, Get followers/following, Get likes, Get comments, Trending tags, Posts by tag
**Protected (Bearer JWT):** Update profile, Create/Delete post, Follow/Unfollow, Timeline, Like/Unlike, Create/Delete comment, All messaging
**WebSocket:** `GET /api/v1/ws?token=<jwt>` for real-time DMs

## Database

PostgreSQL with GORM auto-migration. Tables: `users`, `posts`, `tags`, `post_tags` (join), `follows`, `likes`, `comments`, `messages`.

UUIDs for primary keys (User, Post, Comment, Message). Auto-generated via `gen_random_uuid()`.

## DevOps

- **Docker Compose:** app, postgres, nginx (reverse proxy), prometheus, grafana
- **CI/CD:** `.github/workflows/ci.yml` — lint → test → build
- **Monitoring:** Prometheus metrics at `/metrics`, Grafana at `:3000`
- **Makefile targets:** run, build, test, lint, swagger, docker-up, docker-down

## Adding a New Feature — Checklist

1. Add domain entity in `internal/domain/`
2. Add to auto-migrate list in `internal/infrastructure/database/postgres.go`
3. Define repository interface in `internal/repository/`
4. Implement repository in `internal/infrastructure/database/repository/`
5. Create use case in `internal/usecase/{feature}/usecase.go`
6. Add DTOs in `internal/dto/`
7. Create handler in `internal/handler/`
8. Wire everything in `cmd/server/main.go`
9. Add Swagger annotations, run `make swagger`
10. Write tests

## Testing

- Unit tests exist for: `pkg/auth`, `pkg/hash`, `internal/usecase/user`, `internal/usecase/post`
- Tests use manual mock structs (no external mock libraries)
- Run: `go test ./... -v -race`

## Important Files

| File | Purpose |
|------|---------|
| `cmd/server/main.go` | Entry point, DI wiring, router setup |
| `pkg/config/config.go` | Environment variable loading |
| `.env.example` | All config keys with defaults |
| `internal/usecase/errors/errors.go` | Shared domain error definitions |
| `internal/handler/helpers.go` | getUserID, getPagination, handleError |
| `internal/middleware/auth.go` | AuthMiddleware + OptionalAuth |
| `docs/SocialConnect.postman_collection.json` | Postman collection for all endpoints |
