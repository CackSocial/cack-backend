package repository

import "github.com/CackSocial/cack-backend/internal/domain"

type ConversationPreview struct {
	User        domain.User    `json:"user"`
	LastMessage domain.Message `json:"last_message"`
	UnreadCount int64          `json:"unread_count"`
}

type MessageRepository interface {
	Create(message *domain.Message) error
	GetConversation(userID1, userID2 string, page, limit int) ([]domain.Message, int64, error)
	GetConversations(userID string, page, limit int) ([]ConversationPreview, int64, error)
	MarkAsRead(receiverID, senderID string) error
}
