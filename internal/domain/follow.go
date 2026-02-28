package domain

import "time"

type Follow struct {
	FollowerID  string    `json:"follower_id" gorm:"type:uuid;primaryKey"`
	Follower    User      `json:"-" gorm:"foreignKey:FollowerID"`
	FollowingID string    `json:"following_id" gorm:"type:uuid;primaryKey"`
	Following   User      `json:"-" gorm:"foreignKey:FollowingID"`
	CreatedAt   time.Time `json:"created_at"`
}
