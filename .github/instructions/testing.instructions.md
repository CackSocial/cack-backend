# Testing Guide

## Test Structure

Tests use Go's standard `testing` package with manual mock structs. No external mock libraries.

## Writing Tests for Use Cases

1. Create mock structs that implement repository interfaces in the test file:

```go
type mockUserRepo struct {
    users map[string]*domain.User
}
func (m *mockUserRepo) Create(user *domain.User) error { ... }
func (m *mockUserRepo) GetByID(id string) (*domain.User, error) { ... }
// ... implement all interface methods
```

2. Test happy path AND error cases:

```go
func TestRegister_Success(t *testing.T) { ... }
func TestRegister_DuplicateUsername(t *testing.T) { ... }
func TestLogin_WrongPassword(t *testing.T) { ... }
```

3. Use `t.Errorf` or `t.Fatalf` for assertions. Compare with expected values.

## Existing Test Files

- `pkg/auth/jwt_test.go` — JWT generation, validation, extraction
- `pkg/hash/password_test.go` — bcrypt hash and check
- `internal/usecase/user/usecase_test.go` — register, login, profile, update
- `internal/usecase/post/usecase_test.go` — create, delete, hashtag parsing

## Running Tests

```bash
go test ./... -v -race          # all tests
go test ./pkg/auth/ -v          # specific package
go test -run TestRegister ./... # specific test
```

## Adding Tests for a New Feature

Place test file next to the code: `internal/usecase/feature/usecase_test.go`

Test at minimum:
- Success case
- Not found errors
- Authorization failures (wrong user trying to modify)
- Validation edge cases
