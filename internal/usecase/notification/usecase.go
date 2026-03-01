package notification

import (
	"encoding/json"
	"log"

	"github.com/CackSocial/cack-backend/internal/domain"
	"github.com/CackSocial/cack-backend/internal/dto"
	"github.com/CackSocial/cack-backend/internal/repository"
)

// RealtimePusher abstracts WebSocket push to avoid circular dependencies.
type RealtimePusher interface {
	SendToUser(userID string, data []byte)
}

// NotificationUseCase encapsulates all notification-related business logic.
type NotificationUseCase struct {
	notifRepo repository.NotificationRepository
	userRepo  repository.UserRepository
	pusher    RealtimePusher
}

// NewNotificationUseCase creates a new NotificationUseCase with the given dependencies.
func NewNotificationUseCase(notifRepo repository.NotificationRepository, userRepo repository.UserRepository, pusher RealtimePusher) *NotificationUseCase {
	return &NotificationUseCase{
		notifRepo: notifRepo,
		userRepo:  userRepo,
		pusher:    pusher,
	}
}

// GetNotifications returns a paginated list of notifications for the given user.
func (uc *NotificationUseCase) GetNotifications(userID string, page, limit int) ([]dto.NotificationResponse, error) {
	notifications, err := uc.notifRepo.GetByUserID(userID, page, limit)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.NotificationResponse, 0, len(notifications))
	for _, n := range notifications {
		responses = append(responses, toNotificationResponse(&n))
	}

	return responses, nil
}

// MarkAsRead marks a single notification as read for the given user.
func (uc *NotificationUseCase) MarkAsRead(id, userID string) error {
	return uc.notifRepo.MarkAsRead(id, userID)
}

// MarkAllAsRead marks all notifications as read for the given user.
func (uc *NotificationUseCase) MarkAllAsRead(userID string) error {
	return uc.notifRepo.MarkAllAsRead(userID)
}

// GetUnreadCount returns the number of unread notifications for the given user.
func (uc *NotificationUseCase) GetUnreadCount(userID string) (int64, error) {
	return uc.notifRepo.CountUnread(userID)
}

// CreateNotification creates a notification and pushes it in real-time if the recipient is online.
func (uc *NotificationUseCase) CreateNotification(userID, actorID, notifType, referenceID, referenceType string) error {
	notif := &domain.Notification{
		UserID:        userID,
		ActorID:       actorID,
		Type:          notifType,
		ReferenceID:   referenceID,
		ReferenceType: referenceType,
	}

	if err := uc.notifRepo.Create(notif); err != nil {
		return err
	}

	// Push real-time notification if pusher is available
	if uc.pusher != nil {
		actor, err := uc.userRepo.GetByID(actorID)
		if err == nil && actor != nil {
			notif.Actor = *actor
		}
		resp := toNotificationResponse(notif)
		payload := map[string]interface{}{
			"type": "notification",
			"data": resp,
		}
		data, err := json.Marshal(payload)
		if err == nil {
			uc.pusher.SendToUser(userID, data)
		} else {
			log.Printf("Failed to marshal notification: %v", err)
		}
	}

	return nil
}

func toNotificationResponse(n *domain.Notification) dto.NotificationResponse {
	return dto.NotificationResponse{
		ID: n.ID,
		Actor: dto.UserProfile{
			ID:          n.Actor.ID,
			Username:    n.Actor.Username,
			DisplayName: n.Actor.DisplayName,
			Bio:         n.Actor.Bio,
			AvatarURL:   n.Actor.AvatarURL,
		},
		Type:          n.Type,
		ReferenceID:   n.ReferenceID,
		ReferenceType: n.ReferenceType,
		IsRead:        n.IsRead,
		CreatedAt:     n.CreatedAt,
	}
}
