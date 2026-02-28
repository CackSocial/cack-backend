package repository

import "github.com/CackSocial/cack-backend/internal/domain"

type UserRepository interface {
	Create(user *domain.User) error
	GetByID(id string) (*domain.User, error)
	GetByUsername(username string) (*domain.User, error)
	Update(user *domain.User) error
	Delete(id string) error
	Search(query string, page, limit int) ([]domain.User, int64, error)
}
