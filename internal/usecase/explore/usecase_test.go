package explore

import (
	"errors"
	"testing"
	"time"

	"github.com/CackSocial/cack-backend/internal/domain"
	"github.com/CackSocial/cack-backend/internal/repository"
)

// --- Mock UserRepository ---

type mockUserRepo struct {
	users     map[string]*domain.User
	suggested []repository.SuggestedUser
}

func newMockUserRepo() *mockUserRepo {
	return &mockUserRepo{users: make(map[string]*domain.User)}
}

func (m *mockUserRepo) Create(user *domain.User) error                 { return nil }
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
func (m *mockUserRepo) Update(user *domain.User) error                                       { return nil }
func (m *mockUserRepo) Delete(id string) error                                               { return nil }
func (m *mockUserRepo) Search(query string, page, limit int) ([]domain.User, int64, error)   { return nil, 0, nil }

func (m *mockUserRepo) GetSuggestedUsers(currentUserID string, followingIDs []string, limit int) ([]repository.SuggestedUser, error) {
	if m.suggested != nil {
		if len(m.suggested) > limit {
			return m.suggested[:limit], nil
		}
		return m.suggested, nil
	}
	return []repository.SuggestedUser{}, nil
}

// --- Mock PostRepository ---

type mockPostRepo struct {
	popularPosts  []domain.Post
	popularTotal  int64
	discoverPosts []domain.Post
	discoverTotal int64
}

func (m *mockPostRepo) Create(post *domain.Post) error { return nil }
func (m *mockPostRepo) GetByID(id string) (*domain.Post, error) { return nil, nil }
func (m *mockPostRepo) GetByUserID(userID string, page, limit int) ([]domain.Post, int64, error) { return nil, 0, nil }
func (m *mockPostRepo) Delete(id string) error { return nil }
func (m *mockPostRepo) GetFeed(userIDs []string, page, limit int) ([]domain.Post, int64, error) { return nil, 0, nil }
func (m *mockPostRepo) GetByTagName(tagName string, page, limit int) ([]domain.Post, int64, error) { return nil, 0, nil }
func (m *mockPostRepo) IsReposted(userID, postID string) (bool, error) { return false, nil }
func (m *mockPostRepo) CountReposts(postID string) (int64, error) { return 0, nil }
func (m *mockPostRepo) GetRepostByUser(userID, postID string) (*domain.Post, error) { return nil, nil }

func (m *mockPostRepo) GetPopularPosts(excludeUserIDs []string, page, limit int, since time.Time) ([]domain.Post, int64, error) {
	return m.popularPosts, m.popularTotal, nil
}

func (m *mockPostRepo) GetDiscoverPosts(tagNames []string, excludeUserIDs []string, page, limit int) ([]domain.Post, int64, error) {
	if len(tagNames) == 0 {
		return nil, 0, nil
	}
	return m.discoverPosts, m.discoverTotal, nil
}

// --- Mock FollowRepository ---

type mockFollowRepo struct {
	followingIDs []string
}

func (m *mockFollowRepo) Create(follow *domain.Follow) error { return nil }
func (m *mockFollowRepo) Delete(followerID, followingID string) error { return nil }
func (m *mockFollowRepo) IsFollowing(followerID, followingID string) (bool, error) { return false, nil }
func (m *mockFollowRepo) GetFollowers(userID string, page, limit int) ([]domain.User, int64, error) { return nil, 0, nil }
func (m *mockFollowRepo) GetFollowing(userID string, page, limit int) ([]domain.User, int64, error) { return nil, 0, nil }

func (m *mockFollowRepo) GetFollowingIDs(userID string) ([]string, error) {
	return m.followingIDs, nil
}

// --- Mock LikeRepository ---

type mockLikeRepo struct {
	likedTagNames []string
}

func (m *mockLikeRepo) Create(like *domain.Like) error { return nil }
func (m *mockLikeRepo) Delete(userID, postID string) error { return nil }
func (m *mockLikeRepo) GetByPostID(postID string, page, limit int) ([]domain.User, int64, error) { return nil, 0, nil }
func (m *mockLikeRepo) GetLikedPostsByUserID(userID string, page, limit int) ([]domain.Post, int64, error) { return nil, 0, nil }
func (m *mockLikeRepo) CountByPostID(postID string) (int64, error) { return 0, nil }
func (m *mockLikeRepo) IsLiked(userID, postID string) (bool, error) { return false, nil }

