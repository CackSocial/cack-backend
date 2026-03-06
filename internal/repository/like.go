package repository

import "github.com/CackSocial/cack-backend/internal/domain"

type LikeRepository interface {
	Create(like *domain.Like) error
	Delete(userID, postID string) error
	GetByPostID(postID string, page, limit int) ([]domain.User, int64, error)
	GetLikedPostsByUserID(userID string, page, limit int) ([]domain.Post, int64, error)
	CountByPostID(postID string) (int64, error)
	IsLiked(userID, postID string) (bool, error)
	GetLikedTagNames(userID string, limit int) ([]string, error)
}
