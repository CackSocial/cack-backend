package repository

import (
	"strings"
	"time"

	"github.com/CackSocial/cack-backend/internal/domain"
	"github.com/CackSocial/cack-backend/internal/repository"
	"gorm.io/gorm"
)

type tagRepository struct {
	db *gorm.DB
}

func NewTagRepository(db *gorm.DB) repository.TagRepository {
	return &tagRepository{db: db}
}

func (r *tagRepository) FindOrCreate(name string) (*domain.Tag, error) {
	var tag domain.Tag
	tag.Name = strings.ToLower(name)
	if err := r.db.Where("name = ?", tag.Name).FirstOrCreate(&tag).Error; err != nil {
		return nil, err
	}
	return &tag, nil
}

func (r *tagRepository) GetByPostID(postID string) ([]domain.Tag, error) {
	var tags []domain.Tag
	if err := r.db.
		Joins("JOIN post_tags ON post_tags.tag_id = tags.id").
		Where("post_tags.post_id = ?", postID).
		Find(&tags).Error; err != nil {
		return nil, err
	}
	return tags, nil
}

func (r *tagRepository) GetTrending(limit int, since time.Time) ([]repository.TrendingTag, error) {
	var trending []repository.TrendingTag
	if err := r.db.Raw(`
		SELECT t.name, COUNT(pt.post_id) AS post_count
		FROM tags t
		JOIN post_tags pt ON pt.tag_id = t.id
		JOIN posts p ON p.id = pt.post_id
		WHERE p.created_at >= ?
		GROUP BY t.name
		ORDER BY post_count DESC
		LIMIT ?
	`, since, limit).Scan(&trending).Error; err != nil {
		return nil, err
	}
	return trending, nil
}