func (m *mockLikeRepo) GetLikedTagNames(userID string, limit int) ([]string, error) {
	return m.likedTagNames, nil
}

// --- Mock CommentRepository ---

type mockCommentRepo struct{}

func (m *mockCommentRepo) Create(comment *domain.Comment) error { return nil }
func (m *mockCommentRepo) GetByID(id string) (*domain.Comment, error) { return nil, nil }
func (m *mockCommentRepo) GetByPostID(postID string, page, limit int) ([]domain.Comment, int64, error) { return nil, 0, nil }
func (m *mockCommentRepo) Delete(id string) error { return nil }
func (m *mockCommentRepo) CountByPostID(postID string) (int64, error) { return 0, nil }

// --- Mock BookmarkRepository ---

type mockBookmarkRepo struct{}

func (m *mockBookmarkRepo) Create(bookmark *domain.Bookmark) error { return nil }
func (m *mockBookmarkRepo) Delete(userID, postID string) error { return nil }
func (m *mockBookmarkRepo) GetByUserID(userID string, page, limit int) ([]domain.Bookmark, int64, error) { return nil, 0, nil }
func (m *mockBookmarkRepo) IsBookmarked(userID, postID string) (bool, error) { return false, nil }

// --- Helper ---

func makePost(id, userID, username, content string, tags []domain.Tag) domain.Post {
	return domain.Post{
		ID:        id,
		UserID:    userID,
		User:      domain.User{ID: userID, Username: username, DisplayName: username},
		Content:   content,
		PostType:  "original",
		Tags:      tags,
		CreatedAt: time.Now(),
	}
}

// --- Tests: GetSuggestedUsers ---

func TestGetSuggestedUsers_WithMutuals(t *testing.T) {
	userRepo := newMockUserRepo()
	userRepo.suggested = []repository.SuggestedUser{
		{User: domain.User{ID: "u2", Username: "alice", DisplayName: "Alice"}, MutualCount: 3},
		{User: domain.User{ID: "u3", Username: "bob", DisplayName: "Bob"}, MutualCount: 1},
	}

	uc := NewExploreUseCase(
		userRepo,
		&mockPostRepo{},
		&mockFollowRepo{followingIDs: []string{"f1", "f2"}},
		&mockLikeRepo{},
		&mockCommentRepo{},
		&mockBookmarkRepo{},
	)

	results, err := uc.GetSuggestedUsers("me", 10)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 suggestions, got %d", len(results))
	}
	if results[0].Username != "alice" {
		t.Fatalf("expected first suggestion to be alice, got %q", results[0].Username)
	}
	if results[0].MutualFollowerCount != 3 {
		t.Fatalf("expected mutual count 3, got %d", results[0].MutualFollowerCount)
	}
}

func TestGetSuggestedUsers_Empty(t *testing.T) {
	userRepo := newMockUserRepo()

	uc := NewExploreUseCase(
		userRepo,
		&mockPostRepo{},
		&mockFollowRepo{},
		&mockLikeRepo{},
		&mockCommentRepo{},
		&mockBookmarkRepo{},
	)

	results, err := uc.GetSuggestedUsers("me", 10)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(results) != 0 {
		t.Fatalf("expected 0 suggestions, got %d", len(results))
	}
}

// --- Tests: GetPopularPosts ---

func TestGetPopularPosts_Success(t *testing.T) {
	post1 := makePost("p1", "u1", "alice", "Popular post", nil)
	post2 := makePost("p2", "u2", "bob", "Another popular", nil)

	uc := NewExploreUseCase(
		newMockUserRepo(),
		&mockPostRepo{
			popularPosts: []domain.Post{post1, post2},
			popularTotal: 2,
		},
		&mockFollowRepo{followingIDs: []string{"f1"}},
		&mockLikeRepo{},
		&mockCommentRepo{},
		&mockBookmarkRepo{},
	)

	posts, total, err := uc.GetPopularPosts("me", 1, 20)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if total != 2 {
		t.Fatalf("expected total 2, got %d", total)
	}
	if len(posts) != 2 {
		t.Fatalf("expected 2 posts, got %d", len(posts))
	}
	if posts[0].Author.Username != "alice" {
		t.Fatalf("expected author alice, got %q", posts[0].Author.Username)
	}
	if posts[0].PostType != "original" {
		t.Fatalf("expected post type original, got %q", posts[0].PostType)
	}
}

