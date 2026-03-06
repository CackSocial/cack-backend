package repository

import (
	"time"

	"github.com/CackSocial/cack-backend/internal/domain"
)

type PostRepository interface {
	Create(post *domain.Post) error
	GetByID(id string) (*domain.Post, error)
	GetByUserID(userID string, page, limit int) ([]domain.Post, int64, error)
	Delete(id string) error
	GetFeed(userIDs []string, page, limit int) ([]domain.Post, int64, error)
	GetByTagName(tagName string, page, limit int) ([]domain.Post, int64, error)
	IsReposted(userID, postID string) (bool, error)
	CountReposts(postID string) (int64, error)
	GetRepostByUser(userID, postID string) (*domain.Post, error)
	GetPopularPosts(excludeUserIDs []string, page, limit int, since time.Time) ([]domain.Post, int64, error)
	GetDiscoverPosts(tagNames []string, excludeUserIDs []string, page, limit int) ([]domain.Post, int64, error)
}
