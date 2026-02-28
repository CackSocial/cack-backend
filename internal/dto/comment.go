package dto

import "time"

type CreateCommentRequest struct {
	Content string `json:"content" binding:"required,max=2000"`
}

type CommentResponse struct {
	ID        string      `json:"id"`
	Content   string      `json:"content"`
	Author    UserProfile `json:"author"`
	CreatedAt time.Time   `json:"created_at"`
}
