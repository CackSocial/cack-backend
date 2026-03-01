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
	followRepo   repository.FollowRepository
	postRepo     repository.PostRepository
	likeRepo     repository.LikeRepository
	commentRepo  repository.CommentRepository
	bookmarkRepo repository.BookmarkRepository
}

// NewTimelineUseCase creates a new TimelineUseCase with the given dependencies.
func NewTimelineUseCase(
	followRepo repository.FollowRepository,
	postRepo repository.PostRepository,
	likeRepo repository.LikeRepository,
	commentRepo repository.CommentRepository,
	bookmarkRepo repository.BookmarkRepository,
) *TimelineUseCase {
	return &TimelineUseCase{
		followRepo:   followRepo,
		postRepo:     postRepo,
		likeRepo:     likeRepo,
		commentRepo:  commentRepo,
		bookmarkRepo: bookmarkRepo,
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
		isBookmarked, _ := uc.bookmarkRepo.IsBookmarked(userID, p.ID)
		repostCount, _ := uc.postRepo.CountReposts(p.ID)
		isReposted, _ := uc.postRepo.IsReposted(userID, p.ID)

		tagNames := make([]string, 0, len(p.Tags))
		for _, t := range p.Tags {
			tagNames = append(tagNames, t.Name)
		}

		postType := p.PostType
		if postType == "" {
			postType = "original"
		}

		resp := dto.PostResponse{
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
			PostType:     postType,
			RepostCount:  repostCount,
			IsReposted:   isReposted,
			LikeCount:    likeCount,
			CommentCount: commentCount,
			IsLiked:      isLiked,
			IsBookmarked: isBookmarked,
			CreatedAt:    p.CreatedAt,
		}

		if p.OriginalPost != nil {
			op := p.OriginalPost
			opLikeCount, _ := uc.likeRepo.CountByPostID(op.ID)
			opCommentCount, _ := uc.commentRepo.CountByPostID(op.ID)
			opIsLiked, _ := uc.likeRepo.IsLiked(userID, op.ID)
			opIsBookmarked, _ := uc.bookmarkRepo.IsBookmarked(userID, op.ID)
			opRepostCount, _ := uc.postRepo.CountReposts(op.ID)
			opIsReposted, _ := uc.postRepo.IsReposted(userID, op.ID)
			opTagNames := make([]string, 0, len(op.Tags))
			for _, t := range op.Tags {
				opTagNames = append(opTagNames, t.Name)
			}
			opType := op.PostType
			if opType == "" {
				opType = "original"
			}
			resp.OriginalPost = &dto.PostResponse{
				ID:       op.ID,
				Content:  op.Content,
				ImageURL: op.ImageURL,
				Author: dto.UserProfile{
					ID:          op.User.ID,
					Username:    op.User.Username,
					DisplayName: op.User.DisplayName,
					Bio:         op.User.Bio,
					AvatarURL:   op.User.AvatarURL,
				},
				Tags:         opTagNames,
				PostType:     opType,
				RepostCount:  opRepostCount,
				IsReposted:   opIsReposted,
				LikeCount:    opLikeCount,
				CommentCount: opCommentCount,
				IsLiked:      opIsLiked,
				IsBookmarked: opIsBookmarked,
				CreatedAt:    op.CreatedAt,
			}
		}

		responses = append(responses, resp)
	}

	return responses, total, nil
}
