// Package post implements the business logic for creating, retrieving,
// and deleting posts, including image upload and hashtag extraction.
package post

import (
	"regexp"
	"strings"

	"github.com/CackSocial/cack-backend/internal/domain"
	"github.com/CackSocial/cack-backend/internal/dto"
	"github.com/CackSocial/cack-backend/internal/infrastructure/storage"
	"github.com/CackSocial/cack-backend/internal/repository"
	ucerrors "github.com/CackSocial/cack-backend/internal/usecase/errors"
)

var hashtagRegex = regexp.MustCompile(`#(\w+)`)

// PostUseCase encapsulates all post-related business logic including
// creation with image upload, hashtag parsing, retrieval, and deletion.
type PostUseCase struct {
	postRepo    repository.PostRepository
	tagRepo     repository.TagRepository
	likeRepo    repository.LikeRepository
	commentRepo repository.CommentRepository
	userRepo    repository.UserRepository
	storage     storage.Storage
}

// NewPostUseCase creates a new PostUseCase with the given dependencies.
func NewPostUseCase(
	postRepo repository.PostRepository,
	tagRepo repository.TagRepository,
	likeRepo repository.LikeRepository,
	commentRepo repository.CommentRepository,
	userRepo repository.UserRepository,
	storage storage.Storage,
) *PostUseCase {
	return &PostUseCase{
		postRepo:    postRepo,
		tagRepo:     tagRepo,
		likeRepo:    likeRepo,
		commentRepo: commentRepo,
		userRepo:    userRepo,
		storage:     storage,
	}
}

// Create creates a new post for the given user. If the request includes an
// image, it is uploaded via the storage backend. Hashtags are extracted from
// the content and linked to the post as tags.
func (uc *PostUseCase) Create(userID string, req *dto.CreatePostRequest) (*dto.PostResponse, error) {
	user, err := uc.userRepo.GetByID(userID)
	if err != nil || user == nil {
		return nil, ucerrors.ErrUserNotFound
	}

	var imageURL string
	if req.Image != nil {
		url, err := uc.storage.Upload(req.Image)
		if err != nil {
			return nil, err
		}
		imageURL = url
	}

	// Parse hashtags and find-or-create each tag.
	tagNames := parseHashtags(req.Content)
	var tags []domain.Tag
	for _, name := range tagNames {
		tag, err := uc.tagRepo.FindOrCreate(name)
		if err != nil {
			return nil, err
		}
		tags = append(tags, *tag)
	}

	post := &domain.Post{
		UserID:   userID,
		Content:  req.Content,
		ImageURL: imageURL,
		Tags:     tags,
	}

	if err := uc.postRepo.Create(post); err != nil {
		return nil, err
	}

	// Populate the User field for the response.
	post.User = *user

	return uc.toPostResponse(post, "")
}

// GetByID retrieves a single post by its ID, including like/comment counts
// and whether the current user has liked it.
func (uc *PostUseCase) GetByID(postID string, currentUserID string) (*dto.PostResponse, error) {
	post, err := uc.postRepo.GetByID(postID)
	if err != nil || post == nil {
		return nil, ucerrors.ErrPostNotFound
	}

	return uc.toPostResponse(post, currentUserID)
}

// GetByUserID retrieves paginated posts for a user identified by username.
// Each post includes like/comment counts and the liked status for the
// current user.
func (uc *PostUseCase) GetByUserID(username string, currentUserID string, page, limit int) ([]dto.PostResponse, int64, error) {
	user, err := uc.userRepo.GetByUsername(username)
	if err != nil || user == nil {
		return nil, 0, ucerrors.ErrUserNotFound
	}

	posts, total, err := uc.postRepo.GetByUserID(user.ID, page, limit)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]dto.PostResponse, 0, len(posts))
	for i := range posts {
		resp, err := uc.toPostResponse(&posts[i], currentUserID)
		if err != nil {
			return nil, 0, err
		}
		responses = append(responses, *resp)
	}

	return responses, total, nil
}

// Delete removes a post by its ID. Only the post owner is allowed to delete.
// If the post has an attached image, it is also removed from storage.
func (uc *PostUseCase) Delete(postID string, userID string) error {
	post, err := uc.postRepo.GetByID(postID)
	if err != nil || post == nil {
		return ucerrors.ErrPostNotFound
	}

	if post.UserID != userID {
		return ucerrors.ErrUnauthorized
	}

	// Remove the image from storage if present.
	if post.ImageURL != "" {
		_ = uc.storage.Delete(post.ImageURL)
	}

	return uc.postRepo.Delete(postID)
}

// toPostResponse converts a domain Post into a PostResponse DTO, filling
// in like count, comment count, tag names, and the liked status for the
// current user.
func (uc *PostUseCase) toPostResponse(post *domain.Post, currentUserID string) (*dto.PostResponse, error) {
	likeCount, _ := uc.likeRepo.CountByPostID(post.ID)
	commentCount, _ := uc.commentRepo.CountByPostID(post.ID)

	var isLiked bool
	if currentUserID != "" {
		isLiked, _ = uc.likeRepo.IsLiked(currentUserID, post.ID)
	}

	tagNames := make([]string, 0, len(post.Tags))
	for _, t := range post.Tags {
		tagNames = append(tagNames, t.Name)
	}

	return &dto.PostResponse{
		ID:       post.ID,
		Content:  post.Content,
		ImageURL: post.ImageURL,
		Author: dto.UserProfile{
			ID:          post.User.ID,
			Username:    post.User.Username,
			DisplayName: post.User.DisplayName,
			Bio:         post.User.Bio,
			AvatarURL:   post.User.AvatarURL,
		},
		Tags:         tagNames,
		LikeCount:    likeCount,
		CommentCount: commentCount,
		IsLiked:      isLiked,
		CreatedAt:    post.CreatedAt,
	}, nil
}

// parseHashtags extracts unique, lowercase hashtag names from the given
// content string using the pattern #(\w+).
func parseHashtags(content string) []string {
	matches := hashtagRegex.FindAllStringSubmatch(content, -1)
	seen := make(map[string]struct{})
	var tags []string
	for _, match := range matches {
		name := strings.ToLower(match[1])
		if _, ok := seen[name]; !ok {
			seen[name] = struct{}{}
			tags = append(tags, name)
		}
	}
	return tags
}
