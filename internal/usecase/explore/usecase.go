package explore

import (
	"time"

	"github.com/CackSocial/cack-backend/internal/domain"
	"github.com/CackSocial/cack-backend/internal/dto"
	"github.com/CackSocial/cack-backend/internal/repository"
)

type ExploreUseCase struct {
	userRepo     repository.UserRepository
	postRepo     repository.PostRepository
	followRepo   repository.FollowRepository
	likeRepo     repository.LikeRepository
	commentRepo  repository.CommentRepository
	bookmarkRepo repository.BookmarkRepository
}

func NewExploreUseCase(
	userRepo repository.UserRepository,
	postRepo repository.PostRepository,
	followRepo repository.FollowRepository,
	likeRepo repository.LikeRepository,
	commentRepo repository.CommentRepository,
	bookmarkRepo repository.BookmarkRepository,
) *ExploreUseCase {
	return &ExploreUseCase{
		userRepo:     userRepo,
		postRepo:     postRepo,
		followRepo:   followRepo,
		likeRepo:     likeRepo,
		commentRepo:  commentRepo,
		bookmarkRepo: bookmarkRepo,
	}
}

// GetSuggestedUsers returns users the current user might want to follow,
// ranked by mutual follower count with a fallback to popular users.
func (uc *ExploreUseCase) GetSuggestedUsers(currentUserID string, limit int) ([]dto.SuggestedUserResponse, error) {
	followingIDs, err := uc.followRepo.GetFollowingIDs(currentUserID)
	if err != nil {
		return nil, err
	}

	suggested, err := uc.userRepo.GetSuggestedUsers(currentUserID, followingIDs, limit)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.SuggestedUserResponse, 0, len(suggested))
	for _, s := range suggested {
		responses = append(responses, dto.SuggestedUserResponse{
			ID:                  s.User.ID,
			Username:            s.User.Username,
			DisplayName:         s.User.DisplayName,
			Bio:                 s.User.Bio,
			AvatarURL:           s.User.AvatarURL,
			MutualFollowerCount: s.MutualCount,
		})
	}

	return responses, nil
}

// GetPopularPosts returns high-engagement posts from outside the user's
// network, limited to the last 7 days.
func (uc *ExploreUseCase) GetPopularPosts(currentUserID string, page, limit int) ([]dto.PostResponse, int64, error) {
	followingIDs, err := uc.followRepo.GetFollowingIDs(currentUserID)
	if err != nil {
		return nil, 0, err
	}

	excludeIDs := append(followingIDs, currentUserID)
	since := time.Now().AddDate(0, 0, -7)

	posts, total, err := uc.postRepo.GetPopularPosts(excludeIDs, page, limit, since)
	if err != nil {
		return nil, 0, err
	}

	return uc.enrichPosts(posts, currentUserID), total, nil
}

// GetDiscoverFeed returns posts matching tags from the user's liked posts,
// from users outside their network.
func (uc *ExploreUseCase) GetDiscoverFeed(currentUserID string, page, limit int) ([]dto.PostResponse, int64, error) {
	tagNames, err := uc.likeRepo.GetLikedTagNames(currentUserID, 20)
	if err != nil {
		return nil, 0, err
	}

	if len(tagNames) == 0 {
		return []dto.PostResponse{}, 0, nil
	}

	followingIDs, err := uc.followRepo.GetFollowingIDs(currentUserID)
	if err != nil {
		return nil, 0, err
	}

	excludeIDs := append(followingIDs, currentUserID)

	posts, total, err := uc.postRepo.GetDiscoverPosts(tagNames, excludeIDs, page, limit)
	if err != nil {
		return nil, 0, err
	}

	return uc.enrichPosts(posts, currentUserID), total, nil
}

// enrichPosts adds like/comment/bookmark/repost counts and flags to raw posts.
func (uc *ExploreUseCase) enrichPosts(posts []domain.Post, currentUserID string) []dto.PostResponse {
	responses := make([]dto.PostResponse, 0, len(posts))

	for i := range posts {
		p := &posts[i]
		responses = append(responses, uc.buildPostResponse(p, currentUserID))
	}

	return responses
}

func (uc *ExploreUseCase) buildPostResponse(p *domain.Post, currentUserID string) dto.PostResponse {
	likeCount, _ := uc.likeRepo.CountByPostID(p.ID)
	commentCount, _ := uc.commentRepo.CountByPostID(p.ID)
	isLiked, _ := uc.likeRepo.IsLiked(currentUserID, p.ID)
	isBookmarked, _ := uc.bookmarkRepo.IsBookmarked(currentUserID, p.ID)
	repostCount, _ := uc.postRepo.CountReposts(p.ID)
	isReposted, _ := uc.postRepo.IsReposted(currentUserID, p.ID)

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
		opResp := uc.buildPostResponse(op, currentUserID)
		resp.OriginalPost = &opResp
	}

	return resp
}
