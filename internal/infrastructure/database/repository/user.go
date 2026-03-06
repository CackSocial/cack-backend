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
		// Delete posts with all dependent rows that reference them.
		var postIDs []string
		if err := tx.Model(&domain.Post{}).Where("user_id = ?", id).Pluck("id", &postIDs).Error; err != nil {
			return err
		}
		for _, pid := range postIDs {
			if err := deletePostWithDependencies(tx, pid); err != nil {
				return err
			}
		}
		// Delete follows
		if err := tx.Where("follower_id = ? OR following_id = ?", id, id).Delete(&domain.Follow{}).Error; err != nil {
			return err
		}
		// Delete user
		return tx.Where("id = ?", id).Delete(&domain.User{}).Error
	})
}

func (r *userRepository) GetSuggestedUsers(currentUserID string, followingIDs []string, limit int) ([]repository.SuggestedUser, error) {
	type result struct {
		domain.User
		MutualCount int64
	}
	var results []result

	if len(followingIDs) > 0 {
		// Users followed by people I follow, but not me and not already followed
		excludeIDs := append(followingIDs, currentUserID)
		err := r.db.Raw(`
			SELECT u.*, COUNT(DISTINCT f_mutual.follower_id) as mutual_count
			FROM follows f_mutual
			JOIN users u ON f_mutual.following_id = u.id
			WHERE f_mutual.follower_id IN ?
			  AND f_mutual.following_id NOT IN ?
			GROUP BY u.id
			ORDER BY mutual_count DESC
			LIMIT ?
		`, followingIDs, excludeIDs, limit).Scan(&results).Error
		if err != nil {
			return nil, err
		}
	}

	// Fallback: popular users by follower count if not enough mutual-based suggestions
	if len(results) < limit {
		remaining := limit - len(results)
		existingIDs := make([]string, 0, len(results))
		for _, r := range results {
			existingIDs = append(existingIDs, r.ID)
		}
		excludeIDs := append(followingIDs, currentUserID)
		excludeIDs = append(excludeIDs, existingIDs...)

		var popular []result
		q := r.db.Raw(`
			SELECT u.*, COUNT(f.follower_id) as mutual_count
			FROM users u
			LEFT JOIN follows f ON f.following_id = u.id
			WHERE u.id NOT IN ?
			GROUP BY u.id
			ORDER BY mutual_count DESC
			LIMIT ?
		`, excludeIDs, remaining)
		if err := q.Scan(&popular).Error; err != nil {
			return nil, err
		}
		results = append(results, popular...)
	}

	suggested := make([]repository.SuggestedUser, 0, len(results))
	for _, r := range results {
		suggested = append(suggested, repository.SuggestedUser{
			User:        r.User,
			MutualCount: r.MutualCount,
		})
	}
	return suggested, nil
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
