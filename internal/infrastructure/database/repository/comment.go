package repository

import (
	"fmt"

	"github.com/CackSocial/cack-backend/internal/domain"
	"github.com/CackSocial/cack-backend/internal/repository"
	"gorm.io/gorm"
)

type commentRepository struct {
	db *gorm.DB
}

func NewCommentRepository(db *gorm.DB) repository.CommentRepository {
	return &commentRepository{db: db}
}

func (r *commentRepository) Create(comment *domain.Comment) error {
	return r.db.Create(comment).Error
}

func (r *commentRepository) GetByID(id string) (*domain.Comment, error) {
	var comment domain.Comment
	if err := r.db.Preload("User").Where("id = ?", id).First(&comment).Error; err != nil {
		return nil, fmt.Errorf("comment not found: %w", err)
	}
	return &comment, nil
}

func (r *commentRepository) GetByPostID(postID string, page, limit int) ([]domain.Comment, int64, error) {
	var comments []domain.Comment
	var total int64

	q := r.db.Model(&domain.Comment{}).Where("post_id = ?", postID)

	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := q.Preload("User").Order("created_at ASC").Offset(offset).Limit(limit).Find(&comments).Error; err != nil {
		return nil, 0, err
	}

	return comments, total, nil
}

func (r *commentRepository) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&domain.Comment{}).Error
}

func (r *commentRepository) CountByPostID(postID string) (int64, error) {
	var count int64
	if err := r.db.Model(&domain.Comment{}).Where("post_id = ?", postID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
