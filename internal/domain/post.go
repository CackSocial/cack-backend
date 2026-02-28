package domain

import "time"

type Post struct {
	ID        string    `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID    string    `json:"user_id" gorm:"type:uuid;index;not null"`
	User      User      `json:"user" gorm:"foreignKey:UserID"`
	Content   string    `json:"content" gorm:"type:text"`
	ImageURL  string    `json:"image_url"`
	Tags      []Tag     `json:"tags" gorm:"many2many:post_tags;"`
	CreatedAt time.Time `json:"created_at" gorm:"index"`
	UpdatedAt time.Time `json:"updated_at"`
}
