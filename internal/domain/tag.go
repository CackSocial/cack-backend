package domain

import "time"

type Tag struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"uniqueIndex;size:100;not null"`
	CreatedAt time.Time `json:"created_at"`
}
