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
	tagRepo      repository.TagRepository
	postRepo     repository.PostRepository
	likeRepo     repository.LikeRepository
	commentRepo  repository.CommentRepository
	bookmarkRepo repository.BookmarkRepository
}

// NewTagUseCase creates a new TagUseCase with the given dependencies.
func NewTagUseCase(
	tagRepo repository.TagRepository,
	postRepo repository.PostRepository,
	likeRepo repository.LikeRepository,
	commentRepo repository.CommentRepository,
	bookmarkRepo repository.BookmarkRepository,
) *TagUseCase {
	return &TagUseCase{
		tagRepo:      tagRepo,
		postRepo:     postRepo,
		likeRepo:     likeRepo,
		commentRepo:  commentRepo,
		bookmarkRepo: bookmarkRepo,
	}
}

// GetTrending returns the most-used tags from the last 24 hours, limited
// to the specified count.
func (uc *TagUseCase) GetTrending(limit int) ([]repository.TrendingTag, error) {
	since := time.Now().Add(-24 * time.Hour)
	return uc.tagRepo.GetTrending(limit, since)
}

// GetPostsByTag returns a paginated list of posts associated with the given
// tag name. Each post includes like/comment/repost counts and the liked status
// for the current user.
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
		repostCount, _ := uc.postRepo.CountReposts(p.ID)

		var isLiked bool
		if currentUserID != "" {
			isLiked, _ = uc.likeRepo.IsLiked(currentUserID, p.ID)
		}

		var isBookmarked bool
		if currentUserID != "" {
			isBookmarked, _ = uc.bookmarkRepo.IsBookmarked(currentUserID, p.ID)
		}

		var isReposted bool
		if currentUserID != "" {
			isReposted, _ = uc.postRepo.IsReposted(currentUserID, p.ID)
		}

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
			opRepostCount, _ := uc.postRepo.CountReposts(op.ID)
			var opIsLiked bool
			if currentUserID != "" {
				opIsLiked, _ = uc.likeRepo.IsLiked(currentUserID, op.ID)
			}
			var opIsBookmarked bool
			if currentUserID != "" {
				opIsBookmarked, _ = uc.bookmarkRepo.IsBookmarked(currentUserID, op.ID)
			}
			var opIsReposted bool
			if currentUserID != "" {
				opIsReposted, _ = uc.postRepo.IsReposted(currentUserID, op.ID)
			}
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
