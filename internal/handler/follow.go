package handler

import (
	"net/http"

	"github.com/CackSocial/cack-backend/internal/usecase/follow"
	"github.com/CackSocial/cack-backend/pkg/response"
	"github.com/gin-gonic/gin"
)

type FollowHandler struct {
	followUseCase *follow.FollowUseCase
}

func NewFollowHandler(uc *follow.FollowUseCase) *FollowHandler {
	return &FollowHandler{followUseCase: uc}
}

func (h *FollowHandler) RegisterRoutes(public, protected *gin.RouterGroup) {
	protected.POST("/users/:username/follow", h.Follow)
	protected.DELETE("/users/:username/follow", h.Unfollow)
	public.GET("/users/:username/followers", h.GetFollowers)
	public.GET("/users/:username/following", h.GetFollowing)
}

// Follow godoc
// @Summary Follow a user
// @Description Follow another user by username
// @Tags Follows
// @Produce json
// @Param username path string true "Username to follow"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Security BearerAuth
// @Router /users/{username}/follow [post]
func (h *FollowHandler) Follow(c *gin.Context) {
	userID := getUserID(c)
	username := c.Param("username")

	if err := h.followUseCase.Follow(userID, username); err != nil {
		handleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, gin.H{"message": "followed"})
}

// Unfollow godoc
// @Summary Unfollow a user
// @Description Unfollow a user by username
// @Tags Follows
// @Produce json
// @Param username path string true "Username to unfollow"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Security BearerAuth
// @Router /users/{username}/follow [delete]
func (h *FollowHandler) Unfollow(c *gin.Context) {
	userID := getUserID(c)
	username := c.Param("username")

	if err := h.followUseCase.Unfollow(userID, username); err != nil {
		handleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, gin.H{"message": "unfollowed"})
}

// GetFollowers godoc
// @Summary Get followers
// @Description Get a list of users following the specified user
// @Tags Follows
// @Produce json
// @Param username path string true "Username"
// @Param page query int false "Page number"
// @Param limit query int false "Results per page"
// @Success 200 {object} response.PaginatedResponse
// @Failure 404 {object} response.APIResponse
// @Router /users/{username}/followers [get]
func (h *FollowHandler) GetFollowers(c *gin.Context) {
	username := c.Param("username")
	page, limit := getPagination(c)

	followers, total, err := h.followUseCase.GetFollowers(username, page, limit)
	if err != nil {
		handleError(c, err)
		return
	}

	response.Paginated(c, followers, page, limit, total)
}

// GetFollowing godoc
// @Summary Get following
// @Description Get a list of users the specified user is following
// @Tags Follows
// @Produce json
// @Param username path string true "Username"
// @Param page query int false "Page number"
// @Param limit query int false "Results per page"
// @Success 200 {object} response.PaginatedResponse
// @Failure 404 {object} response.APIResponse
// @Router /users/{username}/following [get]
func (h *FollowHandler) GetFollowing(c *gin.Context) {
	username := c.Param("username")
	page, limit := getPagination(c)

	following, total, err := h.followUseCase.GetFollowing(username, page, limit)
	if err != nil {
		handleError(c, err)
		return
	}

	response.Paginated(c, following, page, limit, total)
}
