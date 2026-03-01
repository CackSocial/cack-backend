package handler

import (
	"net/http"

	"github.com/CackSocial/cack-backend/internal/usecase/notification"
	"github.com/CackSocial/cack-backend/pkg/response"
	"github.com/gin-gonic/gin"
)

type NotificationHandler struct {
	notifUseCase *notification.NotificationUseCase
}

func NewNotificationHandler(uc *notification.NotificationUseCase) *NotificationHandler {
	return &NotificationHandler{notifUseCase: uc}
}

func (h *NotificationHandler) RegisterNotificationRoutes(protected *gin.RouterGroup) {
	protected.GET("/notifications", h.GetNotifications)
	protected.PUT("/notifications/read-all", h.MarkAllAsRead)
	protected.PUT("/notifications/:id/read", h.MarkAsRead)
	protected.GET("/notifications/unread-count", h.GetUnreadCount)
}

// GetNotifications godoc
// @Summary Get notifications
// @Description Get a paginated list of notifications for the authenticated user
// @Tags Notifications
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Results per page"
// @Success 200 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Security BearerAuth
// @Router /notifications [get]
func (h *NotificationHandler) GetNotifications(c *gin.Context) {
	userID := getUserID(c)
	page, limit := getPagination(c)

	notifications, err := h.notifUseCase.GetNotifications(userID, page, limit)
	if err != nil {
		handleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, notifications)
}

// MarkAsRead godoc
// @Summary Mark notification as read
// @Description Mark a single notification as read
// @Tags Notifications
// @Produce json
// @Param id path string true "Notification ID"
// @Success 200 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Security BearerAuth
// @Router /notifications/{id}/read [put]
func (h *NotificationHandler) MarkAsRead(c *gin.Context) {
	userID := getUserID(c)
	id := c.Param("id")

	if err := h.notifUseCase.MarkAsRead(id, userID); err != nil {
		handleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, gin.H{"message": "notification marked as read"})
}

// MarkAllAsRead godoc
// @Summary Mark all notifications as read
// @Description Mark all notifications as read for the authenticated user
// @Tags Notifications
// @Produce json
// @Success 200 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Security BearerAuth
// @Router /notifications/read-all [put]
func (h *NotificationHandler) MarkAllAsRead(c *gin.Context) {
	userID := getUserID(c)

	if err := h.notifUseCase.MarkAllAsRead(userID); err != nil {
		handleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, gin.H{"message": "all notifications marked as read"})
}

// GetUnreadCount godoc
// @Summary Get unread notification count
// @Description Get the number of unread notifications for the authenticated user
// @Tags Notifications
// @Produce json
// @Success 200 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Security BearerAuth
// @Router /notifications/unread-count [get]
func (h *NotificationHandler) GetUnreadCount(c *gin.Context) {
	userID := getUserID(c)

	count, err := h.notifUseCase.GetUnreadCount(userID)
	if err != nil {
		handleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, gin.H{"count": count})
}
