package repository

import (
	"github.com/CackSocial/cack-backend/internal/domain"
	"github.com/CackSocial/cack-backend/internal/repository"
	"gorm.io/gorm"
)

type followRepository struct {
	db *gorm.DB
}

func NewFollowRepository(db *gorm.DB) repository.FollowRepository {
	return &followRepository{db: db}
}

func (r *followRepository) Create(follow *domain.Follow) error {
	return r.db.Create(follow).Error
}

func (r *followRepository) Delete(followerID, followingID string) error {
	return r.db.Where("follower_id = ? AND following_id = ?", followerID, followingID).Delete(&domain.Follow{}).Error
}

func (r *followRepository) IsFollowing(followerID, followingID string) (bool, error) {
	var count int64
	if err := r.db.Model(&domain.Follow{}).Where("follower_id = ? AND following_id = ?", followerID, followingID).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *followRepository) GetFollowers(userID string, page, limit int) ([]domain.User, int64, error) {
	var users []domain.User
	var total int64

	q := r.db.Model(&domain.User{}).
		Joins("JOIN follows ON follows.follower_id = users.id").
		Where("follows.following_id = ?", userID)

	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := q.Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (r *followRepository) GetFollowing(userID string, page, limit int) ([]domain.User, int64, error) {
	var users []domain.User
	var total int64

	q := r.db.Model(&domain.User{}).
		Joins("JOIN follows ON follows.following_id = users.id").
		Where("follows.follower_id = ?", userID)

	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := q.Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (r *followRepository) GetFollowingIDs(userID string) ([]string, error) {
	var ids []string
	if err := r.db.Model(&domain.Follow{}).Where("follower_id = ?", userID).Pluck("following_id", &ids).Error; err != nil {
		return nil, err
	}
	return ids, nil
}
