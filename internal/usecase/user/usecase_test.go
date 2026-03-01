package user

import (
	"errors"
	"mime/multipart"
	"testing"

	"github.com/CackSocial/cack-backend/internal/domain"
	"github.com/CackSocial/cack-backend/internal/dto"
	ucerrors "github.com/CackSocial/cack-backend/internal/usecase/errors"
	"github.com/CackSocial/cack-backend/pkg/hash"
)

// --- Mock UserRepository ---

type mockUserRepo struct {
	users map[string]*domain.User // keyed by ID
}

func newMockUserRepo() *mockUserRepo {
	return &mockUserRepo{users: make(map[string]*domain.User)}
}

func (m *mockUserRepo) Create(user *domain.User) error {
	if user.ID == "" {
		user.ID = "generated-id-" + user.Username
	}
	m.users[user.ID] = user
	return nil
}

func (m *mockUserRepo) GetByID(id string) (*domain.User, error) {
	u, ok := m.users[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return u, nil
}

func (m *mockUserRepo) GetByUsername(username string) (*domain.User, error) {
	for _, u := range m.users {
		if u.Username == username {
			return u, nil
		}
	}
	return nil, errors.New("not found")
}

func (m *mockUserRepo) Update(user *domain.User) error {
	m.users[user.ID] = user
	return nil
}

func (m *mockUserRepo) Search(query string, page, limit int) ([]domain.User, int64, error) {
	return nil, 0, nil
}

func (m *mockUserRepo) Delete(id string) error {
	delete(m.users, id)
	return nil
}

// --- Mock FollowRepository ---

type mockFollowRepo struct {
	follows map[string]map[string]bool // followerID -> followingID -> exists
}

func newMockFollowRepo() *mockFollowRepo {
	return &mockFollowRepo{follows: make(map[string]map[string]bool)}
}

func (m *mockFollowRepo) Create(follow *domain.Follow) error {
	if m.follows[follow.FollowerID] == nil {
		m.follows[follow.FollowerID] = make(map[string]bool)
	}
	m.follows[follow.FollowerID][follow.FollowingID] = true
	return nil
}

func (m *mockFollowRepo) Delete(followerID, followingID string) error {
	if m.follows[followerID] != nil {
		delete(m.follows[followerID], followingID)
	}
	return nil
}

func (m *mockFollowRepo) IsFollowing(followerID, followingID string) (bool, error) {
	if m.follows[followerID] != nil {
		return m.follows[followerID][followingID], nil
	}
	return false, nil
}

func (m *mockFollowRepo) GetFollowers(userID string, page, limit int) ([]domain.User, int64, error) {
	var count int64
	for _, followings := range m.follows {
		if followings[userID] {
			count++
		}
	}
	return nil, count, nil
}

func (m *mockFollowRepo) GetFollowing(userID string, page, limit int) ([]domain.User, int64, error) {
	var count int64
	if m.follows[userID] != nil {
		count = int64(len(m.follows[userID]))
	}
	return nil, count, nil
}

func (m *mockFollowRepo) GetFollowingIDs(userID string) ([]string, error) {
	var ids []string
	if m.follows[userID] != nil {
		for id := range m.follows[userID] {
			ids = append(ids, id)
		}
	}
	return ids, nil
}

// --- Mock Storage ---

type mockStorage struct {
	uploaded string
	deleted  string
}

func (m *mockStorage) Upload(file *multipart.FileHeader) (string, error) {
	m.uploaded = file.Filename
	return "/uploads/" + file.Filename, nil
}

func (m *mockStorage) Delete(filePath string) error {
	m.deleted = filePath
	return nil
}

// --- Tests ---

func TestRegister_Success(t *testing.T) {
	userRepo := newMockUserRepo()
	followRepo := newMockFollowRepo()
	uc := NewUserUseCase(userRepo, followRepo, &mockStorage{}, "testsecret", 1)

	resp, err := uc.Register(&dto.RegisterRequest{
		Username:    "testuser",
		Password:    "password123",
		DisplayName: "Test User",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp.Token == "" {
		t.Fatal("expected non-empty token")
	}
	if resp.User.Username != "testuser" {
		t.Fatalf("expected username 'testuser', got %q", resp.User.Username)
	}
}

func TestRegister_DuplicateUsername(t *testing.T) {
	userRepo := newMockUserRepo()
	followRepo := newMockFollowRepo()
	uc := NewUserUseCase(userRepo, followRepo, &mockStorage{}, "testsecret", 1)

	_, err := uc.Register(&dto.RegisterRequest{
		Username: "testuser",
		Password: "password123",
	})
	if err != nil {
		t.Fatalf("first register failed: %v", err)
	}

	_, err = uc.Register(&dto.RegisterRequest{
		Username: "testuser",
		Password: "password456",
	})
	if !errors.Is(err, ucerrors.ErrUsernameTaken) {
		t.Fatalf("expected ErrUsernameTaken, got %v", err)
	}
}

func TestLogin_Success(t *testing.T) {
	userRepo := newMockUserRepo()
	followRepo := newMockFollowRepo()
	uc := NewUserUseCase(userRepo, followRepo, &mockStorage{}, "testsecret", 1)

	// First register the user.
	hashed, _ := hash.HashPassword("password123")
	userRepo.users["user-1"] = &domain.User{
		ID:       "user-1",
		Username: "testuser",
		Password: hashed,
	}

	resp, err := uc.Login(&dto.LoginRequest{
		Username: "testuser",
		Password: "password123",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp.Token == "" {
		t.Fatal("expected non-empty token")
	}
}

func TestLogin_WrongPassword(t *testing.T) {
	userRepo := newMockUserRepo()
	followRepo := newMockFollowRepo()
	uc := NewUserUseCase(userRepo, followRepo, &mockStorage{}, "testsecret", 1)

	hashed, _ := hash.HashPassword("password123")
	userRepo.users["user-1"] = &domain.User{
		ID:       "user-1",
		Username: "testuser",
		Password: hashed,
	}

	_, err := uc.Login(&dto.LoginRequest{
		Username: "testuser",
		Password: "wrongpassword",
	})
	if !errors.Is(err, ucerrors.ErrInvalidCredentials) {
		t.Fatalf("expected ErrInvalidCredentials, got %v", err)
	}
}

func TestLogin_NonExistentUser(t *testing.T) {
	userRepo := newMockUserRepo()
	followRepo := newMockFollowRepo()
	uc := NewUserUseCase(userRepo, followRepo, &mockStorage{}, "testsecret", 1)

	_, err := uc.Login(&dto.LoginRequest{
		Username: "nouser",
		Password: "password123",
	})
	if !errors.Is(err, ucerrors.ErrInvalidCredentials) {
		t.Fatalf("expected ErrInvalidCredentials, got %v", err)
	}
}

func TestGetProfile_Success(t *testing.T) {
	userRepo := newMockUserRepo()
	followRepo := newMockFollowRepo()
	uc := NewUserUseCase(userRepo, followRepo, &mockStorage{}, "testsecret", 1)

	userRepo.users["user-1"] = &domain.User{
		ID:          "user-1",
		Username:    "testuser",
		DisplayName: "Test User",
		Bio:         "Hello world",
	}

	profile, err := uc.GetProfile("testuser", "")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if profile.Username != "testuser" {
		t.Fatalf("expected username 'testuser', got %q", profile.Username)
	}
	if profile.DisplayName != "Test User" {
		t.Fatalf("expected display name 'Test User', got %q", profile.DisplayName)
	}
	if profile.Bio != "Hello world" {
		t.Fatalf("expected bio 'Hello world', got %q", profile.Bio)
	}
}

func TestUpdateProfile_Success(t *testing.T) {
	userRepo := newMockUserRepo()
	followRepo := newMockFollowRepo()
	uc := NewUserUseCase(userRepo, followRepo, &mockStorage{}, "testsecret", 1)

	userRepo.users["user-1"] = &domain.User{
		ID:          "user-1",
		Username:    "testuser",
		DisplayName: "Old Name",
		Bio:         "Old bio",
	}

	newName := "New Name"
	newBio := "New bio"
	profile, err := uc.UpdateProfile("user-1", &dto.UpdateProfileRequest{
		DisplayName: &newName,
		Bio:         &newBio,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if profile.DisplayName != "New Name" {
		t.Fatalf("expected display name 'New Name', got %q", profile.DisplayName)
	}
	if profile.Bio != "New bio" {
		t.Fatalf("expected bio 'New bio', got %q", profile.Bio)
	}
}
