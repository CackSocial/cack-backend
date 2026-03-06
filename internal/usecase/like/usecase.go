// Package like implements the business logic for liking and unliking posts,
// as well as listing users who liked a post.
package like

import (
	"github.com/CackSocial/cack-backend/internal/domain"
	"github.com/CackSocial/cack-backend/internal/dto"
	"github.com/CackSocial/cack-backend/internal/repository"
	ucerrors "github.com/CackSocial/cack-backend/internal/usecase/errors"
)

// NotificationCreator abstracts notification creation to avoid circular dependencies.
type NotificationCreator interface {
	CreateNotification(userID, actorID, notifType, referenceID, referenceType string) error
}

// LikeUseCase encapsulates all like-related business logic including
// liking, unliking, and retrieving the list of users who liked a post.
type LikeUseCase struct {
	likeRepo     repository.LikeRepository
	postRepo     repository.PostRepository
	userRepo     repository.UserRepository
	commentRepo  repository.CommentRepository
	bookmarkRepo repository.BookmarkRepository
	notifCase    NotificationCreator
}

// NewLikeUseCase creates a new LikeUseCase with the given dependencies.
func NewLikeUseCase(
	likeRepo repository.LikeRepository,
	postRepo repository.PostRepository,
	userRepo repository.UserRepository,
	commentRepo repository.CommentRepository,
	bookmarkRepo repository.BookmarkRepository,
	notifCase NotificationCreator,
) *LikeUseCase {
	return &LikeUseCase{
		likeRepo:     likeRepo,
		postRepo:     postRepo,
		userRepo:     userRepo,
		commentRepo:  commentRepo,
		bookmarkRepo: bookmarkRepo,
		notifCase:    notifCase,
	}
}

// Like adds a like from the given user to the specified post. It checks
// that the post exists and that the user has not already liked it.
func (uc *LikeUseCase) Like(userID, postID string) error {
	post, err := uc.postRepo.GetByID(postID)
	if err != nil || post == nil {
		return ucerrors.ErrPostNotFound
	}

	already, _ := uc.likeRepo.IsLiked(userID, postID)
	if already {
		return ucerrors.ErrAlreadyLiked
	}

	if err := uc.likeRepo.Create(&domain.Like{
		UserID: userID,
		PostID: postID,
	}); err != nil {
		return err
	}

	// Notify the post owner (don't notify if liking own post)
	if uc.notifCase != nil && post.UserID != userID {
		_ = uc.notifCase.CreateNotification(post.UserID, userID, "like", postID, "post")
	}

	return nil
}

// Unlike removes the like from the given user on the specified post.
func (uc *LikeUseCase) Unlike(userID, postID string) error {
	return uc.likeRepo.Delete(userID, postID)
}

// GetPostLikes returns a paginated list of users who liked the specified
// post, along with the total count.
func (uc *LikeUseCase) GetPostLikes(postID string, page, limit int) ([]dto.UserProfile, int64, error) {
	users, total, err := uc.likeRepo.GetByPostID(postID, page, limit)
	if err != nil {
		return nil, 0, err
	}

	profiles := make([]dto.UserProfile, 0, len(users))
	for _, u := range users {
		profiles = append(profiles, dto.UserProfile{
			ID:          u.ID,
			Username:    u.Username,
			DisplayName: u.DisplayName,
			Bio:         u.Bio,
			AvatarURL:   u.AvatarURL,
		})
	}

	return profiles, total, nil
}

// GetLikedPosts returns a paginated list of posts liked by the specified user.
// It enriches each post with like/comment/repost counts and per-viewer status flags.
func (uc *LikeUseCase) GetLikedPosts(username string, currentUserID string, page, limit int) ([]dto.PostResponse, int64, error) {
	user, err := uc.userRepo.GetByUsername(username)
	if err != nil || user == nil {
		return nil, 0, ucerrors.ErrUserNotFound
	}

	posts, total, err := uc.likeRepo.GetLikedPostsByUserID(user.ID, page, limit)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]dto.PostResponse, 0, len(posts))
	for i := range posts {
		post := &posts[i]

		likeCount, _ := uc.likeRepo.CountByPostID(post.ID)
		commentCount, _ := uc.commentRepo.CountByPostID(post.ID)
		repostCount, _ := uc.postRepo.CountReposts(post.ID)

		var isLiked, isBookmarked, isReposted bool
		if currentUserID != "" {
			isLiked, _ = uc.likeRepo.IsLiked(currentUserID, post.ID)
			isBookmarked, _ = uc.bookmarkRepo.IsBookmarked(currentUserID, post.ID)
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

		resp := dto.PostResponse{
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
			op := post.OriginalPost
			opLikeCount, _ := uc.likeRepo.CountByPostID(op.ID)
			opCommentCount, _ := uc.commentRepo.CountByPostID(op.ID)
			opRepostCount, _ := uc.postRepo.CountReposts(op.ID)
			opTagNames := make([]string, 0, len(op.Tags))
			for _, t := range op.Tags {
				opTagNames = append(opTagNames, t.Name)
			}
			opType := op.PostType
			if opType == "" {
				opType = "original"
			}
			var opIsLiked, opIsBookmarked, opIsReposted bool
			if currentUserID != "" {
				opIsLiked, _ = uc.likeRepo.IsLiked(currentUserID, op.ID)
				opIsBookmarked, _ = uc.bookmarkRepo.IsBookmarked(currentUserID, op.ID)
				opIsReposted, _ = uc.postRepo.IsReposted(currentUserID, op.ID)
			}
			resp.OriginalPost = &dto.PostResponse{
				ID:       op.ID,
				Content:  op.Content,
				ImageURL: op.ImageURL,
				Author: dto.UserProfile{
					ID:          op.User.ID,
					Username:    op.User.Username,
					DisplayName: op.User.DisplayName,
					Bio:         op.User.Bio,
					AvatarURL:   op.User.AvatarURL,
				},
				Tags:         opTagNames,
				PostType:     opType,
				RepostCount:  opRepostCount,
				IsReposted:   opIsReposted,
				LikeCount:    opLikeCount,
				CommentCount: opCommentCount,
				IsLiked:      opIsLiked,
				IsBookmarked: opIsBookmarked,
				CreatedAt:    op.CreatedAt,
			}
		}

		responses = append(responses, resp)
	}

	return responses, total, nil
}
