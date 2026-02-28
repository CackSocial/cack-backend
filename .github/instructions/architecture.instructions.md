# Architecture Guide

This codebase follows Clean Architecture with strict layer boundaries.

## Layer Dependency Rules

```
Handler → Use Case → Repository Interface ← Repository Implementation
                   → Storage Interface    ← Storage Implementation
```

- Handlers depend on Use Cases (never on Repositories directly)
- Use Cases depend on Repository/Storage interfaces (never on implementations)
- Implementations depend on GORM/filesystem (infrastructure concerns)
- Domain entities have no dependencies

## Wiring

All dependency injection happens in `cmd/server/main.go`:
1. Config → Database → Repositories → Storage → Use Cases → Handlers → Router

## Adding a New Domain Entity

When adding a new entity (e.g., Notification):

```
internal/domain/notification.go          # GORM model
internal/repository/notification.go      # Interface
internal/infrastructure/database/repository/notification.go  # GORM impl
internal/usecase/notification/usecase.go # Business logic
internal/dto/notification.go             # Request/Response types
internal/handler/notification.go         # HTTP handler
```

Then wire in `cmd/server/main.go` following the existing pattern.

## Key Interfaces

- `repository.UserRepository` — user CRUD + search
- `repository.PostRepository` — post CRUD + feed + tag-based queries
- `repository.TagRepository` — find-or-create + trending
- `repository.FollowRepository` — follow/unfollow + follower/following lists
- `repository.LikeRepository` — like/unlike + counts
- `repository.CommentRepository` — comment CRUD + counts
- `repository.MessageRepository` — messages + conversations + mark-read
- `storage.Storage` — file upload/delete (currently local, designed for swap to cloud)
