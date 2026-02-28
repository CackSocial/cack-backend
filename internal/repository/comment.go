package repository

import "github.com/CackSocial/cack-backend/internal/domain"

type CommentRepository interface {
	Create(comment *domain.Comment) error
	GetByID(id string) (*domain.Comment, error)
	GetByPostID(postID string, page, limit int) ([]domain.Comment, int64, error)
	Delete(id string) error
	CountByPostID(postID string) (int64, error)
}
