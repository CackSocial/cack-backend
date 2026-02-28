package repository

import (
	"fmt"

	"github.com/CackSocial/cack-backend/internal/domain"
	"github.com/CackSocial/cack-backend/internal/repository"
	"gorm.io/gorm"
)

type postRepository struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) repository.PostRepository {
	return &postRepository{db: db}
}

func (r *postRepository) Create(post *domain.Post) error {
	return r.db.Create(post).Error
}

func (r *postRepository) GetByID(id string) (*domain.Post, error) {
	var post domain.Post
	if err := r.db.Preload("User").Preload("Tags").Where("id = ?", id).First(&post).Error; err != nil {
		return nil, fmt.Errorf("post not found: %w", err)
	}
	return &post, nil
}

func (r *postRepository) GetByUserID(userID string, page, limit int) ([]domain.Post, int64, error) {
	var posts []domain.Post
	var total int64

	q := r.db.Model(&domain.Post{}).Where("user_id = ?", userID)

	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := q.Preload("User").Preload("Tags").Order("created_at DESC").Offset(offset).Limit(limit).Find(&posts).Error; err != nil {
		return nil, 0, err
	}

	return posts, total, nil
}

func (r *postRepository) Delete(id string) error {
	// Remove post_tags associations first
	if err := r.db.Model(&domain.Post{ID: id}).Association("Tags").Clear(); err != nil {
		return err
	}
	return r.db.Where("id = ?", id).Delete(&domain.Post{}).Error
}

func (r *postRepository) GetFeed(userIDs []string, page, limit int) ([]domain.Post, int64, error) {
	var posts []domain.Post
	var total int64

	q := r.db.Model(&domain.Post{}).Where("user_id IN ?", userIDs)

	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := q.Preload("User").Preload("Tags").Order("created_at DESC").Offset(offset).Limit(limit).Find(&posts).Error; err != nil {
		return nil, 0, err
	}

	return posts, total, nil
}

func (r *postRepository) GetByTagName(tagName string, page, limit int) ([]domain.Post, int64, error) {
	var posts []domain.Post
	var total int64

	q := r.db.Model(&domain.Post{}).
		Joins("JOIN post_tags ON post_tags.post_id = posts.id").
		Joins("JOIN tags ON tags.id = post_tags.tag_id").
		Where("tags.name = ?", tagName)

	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := q.Preload("User").Preload("Tags").Order("posts.created_at DESC").Offset(offset).Limit(limit).Find(&posts).Error; err != nil {
		return nil, 0, err
	}

	return posts, total, nil
}
