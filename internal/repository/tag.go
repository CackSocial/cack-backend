package repository

import (
	"time"

	"github.com/CackSocial/cack-backend/internal/domain"
)

type TrendingTag struct {
	Name      string `json:"name"`
	PostCount int64  `json:"post_count"`
}

type TagRepository interface {
	FindOrCreate(name string) (*domain.Tag, error)
	GetByPostID(postID string) ([]domain.Tag, error)
	GetTrending(limit int, since time.Time) ([]TrendingTag, error)
}
