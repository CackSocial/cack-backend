// Package timeline implements the business logic for generating a user's
// home feed from the posts of followed users plus the user's own posts.
package timeline

import (
	"github.com/CackSocial/cack-backend/internal/dto"
	"github.com/CackSocial/cack-backend/internal/repository"
)

// TimelineUseCase encapsulates the business logic for building the
// authenticated user's chronological feed.
type TimelineUseCase struct {
	followRepo  repository.FollowRepository
	postRepo    repository.PostRepository
	likeRepo    repository.LikeRepository
	commentRepo repository.CommentRepository
}

// NewTimelineUseCase creates a new TimelineUseCase with the given dependencies.
func NewTimelineUseCase(
	followRepo repository.FollowRepository,
	postRepo repository.PostRepository,
	likeRepo repository.LikeRepository,
	commentRepo repository.CommentRepository,
) *TimelineUseCase {
	return &TimelineUseCase{
		followRepo:  followRepo,
		postRepo:    postRepo,
		likeRepo:    likeRepo,
		commentRepo: commentRepo,
	}
}

// GetFeed returns a paginated feed of posts from users the authenticated
// user follows, plus the user's own posts. Each post includes like/comment
// counts and the liked status for the current user.
func (uc *TimelineUseCase) GetFeed(userID string, page, limit int) ([]dto.PostResponse, int64, error) {
	// Get the IDs of users the current user follows.
	followingIDs, err := uc.followRepo.GetFollowingIDs(userID)
	if err != nil {
		return nil, 0, err
	}

	// Include the user's own posts in the feed.
	feedIDs := append(followingIDs, userID)

	posts, total, err := uc.postRepo.GetFeed(feedIDs, page, limit)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]dto.PostResponse, 0, len(posts))
	for i := range posts {
		p := &posts[i]

		likeCount, _ := uc.likeRepo.CountByPostID(p.ID)
		commentCount, _ := uc.commentRepo.CountByPostID(p.ID)
		isLiked, _ := uc.likeRepo.IsLiked(userID, p.ID)

		tagNames := make([]string, 0, len(p.Tags))
		for _, t := range p.Tags {
			tagNames = append(tagNames, t.Name)
		}

		responses = append(responses, dto.PostResponse{
			ID:       p.ID,
			Content:  p.Content,
			ImageURL: p.ImageURL,
			Author: dto.UserProfile{
				ID:          p.User.ID,
				Username:    p.User.Username,
				DisplayName: p.User.DisplayName,
				Bio:         p.User.Bio,
				AvatarURL:   p.User.AvatarURL,
			},
			Tags:         tagNames,
			LikeCount:    likeCount,
			CommentCount: commentCount,
			IsLiked:      isLiked,
			CreatedAt:    p.CreatedAt,
		})
	}

	return responses, total, nil
}
