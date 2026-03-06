package domain

// PostTag is the explicit join table for the many-to-many relationship
// between Post and Tag. Defining it as a concrete model prevents GORM's
// AutoMigrate from silently recreating the join table on each startup,
// which would wipe all post–tag associations.
type PostTag struct {
	PostID string `gorm:"primaryKey"`
	TagID  uint   `gorm:"primaryKey"`
}

func (PostTag) TableName() string {
	return "post_tags"
}
