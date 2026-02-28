package repository

import (
	"github.com/CackSocial/cack-backend/internal/domain"
	"github.com/CackSocial/cack-backend/internal/repository"
	"gorm.io/gorm"
)

type likeRepository struct {
	db *gorm.DB
}

func NewLikeRepository(db *gorm.DB) repository.LikeRepository {
	return &likeRepository{db: db}
}

func (r *likeRepository) Create(like *domain.Like) error {
	return r.db.Create(like).Error
}

func (r *likeRepository) Delete(userID, postID string) error {
	return r.db.Where("user_id = ? AND post_id = ?", userID, postID).Delete(&domain.Like{}).Error
}

func (r *likeRepository) GetByPostID(postID string, page, limit int) ([]domain.User, int64, error) {
	var users []domain.User
	var total int64

	q := r.db.Model(&domain.User{}).
		Joins("JOIN likes ON likes.user_id = users.id").
		Where("likes.post_id = ?", postID)

	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := q.Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (r *likeRepository) CountByPostID(postID string) (int64, error) {
	var count int64
	if err := r.db.Model(&domain.Like{}).Where("post_id = ?", postID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *likeRepository) IsLiked(userID, postID string) (bool, error) {
	var count int64
	if err := r.db.Model(&domain.Like{}).Where("user_id = ? AND post_id = ?", userID, postID).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
