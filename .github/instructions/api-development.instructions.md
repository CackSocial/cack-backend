# API Development Guide

## Route Registration Pattern

Routes are registered via `RegisterRoutes` methods on handler structs. There are three route groups:

- `public` — no authentication required
- `protected` — `AuthMiddleware` required (Bearer JWT)
- `optionalAuth` — middleware passed as `gin.HandlerFunc`, sets userID if token present

## Handler Method Pattern

```go
func (h *XxxHandler) MethodName(c *gin.Context) {
    // 1. Get authenticated user (if applicable)
    userID := getUserID(c)

    // 2. Parse path params
    id := c.Param("id")

    // 3. Bind request body or query params
    var req dto.SomeRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        response.Error(c, http.StatusBadRequest, err.Error())
        return
    }

    // 4. Call use case
    result, err := h.useCase.Method(userID, &req)
    if err != nil {
        handleError(c, err)
        return
    }

    // 5. Return response
    response.Success(c, http.StatusOK, result)
}
```

## Pagination

All list endpoints use `getPagination(c)` which returns `(page, limit)` from query params with defaults (page=1, limit=20).

Return paginated responses with `response.Paginated(c, data, page, limit, total)`.

## Multipart File Uploads

For endpoints that accept files (posts, messages), use `formdata` content type:
```go
var req dto.CreatePostRequest
if err := c.ShouldBind(&req); err != nil { ... }
// req.Image is *multipart.FileHeader (nil if no file)
```

## Swagger Annotations

Every handler method must have Swagger annotations. After changes, regenerate with:
```bash
make swagger
# or: swag init -g cmd/server/main.go -o docs
```

## Error Handling

Use sentinel errors from `internal/usecase/errors/errors.go`. The `handleError` function in `internal/handler/helpers.go` maps them to HTTP status codes:

| Error | HTTP Status |
|-------|-------------|
| ErrUserNotFound, ErrPostNotFound, ErrCommentNotFound | 404 |
| ErrInvalidCredentials, ErrUnauthorized | 401 |
| ErrUsernameTaken | 409 |
| ErrSelfFollow, ErrAlreadyFollowing, ErrAlreadyLiked | 400 |
| (default) | 500 |
