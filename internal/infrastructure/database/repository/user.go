package repository

import (
	"fmt"

	"github.com/CackSocial/cack-backend/internal/domain"
	"github.com/CackSocial/cack-backend/internal/repository"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) repository.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *domain.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) GetByID(id string) (*domain.User, error) {
	var user domain.User
	if err := r.db.Where("id = ?", id).First(&user).Error; err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	return &user, nil
}

func (r *userRepository) GetByUsername(username string) (*domain.User, error) {
	var user domain.User
	if err := r.db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	return &user, nil
}

func (r *userRepository) Update(user *domain.User) error {
	return r.db.Save(user).Error
}

func (r *userRepository) Delete(id string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Delete messages where user is sender or receiver
		if err := tx.Where("sender_id = ? OR receiver_id = ?", id, id).Delete(&domain.Message{}).Error; err != nil {
			return err
		}
		// Delete bookmarks
		if err := tx.Where("user_id = ?", id).Delete(&domain.Bookmark{}).Error; err != nil {
			return err
		}
		// Delete likes
		if err := tx.Where("user_id = ?", id).Delete(&domain.Like{}).Error; err != nil {
			return err
		}
		// Delete comments
		if err := tx.Where("user_id = ?", id).Delete(&domain.Comment{}).Error; err != nil {
			return err
		}
		// Delete post_tags associations and posts
		var postIDs []string
		tx.Model(&domain.Post{}).Where("user_id = ?", id).Pluck("id", &postIDs)
		for _, pid := range postIDs {
			tx.Model(&domain.Post{ID: pid}).Association("Tags").Clear()
		}
		if err := tx.Where("user_id = ?", id).Delete(&domain.Post{}).Error; err != nil {
			return err
		}
		// Delete follows
		if err := tx.Where("follower_id = ? OR following_id = ?", id, id).Delete(&domain.Follow{}).Error; err != nil {
			return err
		}
		// Delete user
		return tx.Where("id = ?", id).Delete(&domain.User{}).Error
	})
}

func (r *userRepository) Search(query string, page, limit int) ([]domain.User, int64, error) {
	var users []domain.User
	var total int64

	search := "%" + query + "%"
	q := r.db.Model(&domain.User{}).Where("username ILIKE ? OR display_name ILIKE ?", search, search)

	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := q.Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}
