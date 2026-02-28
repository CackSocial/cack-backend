package dto

import (
	"mime/multipart"
	"time"
)

type SendMessageRequest struct {
	Content string                `form:"content" binding:"required_without=Image,max=5000"`
	Image   *multipart.FileHeader `form:"image"`
}

type MessageResponse struct {
	ID         string     `json:"id"`
	SenderID   string     `json:"sender_id"`
	ReceiverID string     `json:"receiver_id"`
	Content    string     `json:"content"`
	ImageURL   string     `json:"image_url,omitempty"`
	ReadAt     *time.Time `json:"read_at"`
	CreatedAt  time.Time  `json:"created_at"`
}

type ConversationListResponse struct {
	User        UserProfile     `json:"user"`
	LastMessage MessageResponse `json:"last_message"`
	UnreadCount int64           `json:"unread_count"`
}

type WSMessage struct {
	Type       string `json:"type"`
	ReceiverID string `json:"receiver_id,omitempty"`
	Content    string `json:"content,omitempty"`
	ImageURL   string `json:"image_url,omitempty"`
}
