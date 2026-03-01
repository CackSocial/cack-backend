package post

import (
	"errors"
	"mime/multipart"
	"testing"
	"time"

	"github.com/CackSocial/cack-backend/internal/domain"
	"github.com/CackSocial/cack-backend/internal/dto"
	"github.com/CackSocial/cack-backend/internal/repository"
	ucerrors "github.com/CackSocial/cack-backend/internal/usecase/errors"
)

// --- Mock PostRepository ---

type mockPostRepo struct {
	posts map[string]*domain.Post
}

func newMockPostRepo() *mockPostRepo {
	return &mockPostRepo{posts: make(map[string]*domain.Post)}
}

func (m *mockPostRepo) Create(post *domain.Post) error {
	if post.ID == "" {
		post.ID = "post-" + post.UserID + "-" + post.Content[:min(8, len(post.Content))]
	}
	post.CreatedAt = time.Now()
	m.posts[post.ID] = post
	return nil
}

func (m *mockPostRepo) GetByID(id string) (*domain.Post, error) {
	p, ok := m.posts[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return p, nil
}

func (m *mockPostRepo) GetByUserID(userID string, page, limit int) ([]domain.Post, int64, error) {
	var result []domain.Post
	for _, p := range m.posts {
		if p.UserID == userID {
			result = append(result, *p)
		}
	}
	return result, int64(len(result)), nil
}

func (m *mockPostRepo) Delete(id string) error {
	delete(m.posts, id)
	return nil
}

func (m *mockPostRepo) GetFeed(userIDs []string, page, limit int) ([]domain.Post, int64, error) {
	return nil, 0, nil
}

func (m *mockPostRepo) GetByTagName(tagName string, page, limit int) ([]domain.Post, int64, error) {
	return nil, 0, nil
}

func (m *mockPostRepo) IsReposted(userID, postID string) (bool, error) {
	for _, p := range m.posts {
		if p.UserID == userID && p.OriginalPostID != nil && *p.OriginalPostID == postID && p.PostType == "repost" {
			return true, nil
		}
	}
	return false, nil
}

func (m *mockPostRepo) CountReposts(postID string) (int64, error) {
	var count int64
	for _, p := range m.posts {
		if p.OriginalPostID != nil && *p.OriginalPostID == postID && (p.PostType == "repost" || p.PostType == "quote") {
			count++
		}
	}
	return count, nil
}

func (m *mockPostRepo) GetRepostByUser(userID, postID string) (*domain.Post, error) {
	for _, p := range m.posts {
		if p.UserID == userID && p.OriginalPostID != nil && *p.OriginalPostID == postID && p.PostType == "repost" {
			return p, nil
		}
	}
	return nil, errors.New("repost not found")
}

// --- Mock TagRepository ---

type mockTagRepo struct {
	tags   map[string]*domain.Tag
	nextID uint
}

func newMockTagRepo() *mockTagRepo {
	return &mockTagRepo{tags: make(map[string]*domain.Tag), nextID: 1}
}

func (m *mockTagRepo) FindOrCreate(name string) (*domain.Tag, error) {
	if t, ok := m.tags[name]; ok {
		return t, nil
	}
	tag := &domain.Tag{ID: m.nextID, Name: name}
	m.nextID++
	m.tags[name] = tag
	return tag, nil
}

func (m *mockTagRepo) GetByPostID(postID string) ([]domain.Tag, error) {
	return nil, nil
}

func (m *mockTagRepo) GetTrending(limit int, since time.Time) ([]repository.TrendingTag, error) {
	return nil, nil
}

// --- Mock LikeRepository ---

type mockLikeRepo struct{}

func (m *mockLikeRepo) Create(like *domain.Like) error              { return nil }
func (m *mockLikeRepo) Delete(userID, postID string) error          { return nil }
func (m *mockLikeRepo) GetByPostID(postID string, page, limit int) ([]domain.User, int64, error) {
	return nil, 0, nil
}
func (m *mockLikeRepo) CountByPostID(postID string) (int64, error) { return 0, nil }
func (m *mockLikeRepo) IsLiked(userID, postID string) (bool, error) { return false, nil }

// --- Mock CommentRepository ---

type mockCommentRepo struct{}

func (m *mockCommentRepo) Create(comment *domain.Comment) error { return nil }
func (m *mockCommentRepo) GetByID(id string) (*domain.Comment, error) { return nil, nil }
func (m *mockCommentRepo) GetByPostID(postID string, page, limit int) ([]domain.Comment, int64, error) {
	return nil, 0, nil
}
func (m *mockCommentRepo) Delete(id string) error                  { return nil }
func (m *mockCommentRepo) CountByPostID(postID string) (int64, error) { return 0, nil }

// --- Mock BookmarkRepository ---

type mockBookmarkRepo struct{}

func (m *mockBookmarkRepo) Create(bookmark *domain.Bookmark) error              { return nil }
func (m *mockBookmarkRepo) Delete(userID, postID string) error                  { return nil }
func (m *mockBookmarkRepo) GetByUserID(userID string, page, limit int) ([]domain.Bookmark, int64, error) {
	return nil, 0, nil
}
func (m *mockBookmarkRepo) IsBookmarked(userID, postID string) (bool, error) { return false, nil }

// --- Mock UserRepository ---

type mockUserRepo struct {
	users map[string]*domain.User
}

func newMockUserRepo() *mockUserRepo {
	return &mockUserRepo{users: make(map[string]*domain.User)}
}

func (m *mockUserRepo) Create(user *domain.User) error {
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

// --- Mock Storage ---

type mockStorage struct{}

func (m *mockStorage) Upload(file *multipart.FileHeader) (string, error) {
	return "http://example.com/image.jpg", nil
}

func (m *mockStorage) Delete(filePath string) error {
	return nil
}

// --- Helper ---

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// --- Tests ---

func TestParseHashtags_ViaCreate(t *testing.T) {
	// parseHashtags is unexported, so we test it indirectly through Create.
	postRepo := newMockPostRepo()
	tagRepo := newMockTagRepo()
	userRepo := newMockUserRepo()

	userRepo.users["user-1"] = &domain.User{
		ID:       "user-1",
		Username: "testuser",
	}

	uc := NewPostUseCase(postRepo, tagRepo, &mockLikeRepo{}, &mockCommentRepo{}, userRepo, &mockBookmarkRepo{}, &mockStorage{}, nil)

	resp, err := uc.Create("user-1", &dto.CreatePostRequest{
		Content: "Hello #world #golang #world",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Should have 2 unique tags: world, golang
	if len(resp.Tags) != 2 {
		t.Fatalf("expected 2 tags, got %d: %v", len(resp.Tags), resp.Tags)
	}

	tagSet := make(map[string]bool)
	for _, tag := range resp.Tags {
		tagSet[tag] = true
	}
	if !tagSet["world"] || !tagSet["golang"] {
		t.Fatalf("expected tags 'world' and 'golang', got %v", resp.Tags)
	}
}

func TestCreate_TextOnly(t *testing.T) {
	postRepo := newMockPostRepo()
	tagRepo := newMockTagRepo()
	userRepo := newMockUserRepo()

	userRepo.users["user-1"] = &domain.User{
		ID:       "user-1",
		Username: "testuser",
	}

	uc := NewPostUseCase(postRepo, tagRepo, &mockLikeRepo{}, &mockCommentRepo{}, userRepo, &mockBookmarkRepo{}, &mockStorage{}, nil)

	resp, err := uc.Create("user-1", &dto.CreatePostRequest{
		Content: "Just a simple post without hashtags",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp.Content != "Just a simple post without hashtags" {
		t.Fatalf("unexpected content: %q", resp.Content)
	}
	if resp.ImageURL != "" {
		t.Fatalf("expected empty image URL, got %q", resp.ImageURL)
	}
	if len(resp.Tags) != 0 {
		t.Fatalf("expected 0 tags, got %d", len(resp.Tags))
	}
	if resp.Author.Username != "testuser" {
		t.Fatalf("expected author username 'testuser', got %q", resp.Author.Username)
	}
}

func TestCreate_WithHashtags(t *testing.T) {
	postRepo := newMockPostRepo()
	tagRepo := newMockTagRepo()
	userRepo := newMockUserRepo()

	userRepo.users["user-1"] = &domain.User{
		ID:       "user-1",
		Username: "testuser",
	}

	uc := NewPostUseCase(postRepo, tagRepo, &mockLikeRepo{}, &mockCommentRepo{}, userRepo, &mockBookmarkRepo{}, &mockStorage{}, nil)

	resp, err := uc.Create("user-1", &dto.CreatePostRequest{
		Content: "Check out #Go and #Programming",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(resp.Tags) != 2 {
		t.Fatalf("expected 2 tags, got %d: %v", len(resp.Tags), resp.Tags)
	}
}

func TestDelete_OwnPost(t *testing.T) {
	postRepo := newMockPostRepo()
	tagRepo := newMockTagRepo()
	userRepo := newMockUserRepo()

	userRepo.users["user-1"] = &domain.User{ID: "user-1", Username: "testuser"}

	uc := NewPostUseCase(postRepo, tagRepo, &mockLikeRepo{}, &mockCommentRepo{}, userRepo, &mockBookmarkRepo{}, &mockStorage{}, nil)

	postRepo.posts["post-1"] = &domain.Post{
		ID:     "post-1",
		UserID: "user-1",
	}

	err := uc.Delete("post-1", "user-1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if _, ok := postRepo.posts["post-1"]; ok {
		t.Fatal("expected post to be deleted")
	}
}

func TestDelete_OtherUserPost_Unauthorized(t *testing.T) {
	postRepo := newMockPostRepo()
	tagRepo := newMockTagRepo()
	userRepo := newMockUserRepo()

	userRepo.users["user-1"] = &domain.User{ID: "user-1", Username: "testuser"}
	userRepo.users["user-2"] = &domain.User{ID: "user-2", Username: "otheruser"}

	uc := NewPostUseCase(postRepo, tagRepo, &mockLikeRepo{}, &mockCommentRepo{}, userRepo, &mockBookmarkRepo{}, &mockStorage{}, nil)

	postRepo.posts["post-1"] = &domain.Post{
		ID:     "post-1",
		UserID: "user-1",
	}

	err := uc.Delete("post-1", "user-2")
	if !errors.Is(err, ucerrors.ErrUnauthorized) {
		t.Fatalf("expected ErrUnauthorized, got %v", err)
	}
}
