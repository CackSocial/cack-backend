package handler

import (
	"net/http"

	"github.com/CackSocial/cack-backend/internal/usecase/like"
	"github.com/CackSocial/cack-backend/pkg/response"
	"github.com/gin-gonic/gin"
)

type LikeHandler struct {
	likeUseCase *like.LikeUseCase
}

func NewLikeHandler(uc *like.LikeUseCase) *LikeHandler {
	return &LikeHandler{likeUseCase: uc}
}

func (h *LikeHandler) RegisterRoutes(public, protected *gin.RouterGroup, optionalAuth gin.HandlerFunc) {
	protected.POST("/posts/:id/like", h.Like)
	protected.DELETE("/posts/:id/like", h.Unlike)
	public.GET("/posts/:id/likes", h.GetPostLikes)
	public.GET("/users/:username/likes", optionalAuth, h.GetUserLikedPosts)
}

// Like godoc
// @Summary Like a post
// @Description Like a post by its ID
// @Tags Likes
// @Produce json
// @Param id path string true "Post ID"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Security BearerAuth
// @Router /posts/{id}/like [post]
func (h *LikeHandler) Like(c *gin.Context) {
	userID := getUserID(c)
	postID := c.Param("id")

	if err := h.likeUseCase.Like(userID, postID); err != nil {
		handleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, gin.H{"message": "liked"})
}

// Unlike godoc
// @Summary Unlike a post
// @Description Remove a like from a post by its ID
// @Tags Likes
// @Produce json
// @Param id path string true "Post ID"
// @Success 200 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Security BearerAuth
// @Router /posts/{id}/like [delete]
func (h *LikeHandler) Unlike(c *gin.Context) {
	userID := getUserID(c)
	postID := c.Param("id")

	if err := h.likeUseCase.Unlike(userID, postID); err != nil {
		handleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, gin.H{"message": "unliked"})
}

// GetPostLikes godoc
// @Summary Get post likes
// @Description Get a list of users who liked a post
// @Tags Likes
// @Produce json
// @Param id path string true "Post ID"
// @Param page query int false "Page number"
// @Param limit query int false "Results per page"
// @Success 200 {object} response.PaginatedResponse
// @Failure 404 {object} response.APIResponse
// @Router /posts/{id}/likes [get]
func (h *LikeHandler) GetPostLikes(c *gin.Context) {
	postID := c.Param("id")
	page, limit := getPagination(c)

	users, total, err := h.likeUseCase.GetPostLikes(postID, page, limit)
	if err != nil {
		handleError(c, err)
		return
	}

	response.Paginated(c, users, page, limit, total)
}

// GetUserLikedPosts godoc
// @Summary Get posts liked by a user
// @Description Get a paginated list of posts that a user has liked
// @Tags Likes
// @Produce json
// @Param username path string true "Username"
// @Param page query int false "Page number"
// @Param limit query int false "Results per page"
// @Success 200 {object} response.PaginatedResponse
// @Failure 404 {object} response.APIResponse
// @Router /users/{username}/likes [get]
func (h *LikeHandler) GetUserLikedPosts(c *gin.Context) {
	username := c.Param("username")
	currentUserID := getUserID(c)
	page, limit := getPagination(c)

	posts, total, err := h.likeUseCase.GetLikedPosts(username, currentUserID, page, limit)
	if err != nil {
		handleError(c, err)
		return
	}

	response.Paginated(c, posts, page, limit, total)
}
