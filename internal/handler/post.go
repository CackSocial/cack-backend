package handler

import (
	"net/http"

	"github.com/CackSocial/cack-backend/internal/dto"
	"github.com/CackSocial/cack-backend/internal/usecase/post"
	"github.com/CackSocial/cack-backend/pkg/response"
	"github.com/gin-gonic/gin"
)

type PostHandler struct {
	postUseCase *post.PostUseCase
}

func NewPostHandler(uc *post.PostUseCase) *PostHandler {
	return &PostHandler{postUseCase: uc}
}

func (h *PostHandler) RegisterRoutes(public, protected *gin.RouterGroup, optionalAuth gin.HandlerFunc) {
	protected.POST("/posts", h.Create)
	public.GET("/posts/:id", optionalAuth, h.GetByID)
	protected.DELETE("/posts/:id", h.Delete)
	public.GET("/users/:username/posts", optionalAuth, h.GetByUserID)
}

// Create godoc
// @Summary Create a new post
// @Description Create a new post with content and optional image
// @Tags Posts
// @Accept multipart/form-data
// @Produce json
// @Param content formData string true "Post content"
// @Param image formData file false "Image file"
// @Success 201 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Security BearerAuth
// @Router /posts [post]
func (h *PostHandler) Create(c *gin.Context) {
	var req dto.CreatePostRequest
	if err := c.ShouldBind(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	userID := getUserID(c)
	resp, err := h.postUseCase.Create(userID, &req)
	if err != nil {
		handleError(c, err)
		return
	}

	response.Success(c, http.StatusCreated, resp)
}

// GetByID godoc
// @Summary Get post by ID
// @Description Get a single post by its ID
// @Tags Posts
// @Produce json
// @Param id path string true "Post ID"
// @Success 200 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Router /posts/{id} [get]
func (h *PostHandler) GetByID(c *gin.Context) {
	postID := c.Param("id")
	currentUserID := getUserID(c)

	resp, err := h.postUseCase.GetByID(postID, currentUserID)
	if err != nil {
		handleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, resp)
}

// Delete godoc
// @Summary Delete a post
// @Description Delete a post by its ID (must be the author)
// @Tags Posts
// @Produce json
// @Param id path string true "Post ID"
// @Success 200 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Security BearerAuth
// @Router /posts/{id} [delete]
func (h *PostHandler) Delete(c *gin.Context) {
	postID := c.Param("id")
	userID := getUserID(c)

	if err := h.postUseCase.Delete(postID, userID); err != nil {
		handleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, gin.H{"message": "post deleted"})
}

// GetByUserID godoc
// @Summary Get posts by user
// @Description Get all posts by a specific user
// @Tags Posts
// @Produce json
// @Param username path string true "Username"
// @Param page query int false "Page number"
// @Param limit query int false "Results per page"
// @Success 200 {object} response.PaginatedResponse
// @Failure 404 {object} response.APIResponse
// @Router /users/{username}/posts [get]
func (h *PostHandler) GetByUserID(c *gin.Context) {
	username := c.Param("username")
	currentUserID := getUserID(c)
	page, limit := getPagination(c)

	posts, total, err := h.postUseCase.GetByUserID(username, currentUserID, page, limit)
	if err != nil {
		handleError(c, err)
		return
	}

	response.Paginated(c, posts, page, limit, total)
}
