package repository

import (
	"github.com/CackSocial/cack-backend/internal/domain"
	"github.com/CackSocial/cack-backend/internal/repository"
	"gorm.io/gorm"
)

type notificationRepository struct {
	db *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) repository.NotificationRepository {
	return &notificationRepository{db: db}
}

func (r *notificationRepository) Create(notification *domain.Notification) error {
	return r.db.Create(notification).Error
}

func (r *notificationRepository) GetByUserID(userID string, page, limit int) ([]domain.Notification, error) {
	var notifications []domain.Notification
	offset := (page - 1) * limit
	err := r.db.Preload("Actor").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Offset(offset).Limit(limit).
		Find(&notifications).Error
	return notifications, err
}

func (r *notificationRepository) MarkAsRead(id string, userID string) error {
	return r.db.Model(&domain.Notification{}).
		Where("id = ? AND user_id = ?", id, userID).
		Update("is_read", true).Error
}

func (r *notificationRepository) MarkAllAsRead(userID string) error {
	return r.db.Model(&domain.Notification{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Update("is_read", true).Error
}

func (r *notificationRepository) CountUnread(userID string) (int64, error) {
	var count int64
	err := r.db.Model(&domain.Notification{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Count(&count).Error
	return count, err
}

func (r *notificationRepository) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&domain.Notification{}).Error
}
