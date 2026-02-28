package handler

import (
	"net/http"

	"github.com/CackSocial/cack-backend/internal/dto"
	"github.com/CackSocial/cack-backend/internal/usecase/comment"
	"github.com/CackSocial/cack-backend/pkg/response"
	"github.com/gin-gonic/gin"
)

type CommentHandler struct {
	commentUseCase *comment.CommentUseCase
}

func NewCommentHandler(uc *comment.CommentUseCase) *CommentHandler {
	return &CommentHandler{commentUseCase: uc}
}

func (h *CommentHandler) RegisterRoutes(public, protected *gin.RouterGroup) {
	protected.POST("/posts/:id/comments", h.Create)
	public.GET("/posts/:id/comments", h.GetByPostID)
	protected.DELETE("/comments/:id", h.Delete)
}

// Create godoc
// @Summary Create a comment
// @Description Add a comment to a post
// @Tags Comments
// @Accept json
// @Produce json
// @Param id path string true "Post ID"
// @Param body body dto.CreateCommentRequest true "Comment request"
// @Success 201 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Security BearerAuth
// @Router /posts/{id}/comments [post]
func (h *CommentHandler) Create(c *gin.Context) {
	var req dto.CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	userID := getUserID(c)
	postID := c.Param("id")

	resp, err := h.commentUseCase.Create(userID, postID, &req)
	if err != nil {
		handleError(c, err)
		return
	}

	response.Success(c, http.StatusCreated, resp)
}

// GetByPostID godoc
// @Summary Get comments by post
// @Description Get all comments for a specific post
// @Tags Comments
// @Produce json
// @Param id path string true "Post ID"
// @Param page query int false "Page number"
// @Param limit query int false "Results per page"
// @Success 200 {object} response.PaginatedResponse
// @Failure 404 {object} response.APIResponse
// @Router /posts/{id}/comments [get]
func (h *CommentHandler) GetByPostID(c *gin.Context) {
	postID := c.Param("id")
	page, limit := getPagination(c)

	comments, total, err := h.commentUseCase.GetByPostID(postID, page, limit)
	if err != nil {
		handleError(c, err)
		return
	}

	response.Paginated(c, comments, page, limit, total)
}

// Delete godoc
// @Summary Delete a comment
// @Description Delete a comment by its ID (must be the author)
// @Tags Comments
// @Produce json
// @Param id path string true "Comment ID"
// @Success 200 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Security BearerAuth
// @Router /comments/{id} [delete]
func (h *CommentHandler) Delete(c *gin.Context) {
	commentID := c.Param("id")
	userID := getUserID(c)

	if err := h.commentUseCase.Delete(commentID, userID); err != nil {
		handleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, gin.H{"message": "comment deleted"})
}
