package domain

import "time"

type Comment struct {
	ID        string    `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	PostID    string    `json:"post_id" gorm:"type:uuid;index;not null"`
	Post      Post      `json:"-" gorm:"foreignKey:PostID"`
	UserID    string    `json:"user_id" gorm:"type:uuid;not null"`
	User      User      `json:"user" gorm:"foreignKey:UserID"`
	Content   string    `json:"content" gorm:"type:text;not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
