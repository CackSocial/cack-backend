package handler

import (
	"net/http"

	"github.com/CackSocial/cack-backend/internal/usecase/tag"
	"github.com/CackSocial/cack-backend/pkg/response"
	"github.com/gin-gonic/gin"
)

type TagHandler struct {
	tagUseCase *tag.TagUseCase
}

func NewTagHandler(uc *tag.TagUseCase) *TagHandler {
	return &TagHandler{tagUseCase: uc}
}

func (h *TagHandler) RegisterRoutes(public *gin.RouterGroup, optionalAuth gin.HandlerFunc) {
	public.GET("/tags/trending", h.GetTrending)
	public.GET("/tags/:name/posts", optionalAuth, h.GetPostsByTag)
}

// GetTrending godoc
// @Summary Get trending tags
// @Description Get the top trending tags
// @Tags Tags
// @Produce json
// @Success 200 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /tags/trending [get]
func (h *TagHandler) GetTrending(c *gin.Context) {
	tags, err := h.tagUseCase.GetTrending(10)
	if err != nil {
		handleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, tags)
}

// GetPostsByTag godoc
// @Summary Get posts by tag
// @Description Get all posts associated with a specific tag
// @Tags Tags
// @Produce json
// @Param name path string true "Tag name"
// @Param page query int false "Page number"
// @Param limit query int false "Results per page"
// @Success 200 {object} response.PaginatedResponse
// @Failure 404 {object} response.APIResponse
// @Router /tags/{name}/posts [get]
func (h *TagHandler) GetPostsByTag(c *gin.Context) {
	tagName := c.Param("name")
	currentUserID := getUserID(c)
	page, limit := getPagination(c)

	posts, total, err := h.tagUseCase.GetPostsByTag(tagName, currentUserID, page, limit)
	if err != nil {
		handleError(c, err)
		return
	}

	response.Paginated(c, posts, page, limit, total)
}
