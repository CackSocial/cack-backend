package handler

import (
	"github.com/CackSocial/cack-backend/internal/usecase/timeline"
	"github.com/CackSocial/cack-backend/pkg/response"
	"github.com/gin-gonic/gin"
)

type TimelineHandler struct {
	timelineUseCase *timeline.TimelineUseCase
}

func NewTimelineHandler(uc *timeline.TimelineUseCase) *TimelineHandler {
	return &TimelineHandler{timelineUseCase: uc}
}

func (h *TimelineHandler) RegisterRoutes(protected *gin.RouterGroup) {
	protected.GET("/timeline", h.GetFeed)
}

// GetFeed godoc
// @Summary Get timeline feed
// @Description Get the authenticated user's timeline feed with posts from followed users
// @Tags Timeline
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Results per page"
// @Success 200 {object} response.PaginatedResponse
// @Failure 401 {object} response.APIResponse
// @Security BearerAuth
// @Router /timeline [get]
func (h *TimelineHandler) GetFeed(c *gin.Context) {
	userID := getUserID(c)
	page, limit := getPagination(c)

	posts, total, err := h.timelineUseCase.GetFeed(userID, page, limit)
	if err != nil {
		handleError(c, err)
		return
	}

	response.Paginated(c, posts, page, limit, total)
}
