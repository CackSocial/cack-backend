package dto

import "time"

type NotificationResponse struct {
	ID            string      `json:"id"`
	Actor         UserProfile `json:"actor"`
	Type          string      `json:"type"`
	ReferenceID   string      `json:"reference_id"`
	ReferenceType string      `json:"reference_type"`
	IsRead        bool        `json:"is_read"`
	CreatedAt     time.Time   `json:"created_at"`
}

type UnreadCountResponse struct {
	Count int64 `json:"count"`
}
