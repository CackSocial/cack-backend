package handler

import (
	"net/http"

	"github.com/CackSocial/cack-backend/internal/usecase/bookmark"
	"github.com/CackSocial/cack-backend/pkg/response"
	"github.com/gin-gonic/gin"
)

type BookmarkHandler struct {
	bookmarkUseCase *bookmark.BookmarkUseCase
}

func NewBookmarkHandler(uc *bookmark.BookmarkUseCase) *BookmarkHandler {
	return &BookmarkHandler{bookmarkUseCase: uc}
}

func (h *BookmarkHandler) RegisterRoutes(protected *gin.RouterGroup) {
	protected.POST("/posts/:id/bookmark", h.Bookmark)
	protected.DELETE("/posts/:id/bookmark", h.Unbookmark)
	protected.GET("/bookmarks", h.GetBookmarks)
}

// Bookmark godoc
// @Summary Bookmark a post
// @Description Add a post to the authenticated user's bookmarks
// @Tags Bookmarks
// @Produce json
// @Param id path string true "Post ID"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Security BearerAuth
// @Router /posts/{id}/bookmark [post]
func (h *BookmarkHandler) Bookmark(c *gin.Context) {
	userID := getUserID(c)
	postID := c.Param("id")

	if err := h.bookmarkUseCase.Bookmark(userID, postID); err != nil {
		handleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, gin.H{"message": "bookmarked"})
}

// Unbookmark godoc
// @Summary Remove bookmark from a post
// @Description Remove a post from the authenticated user's bookmarks
// @Tags Bookmarks
// @Produce json
// @Param id path string true "Post ID"
// @Success 200 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Security BearerAuth
// @Router /posts/{id}/bookmark [delete]
func (h *BookmarkHandler) Unbookmark(c *gin.Context) {
	userID := getUserID(c)
	postID := c.Param("id")

	if err := h.bookmarkUseCase.Unbookmark(userID, postID); err != nil {
		handleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, gin.H{"message": "unbookmarked"})
}

// GetBookmarks godoc
// @Summary Get bookmarked posts
// @Description Get a paginated list of the authenticated user's bookmarked posts
// @Tags Bookmarks
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Results per page"
// @Success 200 {object} response.PaginatedResponse
// @Failure 401 {object} response.APIResponse
// @Security BearerAuth
// @Router /bookmarks [get]
func (h *BookmarkHandler) GetBookmarks(c *gin.Context) {
	userID := getUserID(c)
	page, limit := getPagination(c)

	posts, total, err := h.bookmarkUseCase.GetBookmarks(userID, page, limit)
	if err != nil {
		handleError(c, err)
		return
	}

	response.Paginated(c, posts, page, limit, total)
}
