// Package like implements the business logic for liking and unliking posts,
// as well as listing users who liked a post.
package like

import (
	"github.com/CackSocial/cack-backend/internal/domain"
	"github.com/CackSocial/cack-backend/internal/dto"
	"github.com/CackSocial/cack-backend/internal/repository"
	ucerrors "github.com/CackSocial/cack-backend/internal/usecase/errors"
)

// LikeUseCase encapsulates all like-related business logic including
// liking, unliking, and retrieving the list of users who liked a post.
type LikeUseCase struct {
	likeRepo repository.LikeRepository
	postRepo repository.PostRepository
	userRepo repository.UserRepository
}

// NewLikeUseCase creates a new LikeUseCase with the given dependencies.
func NewLikeUseCase(likeRepo repository.LikeRepository, postRepo repository.PostRepository, userRepo repository.UserRepository) *LikeUseCase {
	return &LikeUseCase{
		likeRepo: likeRepo,
		postRepo: postRepo,
		userRepo: userRepo,
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

	return uc.likeRepo.Create(&domain.Like{
		UserID: userID,
		PostID: postID,
	})
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
