// Package follow implements the business logic for following and
// unfollowing users, as well as listing followers and following.
package follow

import (
	"github.com/CackSocial/cack-backend/internal/domain"
	"github.com/CackSocial/cack-backend/internal/dto"
	"github.com/CackSocial/cack-backend/internal/repository"
	ucerrors "github.com/CackSocial/cack-backend/internal/usecase/errors"
)

// FollowUseCase encapsulates all follow-related business logic including
// following, unfollowing, and retrieving follower/following lists.
type FollowUseCase struct {
	followRepo repository.FollowRepository
	userRepo   repository.UserRepository
}

// NewFollowUseCase creates a new FollowUseCase with the given dependencies.
func NewFollowUseCase(followRepo repository.FollowRepository, userRepo repository.UserRepository) *FollowUseCase {
	return &FollowUseCase{
		followRepo: followRepo,
		userRepo:   userRepo,
	}
}

// Follow makes the authenticated user (followerID) follow the user
// identified by username. Self-follows and duplicate follows are rejected.
func (uc *FollowUseCase) Follow(followerID string, username string) error {
	target, err := uc.userRepo.GetByUsername(username)
	if err != nil || target == nil {
		return ucerrors.ErrUserNotFound
	}

	if followerID == target.ID {
		return ucerrors.ErrSelfFollow
	}

	already, _ := uc.followRepo.IsFollowing(followerID, target.ID)
	if already {
		return ucerrors.ErrAlreadyFollowing
	}

	return uc.followRepo.Create(&domain.Follow{
		FollowerID:  followerID,
		FollowingID: target.ID,
	})
}

// Unfollow removes the follow relationship between the authenticated user
// (followerID) and the user identified by username.
func (uc *FollowUseCase) Unfollow(followerID string, username string) error {
	target, err := uc.userRepo.GetByUsername(username)
	if err != nil || target == nil {
		return ucerrors.ErrUserNotFound
	}

	return uc.followRepo.Delete(followerID, target.ID)
}

// GetFollowers returns a paginated list of users who follow the user
// identified by username.
func (uc *FollowUseCase) GetFollowers(username string, page, limit int) ([]dto.UserProfile, int64, error) {
	user, err := uc.userRepo.GetByUsername(username)
	if err != nil || user == nil {
		return nil, 0, ucerrors.ErrUserNotFound
	}

	followers, total, err := uc.followRepo.GetFollowers(user.ID, page, limit)
	if err != nil {
		return nil, 0, err
	}

	return toUserProfiles(followers), total, nil
}

// GetFollowing returns a paginated list of users that the user identified
// by username is following.
func (uc *FollowUseCase) GetFollowing(username string, page, limit int) ([]dto.UserProfile, int64, error) {
	user, err := uc.userRepo.GetByUsername(username)
	if err != nil || user == nil {
		return nil, 0, ucerrors.ErrUserNotFound
	}

	following, total, err := uc.followRepo.GetFollowing(user.ID, page, limit)
	if err != nil {
		return nil, 0, err
	}

	return toUserProfiles(following), total, nil
}

// toUserProfiles converts a slice of domain Users to minimal UserProfile DTOs.
func toUserProfiles(users []domain.User) []dto.UserProfile {
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
	return profiles
}
