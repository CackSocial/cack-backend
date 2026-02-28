package repository

import "github.com/CackSocial/cack-backend/internal/domain"

type BookmarkRepository interface {
	Create(bookmark *domain.Bookmark) error
	Delete(userID, postID string) error
	GetByUserID(userID string, page, limit int) ([]domain.Bookmark, int64, error)
	IsBookmarked(userID, postID string) (bool, error)
}
