// Package tag implements the business logic for retrieving trending tags
// and posts associated with a specific tag.
package tag

import (
	"time"

	"github.com/CackSocial/cack-backend/internal/dto"
	"github.com/CackSocial/cack-backend/internal/repository"
)

// TagUseCase encapsulates all tag-related business logic including
// trending tag retrieval and fetching posts by tag name.
type TagUseCase struct {
	tagRepo     repository.TagRepository
	postRepo    repository.PostRepository
	likeRepo    repository.LikeRepository
	commentRepo repository.CommentRepository
}

// NewTagUseCase creates a new TagUseCase with the given dependencies.
func NewTagUseCase(
	tagRepo repository.TagRepository,
	postRepo repository.PostRepository,
	likeRepo repository.LikeRepository,
	commentRepo repository.CommentRepository,
) *TagUseCase {
	return &TagUseCase{
		tagRepo:     tagRepo,
		postRepo:    postRepo,
		likeRepo:    likeRepo,
		commentRepo: commentRepo,
	}
}

// GetTrending returns the most-used tags from the last 24 hours, limited
// to the specified count.
func (uc *TagUseCase) GetTrending(limit int) ([]repository.TrendingTag, error) {
	since := time.Now().Add(-24 * time.Hour)
	return uc.tagRepo.GetTrending(limit, since)
}

// GetPostsByTag returns a paginated list of posts associated with the given
// tag name. Each post includes like/comment counts and the liked status for
// the current user.
func (uc *TagUseCase) GetPostsByTag(tagName string, currentUserID string, page, limit int) ([]dto.PostResponse, int64, error) {
	posts, total, err := uc.postRepo.GetByTagName(tagName, page, limit)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]dto.PostResponse, 0, len(posts))
	for i := range posts {
		p := &posts[i]

		likeCount, _ := uc.likeRepo.CountByPostID(p.ID)
		commentCount, _ := uc.commentRepo.CountByPostID(p.ID)

		var isLiked bool
		if currentUserID != "" {
			isLiked, _ = uc.likeRepo.IsLiked(currentUserID, p.ID)
		}

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
