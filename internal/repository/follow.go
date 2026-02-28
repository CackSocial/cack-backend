package repository

import "github.com/CackSocial/cack-backend/internal/domain"

type FollowRepository interface {
	Create(follow *domain.Follow) error
	Delete(followerID, followingID string) error
	IsFollowing(followerID, followingID string) (bool, error)
	GetFollowers(userID string, page, limit int) ([]domain.User, int64, error)
	GetFollowing(userID string, page, limit int) ([]domain.User, int64, error)
	GetFollowingIDs(userID string) ([]string, error)
}
