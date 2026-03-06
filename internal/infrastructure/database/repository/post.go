package repository

import (
	"fmt"
	"time"

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
	if err := r.db.Preload("User").Preload("Tags").
		Preload("OriginalPost").Preload("OriginalPost.User").Preload("OriginalPost.Tags").
		Where("id = ?", id).First(&post).Error; err != nil {
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
	if err := q.Preload("User").Preload("Tags").
		Preload("OriginalPost").Preload("OriginalPost.User").Preload("OriginalPost.Tags").
		Order("created_at DESC").Offset(offset).Limit(limit).Find(&posts).Error; err != nil {
		return nil, 0, err
	}

	return posts, total, nil
}

func (r *postRepository) Delete(id string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		return deletePostWithDependencies(tx, id)
	})
}

func (r *postRepository) GetFeed(userIDs []string, page, limit int) ([]domain.Post, int64, error) {
	var posts []domain.Post
	var total int64

	q := r.db.Model(&domain.Post{}).Where("user_id IN ?", userIDs)

	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := q.Preload("User").Preload("Tags").
		Preload("OriginalPost").Preload("OriginalPost.User").Preload("OriginalPost.Tags").
		Order("created_at DESC").Offset(offset).Limit(limit).Find(&posts).Error; err != nil {
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
	if err := q.Preload("User").Preload("Tags").
		Preload("OriginalPost").Preload("OriginalPost.User").Preload("OriginalPost.Tags").
		Order("posts.created_at DESC").Offset(offset).Limit(limit).Find(&posts).Error; err != nil {
		return nil, 0, err
	}

	return posts, total, nil
}

func (r *postRepository) IsReposted(userID, postID string) (bool, error) {
	var count int64
	err := r.db.Model(&domain.Post{}).
		Where("user_id = ? AND original_post_id = ? AND post_type = 'repost'", userID, postID).
		Count(&count).Error
	return count > 0, err
}

func (r *postRepository) CountReposts(postID string) (int64, error) {
	var count int64
	err := r.db.Model(&domain.Post{}).
		Where("original_post_id = ? AND (post_type = 'repost' OR post_type = 'quote')", postID).
		Count(&count).Error
	return count, err
}

func (r *postRepository) GetRepostByUser(userID, postID string) (*domain.Post, error) {
	var post domain.Post
	if err := r.db.Where("user_id = ? AND original_post_id = ? AND post_type = 'repost'", userID, postID).
		First(&post).Error; err != nil {
		return nil, fmt.Errorf("repost not found: %w", err)
	}
	return &post, nil
}

func (r *postRepository) GetPopularPosts(excludeUserIDs []string, page, limit int, since time.Time) ([]domain.Post, int64, error) {
	var posts []domain.Post
	var total int64

	q := r.db.Model(&domain.Post{}).
		Where("post_type = 'original'").
		Where("created_at > ?", since)

	if len(excludeUserIDs) > 0 {
		q = q.Where("user_id NOT IN ?", excludeUserIDs)
	}

	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	fetchQ := r.db.
		Select("posts.*, "+
			"(SELECT COUNT(*) FROM likes WHERE likes.post_id = posts.id) + "+
			"(SELECT COUNT(*) FROM comments WHERE comments.post_id = posts.id) AS engagement").
		Table("posts").
		Where("posts.post_type = 'original'").
		Where("posts.created_at > ?", since)

	if len(excludeUserIDs) > 0 {
		fetchQ = fetchQ.Where("posts.user_id NOT IN ?", excludeUserIDs)
	}

	if err := fetchQ.
		Preload("User").Preload("Tags").
		Preload("OriginalPost").Preload("OriginalPost.User").Preload("OriginalPost.Tags").
		Order("engagement DESC, posts.created_at DESC").
		Offset(offset).Limit(limit).
		Find(&posts).Error; err != nil {
		return nil, 0, err
	}

	return posts, total, nil
}

func (r *postRepository) GetDiscoverPosts(tagNames []string, excludeUserIDs []string, page, limit int) ([]domain.Post, int64, error) {
	var posts []domain.Post
	var total int64

	if len(tagNames) == 0 {
		return posts, 0, nil
	}

	q := r.db.Model(&domain.Post{}).
		Joins("JOIN post_tags ON post_tags.post_id = posts.id").
		Joins("JOIN tags ON tags.id = post_tags.tag_id").
		Where("tags.name IN ?", tagNames).
		Where("posts.post_type = 'original'")

	if len(excludeUserIDs) > 0 {
		q = q.Where("posts.user_id NOT IN ?", excludeUserIDs)
	}

	if err := q.Distinct("posts.id").Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	fetchQ := r.db.
		Joins("JOIN post_tags ON post_tags.post_id = posts.id").
		Joins("JOIN tags ON tags.id = post_tags.tag_id").
		Where("tags.name IN ?", tagNames).
		Where("posts.post_type = 'original'")

	if len(excludeUserIDs) > 0 {
		fetchQ = fetchQ.Where("posts.user_id NOT IN ?", excludeUserIDs)
	}

	if err := fetchQ.
		Preload("User").Preload("Tags").
		Preload("OriginalPost").Preload("OriginalPost.User").Preload("OriginalPost.Tags").
		Group("posts.id").
		Order("posts.created_at DESC").
		Offset(offset).Limit(limit).
		Find(&posts).Error; err != nil {
		return nil, 0, err
	}

	return posts, total, nil
}
