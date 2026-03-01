package dto

import (
	"mime/multipart"
	"time"
)

type CreatePostRequest struct {
	Content string                `form:"content" binding:"required_without=Image,max=5000"`
	Image   *multipart.FileHeader `form:"image"`
}

type PostResponse struct {
	ID           string        `json:"id"`
	Content      string        `json:"content"`
	ImageURL     string        `json:"image_url,omitempty"`
	Author       UserProfile   `json:"author"`
	Tags         []string      `json:"tags"`
	PostType     string        `json:"post_type"`
	OriginalPost *PostResponse `json:"original_post,omitempty"`
	RepostCount  int64         `json:"repost_count"`
	IsReposted   bool          `json:"is_reposted"`
	LikeCount    int64         `json:"like_count"`
	CommentCount int64         `json:"comment_count"`
	IsLiked      bool          `json:"is_liked"`
	IsBookmarked bool          `json:"is_bookmarked"`
	CreatedAt    time.Time     `json:"created_at"`
}
