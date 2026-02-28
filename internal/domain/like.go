package domain

import "time"

type Like struct {
	UserID    string    `json:"user_id" gorm:"type:uuid;primaryKey"`
	User      User      `json:"-" gorm:"foreignKey:UserID"`
	PostID    string    `json:"post_id" gorm:"type:uuid;primaryKey"`
	Post      Post      `json:"-" gorm:"foreignKey:PostID"`
	CreatedAt time.Time `json:"created_at"`
}
