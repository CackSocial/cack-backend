package domain

import "time"

type Notification struct {
	ID            string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID        string    `gorm:"type:uuid;not null;index"` // recipient
	ActorID       string    `gorm:"type:uuid;not null"`       // who triggered it
	Actor         User      `gorm:"foreignKey:ActorID"`
	Type          string    `gorm:"type:varchar(50);not null;index"` // "like", "follow", "comment", "mention", "repost"
	ReferenceID   string    `gorm:"type:uuid"`                       // post ID, comment ID, etc.
	ReferenceType string    `gorm:"type:varchar(50)"`                // "post", "comment", "user"
	IsRead        bool      `gorm:"default:false;index"`
	CreatedAt     time.Time
}
