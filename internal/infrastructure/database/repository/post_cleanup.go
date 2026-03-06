package repository

import (
	"errors"

	"github.com/CackSocial/cack-backend/internal/domain"
	"gorm.io/gorm"
)

// deletePostWithDependencies removes a post together with rows that still
// hold foreign keys to it. Reposts are deleted recursively, while quotes are
// detached so the quoting post can remain visible without the original.
func deletePostWithDependencies(tx *gorm.DB, postID string) error {
	var post domain.Post
	if err := tx.Select("id").First(&post, "id = ?", postID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	}

	repostIDs, err := getChildPostIDs(tx, postID, "repost")
	if err != nil {
		return err
	}
	for _, repostID := range repostIDs {
		if err := deletePostWithDependencies(tx, repostID); err != nil {
			return err
		}
	}

	if err := tx.Model(&domain.Post{}).
		Where("original_post_id = ? AND post_type = ?", postID, "quote").
		Update("original_post_id", nil).Error; err != nil {
		return err
	}

	if err := tx.Where("post_id = ?", postID).Delete(&domain.Like{}).Error; err != nil {
		return err
	}
	if err := tx.Where("post_id = ?", postID).Delete(&domain.Comment{}).Error; err != nil {
		return err
	}
	if err := tx.Where("post_id = ?", postID).Delete(&domain.Bookmark{}).Error; err != nil {
		return err
	}

	if err := tx.Model(&post).Association("Tags").Clear(); err != nil {
		return err
	}

	return tx.Where("id = ?", postID).Delete(&domain.Post{}).Error
}

func getChildPostIDs(tx *gorm.DB, originalPostID, postType string) ([]string, error) {
	var postIDs []string
	if err := tx.Model(&domain.Post{}).
		Where("original_post_id = ? AND post_type = ?", originalPostID, postType).
		Pluck("id", &postIDs).Error; err != nil {
		return nil, err
	}

	return postIDs, nil
}
