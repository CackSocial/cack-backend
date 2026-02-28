package repository

import (
	"github.com/CackSocial/cack-backend/internal/domain"
	"github.com/CackSocial/cack-backend/internal/repository"
	"gorm.io/gorm"
)

type bookmarkRepository struct {
	db *gorm.DB
}

func NewBookmarkRepository(db *gorm.DB) repository.BookmarkRepository {
	return &bookmarkRepository{db: db}
}

func (r *bookmarkRepository) Create(bookmark *domain.Bookmark) error {
	return r.db.Create(bookmark).Error
}

func (r *bookmarkRepository) Delete(userID, postID string) error {
	return r.db.Where("user_id = ? AND post_id = ?", userID, postID).Delete(&domain.Bookmark{}).Error
}

func (r *bookmarkRepository) GetByUserID(userID string, page, limit int) ([]domain.Bookmark, int64, error) {
	var bookmarks []domain.Bookmark
	var total int64

	q := r.db.Model(&domain.Bookmark{}).Where("user_id = ?", userID)
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := q.Order("created_at DESC").Offset(offset).Limit(limit).Find(&bookmarks).Error; err != nil {
		return nil, 0, err
	}

	return bookmarks, total, nil
}

func (r *bookmarkRepository) IsBookmarked(userID, postID string) (bool, error) {
	var count int64
	err := r.db.Model(&domain.Bookmark{}).Where("user_id = ? AND post_id = ?", userID, postID).Count(&count).Error
	return count > 0, err
}
