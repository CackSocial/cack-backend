package domain

import "time"

type Post struct {
	ID             string    `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID         string    `json:"user_id" gorm:"type:uuid;index;not null"`
	User           User      `json:"user" gorm:"foreignKey:UserID"`
	Content        string    `json:"content" gorm:"type:text"`
	ImageURL       string    `json:"image_url"`
	PostType       string    `json:"post_type" gorm:"type:varchar(20);default:'original'"`
	OriginalPostID *string   `json:"original_post_id" gorm:"type:uuid;index"`
	OriginalPost   *Post     `json:"original_post" gorm:"foreignKey:OriginalPostID"`
	Tags           []Tag     `json:"tags" gorm:"many2many:post_tags;"`
	CreatedAt      time.Time `json:"created_at" gorm:"index"`
	UpdatedAt      time.Time `json:"updated_at"`
}
