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

func (r *likeRepository) GetLikedPostsByUserID(userID string, page, limit int) ([]domain.Post, int64, error) {
	var posts []domain.Post
	var total int64

	q := r.db.Model(&domain.Post{}).
		Joins("JOIN likes ON likes.post_id = posts.id").
		Where("likes.user_id = ?", userID)

	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := r.db.
		Joins("JOIN likes ON likes.post_id = posts.id").
		Where("likes.user_id = ?", userID).
		Preload("User").Preload("Tags").
		Preload("OriginalPost").Preload("OriginalPost.User").Preload("OriginalPost.Tags").
		Order("likes.created_at DESC").
		Offset(offset).Limit(limit).
		Find(&posts).Error; err != nil {
		return nil, 0, err
	}

	return posts, total, nil
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

func (r *likeRepository) GetLikedTagNames(userID string, limit int) ([]string, error) {
	var tagNames []string
	err := r.db.Raw(`
		SELECT DISTINCT t.name
		FROM tags t
		JOIN post_tags pt ON t.id = pt.tag_id
		JOIN likes l ON l.post_id = pt.post_id
		WHERE l.user_id = ?
		LIMIT ?
	`, userID, limit).Scan(&tagNames).Error
	if err != nil {
		return nil, err
	}
	return tagNames, nil
}
