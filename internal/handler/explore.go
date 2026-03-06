package handler

import (
	"net/http"
	"strconv"

	"github.com/CackSocial/cack-backend/internal/usecase/explore"
	"github.com/CackSocial/cack-backend/pkg/response"
	"github.com/gin-gonic/gin"
)

type ExploreHandler struct {
	useCase *explore.ExploreUseCase
}

func NewExploreHandler(uc *explore.ExploreUseCase) *ExploreHandler {
	return &ExploreHandler{useCase: uc}
}

func (h *ExploreHandler) RegisterRoutes(protected *gin.RouterGroup) {
	protected.GET("/explore/suggested-users", h.GetSuggestedUsers)
	protected.GET("/explore/popular", h.GetPopularPosts)
	protected.GET("/explore/discover", h.GetDiscoverFeed)
}

// GetSuggestedUsers godoc
// @Summary Get suggested users to follow
// @Description Returns users the current user might want to follow, ranked by mutual followers with a fallback to popular users
// @Tags Explore
// @Produce json
// @Param limit query int false "Maximum number of suggestions (default 10)"
// @Success 200 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Security BearerAuth
// @Router /explore/suggested-users [get]
func (h *ExploreHandler) GetSuggestedUsers(c *gin.Context) {
	userID := getUserID(c)

	limit := 10
	if l, err := strconv.Atoi(c.Query("limit")); err == nil && l > 0 && l <= 50 {
		limit = l
	}

	users, err := h.useCase.GetSuggestedUsers(userID, limit)
	if err != nil {
		handleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, users)
}

// GetPopularPosts godoc
// @Summary Get popular posts
// @Description Returns high-engagement posts from outside the user's network, limited to the last 7 days
// @Tags Explore
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Results per page"
// @Success 200 {object} response.PaginatedResponse
// @Failure 401 {object} response.APIResponse
// @Security BearerAuth
// @Router /explore/popular [get]
func (h *ExploreHandler) GetPopularPosts(c *gin.Context) {
	userID := getUserID(c)
	page, limit := getPagination(c)

	posts, total, err := h.useCase.GetPopularPosts(userID, page, limit)
	if err != nil {
		handleError(c, err)
		return
	}

	response.Paginated(c, posts, page, limit, total)
}

// GetDiscoverFeed godoc
// @Summary Get discover feed
// @Description Returns posts matching tags from the user's liked posts, from users outside their network
// @Tags Explore
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Results per page"
// @Success 200 {object} response.PaginatedResponse
// @Failure 401 {object} response.APIResponse
// @Security BearerAuth
// @Router /explore/discover [get]
func (h *ExploreHandler) GetDiscoverFeed(c *gin.Context) {
	userID := getUserID(c)
	page, limit := getPagination(c)

	posts, total, err := h.useCase.GetDiscoverFeed(userID, page, limit)
	if err != nil {
		handleError(c, err)
		return
	}

	response.Paginated(c, posts, page, limit, total)
}
