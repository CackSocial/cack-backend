package handler

import (
	"encoding/json"
	"net/http"

	"github.com/CackSocial/cack-backend/internal/dto"
	"github.com/CackSocial/cack-backend/internal/usecase/message"
	"github.com/CackSocial/cack-backend/pkg/response"
	"github.com/gin-gonic/gin"
)

// MessagePusher abstracts the WebSocket hub for real-time message delivery.
type MessagePusher interface {
	SendToUser(userID string, data []byte)
}

type MessageHandler struct {
	messageUseCase *message.MessageUseCase
	hub            MessagePusher
}

func NewMessageHandler(uc *message.MessageUseCase, hub MessagePusher) *MessageHandler {
	return &MessageHandler{messageUseCase: uc, hub: hub}
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

	// Push the new message to both parties in real-time via WebSocket.
	if h.hub != nil {
		wsPayload := map[string]interface{}{
			"type":        "message",
			"id":          resp.ID,
			"sender_id":   resp.SenderID,
			"receiver_id": resp.ReceiverID,
			"content":     resp.Content,
			"image_url":   resp.ImageURL,
			"created_at":  resp.CreatedAt,
		}
		if data, err := json.Marshal(wsPayload); err == nil {
			h.hub.SendToUser(resp.ReceiverID, data)
			h.hub.SendToUser(userID, data)
		}
	}

	response.Success(c, http.StatusCreated, resp)
}
