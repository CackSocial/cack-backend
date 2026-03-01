package repository

import "github.com/CackSocial/cack-backend/internal/domain"

type NotificationRepository interface {
	Create(notification *domain.Notification) error
	GetByUserID(userID string, page, limit int) ([]domain.Notification, error)
	MarkAsRead(id string, userID string) error
	MarkAllAsRead(userID string) error
	CountUnread(userID string) (int64, error)
	Delete(id string) error
}
