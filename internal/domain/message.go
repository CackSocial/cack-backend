package domain

import "time"

type Message struct {
	ID         string     `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	SenderID   string     `json:"sender_id" gorm:"type:uuid;not null;index:idx_messages_conversation"`
	Sender     User       `json:"sender" gorm:"foreignKey:SenderID"`
	ReceiverID string     `json:"receiver_id" gorm:"type:uuid;not null;index:idx_messages_conversation"`
	Receiver   User       `json:"-" gorm:"foreignKey:ReceiverID"`
	Content    string     `json:"content" gorm:"type:text"`
	ImageURL   string     `json:"image_url"`
	ReadAt     *time.Time `json:"read_at"`
	CreatedAt  time.Time  `json:"created_at" gorm:"index:idx_messages_conversation"`
}
