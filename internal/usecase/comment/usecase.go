// Package comment implements the business logic for creating, listing,
// and deleting comments on posts.
package comment

import (
	"log"

	"github.com/CackSocial/cack-backend/internal/domain"
	"github.com/CackSocial/cack-backend/internal/dto"
	"github.com/CackSocial/cack-backend/internal/repository"
	ucerrors "github.com/CackSocial/cack-backend/internal/usecase/errors"
	"github.com/CackSocial/cack-backend/pkg/mentions"
)

// NotificationCreator abstracts notification creation to avoid circular dependencies.
type NotificationCreator interface {
	CreateNotification(userID, actorID, notifType, referenceID, referenceType string) error
}

// CommentUseCase encapsulates all comment-related business logic including
// creating comments, listing comments by post, and deleting comments.
type CommentUseCase struct {
	commentRepo repository.CommentRepository
	postRepo    repository.PostRepository
	userRepo    repository.UserRepository
	notifCase   NotificationCreator
}

// NewCommentUseCase creates a new CommentUseCase with the given dependencies.
func NewCommentUseCase(commentRepo repository.CommentRepository, postRepo repository.PostRepository, userRepo repository.UserRepository, notifCase NotificationCreator) *CommentUseCase {
	return &CommentUseCase{
		commentRepo: commentRepo,
		postRepo:    postRepo,
		userRepo:    userRepo,
		notifCase:   notifCase,
	}
}

// Create adds a new comment to the specified post. It validates that the
// post exists, persists the comment, and returns the comment with author
// information.
func (uc *CommentUseCase) Create(userID, postID string, req *dto.CreateCommentRequest) (*dto.CommentResponse, error) {
	post, err := uc.postRepo.GetByID(postID)
	if err != nil || post == nil {
		return nil, ucerrors.ErrPostNotFound
	}

	comment := &domain.Comment{
		PostID:  postID,
		UserID:  userID,
		Content: req.Content,
	}

	if err := uc.commentRepo.Create(comment); err != nil {
		return nil, err
	}

	// Fetch the created comment with the user preloaded.
	created, err := uc.commentRepo.GetByID(comment.ID)
	if err != nil {
		return nil, err
	}

	// Notify the post owner (don't notify if commenting on own post)
	if uc.notifCase != nil && post.UserID != userID {
		_ = uc.notifCase.CreateNotification(post.UserID, userID, "comment", postID, "post")
	}

	// Send mention notifications
	if uc.notifCase != nil {
		for _, username := range mentions.ExtractMentions(req.Content) {
			mentioned, err := uc.userRepo.GetByUsername(username)
			if err != nil || mentioned == nil || mentioned.ID == userID {
				continue
			}
			if err := uc.notifCase.CreateNotification(mentioned.ID, userID, "mention", postID, "post"); err != nil {
				log.Printf("Failed to create mention notification for @%s: %v", username, err)
			}
		}
	}

	return toCommentResponse(created), nil
}

// GetByPostID returns a paginated list of comments for the specified post.
func (uc *CommentUseCase) GetByPostID(postID string, page, limit int) ([]dto.CommentResponse, int64, error) {
	comments, total, err := uc.commentRepo.GetByPostID(postID, page, limit)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]dto.CommentResponse, 0, len(comments))
	for i := range comments {
		responses = append(responses, *toCommentResponse(&comments[i]))
	}

	return responses, total, nil
}

// Delete removes a comment by its ID. Only the comment owner is allowed
// to delete their comment.
func (uc *CommentUseCase) Delete(commentID, userID string) error {
	comment, err := uc.commentRepo.GetByID(commentID)
	if err != nil || comment == nil {
		return ucerrors.ErrCommentNotFound
	}

	if comment.UserID != userID {
		return ucerrors.ErrUnauthorized
	}

	return uc.commentRepo.Delete(commentID)
}

// toCommentResponse converts a domain Comment into a CommentResponse DTO.
func toCommentResponse(c *domain.Comment) *dto.CommentResponse {
	return &dto.CommentResponse{
		ID:      c.ID,
		Content: c.Content,
		Author: dto.UserProfile{
			ID:          c.User.ID,
			Username:    c.User.Username,
			DisplayName: c.User.DisplayName,
			Bio:         c.User.Bio,
			AvatarURL:   c.User.AvatarURL,
		},
		CreatedAt: c.CreatedAt,
	}
}
