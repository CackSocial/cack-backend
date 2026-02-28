// Package bookmark implements the business logic for bookmarking posts.
package bookmark

import (
	"github.com/CackSocial/cack-backend/internal/domain"
	"github.com/CackSocial/cack-backend/internal/dto"
	"github.com/CackSocial/cack-backend/internal/repository"
	ucerrors "github.com/CackSocial/cack-backend/internal/usecase/errors"
)

// BookmarkUseCase encapsulates all bookmark-related business logic.
type BookmarkUseCase struct {
	bookmarkRepo repository.BookmarkRepository
	postRepo     repository.PostRepository
	likeRepo     repository.LikeRepository
	commentRepo  repository.CommentRepository
	userRepo     repository.UserRepository
}

// NewBookmarkUseCase creates a new BookmarkUseCase with the given dependencies.
func NewBookmarkUseCase(
	bookmarkRepo repository.BookmarkRepository,
	postRepo repository.PostRepository,
	likeRepo repository.LikeRepository,
	commentRepo repository.CommentRepository,
	userRepo repository.UserRepository,
) *BookmarkUseCase {
	return &BookmarkUseCase{
		bookmarkRepo: bookmarkRepo,
		postRepo:     postRepo,
		likeRepo:     likeRepo,
		commentRepo:  commentRepo,
		userRepo:     userRepo,
	}
}

// Bookmark adds a post to the user's bookmarks.
func (uc *BookmarkUseCase) Bookmark(userID, postID string) error {
	post, err := uc.postRepo.GetByID(postID)
	if err != nil || post == nil {
		return ucerrors.ErrPostNotFound
	}

	existing, _ := uc.bookmarkRepo.IsBookmarked(userID, postID)
	if existing {
		return ucerrors.ErrAlreadyBookmarked
	}

	return uc.bookmarkRepo.Create(&domain.Bookmark{
		UserID: userID,
		PostID: postID,
	})
}

// Unbookmark removes a post from the user's bookmarks.
func (uc *BookmarkUseCase) Unbookmark(userID, postID string) error {
	return uc.bookmarkRepo.Delete(userID, postID)
}

// GetBookmarks returns a paginated list of the user's bookmarked posts.
func (uc *BookmarkUseCase) GetBookmarks(userID string, page, limit int) ([]dto.PostResponse, int64, error) {
	bookmarks, total, err := uc.bookmarkRepo.GetByUserID(userID, page, limit)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]dto.PostResponse, 0, len(bookmarks))
	for _, b := range bookmarks {
		post, err := uc.postRepo.GetByID(b.PostID)
		if err != nil || post == nil {
			continue
		}

		likeCount, _ := uc.likeRepo.CountByPostID(post.ID)
		commentCount, _ := uc.commentRepo.CountByPostID(post.ID)
		isLiked, _ := uc.likeRepo.IsLiked(userID, post.ID)

		tagNames := make([]string, 0, len(post.Tags))
		for _, t := range post.Tags {
			tagNames = append(tagNames, t.Name)
		}

		responses = append(responses, dto.PostResponse{
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
			IsBookmarked: true,
			CreatedAt:    post.CreatedAt,
		})
	}

	return responses, total, nil
}
