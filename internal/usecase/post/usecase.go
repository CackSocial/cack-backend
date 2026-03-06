// Package post implements the business logic for creating, retrieving,
// and deleting posts, including image upload and hashtag extraction.
package post

import (
	"log"
	"regexp"
	"strings"

	"github.com/CackSocial/cack-backend/internal/domain"
	"github.com/CackSocial/cack-backend/internal/dto"
	"github.com/CackSocial/cack-backend/internal/infrastructure/storage"
	"github.com/CackSocial/cack-backend/internal/repository"
	ucerrors "github.com/CackSocial/cack-backend/internal/usecase/errors"
	"github.com/CackSocial/cack-backend/pkg/mentions"
)

var hashtagRegex = regexp.MustCompile(`#(\w+)`)

// NotificationCreator abstracts notification creation to avoid circular dependencies.
type NotificationCreator interface {
	CreateNotification(userID, actorID, notifType, referenceID, referenceType string) error
}

// PostUseCase encapsulates all post-related business logic including
// creation with image upload, hashtag parsing, retrieval, and deletion.
type PostUseCase struct {
	postRepo     repository.PostRepository
	tagRepo      repository.TagRepository
	likeRepo     repository.LikeRepository
	commentRepo  repository.CommentRepository
	userRepo     repository.UserRepository
	bookmarkRepo repository.BookmarkRepository
	storage      storage.Storage
	notifCase    NotificationCreator
}