func TestGetPopularPosts_Empty(t *testing.T) {
	uc := NewExploreUseCase(
		newMockUserRepo(),
		&mockPostRepo{},
		&mockFollowRepo{},
		&mockLikeRepo{},
		&mockCommentRepo{},
		&mockBookmarkRepo{},
	)

	posts, total, err := uc.GetPopularPosts("me", 1, 20)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if total != 0 {
		t.Fatalf("expected total 0, got %d", total)
	}
	if len(posts) != 0 {
		t.Fatalf("expected 0 posts, got %d", len(posts))
	}
}

// --- Tests: GetDiscoverFeed ---

func TestGetDiscoverFeed_WithAffinityTags(t *testing.T) {
	tags := []domain.Tag{{ID: 1, Name: "golang"}}
	post1 := makePost("p1", "u1", "alice", "Go is great #golang", tags)

	uc := NewExploreUseCase(
		newMockUserRepo(),
		&mockPostRepo{
			discoverPosts: []domain.Post{post1},
			discoverTotal: 1,
		},
		&mockFollowRepo{followingIDs: []string{"f1"}},
		&mockLikeRepo{likedTagNames: []string{"golang", "rust"}},
		&mockCommentRepo{},
		&mockBookmarkRepo{},
	)

	posts, total, err := uc.GetDiscoverFeed("me", 1, 20)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if total != 1 {
		t.Fatalf("expected total 1, got %d", total)
	}
	if len(posts) != 1 {
		t.Fatalf("expected 1 post, got %d", len(posts))
	}
	if len(posts[0].Tags) != 1 || posts[0].Tags[0] != "golang" {
		t.Fatalf("expected tag golang, got %v", posts[0].Tags)
	}
}

func TestGetDiscoverFeed_NoLikedTags(t *testing.T) {
	uc := NewExploreUseCase(
		newMockUserRepo(),
		&mockPostRepo{
			discoverPosts: []domain.Post{makePost("p1", "u1", "alice", "post", nil)},
			discoverTotal: 1,
		},
		&mockFollowRepo{},
		&mockLikeRepo{likedTagNames: nil},
		&mockCommentRepo{},
		&mockBookmarkRepo{},
	)

	posts, total, err := uc.GetDiscoverFeed("me", 1, 20)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if total != 0 {
		t.Fatalf("expected total 0 (no affinity tags), got %d", total)
	}
	if len(posts) != 0 {
		t.Fatalf("expected 0 posts, got %d", len(posts))
	}
}

func TestGetPopularPosts_EnrichesPosts(t *testing.T) {
	tags := []domain.Tag{{ID: 1, Name: "tech"}, {ID: 2, Name: "news"}}
	post := makePost("p1", "u1", "alice", "Tech news #tech #news", tags)

	uc := NewExploreUseCase(
		newMockUserRepo(),
		&mockPostRepo{
			popularPosts: []domain.Post{post},
			popularTotal: 1,
		},
		&mockFollowRepo{},
		&mockLikeRepo{},
		&mockCommentRepo{},
		&mockBookmarkRepo{},
	)

	posts, _, err := uc.GetPopularPosts("me", 1, 20)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(posts) != 1 {
		t.Fatalf("expected 1 post, got %d", len(posts))
	}

	p := posts[0]
	if p.ID != "p1" {
		t.Fatalf("expected post ID p1, got %q", p.ID)
	}
	if len(p.Tags) != 2 {
		t.Fatalf("expected 2 tags, got %d", len(p.Tags))
	}
	// Mock returns 0 for all counts and false for all flags
	if p.LikeCount != 0 || p.CommentCount != 0 {
		t.Fatalf("expected 0 counts from mocks, got like=%d comment=%d", p.LikeCount, p.CommentCount)
	}
	if p.IsLiked || p.IsBookmarked || p.IsReposted {
		t.Fatal("expected all flags to be false from mocks")
	}
}
