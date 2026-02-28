package handler

import (
	"net/http"

	"github.com/CackSocial/cack-backend/internal/dto"
	"github.com/CackSocial/cack-backend/internal/usecase/message"
	"github.com/CackSocial/cack-backend/pkg/response"
	"github.com/gin-gonic/gin"
)

type MessageHandler struct {
	messageUseCase *message.MessageUseCase
}

func NewMessageHandler(uc *message.MessageUseCase) *MessageHandler {
	return &MessageHandler{messageUseCase: uc}
}

func (h *MessageHandler) RegisterRoutes(protected *gin.RouterGroup) {
	protected.GET("/messages/conversations", h.GetConversations)
	protected.GET("/messages/:username", h.GetConversation)
	protected.POST("/messages/:username", h.Send)
}

// GetConversations godoc
// @Summary Get conversations
// @Description Get all conversations for the authenticated user
// @Tags Messages
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Results per page"
// @Success 200 {object} response.PaginatedResponse
// @Failure 401 {object} response.APIResponse
// @Security BearerAuth
// @Router /messages/conversations [get]
func (h *MessageHandler) GetConversations(c *gin.Context) {
	userID := getUserID(c)
	page, limit := getPagination(c)

	conversations, total, err := h.messageUseCase.GetConversations(userID, page, limit)
	if err != nil {
		handleError(c, err)
		return
	}

	response.Paginated(c, conversations, page, limit, total)
}

// GetConversation godoc
// @Summary Get conversation with user
// @Description Get messages in a conversation with a specific user
// @Tags Messages
// @Produce json
// @Param username path string true "Username"
// @Param page query int false "Page number"
// @Param limit query int false "Results per page"
// @Success 200 {object} response.PaginatedResponse
// @Failure 401 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Security BearerAuth
// @Router /messages/{username} [get]
func (h *MessageHandler) GetConversation(c *gin.Context) {
	userID := getUserID(c)
	username := c.Param("username")
	page, limit := getPagination(c)

	messages, total, err := h.messageUseCase.GetConversation(userID, username, page, limit)
	if err != nil {
		handleError(c, err)
		return
	}

	response.Paginated(c, messages, page, limit, total)
}

// Send godoc
// @Summary Send a message
// @Description Send a message to another user with content and optional image
// @Tags Messages
// @Accept multipart/form-data
// @Produce json
// @Param username path string true "Recipient username"
// @Param content formData string true "Message content"
// @Param image formData file false "Image file"
// @Success 201 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Security BearerAuth
// @Router /messages/{username} [post]
func (h *MessageHandler) Send(c *gin.Context) {
	var req dto.SendMessageRequest
	if err := c.ShouldBind(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	userID := getUserID(c)
	username := c.Param("username")

	resp, err := h.messageUseCase.Send(userID, username, &req)
	if err != nil {
		handleError(c, err)
		return
	}

	response.Success(c, http.StatusCreated, resp)
}
