package repository

import (
	"time"

	"github.com/CackSocial/cack-backend/internal/domain"
	"github.com/CackSocial/cack-backend/internal/repository"
	"gorm.io/gorm"
)

type messageRepository struct {
	db *gorm.DB
}

func NewMessageRepository(db *gorm.DB) repository.MessageRepository {
	return &messageRepository{db: db}
}

func (r *messageRepository) Create(message *domain.Message) error {
	return r.db.Create(message).Error
}

func (r *messageRepository) GetConversation(userID1, userID2 string, page, limit int) ([]domain.Message, int64, error) {
	var messages []domain.Message
	var total int64

	q := r.db.Model(&domain.Message{}).
		Where("(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)",
			userID1, userID2, userID2, userID1)

	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := q.Preload("Sender").Order("created_at ASC").Offset(offset).Limit(limit).Find(&messages).Error; err != nil {
		return nil, 0, err
	}

	return messages, total, nil
}

func (r *messageRepository) GetConversations(userID string, page, limit int) ([]repository.ConversationPreview, int64, error) {
	// Subquery: get distinct conversation partner IDs
	type partnerRow struct {
		PartnerID string
	}
	var partners []partnerRow
	if err := r.db.Raw(`
		SELECT DISTINCT partner_id FROM (
			SELECT receiver_id AS partner_id FROM messages WHERE sender_id = ?
			UNION
			SELECT sender_id AS partner_id FROM messages WHERE receiver_id = ?
		) AS conversations
	`, userID, userID).Scan(&partners).Error; err != nil {
		return nil, 0, err
	}

	total := int64(len(partners))

	offset := (page - 1) * limit
	end := offset + limit
	if end > int(total) {
		end = int(total)
	}
	if offset >= int(total) {
		return []repository.ConversationPreview{}, total, nil
	}
	pagedPartners := partners[offset:end]

	var previews []repository.ConversationPreview
	for _, p := range pagedPartners {
		var preview repository.ConversationPreview

		// Get partner user info
		if err := r.db.Where("id = ?", p.PartnerID).First(&preview.User).Error; err != nil {
			continue
		}

		// Get last message
		if err := r.db.Preload("Sender").
			Where("(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)",
				userID, p.PartnerID, p.PartnerID, userID).
			Order("created_at DESC").
			First(&preview.LastMessage).Error; err != nil {
			continue
		}

		// Count unread messages from this partner
		r.db.Model(&domain.Message{}).
			Where("sender_id = ? AND receiver_id = ? AND read_at IS NULL", p.PartnerID, userID).
			Count(&preview.UnreadCount)

		previews = append(previews, preview)
	}

	return previews, total, nil
}

func (r *messageRepository) MarkAsRead(receiverID, senderID string) error {
	now := time.Now()
	return r.db.Model(&domain.Message{}).
		Where("receiver_id = ? AND sender_id = ? AND read_at IS NULL", receiverID, senderID).
		Update("read_at", now).Error
}
