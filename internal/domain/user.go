package domain

import "time"

type User struct {
	ID          string    `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Username    string    `json:"username" gorm:"uniqueIndex;size:50;not null"`
	Password    string    `json:"-" gorm:"not null"`
	DisplayName string    `json:"display_name" gorm:"size:100"`
	Bio         string    `json:"bio" gorm:"size:500"`
	AvatarURL   string    `json:"avatar_url"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