// NewPostUseCase creates a new PostUseCase with the given dependencies.
func NewPostUseCase(
	postRepo repository.PostRepository,
	tagRepo repository.TagRepository,
	likeRepo repository.LikeRepository,
	commentRepo repository.CommentRepository,
	userRepo repository.UserRepository,
	bookmarkRepo repository.BookmarkRepository,
	storage storage.Storage,
	notifCase NotificationCreator,
) *PostUseCase {
	return &PostUseCase{
		postRepo:     postRepo,
		tagRepo:      tagRepo,
		likeRepo:     likeRepo,
		commentRepo:  commentRepo,
		userRepo:     userRepo,
		bookmarkRepo: bookmarkRepo,
		storage:      storage,
		notifCase:    notifCase,
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
			if imageURL != "" {
				_ = uc.storage.Delete(imageURL)
			}
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
		if imageURL != "" {
			_ = uc.storage.Delete(imageURL)
		}
		return nil, err
	}

	// Populate the User field for the response.
	post.User = *user

	// Send mention notifications
	if uc.notifCase != nil {
		for _, username := range mentions.ExtractMentions(req.Content) {
			mentioned, err := uc.userRepo.GetByUsername(username)
			if err != nil || mentioned == nil || mentioned.ID == userID {
				continue
			}
			if err := uc.notifCase.CreateNotification(mentioned.ID, userID, "mention", post.ID, "post"); err != nil {
				log.Printf("Failed to create mention notification for @%s: %v", username, err)
			}
		}
	}

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
// in like count, comment count, repost count, tag names, and the liked/reposted
// status for the current user.
func (uc *PostUseCase) toPostResponse(post *domain.Post, currentUserID string) (*dto.PostResponse, error) {
	likeCount, _ := uc.likeRepo.CountByPostID(post.ID)
	commentCount, _ := uc.commentRepo.CountByPostID(post.ID)
	repostCount, _ := uc.postRepo.CountReposts(post.ID)

	var isLiked bool
	if currentUserID != "" {
		isLiked, _ = uc.likeRepo.IsLiked(currentUserID, post.ID)
	}

	var isBookmarked bool
	if currentUserID != "" {
		isBookmarked, _ = uc.bookmarkRepo.IsBookmarked(currentUserID, post.ID)
	}

	var isReposted bool
	if currentUserID != "" {
		isReposted, _ = uc.postRepo.IsReposted(currentUserID, post.ID)
	}

	tagNames := make([]string, 0, len(post.Tags))
	for _, t := range post.Tags {
		tagNames = append(tagNames, t.Name)
	}

	postType := post.PostType
	if postType == "" {
		postType = "original"
	}

	resp := &dto.PostResponse{
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
		PostType:     postType,
		RepostCount:  repostCount,
		IsReposted:   isReposted,
		LikeCount:    likeCount,
		CommentCount: commentCount,
		IsLiked:      isLiked,
		IsBookmarked: isBookmarked,
		CreatedAt:    post.CreatedAt,
	}

	if post.OriginalPost != nil {
		origResp, err := uc.toPostResponse(post.OriginalPost, currentUserID)
		if err == nil {
			resp.OriginalPost = origResp
		}
	}

	return resp, nil
}

// Repost creates a repost of the given post for the user. Only original posts
// and quotes can be reposted.
func (uc *PostUseCase) Repost(userID, postID string) (*dto.PostResponse, error) {
	original, err := uc.postRepo.GetByID(postID)
	if err != nil || original == nil {
		return nil, ucerrors.ErrPostNotFound
	}

	if original.PostType == "repost" {
		return nil, ucerrors.ErrCannotRepost
	}

	already, _ := uc.postRepo.IsReposted(userID, postID)
	if already {
		return nil, ucerrors.ErrAlreadyReposted
	}

	user, err := uc.userRepo.GetByID(userID)
	if err != nil || user == nil {
		return nil, ucerrors.ErrUserNotFound
	}

	post := &domain.Post{
		UserID:         userID,
		PostType:       "repost",
		OriginalPostID: &postID,
	}

	if err := uc.postRepo.Create(post); err != nil {
		return nil, err
	}

	post.User = *user
	post.OriginalPost = original

	if uc.notifCase != nil && original.UserID != userID {
		_ = uc.notifCase.CreateNotification(original.UserID, userID, "repost", postID, "post")
	}

	return uc.toPostResponse(post, userID)
}

// DeleteRepost removes the user's repost of the given post.
func (uc *PostUseCase) DeleteRepost(userID, postID string) error {
	repost, err := uc.postRepo.GetRepostByUser(userID, postID)
	if err != nil || repost == nil {
		return ucerrors.ErrRepostNotFound
	}

	return uc.postRepo.Delete(repost.ID)
}

// QuotePost creates a quote post referencing the original. Content is required.
func (uc *PostUseCase) QuotePost(userID, postID string, req *dto.CreatePostRequest) (*dto.PostResponse, error) {
	original, err := uc.postRepo.GetByID(postID)
	if err != nil || original == nil {
		return nil, ucerrors.ErrPostNotFound
	}

	if original.PostType == "repost" {
		return nil, ucerrors.ErrCannotRepost
	}

	if strings.TrimSpace(req.Content) == "" {
		return nil, ucerrors.ErrContentRequired
	}

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

	tagNames := parseHashtags(req.Content)
	var tags []domain.Tag
	for _, name := range tagNames {
		tag, err := uc.tagRepo.FindOrCreate(name)
		if err != nil {
			if imageURL != "" {
				_ = uc.storage.Delete(imageURL)
			}
			return nil, err
		}
		tags = append(tags, *tag)
	}

	post := &domain.Post{
		UserID:         userID,
		Content:        req.Content,
		ImageURL:       imageURL,
		PostType:       "quote",
		OriginalPostID: &postID,
		Tags:           tags,
	}

	if err := uc.postRepo.Create(post); err != nil {
		if imageURL != "" {
			_ = uc.storage.Delete(imageURL)
		}
		return nil, err
	}

	post.User = *user
	post.OriginalPost = original

	if uc.notifCase != nil && original.UserID != userID {
		_ = uc.notifCase.CreateNotification(original.UserID, userID, "quote", postID, "post")
	}

	// Send mention notifications for quote content
	if uc.notifCase != nil {
		for _, username := range mentions.ExtractMentions(req.Content) {
			mentioned, err := uc.userRepo.GetByUsername(username)
			if err != nil || mentioned == nil || mentioned.ID == userID {
				continue
			}
			if err := uc.notifCase.CreateNotification(mentioned.ID, userID, "mention", post.ID, "post"); err != nil {
				log.Printf("Failed to create mention notification for @%s: %v", username, err)
			}
		}
	}

	return uc.toPostResponse(post, userID)
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
