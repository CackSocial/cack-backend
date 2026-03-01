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
		repostCount, _ := uc.postRepo.CountReposts(post.ID)
		isReposted, _ := uc.postRepo.IsReposted(userID, post.ID)

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
			IsBookmarked: true,
			CreatedAt:    post.CreatedAt,
		}

		if post.OriginalPost != nil {
			op := post.OriginalPost
			opLikeCount, _ := uc.likeRepo.CountByPostID(op.ID)
			opCommentCount, _ := uc.commentRepo.CountByPostID(op.ID)
			opRepostCount, _ := uc.postRepo.CountReposts(op.ID)
			opIsLiked, _ := uc.likeRepo.IsLiked(userID, op.ID)
			opIsBookmarked, _ := uc.bookmarkRepo.IsBookmarked(userID, op.ID)
			opIsReposted, _ := uc.postRepo.IsReposted(userID, op.ID)
			opTagNames := make([]string, 0, len(op.Tags))
			for _, t := range op.Tags {
				opTagNames = append(opTagNames, t.Name)
			}
			opType := op.PostType
			if opType == "" {
				opType = "original"
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
